FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN make build-app

EXPOSE 8080

CMD ["build/url-collector"]
