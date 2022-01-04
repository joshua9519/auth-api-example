FROM golang:alpine AS builder

WORKDIR /api

COPY . .

RUN go mod download
RUN go build .
RUN ls /api

FROM alpine
COPY --from=builder /api/auth /
CMD ["/auth"]