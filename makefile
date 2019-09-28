build:
	export http_proxy=http://127.0.0.1:1087
	export https_proxy=http://127.0.0.1:1087
	go get -u -v github.com/go-ini/ini
	go get -u -v -tags "etcd" github.com/smallnest/rpcx/...
	go get -u -v github.com/gin-gonic/gin
	go get -u -v github.com/go-sql-driver/mysql
	go get -u -v github.com/jmoiron/sqlx
	go get -u -v github.com/dgrijalva/jwt-go
	go get -u -v github.com/robfig/cron