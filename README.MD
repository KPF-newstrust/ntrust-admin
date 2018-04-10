## go 설치
https://golang.org/dl/

## gopath 설정
$ vi ~/.vash_profile
```
#add
export PATH = $PATH:/usr/local/go/bin
export GOPATH = /Users/ben/.go
```
$ cd /Users/ben  
$ mkdir .go && cd .go && mkdir src  
$ .source ~/.bash_profile  
$ echo $GOPATH

## 서버 설정
* config.toml에 각종 서버 설정이 정리되어 있고, 내용은 모두 임의의 값입니다.
* mongodb, redis, rabbitmq 등을 연동하고 있으며 해당 서버들은 직접 구성한 서버에 맞게 작성해서 연동해야 합니다.

## 프로젝트 실행
gopath 링크
```
$ ln -s `pwd` $GOPATH/src/ntrust_admin
```

패키지 설치
```
$ go get -u ntrust_admin
```

실행
```
$ go run main.go (--dev)
```