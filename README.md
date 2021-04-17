# 한국 코로나19 상황판
## 사용법
<https://t.me/KOR_corona19_status_robot> 에서 서비스를 진행했었습니다. 
현재는 서버 자원이 없어 운영하지 않습니다.

## 소스코드 직접 컴파일 하기
### 필요사항
- Go 1.4
프로젝트를 다음의 명령어로 클론 하세요.
```bash
git clone https://github.com/code-yeongyu/corona19-telegram-notifier-bot
```

그 다음 프로젝트 폴더에 들어간 뒤에 다음의 명령어로 의존성이 있는 모듈을 설치해주세요.
```bash
go mod tidy
```

끝입니다! 실행하고 싶으시다면, 다음과 같은 내용을 환경변수에 지정해주세요.
TELEGRAM_CODE = 봇 API KEY  
DB_INFO = 데이터를 저장할 DB정보 (DB정보 양식은 Golang 에서의 sql 사용방법을 검색해보세요.)
