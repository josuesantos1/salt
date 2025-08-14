ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /run-app .

FROM alpine:latest
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=builder /run-app /usr/local/bin/
COPY public public
CMD ["/usr/local/bin/run-app"]
