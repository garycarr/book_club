# Book Club Standing Be

A REST service for a book club

To run

```
go build && book_club
```

To run from the dockerfile

```
docker build  -t book_club . && docker run -p 8080:8080 book_club
```

To see in action (Major TODO - swagger)

```
curl -d '{"email":"gcarr", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/login
```

To deploy to elastic beanstalk, zip the file (not the parent directory) and upload.
