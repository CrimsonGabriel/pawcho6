# syntax=docker/dockerfile:1.6

# === Etap 1: Budowanie aplikacji Go ===
FROM golang:alpine AS builder

WORKDIR /app

COPY server.go ./
COPY index.html ./

RUN go mod init myapp

ARG VERSION=1.0.0
ENV VERSION=${VERSION}

RUN go build -o aprogram -ldflags "-X main.Version=${VERSION}"

# === Etap 2: Konfiguracja Nginx jako reverse proxy ===
FROM nginx:alpine

# Usunięcie domyślnej konfiguracji Nginx
RUN rm /etc/nginx/conf.d/default.conf

# Kopiowanie własnej konfiguracji Nginx
COPY nginx.conf /etc/nginx/conf.d/

# Kopiowanie aplikacji Go do kontenera
COPY --from=builder /app/aprogram /usr/bin/aprogram
COPY --from=builder /app/index.html /usr/share/nginx/html/

# Ustawienie uprawnień do uruchamiania aplikacji
RUN chmod +x /usr/bin/aprogram

# Otwieramy port 80 dla Nginx
EXPOSE 80

# Healthcheck - sprawdzanie działania serwera
HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -f http://localhost/ || exit 1

# Uruchomienie zarówno aplikacji Go, jak i Nginx
CMD ["/bin/sh", "-c", "/usr/bin/aprogram & nginx -g 'daemon off;'"]
