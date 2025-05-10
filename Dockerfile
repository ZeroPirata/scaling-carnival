FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd

FROM debian:stable-slim

RUN apt update && apt install -y nginx openssl

COPY --from=builder /app/app /usr/local/bin/app
COPY nginx.conf /etc/nginx/conf.d/myapp.conf

RUN mkdir -p /etc/nginx/ssl && chmod 700 /etc/nginx/ssl
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/nginx/ssl/nginx-selfsigned.key \
    -out /etc/nginx/ssl/nginx-selfsigned.pem \
    -subj "/C=BR/ST=Sao Paulo/L=Sao Jose dos Campos/O=BichoCorporacoes/OU=itau-desafio/CN=localhost"

EXPOSE 443
CMD ["sh", "-c", "nginx && /usr/local/bin/app"]
