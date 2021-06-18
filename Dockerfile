FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN make build-app

CMD ["build/url-collector"]
