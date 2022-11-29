# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o main

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .

CMD [ "./main" ]