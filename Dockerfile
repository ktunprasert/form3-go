FROM golang:1.18

WORKDIR /go/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
