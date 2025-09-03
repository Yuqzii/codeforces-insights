FROM golang:1.24

WORKDIR /app
COPY . .
RUN go mod download

RUN GOOS=linux go build -o /server ./cmd/server/

EXPOSE 8080

ENTRYPOINT ["/server"]
