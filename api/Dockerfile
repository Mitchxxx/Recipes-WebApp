FROM golang:1.23-alpine AS builder
WORKDIR /go/src/github.com/api
COPY . .
RUN go mod download
RUN CGO_Enabled=0 GOOS=linux go build -a -installsuffix cgo -o app .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/api/app .
CMD ["./app"]