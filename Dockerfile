# Use an official Golang runtime as a parent image
FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go build -o app .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["./app"]




