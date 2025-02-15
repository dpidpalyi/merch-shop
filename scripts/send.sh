SEND1='{"toUser":"Alice","amount":100}'
SEND2='{"toUser":"Alice","amount":-100}'
SEND3='{"toUser":"Alice","amount":200}'
SEND4='{"toUser":"Bob","amount":200}'
SEND5='{"toUser":"Bob","amount":10}'
SEND6='{"toUser":"Bob","amount":-10}'
SEND7='{"toUser":"Bob","amount":300}'
SEND8='{"toUser":"Alice","amount":20}'
HEADER1='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjF9.CI3aWDccp4DKqkuG-FjQR-C9VIGo187cshHYwc8vSsk'
HEADER2='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjJ9.nA6a4HBh-MS3Fa31gc4es4SKlz28__3E5I61QYdBqfk'
HEADER3='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2Nzc5ODksInVzZXJfaWQiOjN9.Oh8LQVjsGxruG2fi3gEaowm4vyMaeEP9yMN8gW9ThD4'
curl -i -XPOST -d "$SEND1" --header "$HEADER1" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND2" --header "$HEADER1" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND3" --header "$HEADER1" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND4" --header "$HEADER1" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND5" --header "$HEADER2" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND6" --header "$HEADER2" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND7" --header "$HEADER2" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND8" --header "$HEADER2" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND5" --header "$HEADER3" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND6" --header "$HEADER3" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND7" --header "$HEADER3" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND8" --header "$HEADER3" localhost:8080/api/sendCoin
