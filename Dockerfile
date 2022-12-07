FROM alpine:edge AS builder
RUN apk add --no-cache build-base go git
WORKDIR /src

COPY ./go.mod ./go.sum ./tie/
RUN cd tie && go mod download -x all

COPY . ./
RUN go build -v -o /app/tie tie.prodigy9.co

# ---

FROM alpine:edge
LABEL org.opencontainers.image.source=https://github.com/prod9/tie

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /src/data/migrations /app/data/migrations
COPY --from=builder /app/tie /app

EXPOSE 4000
CMD ["/app/tie", "serve"]

