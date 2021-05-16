FROM golang:1.15 AS build

WORKDIR /src/

ADD . /src/

RUN go mod download
RUN go build -o main .

CMD ["/src/main"]