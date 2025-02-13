SEND1='{"toUser":"Alice","amount":100}'
SEND2='{"toUser":"Alice","amount":-100}'
SEND3='{"toUser":"Alice","amount":2000}'
SEND4='{"toUser":"Bob","amount":200}'
HEADER='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk1NTE5MzAsInVzZXJfaWQiOjF9.1BtNlrIFInZkD5DfCVZyM9ZhghDWDdeT1OwSCSfpP8g'
curl -i -XPOST -d "$SEND1" --header "$HEADER" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND2" --header "$HEADER" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND3" --header "$HEADER" localhost:8080/api/sendCoin
curl -i -XPOST -d "$SEND4" --header "$HEADER" localhost:8080/api/sendCoin
