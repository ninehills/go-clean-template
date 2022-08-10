FROM debian:11

WORKDIR /app
ADD ./output/go-web-template_*_Linux-x86_64.tar.gz /app/
EXPOSE 8080
CMD ["/app/go-web-template"]