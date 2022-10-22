#go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#go build ./main.go
ssh root@70.34.209.206 "systemctl stop crypto-notice.service"
scp ./main root@70.34.209.206:~/
ssh root@70.34.209.206 "chmod +x ./main; systemctl start crypto-notice.service"