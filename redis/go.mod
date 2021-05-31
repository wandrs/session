module go.wandrs.dev/session/redis

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-redis/redis/v8 v8.9.0
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/unknwon/com v1.0.1
	go.wandrs.dev/session v0.0.0-20210531075027-e2a8cbd582a7
	gopkg.in/ini.v1 v1.62.0
)

replace go.wandrs.dev/session => ../
