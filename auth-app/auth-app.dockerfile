FROM golang:1.21

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install inetutils-ping

ENV LOG_LEVEL="prod"
ENV APP_NAME="auth-app"
ENV APP_HOST="0.0.0.0"
ENV APP_PORT="8080"
ENV TIMEOUT="4s"
ENV IDLE_TIMEOUT="120s"
ENV KEYCLOAK_HOST="http://keycloak:8080"
ENV KEYCLOAK_REALM="demo"
ENV KEYCLOAK_CLIENT_ID="go-client-test"
ENV KEYCLOAK_CLIENT_SECRET="1VGIiII5tCz0SxMGEFfBjMpc3cYd24Ll"

RUN go build -o auth-app ./cmd/auth-app

CMD [ "/app/auth-app" ]
