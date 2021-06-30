FROM golang:1.16.4-alpine as builder
COPY . /app
WORKDIR /app
RUN apk update && \
  apk upgrade && \
  apk add --no-cache ca-certificates && \
  apk add --update-cache tzdata && \
  update-ca-certificates 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags '-w -s' -a -installsuffix cgo  -o /app/bin/main ./main.go
  
FROM alpine:3.14
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /app/bin/main /main

COPY ./bin/stock-write /stock/stock-write
COPY ./bin/.env.local /stock/.env.local

ENV TZ=Asia/Seoul \
    ZONEINFO=/zoneinfo.zip  
CMD ["/main"]
# docker run --name ticker stock-write-ticker:latest 
# docker exec -it ticker sh