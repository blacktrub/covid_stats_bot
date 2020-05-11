FROM golang:1.14.2-alpine

RUN apk add git

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

CMD ["app"]
