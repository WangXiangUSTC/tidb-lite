// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package example

import (
	"context"
	"database/sql"
	"testing"
	"time"

	tidblite "github.com/WangXiangUSTC/tidb-lite"
	. "github.com/pingcap/check"
)

func TestClient(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testExampleSuite{})

type testExampleSuite struct{}

func (t *testExampleSuite) TestGetRowCount(c *C) {
	tidbServer, err := tidblite.NewTiDBServer(tidblite.NewOptions(c.MkDir()).WithPort(4040))
	c.Assert(err, IsNil)
	defer tidbServer.Close()

	var dbConn *sql.DB
	for i := 0; i < 5; i++ {
		dbConn, err = tidbServer.CreateConn()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	c.Assert(err, IsNil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = dbConn.ExecContext(ctx, "create database example_test")
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "create table example_test.t(id int primary key, name varchar(24))")
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "insert into example_test.t values(1, 'a'),(2, 'b'),(3, 'c')")
	c.Assert(err, IsNil)

	count, err := GetRowCount(ctx, dbConn, "example_test", "t", "id > 2")
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))

	count, err = GetRowCount(ctx, dbConn, "example_test", "t", "")
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(3))
}
