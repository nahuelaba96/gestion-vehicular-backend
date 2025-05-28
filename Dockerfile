# Etapa de build
FROM golang:1.23-alpine AS builder

# Instalamos herramientas necesarias para compilar
RUN apk add --no-cache git

WORKDIR /app

# Copiamos los archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilamos el binario estático
RUN go build -ldflags="-s -w" -o app

# Etapa final: imagen mínima
FROM alpine:latest

# 🛠 Instalá certificados TLS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiamos el binario desde la etapa anterior
COPY --from=builder /app/app .

# Exponemos el puerto en el que escucha tu app
EXPOSE 8080

# Comando por defecto al iniciar el contenedor
CMD ["./app"]
