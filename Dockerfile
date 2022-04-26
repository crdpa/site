FROM golang:1.18.1-alpine AS builder
WORKDIR /build
COPY go.mod go.sum main.go /build/
COPY blogposts /build/blogposts
RUN go mod download
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/site /app/
COPY static/ /app/static
WORKDIR /app
CMD ["./site"]

