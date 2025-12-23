package config



import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	MONGO_CONNECTION_STRING string
	

	USER_DB              string
	USER_COLLECTION_NAME string

	

	USER_OTP_COLLECTION_NAME           string
	USER_REFRESH_TOKEN_COLLECTION_NAME string

	JWTSECRET        string
	JWTREFRESHSECRET string
	CURR_USER        string

	FROM       string
	APPPASS    string
	SMTPSERVER string
	SMTPPORT   string
	SMTPUSER   string

	LOG_ENC_KEY string

	MIN_PASSWORD_LENGTH   int
	LOCKOUT_THRESHOLD     int
	LOCKOUT_WINDOW_MIN    int
	LOCKOUT_DURATION_MIN  int
	BACKUP_PATH           string

	CLIENT_ID           string
	CLIENT_SECRET       string
	CLIENT_CALLBACK_URL string

	// PORT 				string
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can not load .env file")
	}
	MONGO_CONNECTION_STRING = getEnv("MONGO_CONNECTION_STRING")
	
	USER_DB = getEnv("USER_DB")
	USER_COLLECTION_NAME = getEnv("USER_COLLECTION_NAME")

	JWTSECRET = getEnv("JWTSECRET")
	
	
	FROM = getEnv("FROM")
	APPPASS = getEnv("APPPASS")
	SMTPSERVER = getEnv("SMTPSERVER")
	SMTPPORT = getEnv("SMTPPORT")
	SMTPUSER = getEnv("SMTPUSER")
	CLIENT_ID = getEnv("CLIENT_ID")
	CLIENT_SECRET = getEnv("CLIENT_SECRET")
	CLIENT_CALLBACK_URL = getEnv("CLIENT_CALLBACK_URL")
	USER_OTP_COLLECTION_NAME = getEnv("USER_OTP_COLLECTION_NAME")
	USER_REFRESH_TOKEN_COLLECTION_NAME = getEnv("USER_REFRESH_TOKEN_COLLECTION_NAME")
	JWTREFRESHSECRET = getEnv("JWTREFRESHSECRET")
	LOG_ENC_KEY = getEnv("LOG_ENC_KEY")
	MIN_PASSWORD_LENGTH = getEnvInt("MIN_PASSWORD_LENGTH", 8)
	LOCKOUT_THRESHOLD = getEnvInt("LOCKOUT_THRESHOLD", 5)
	LOCKOUT_WINDOW_MIN = getEnvInt("LOCKOUT_WINDOW_MIN", 15)
	LOCKOUT_DURATION_MIN = getEnvInt("LOCKOUT_DURATION_MIN", 15)
	BACKUP_PATH = getEnv("BACKUP_PATH")
	
	// PORT = getEnv("PORT")
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Environment variable %s must be int", key)
	}
	return parsed
}
