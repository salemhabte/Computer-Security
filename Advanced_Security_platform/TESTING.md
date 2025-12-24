# Quick Test Guide (curl)

Base URL assumes `localhost:8080`.

## Registration & Verification
```bash
curl -X POST http://localhost:8080/registration \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password!234","username":"alice"}'

curl -X POST http://localhost:8080/registration/verification \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","otp":"123456"}'
```

## Login (MFA + Captcha stub)
Step 1: password + captcha (sends OTP to email)
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password!234","captcha_token":"stub"}'
```
Step 2: include OTP from email
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password!234","otp":"123456","captcha_token":"stub"}'
```

## Refresh token
```bash
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"token":"<refresh_token>"}'
```

## Logout (auth required)
```bash
curl -X POST http://localhost:8080/user/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

## Profile update (auth required)
```bash
curl -X PUT http://localhost:8080/user/edit_profile \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"username":"Alice","department":"Payroll","device_trust":"HIGH"}'
```

## Password reset
```bash
curl -X POST http://localhost:8080/forgot_password \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

curl -X PUT http://localhost:8080/reset_password \
  -H "Content-Type: application/json" \
  -d '{"token":"<reset_token>","new_password":"NewPass!234"}'
```

## Admin (requires SUPER_ADMIN access token)
Upsert role:
```bash
curl -X POST http://localhost:8080/admin/role \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"HR_MANAGER","permissions":["approve_leave","document:read"]}'
```

## DAC Permissions (requires authentication - resource owner OR admin)
Grant DAC (resource owners can grant permissions on their own resources, admins can grant on any resource):
```bash
curl -X POST http://localhost:8080/dac/grant \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"resource_id":"doc-123","subject":"user@example.com","actions":["document:read","document:write"]}'
```

Revoke DAC:
```bash
curl -X POST http://localhost:8080/dac/revoke \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"resource_id":"doc-123","subject":"user@example.com"}'
```
Backup:
```bash
curl -X POST http://localhost:8080/admin/backup/run \
  -H "Authorization: Bearer <admin_access_token>"
```

**Note on DAC**: 
- Resource owners (users who created/own the resource) can grant/revoke permissions on their own resources
- Admins (SUPER_ADMIN/ADMIN) can grant/revoke permissions on any resource
- The system checks resource ownership automatically via the resource repository

