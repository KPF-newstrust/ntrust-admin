모든 시스템은 [도커](https://www.docker.com/) 환경 위에서 구현됩니다. 도커를 설치한 후 아래 명령어로 배포합니다.
```
docker-compose up
```

1. admin 외에 rabbitmq, mongo, redis 등이 모두 함께 실행됩니다.
2. 도커 컨테이너가 실행되면 `{호스트}:50000`로 접속 가능합니다.
3. 초기 관리자계정 로그인은 `admin / test`로 가능합니다.

앱 소스를 수정한뒤 적용하려면 아래 순서로 명령어를 입력합니다.
```
docker-compose stop app
docker-compose build app
docker-compose up -d --no-deps app
```
