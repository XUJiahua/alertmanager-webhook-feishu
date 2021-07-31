FROM golang:1.16
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o alertmanager-webhook-feishu .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=0 /app/alertmanager-webhook-feishu .
EXPOSE 8000
ENTRYPOINT ["./alertmanager-webhook-feishu"]
