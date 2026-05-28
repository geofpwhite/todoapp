FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN go build -o todo .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /build/todo .
COPY backend/active_records.json .
COPY backend/completed_records.json .
COPY backend/num_records.txt .

RUN mkdir -p static

EXPOSE 8080

CMD ["./todo"]
