FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o myapp ./cmd/api

FROM alpine
WORKDIR /app
COPY --from=builder /app/myapp /app/myapp
COPY --from=builder /app/database /app/database
ENTRYPOINT ["/app/myapp"]