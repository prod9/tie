FROM alpine:edge AS builder
RUN apk add --no-cache build-base go git
WORKDIR /src

COPY ./tie/go.mod ./tie/go.sum ./tie/
RUN cd tie && go mod download -x all

COPY . ./
RUN go build -v -o /app/tie tie.prodigy9.co

# ---

FROM alpine:edge
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /src/tie/data/migrations /app/data/migrations
COPY --from=builder /app/tie /app

EXPOSE 4000
CMD ["/app/tie", "serve"]

