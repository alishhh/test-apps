# environments for service

# envs for setting up service mode:
# LOCAL, DEV, PROD
export LOG_LEVEL="dev"

# app envs
export APP_NAME="auth-app"
export APP_HOST="localhost"
export APP_PORT="8084"

# database envs
export DB_PATH="./storage/storage.db"

# client envs
export CLIENT_URL="http://localhost:8083"

# envs for timeouts
export TIMEOUT="4s"
export IDLE_TIMEOUT="120s"

# command for start service
go run ./cmd/url-shortener-app