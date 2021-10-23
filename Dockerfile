FROM golang:1.16 AS builder
WORKDIR /go/src/passbase
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /go/src/passbase/app ./app
RUN ls /root/
COPY migrations migrations
CMD ["./app"]    
