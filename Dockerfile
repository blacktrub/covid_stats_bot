FROM golang:1.14.2-alpine

RUN apk add git

WORKDIR /go/src/app

COPY . .

RUN go install

RUN go build -o main .

CMD ["app"]
