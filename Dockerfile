# Build aşaması
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Build için gerekli paketleri yükleme
RUN apk add --no-cache git

# Go modüllerini kopyalama ve indirme
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyalama
COPY . .

# Uygulamayı derleme
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Çalışma aşaması
FROM alpine:latest

WORKDIR /app

# SSL sertifikaları için gerekli paket
RUN apk --no-cache add ca-certificates

# Builder aşamasından derlenmiş uygulamayı kopyalama
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Uygulama için gerekli portu açma
EXPOSE 8080

# Uygulamayı çalıştırma
CMD ["./main"] 