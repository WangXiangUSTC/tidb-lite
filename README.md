# tidb-lite

[English README](./README_EN.md)

tidb-lite 是一个库，可以方便地使用该库在 golang 应用或者单元测试中使用 `mocktikv` 模式的 TiDB。

## 提供的接口

- func NewTiDBServer(options *Options) (*TiDBServer, error)
  
  创建一个 TiDB Server，使用 options 来设置 TiDB 存储数据的路径和服务的端口号。

- func (t *TiDBServer) CreateConn() (*sql.DB, error)
  
  获取一个 TiDB 的链接。

- func (t *TiDBServer) Close()
  
  关闭 TiDB 服务。

- func (t *TiDBServer) CloseGracefully()
  
  优雅地关闭 TiDB 服务。

## 使用示例

可以通过 `example` 路径下的示例代码了解 tidb-lite 的使用方法。

在 [example.go](./example/example.go) 中定义了一个函数 `GetRowCount` 计算指定表中符合条件的数据的行数。
在 [example_test.go](./example/example_test.go) 中使用以下代码创建 TiDB Server 并获取数据库链接：

```
tidbServer, err := tidblite.NewTiDBServer(tidblite.NewOptions(c.MkDir()).WithPort(4040))
c.Assert(err, IsNil)
defer tidbServer.Close()

var dbConn *sql.DB
for i := 0; i< 5; i++ {
	dbConn, err = tidbServer.CreateConn()
	if err != nil {
		time.Sleep(100*time.Millisecond)
        continue
	}
    break
}
c.Assert(err, IsNil)
```

然后就可以使用链接 `dbConn` 生成测试数据，对函数 `GetRowCount` 进行测试。

