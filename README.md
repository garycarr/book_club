# Book Club Standing Be

A REST service for a book club

To run from the dockerfile

```
docker build  -t book_club . && docker run -p 8080:8080 book_club

curl -d '{"username":"gcarr", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/login
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDkzMDQ0OTksImlhdCI6MTUwOTMwMDg5OSwiaXNzIjoiTWUifQ.nd6Q_IHZagWHxAinYYCk3aAUp-5uuV5luq_smwUL6lo"}
```
