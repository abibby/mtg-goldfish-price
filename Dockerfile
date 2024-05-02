FROM golang:1.22 as go-build

WORKDIR /go/src/github.com/abibby/mtg-goldfish-price

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /mtg-goldfish-price

# Now copy it into our base image.
FROM alpine:latest

RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates

COPY --from=go-build /mtg-goldfish-price /mtg-goldfish-price

EXPOSE 3335/tcp

CMD ["/mtg-goldfish-price"]