FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /go/src/pgrebase
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN apk del git

CMD ["pgrebase"]