# stock-write-ticker
1분 마다 체크함. 


```
빌드
$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
go build -o bin/release main.go

```

```
실행
nohup ticker    > ticker.out &
nohup ticker test   > ticker.out &
```

```
chmod +x ./stock-write-ticker/release

ec2 업로드 전 기존 프로세스 kill -9 p_id 시키기.
1. ssh -i "highserpot_stock.pem" ec2-user@ec2-3-35-30-100.ap-northeast-2.compute.amazonaws.com

scp -i "highserpot_stock.pem" bin/release  ec2-user@3.35.30.100:~/ stock-write-ticker/release
nohup ./stock-write-ticker/release    > stock-write-ticker/nohup.out &

```
### ticker  안될때
``` 
// ‘Seoul’ 파일 확인
$ ls /usr/share/zoneinfo/Asia

// Localtime 심볼릭 링크 재설정
$ sudo ln -sf /usr/share/zoneinfo/Asia/Seoul /etc/localtime

// 적용 확인
$ date
```