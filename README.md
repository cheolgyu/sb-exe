# sb-exe

    체크주기: 1분 마다 평일 주식시간
    실행절차
        input
        line
        sbp-line-next
        sbp-stat-volume




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
chmod +x ./sb-exe/release

ec2 업로드 전 기존 프로세스 kill -9 p_id 시키기.
1. ssh -i "highserpot_stock.pem" ec2-user@ec2-3-35-30-100.ap-northeast-2.compute.amazonaws.com

scp -i "highserpot_stock.pem" bin/release  ec2-user@3.35.30.100:~/ sb-exe/release
nohup ./sb-exe/release    > sb-exe/nohup.out &

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