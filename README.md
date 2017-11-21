# Book Club Standing Be

A REST service for a book club

To run from the dockerfile

```
docker build  -t book_club . && docker run -p 8080:8080 book_club

curl -d '{"username":"gcarr", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/login
```
