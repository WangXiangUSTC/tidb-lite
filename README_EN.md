# tidb-lite

tidb-lite is a package, we can use this package to create a TiDB server with `mocktikv` mode in your application or unit test.

## Interface

- func NewTiDBServer(options *Options) (*TiDBServer, error)
  
  Create a TiDB Server, can use options to set the path which used to save db's data and this server's port.

- func (t *TiDBServer) CreateConn() (*sql.DB, error)
  
  Create a database connection.

- func (t *TiDBServer) Close()
  
  Close TiDB Server.

- func (t *TiDBServer) CloseGracefully()
  
  Close TiDB server gracefully.

## Example

We can read the code under [example.go](./example/example.go) to know how to use tidb-lite.

In [example.go](./example/example.go) defines a function `GetRowCount` to calculate the count of a table with condition.

In [example_test.go](./example/example_test.go) use code below to create a TiDB server and get the database's connection used for unit test.

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

And then we can use the `dbConn` to generate test data, and then check function `GetRowCount`'s correctness.
