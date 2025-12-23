# Environment Sample

Copy these into a local `.env` (adjust values):

```
MONGO_CONNECTION_STRING=mongodb://localhost:27017
USER_DB=securitydb
USER_COLLECTION_NAME=users
USER_OTP_COLLECTION_NAME=user_otps
USER_REFRESH_TOKEN_COLLECTION_NAME=refresh_tokens

JWTSECRET=change_me_access_secret
JWTREFRESHSECRET=change_me_refresh_secret

FROM=you@example.com
APPPASS=your_smtp_app_password
SMTPSERVER=smtp.example.com
SMTPPORT=587
SMTPUSER=you@example.com

CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
CLIENT_CALLBACK_URL=http://localhost:8080/auth/google/callback

LOG_ENC_KEY=32_bytes_hex_or_plain_key_____________

MIN_PASSWORD_LENGTH=12
LOCKOUT_THRESHOLD=5
LOCKOUT_WINDOW_MIN=15
LOCKOUT_DURATION_MIN=15

BACKUP_PATH=./backups
```

