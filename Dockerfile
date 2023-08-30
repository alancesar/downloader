FROM golang:1.19
WORKDIR /app
COPY . ./
RUN go mod download
RUN GOOS=linux go build ./cmd/downloader
CMD ["./downloader"]