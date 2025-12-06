FROM golang:alpine
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["sh", "-c", "sleep 2 && goose -dir sql/schema postgres \"$DATABASE_URL\" up && ./main"]
