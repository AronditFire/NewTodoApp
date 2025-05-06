FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod downlaod
RUN go build -o myapp ./cmd/app/main.go

CMD ["./myapp"]
