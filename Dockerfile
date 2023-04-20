FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]