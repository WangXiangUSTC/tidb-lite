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

package tidblite

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	. "github.com/pingcap/check"
)

func TestClient(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testTiDBSuite{})

type testTiDBSuite struct{}

func (t *testTiDBSuite) TestTiDBServer(c *C) {
	tidbServer1, err := NewTiDBServer(NewOptions(c.MkDir()).WithPort(4040))
	c.Assert(err, IsNil)
	defer tidbServer.Close()

	tidbServer2, err := NewTiDBServer(NewOptions(c.MkDir()).WithPort(4041))
	c.Assert(err, IsNil)
	// only support create one tidb server
	c.Assert(tidbServer1, Equals, tidbServer2)

	var dbConn *sql.DB
	for i := 0; i < 5; i++ {
		dbConn, err = tidbServer1.CreateConn()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	c.Assert(err, IsNil)

	result, err := dbConn.Query("SELECT version()")
	c.Assert(err, IsNil)
	defer result.Close()

	var version sql.NullString
	for result.Next() {
		err := result.Scan(&version)
		c.Assert(err, IsNil)
		break
	}

	c.Assert(result.Err(), IsNil)
	c.Assert(version.Valid, IsTrue)
	c.Log(version.String)
	// version looks like 5.7.25-TiDB-None
	c.Assert(strings.Contains(version.String, "TiDB"), IsTrue)
}
