FROM golang:1.20.2-alpine3.17
WORKDIR /app
COPY . .
RUN go build -o main main.go
EXPOSE 8080
CMD ["/app/main"]
