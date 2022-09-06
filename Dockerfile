FROM alpine:edge AS builder
RUN apk add --no-cache build-base go git
WORKDIR /src/tie

COPY go.mod go.sum ./
RUN go mod download -x all

COPY . ./
RUN go build -o /app/tie .

# ---

FROM alpine:edge
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /src/tie/data/migrations /app/data/migrations
COPY --from=builder /app/tie /app

EXPOSE 4000
CMD ["/app/tie", "serve"]

