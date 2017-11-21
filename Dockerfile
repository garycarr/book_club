FROM golang:1.8.3
EXPOSE 8080

WORKDIR /go/src/github/gcarr/book_club
COPY . /go/src/github/gcarr/book_club

RUN go-wrapper install

CMD ["go-wrapper", "run"]
