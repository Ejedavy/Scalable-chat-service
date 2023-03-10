FROM golang:1.20.2-alpine3.17
WORKDIR /app
COPY . .
RUN chmod 777 /app/wait.sh
RUN chmod 777 /app/start.sh
RUN go build -o main /app/main.go
EXPOSE 8080
CMD ["/app/main"]
