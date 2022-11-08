FROM ubuntu:22.04

COPY go-webapp-template /app/
COPY config/config.yml /app/config/config.yml

EXPOSE 8080
WORKDIR /app

CMD ["/app/go-webapp-template"]