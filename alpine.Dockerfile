FROM golang:1.12-alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git ca-certificates libcap && update-ca-certificates
WORKDIR /app
ADD . .
RUN go mod download
RUN go mod verify
RUN mkdir /kubedemo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /kubedemo/kubedemo ./cmd
RUN addgroup -S -g 10101 appuser
RUN adduser -S -D -u 10101 -s /sbin/nologin -h /appuser -G appuser appuser
RUN chown -R appuser:appuser /kubedemo/kubedemo
RUN setcap 'cap_net_bind_service=+ep' /kubedemo/kubedemo

FROM alpine:3.10
EXPOSE 8080 
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/group /etc/passwd /etc/
COPY --from=builder /kubedemo /
USER appuser
ENTRYPOINT ["/kubedemo"]