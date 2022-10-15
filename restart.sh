go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
go build ./main.go
ssh root@70.34.209.206 "kill \$(ps aux |grep main |awk 'NR==1{print \$2}')"
scp ./main root@70.34.209.206:~/
ssh root@70.34.209.206 "chmod +x ./main; nohup ./main &"