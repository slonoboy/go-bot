FROM golang:latest

WORKDIR /usr/local/go/src/app

RUN go mod download

COPY . .

RUN mkdir -p logs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/app

ENTRYPOINT ["/usr/local/go/src/app/app"]
