# Complete Implementation Details - Computer Security Platform

## Overview
This is a comprehensive security platform implementing multiple access control models, authentication mechanisms, audit logging, and backup systems following clean architecture principles.

---

## 🏗️ Architecture Layers

### 1. **Domain Layer** (`domain/`)
Core business entities and interfaces (no external dependencies)

#### **Entities:**
- **User**: Stores user data including new security attributes:
  - `Department`, `EmploymentStatus`, `DeviceTrust`, `BiometricVerified`
- **Resource**: Protected objects with sensitivity labels (PUBLIC, INTERNAL, CONFIDENTIAL)
- **Role**: Maps role names to permission lists
- **AccessControlEntry**: DAC permissions (who can do what on which resource)
- **AttributeSet**: ABAC attributes (department, location, device trust, biometric status)
- **PolicyDecision**: Result of access control evaluation (allowed/denied + reason)

#### **Key Constants:**
- Roles: `SUPER_ADMIN`, `ADMIN`, `USER`
- Sensitivity Levels: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`

---

## 🔐 Access Control Models Implementation

### **1. Mandatory Access Control (MAC)**

**Location:** `infrastructure/policy_engine.go` - `macCheck()` method

**How it works:**
- Users have clearance levels based on their role
- Resources have sensitivity labels (PUBLIC, INTERNAL, CONFIDENTIAL)
- Access rules:
  - `SUPER_ADMIN`: Can access everything
  - `ADMIN`: Can access PUBLIC and INTERNAL, NOT CONFIDENTIAL
  - `USER`: Can only access PUBLIC and INTERNAL
- Access is **denied by default** if clearance is insufficient

**Example:**
```go
// A CONFIDENTIAL document can only be accessed by SUPER_ADMIN
// ADMIN and USER will be denied with reason "MAC: insufficient clearance"
```

---

### **2. Discretionary Access Control (DAC)**

**Location:** `infrastructure/policy_engine.go` - `Decide()` method (DAC checks)

**How it works:**
1. **Resource Ownership**: If user owns the resource → automatic access
2. **ACL Entries**: Admins can grant specific permissions via `/admin/dac/grant`
3. **Storage**: ACL entries stored in MongoDB (`acls` collection)
4. **Admin Endpoints:**
   - `POST /admin/dac/grant`: Grant permissions
   - `POST /admin/dac/revoke`: Remove permissions

**Example:**
```json
POST /admin/dac/grant
{
  "resource_id": "doc123",
  "subject": "user@example.com",
  "actions": ["document:read", "document:write"]
}
```

**Data Flow:**
- ACL entries stored in `Repository/acl_repo.go`
- Policy engine checks ACL after MAC check
- Each entry tracks: ResourceID, Subject (user email), Actions, GrantedBy, GrantedAt

---

### **3. Role-Based Access Control (RBAC)**

**Location:** `infrastructure/policy_engine.go` - `Decide()` method (RBAC checks)

**How it works:**
1. **Role Definitions**: Roles mapped to permission lists (e.g., "HR_MANAGER" → ["approve_leave", "view_hr_data"])
2. **Role Storage**: MongoDB `roles` collection
3. **Permission Format**: `"resource:action"` (e.g., "document:read", "leave:approve")
4. **Admin Endpoint**: `POST /admin/role` to create/update roles

**Example:**
```json
POST /admin/role
{
  "name": "HR_MANAGER",
  "permissions": ["leave:approve", "hr_data:view", "employee:update"]
}
```

**Evaluation:**
- Policy engine checks if user's role contains the required permission
- Checks happen after MAC and DAC

---

### **4. Rule-Based Access Control (RuBAC)**

**Location:** `infrastructure/policy_engine.go` - `rubacCheck()` method

**Implemented Rules:**
1. **Time-Based Access**: Only allows access between 06:00 - 22:00 (configurable)
2. **Device Trust**: Blocks access if device trust is "LOW" (via ABAC)
3. **Location/IP Rules**: Framework exists, can be extended

**How it works:**
- `RequestContext` contains: IP, UserAgent, Timestamp
- Policy engine evaluates time window before allowing access
- Returns: `PolicyDecision{Allowed: false, Reason: "RuBAC: outside allowed time or context"}`

**Future Extensions:**
- IP whitelist/blacklist
- Geographic restrictions
- Device fingerprinting

---

### **5. Attribute-Based Access Control (ABAC)**

**Location:** `infrastructure/policy_engine.go` - `abacCheck()` method

**Attributes Used:**
- `Department`: User's department (e.g., "Payroll", "IT", "Finance")
- `EmploymentStatus`: Active, Suspended, etc.
- `DeviceTrust`: HIGH, MEDIUM, LOW
- `BiometricVerified`: Boolean flag
- `Role`: User's role
- `Location`: (can be extended)

**Implemented Policies:**
1. **Department Isolation**: 
   - Payroll data only accessible by Payroll department employees
   - Example: `if resource.Department == "Payroll" && user.Department != "Payroll" → DENY`

2. **Role-Based Action Restriction**:
   - Leave approval (>10 days) only by HR_MANAGER
   - Example: `if action == "approve_leave" && attrs.Role != "HR_MANAGER" → DENY`

3. **Device Trust Requirement**:
   - LOW trust devices blocked
   - Example: `if attrs.DeviceTrust == "LOW" → DENY`

4. **Biometric Requirement for Confidential Data**:
   - CONFIDENTIAL resources require `BiometricVerified == true`
   - Example: `if resource.Label == CONFIDENTIAL && !attrs.BiometricVerified → DENY`

**Dynamic Evaluation:**
- All attributes evaluated in real-time per request
- Can combine multiple attributes (e.g., "Manager in Finance Department during working hours")

---

## 🔒 Authentication & Authorization

### **1. Password Authentication**

**Location:** `infrastructure/passwordService.go`, `usecase/userUsecase.go`

**Features:**
- **Password Hashing**: Bcrypt with default cost (protects against rainbow tables)
- **Password Policy**: Configurable via `MIN_PASSWORD_LENGTH` (default 8)
  - Requirements: uppercase, lowercase, number, special character
  - Guidance message provided to users
- **Secure Transmission**: Passwords sent over HTTPS (assumed in production)
- **Password Change**: Via profile update endpoint (`PUT /user/edit_profile`)

**Validation Flow:**
1. Check password strength during registration
2. Hash with bcrypt before storage
3. Compare using `bcrypt.CompareHashAndPassword` on login

---

### **2. Account Lockout Policy**

**Location:** `usecase/userUsecase.go` - `markFail()`, `isLocked()`, `resetFails()`

**How it works:**
- **In-Memory Tracking**: `failedLogins` map stores timestamps of failed attempts
- **Configurable Threshold**: `LOCKOUT_THRESHOLD` (default: 5 attempts)
- **Time Window**: `LOCKOUT_WINDOW_MIN` (default: 15 minutes)
- **Lock Duration**: `LOCKOUT_DURATION_MIN` (default: 15 minutes)

**Flow:**
1. Failed login → add timestamp to `failedLogins[email]`
2. Clean old attempts outside the window
3. If count ≥ threshold → lock account until `LOCKOUT_DURATION_MIN` passes
4. Successful login → reset counter

**Thread-Safe**: Uses `sync.Mutex` for concurrent access protection

**Example:**
- User fails login 5 times in 15 minutes → account locked for 15 minutes
- After 15 minutes, user can try again
- If they succeed, counter resets

---

### **3. Token-Based Authentication**

**Location:** `infrastructure/jwt_service.go`, `usecase/userUsecase.go`

**Features:**
- **JWT Access Tokens**: Short-lived (15 minutes), contains: user_id, email, role
- **Refresh Tokens**: Long-lived (7 days), stored in MongoDB
- **Token Rotation**: Old refresh token deleted when new one issued
- **Secure Storage**: Refresh tokens stored in database, not cookies

**Token Claims:**
```json
{
  "user_id": "...",
  "email": "user@example.com",
  "role": "USER",
  "exp": 1234567890
}
```

**Flow:**
1. Login → Generate access + refresh tokens
2. Store refresh token in DB
3. Client stores both tokens
4. Access token expires → Use refresh token to get new tokens
5. Logout → Delete refresh token from DB

---

### **4. Multi-Factor Authentication (MFA)**

**Location:** `usecase/userUsecase.go` - `Login()`, `storeMFA()`, `validateMFA()`

**Implementation: Email OTP**

**Flow:**
1. **First Request**: User sends email, password, captcha_token (NO OTP)
   - Password validated
   - 6-digit OTP generated and emailed
   - OTP stored in-memory with 5-minute expiry
   - Response: `202 Accepted` with message "OTP sent to email"

2. **Second Request**: User sends email, password, captcha_token, AND OTP
   - Password validated again
   - OTP validated against stored value
   - If valid → tokens generated and returned
   - If invalid → fail counter incremented (lockout possible)

**Storage:**
- In-memory map `mfaPending` with expiry timestamps
- Thread-safe with mutex

**Future Extensions:**
- SMS OTP
- TOTP (Google Authenticator, Authy)
- Hardware tokens

---

### **5. CAPTCHA Integration**

**Location:** `infrastructure/captcha_service.go`

**Current Implementation:**
- **Placeholder**: Accepts any non-empty token
- **Hook for Real Provider**: Easy to integrate reCAPTCHA, hCaptcha, etc.

**Usage:**
- Required on registration and login
- Prevents automated bot attacks
- Returns `false` if token is empty

**Future:**
```go
func (c *CaptchaValidator) Validate(token string) bool {
    // Call Google reCAPTCHA API
    resp, err := http.Post("https://www.google.com/recaptcha/api/siteverify", ...)
    // Parse response and return result
}
```

---

## 📊 Audit Logging & Monitoring

### **1. Encrypted Audit Logging**

**Location:** `infrastructure/audit_logger.go`

**Features:**
- **Encryption**: AES-256-GCM using `LOG_ENC_KEY` (32 bytes)
- **Log Format**: JSON lines (one event per line)
- **Encryption Method**: 
  - Generate random nonce per log entry
  - Encrypt JSON payload with AES-GCM
  - Append nonce + ciphertext to file

**Audit Event Structure:**
```json
{
  "user_email": "user@example.com",
  "action": "document:read",
  "resource": "doc123",
  "result": "allow|deny",
  "reason": "RBAC: role permit",
  "ip": "192.168.1.1",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**What Gets Logged:**
- All policy decisions (allow/deny with reason)
- User actions (via policy middleware)
- Failed login attempts (implicit via lockout)
- Admin actions (role changes, DAC grants/revokes)

**File Storage:**
- Default: `audit.log` in project root
- Encrypted, readable only with `LOG_ENC_KEY`

---

### **2. Policy Middleware Logging**

**Location:** `infrastructure/policy_middleware.go`

**How it works:**
1. Middleware wraps protected routes
2. Extracts user info from JWT (email, role)
3. Calls policy engine to evaluate access
4. Logs decision to audit logger
5. Returns 403 if denied, continues if allowed

**Usage:**
```go
router.Use(PolicyMiddleware(policyEngine, auditLogger, resourceProvider, "document:read"))
```

**Logged Information:**
- User email
- Action attempted
- Resource ID (if applicable)
- Decision (allow/deny)
- Reason (which policy allowed/denied)
- IP address
- Timestamp

---

### **3. Centralized Logging**

**Current State:**
- Logs written to local encrypted file
- **Future Extension**: Can send to centralized log aggregation service (ELK, Splunk, etc.)

**Alerting Mechanism:**
- Framework exists in audit logger
- Can trigger alerts on:
  - Multiple failed logins (lockout)
  - Policy denials for sensitive resources
  - Unusual access patterns
- **TODO**: Implement actual alerting service (email, Slack, etc.)

---

## 💾 Data Backups

**Location:** `infrastructure/backup_service.go`, `Delivery/Controller/backup_controller.go`

**Features:**
- **Manual Backup Trigger**: `POST /admin/backup/run` (requires SUPER_ADMIN)
- **Timestamped Files**: Format: `backup-YYYYMMDD-HHMMSS.txt`
- **Storage Location**: Configurable via `BACKUP_PATH` env variable

**Current Implementation:**
- **Placeholder**: Writes text file with timestamp
- **Future**: Should dump MongoDB collections to backup file

**Example Backup File:**
```
backup-20240101-120000.txt
Content: "backup placeholder at 2024-01-01 12:00:00"
```

**Future Enhancements:**
- Automated scheduled backups (cron-like)
- Full MongoDB dump
- Incremental backups
- Backup encryption
- Cloud storage integration (S3, Azure Blob)

---

## 📝 User Registration & Profile Management

### **Registration Flow**

**Location:** `usecase/userUsecase.go` - `HandleRegistration()`, `VerifyOTP()`

**Steps:**
1. User submits registration (email, password, username, etc.)
2. Validate email format and password strength
3. Check if user already exists
4. Hash password with bcrypt
5. Generate 6-digit OTP
6. Email OTP to user
7. Store unverified user data in MongoDB (with OTP, expires in 5 minutes)
8. User submits OTP via `/registration/verification`
9. Validate OTP (correct and not expired)
10. Create verified user in database
11. Delete OTP entry

**Security Features:**
- Email verification required
- OTP expiration (5 minutes)
- Password strength validation
- Duplicate email check

---

### **Profile Management**

**Location:** `usecase/userUsecase.go` - `UpdateProfile()`

**Features:**
- Update username, bio, profile picture, phone, telegram handle
- Secure password update (hashed before storage)
- Requires authentication (JWT middleware)
- Updates stored in MongoDB

**Endpoint:** `PUT /user/edit_profile`

---

## 🔄 Session Management

**Location:** `usecase/userUsecase.go` - `Login()`, `Refresh()`, `Logout()`

**Features:**
- **Refresh Token Rotation**: Old token deleted when new one issued
- **Token Revocation**: Logout deletes refresh token from DB
- **Session Persistence**: Refresh tokens stored in MongoDB (not stateless)
- **Expiry Management**: 
  - Access token: 15 minutes
  - Refresh token: 7 days

**Flow:**
1. Login → Access token (15 min) + Refresh token (7 days) stored in DB
2. Access token expires → Call `/refresh` with refresh token
3. New tokens issued, old refresh token deleted (rotation)
4. Logout → Delete refresh token, user must login again

---

## 🛡️ Security Enhancements

### **1. Password Hashing**
- Algorithm: Bcrypt (adaptive hashing)
- Cost: Default (10 rounds, ~100ms per hash)
- Protection: Against rainbow tables (salts generated automatically)

### **2. Secure Password Transmission**
- Assumes HTTPS in production (not implemented in code, deployment responsibility)
- Passwords never logged
- Only hashed values stored

### **3. Account Lockout**
- Prevents brute-force attacks
- Configurable thresholds
- Automatic unlock after duration

### **4. MFA**
- Adds second factor (email OTP)
- Prevents account takeover even if password compromised

### **5. CAPTCHA**
- Prevents automated registration/login attacks
- Placeholder ready for real provider integration

---

## 🔌 API Endpoints Summary

### **Public Endpoints:**
- `POST /registration` - User registration
- `POST /registration/verification` - Verify OTP
- `POST /login` - Login (requires MFA flow)
- `POST /refresh` - Refresh access token
- `POST /forgot_password` - Request password reset
- `PUT /reset_password` - Reset password with token
- `GET /auth/:provider` - OAuth login
- `GET /auth/:provider/callback` - OAuth callback

### **Authenticated Endpoints (require JWT):**
- `POST /user/logout` - Logout
- `PUT /user/edit_profile` - Update profile

### **Admin Endpoints (require SUPER_ADMIN role):**
- `POST /admin/role` - Create/update role
- `POST /admin/dac/grant` - Grant DAC permission
- `POST /admin/dac/revoke` - Revoke DAC permission
- `POST /admin/backup/run` - Trigger backup

---

## 🗄️ Database Collections

**MongoDB Collections:**
1. **users** - User accounts
2. **unverified_users** - Pending registrations (with OTP)
3. **refresh_tokens** - Active refresh tokens
4. **roles** - Role definitions (name + permissions)
5. **acls** - DAC access control entries
6. **audit_log** - Encrypted audit events (file-based, not DB)

---

## 📋 Configuration Variables

See `ENV_SAMPLE.md` for complete list. Key variables:

- **MongoDB**: `MONGO_CONNECTION_STRING`, `USER_DB`, collection names
- **JWT**: `JWTSECRET`, `JWTREFRESHSECRET`
- **SMTP**: Email service credentials
- **Security**: `LOG_ENC_KEY` (32 bytes), `MIN_PASSWORD_LENGTH`
- **Lockout**: `LOCKOUT_THRESHOLD`, `LOCKOUT_WINDOW_MIN`, `LOCKOUT_DURATION_MIN`
- **Backup**: `BACKUP_PATH`

---

## 🏃 Running the System

1. **Setup Environment**: Copy `ENV_SAMPLE.md` to `.env` and fill values
2. **Start MongoDB**: Ensure MongoDB is running
3. **Run Application**: `go run cmd/main.go`
4. **Test**: Use `TESTING.md` for curl examples

---

## 🔮 Future Enhancements

1. **Real CAPTCHA Provider**: Integrate reCAPTCHA/hCaptcha
2. **SMS/TOTP MFA**: Add alternative MFA methods
3. **Centralized Logging**: Send logs to ELK/Splunk
4. **Alerting Service**: Email/Slack alerts on security events
5. **Real Backups**: MongoDB dump functionality
6. **Biometric Integration**: Actual fingerprint/face recognition
7. **Device Fingerprinting**: Track and verify devices
8. **IP Geolocation**: Location-based access rules
9. **Rate Limiting**: Prevent API abuse
10. **Password History**: Prevent password reuse

---

## 📚 Architecture Benefits

**Clean Architecture Principles:**
- **Separation of Concerns**: Domain, UseCase, Infrastructure, Delivery layers
- **Dependency Inversion**: Interfaces in domain, implementations in infrastructure
- **Testability**: Easy to mock dependencies
- **Maintainability**: Clear structure, easy to extend

**Security Best Practices:**
- Defense in depth (multiple access control models)
- Least privilege (default deny)
- Audit everything (encrypted logs)
- Secure defaults (strong passwords, MFA, lockout)

---

This implementation provides a comprehensive security platform covering all major access control models, authentication mechanisms, and security best practices.

