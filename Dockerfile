# Builder image go
FROM golang:1.17 AS builder

# Build fruits binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /fruits
WORKDIR /fruits
COPY . /fruits
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/fruits-service ./cmd/fruitsd/main.go 

# Runnable image
FROM alpine
COPY --from=builder /fruits/bin/fruits-service /bin/fruits-service
RUN ls /bin/fruits-service
WORKDIR /bin
ENTRYPOINT [ "./fruits-service" ]