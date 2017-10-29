FROM golang:1.8.3
EXPOSE 8080

WORKDIR /go/src/github/gcarr/lastman_standing_be
COPY . /go/src/github/gcarr/lastman_standing_be

RUN go-wrapper install

CMD ["go-wrapper", "run"]
