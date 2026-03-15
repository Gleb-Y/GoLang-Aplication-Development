FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY . .
 
RUN go build -o myapp ./cmd/api/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/myapp /app/myapp
COPY --from=builder /app/database /app/database
ENTRYPOINT ["/app/myapp"]