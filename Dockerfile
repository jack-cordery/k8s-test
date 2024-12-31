FROM alpine:3.16 AS root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home

FROM golang:1.23.3-alpine AS builder
WORKDIR /go-api
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./...

FROM scratch AS final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --chown=1001:1001 --from=builder /go-api/main /main
EXPOSE 8080
USER app 
ENTRYPOINT [ "/main" ]