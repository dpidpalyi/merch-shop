HEADER1='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk1OTk3OTIsInVzZXJfaWQiOjF9.p9L-Z63fLjFCd_93s3dd0WJaquRY5SW-Y9RN1JT8EoM'
HEADER2='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk1OTk3OTIsInVzZXJfaWQiOjJ9.vyhQ_rqGynJmfAvy1GvCh7Y0IgYOgRBG2Mq4xKk_ZiA'
HEADER3='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MDE3MjUsInVzZXJfaWQiOjN9.C4YdKhvo0EEJNAvc5eC9-9DnR4l_UlzzSpH6pg79C1A'
curl -i -XGET --header "$HEADER1" localhost:8080/api/info
curl -i -XGET --header "$HEADER2" localhost:8080/api/info
curl -i -XGET --header "$HEADER3" localhost:8080/api/info
