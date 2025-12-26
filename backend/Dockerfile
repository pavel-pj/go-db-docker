FROM golang:1.25-alpine

WORKDIR /app
 
RUN apk add --no-cache curl tar postgresql-client && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# Устанавливаем Air
RUN go install github.com/air-verse/air@latest

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download
 
COPY . .

EXPOSE 8100

CMD ["air", "-c", ".air.toml"]