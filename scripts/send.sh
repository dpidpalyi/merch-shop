SEND1='{"toUser":"Alice","amount":100}'
SEND2='{"toUser":"Alice","amount":-100}'
SEND3='{"toUser":"Alice","amount":200}'
SEND4='{"toUser":"Bob","amount":200}'
SEND5='{"toUser":"Bob","amount":10}'
SEND6='{"toUser":"Bob","amount":-10}'
SEND7='{"toUser":"Bob","amount":300}'
SEND8='{"toUser":"Alice","amount":20}'
HEADER1='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk1OTk3OTIsInVzZXJfaWQiOjF9.p9L-Z63fLjFCd_93s3dd0WJaquRY5SW-Y9RN1JT8EoM'
HEADER2='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk1OTk3OTIsInVzZXJfaWQiOjJ9.vyhQ_rqGynJmfAvy1GvCh7Y0IgYOgRBG2Mq4xKk_ZiA'
HEADER3='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MDE3MjUsInVzZXJfaWQiOjN9.C4YdKhvo0EEJNAvc5eC9-9DnR4l_UlzzSpH6pg79C1A'
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
