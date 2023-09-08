FROM golang:latest

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go env -w GOFLAGS="-buildvcs=false"
RUN go mod tidy