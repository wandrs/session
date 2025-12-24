module go.wandrs.dev/session/ledis

go 1.16

require (
	github.com/go-chi/chi/v5 v5.2.2
	github.com/ledisdb/ledisdb v0.0.0-20200510135210-d35789ec47e6
	github.com/smartystreets/goconvey v1.6.4
	github.com/unknwon/com v1.0.1
	go.wandrs.dev/session v0.0.0-20210531074458-568b38307750
	gopkg.in/ini.v1 v1.62.0
)

replace go.wandrs.dev/session => ../
