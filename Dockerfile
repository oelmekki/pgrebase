FROM golang:alpine as builder

RUN apk add --no-cache git

WORKDIR /go/src/pgrebase
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest

COPY --from=builder /go/bin/pgrebase /usr/local/bin/

ENTRYPOINT ["pgrebase"]