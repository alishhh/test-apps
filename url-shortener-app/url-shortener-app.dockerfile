FROM golang:1.21

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install inetutils-ping

ENV LOG_LEVEL="prod"
ENV APP_NAME="url-shortener-app"
ENV APP_HOST="0.0.0.0"
ENV APP_PORT="8080"
ENV TIMEOUT="4s"
ENV IDLE_TIMEOUT="120s"
ENV DB_PATH="./storage/storage.db"
ENV CLIENT_URL="http://auth-app:8080"

RUN go build -o url-shortener-app ./cmd/url-shortener-app

CMD [ "/app/url-shortener-app" ]
