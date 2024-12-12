FROM golang:1.23.3

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go clean -modcache

RUN go mod tidy

CMD ["air"]