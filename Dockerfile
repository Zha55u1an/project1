FROM golang:1.21.6

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o x ./cmd/main.go

EXPOSE 8000

CMD ["./x"]
