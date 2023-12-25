# environments for service

# envs for setting up service mode:
# LOCAL, DEV, PROD
export LOG_LEVEL="prod"

# app envs
export APP_NAME="auth-app"
export APP_HOST="localhost"
export APP_PORT="8081"

# envs for timeouts
export TIMEOUT="4s"
export IDLE_TIMEOUT="120s"

# keycloak
export KEYCLOAK_HOST="http://localhost:8082"
export KEYCLOAK_REALM="demo"
export KEYCLOAK_CLIENT_ID="go-client-test"
export KEYCLOAK_CLIENT_SECRET="1VGIiII5tCz0SxMGEFfBjMpc3cYd24Ll"

# command for start service
go run ./cmd/auth-app