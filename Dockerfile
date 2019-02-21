FROM golang:1.10

WORKDIR /go/src/github.com/wbgalvao/tracking_server/
COPY . .

RUN make build
EXPOSE 8080

CMD ["./tracking_server"]
