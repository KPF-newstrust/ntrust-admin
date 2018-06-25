FROM golang:1.8

WORKDIR /go/src/ntrust_admin
COPY . .

RUN go get ./
