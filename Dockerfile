# Builder image go
FROM golang:1.17 AS servicebuilder

# Build fruits binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /fruits
WORKDIR /fruits
COPY . /fruits
RUN go mod download
WORKDIR /fruits/cmd/fruitsd
RUN cd /fruits/cmd/fruitsd && go build -o /fruits-service
CMD [ "/fruits-service" ]