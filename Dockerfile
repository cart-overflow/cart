FROM golang:1.25.0
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY migrations /migrations
COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/cart
CMD ["app"]
