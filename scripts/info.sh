HEADER1='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjF9.CI3aWDccp4DKqkuG-FjQR-C9VIGo187cshHYwc8vSsk'
HEADER2='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjJ9.nA6a4HBh-MS3Fa31gc4es4SKlz28__3E5I61QYdBqfk'
HEADER3='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjN9.Oh8LQVjsGxruG2fi3gEaowm4vyMaeEP9yMN8gW9ThD4'
curl -i -XGET --header "$HEADER1" localhost:8080/api/info
curl -i -XGET --header "$HEADER2" localhost:8080/api/info
curl -i -XGET --header "$HEADER3" localhost:8080/api/info
