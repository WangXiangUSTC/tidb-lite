module github.com/WangXiangUSTC/tidb-lite

go 1.12

require (
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/go-sql-driver/mysql v0.0.0-20170715192408-3955978caca4
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/myesui/uuid v1.0.0 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pingcap/check v0.0.0-20191107115940-caf2b9e6ccf4
	github.com/pingcap/errors v0.11.4
	github.com/pingcap/log v0.0.0-20191012051959-b742a5d432e9
	github.com/pingcap/parser v0.0.0-20191210055545-753e13bfdbf0
	github.com/pingcap/tidb v1.1.0-beta.0.20191211070559-a94cff903cd1
	github.com/twinj/uuid v1.0.0 // indirect
	go.uber.org/zap v1.12.0
)

replace github.com/pingcap/parser v0.0.0-20191210055545-753e13bfdbf0 => github.com/WangXiangUSTC/parser v0.0.0-20210518141443-57a13baea68a
