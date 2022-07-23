FROM golang:1.15-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o greeter

CMD ["/app/greeter"]
