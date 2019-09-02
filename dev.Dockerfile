FROM alpine:3.10
EXPOSE 8080
ADD kubedemo /kubedemo
ADD web /web
ENTRYPOINT ["/kubedemo"]