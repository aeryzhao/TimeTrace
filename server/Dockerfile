FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /timetrace-server ./main.go

FROM golang:1.23-bookworm AS runtime

WORKDIR /app

COPY --from=builder /timetrace-server /usr/local/bin/timetrace-server

ENV GIN_MODE=release
ENV TIMETRACE_DB_PATH=/app/data/timetrace.db

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["timetrace-server"]
