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
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
)

func TestClient(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testTiDBSuite{})

type testTiDBSuite struct{}

func (t *testTiDBSuite) TestTiDBServer(c *C) {
	tidbServer1, err := NewTiDBServer(NewOptions(c.MkDir()))
	c.Assert(err, IsNil)
	defer tidbServer1.Close()

	_, err = GetTiDBServer()
	c.Assert(err, IsNil)

	_, err = NewTiDBServer(NewOptions(c.MkDir()).WithPort(4041))
	// only support exist one tidb server
	c.Assert(err, NotNil)

	dbConn, err := tidbServer1.CreateConn()
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

	// test for set meta
	tblInfo := &model.TableInfo{
		ID:    1000,
		Name:  model.NewCIStr("my_table"),
		State: model.StatePublic,
		Columns: []*model.ColumnInfo{
			{
				Name:         model.NewCIStr("c01"),
				ID:           1,
				Offset:       0,
				DefaultValue: 0,
				State:        model.StatePublic,
				FieldType:    *types.NewFieldType(mysql.TypeLong),
			},
		},
		Indices: []*model.IndexInfo{
			{
				ID:      10,
				Name:    model.NewCIStr("i1"),
				Table:   model.NewCIStr("my_table"),
				Columns: []*model.IndexColumn{{Name: model.NewCIStr("c01"), Offset: 0}},
				Unique:  true,
				State:   model.StatePublic,
				Tp:      model.IndexTypeBtree,
			},
		},
	}
	dbInfo := &model.DBInfo{
		ID:     100,
		Name:   model.NewCIStr("my_db"),
		Tables: []*model.TableInfo{tblInfo},
		State:  model.StatePublic,
	}
	err = tidbServer1.SetDBInfoMetaAndReload([]*model.DBInfo{dbInfo})
	c.Assert(err, IsNil)

	err = tidbServer.dom.Reload()
	c.Assert(err, IsNil)

	_, err = dbConn.Query("use my_db;")
	c.Assert(err, IsNil)

	result, err = dbConn.Query("show create table my_db.my_table")
	c.Assert(err, IsNil)
	var nameString sql.NullString
	var createTableString sql.NullString
	for result.Next() {
		err := result.Scan(&nameString, &createTableString)
		c.Assert(err, IsNil)
		break
	}
	c.Assert(nameString.String, Equals, "my_table")
	c.Assert(createTableString.String, Equals, "CREATE TABLE `my_table` (\n"+
		"  `c01` int(11) CHARACTER SET  COLLATE  DEFAULT '0',\n"+
		"  UNIQUE KEY `i1` (`c01`(0))\n"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")

	c.Assert(err, IsNil)
	result, err = dbConn.Query("SELECT non_unique, index_name, seq_in_index, column_name FROM information_schema.statistics WHERE table_schema = \"my_db\" AND table_name = \"my_table\" ORDER BY seq_in_index ASC")
	c.Assert(err, IsNil)
	var (
		nonUniqueString  sql.NullString
		indexNameString  sql.NullString
		seqInIndexString sql.NullString
		columnString     sql.NullString
	)
	for result.Next() {
		err := result.Scan(&nonUniqueString, &indexNameString, &seqInIndexString, &columnString)
		c.Assert(err, IsNil)
		break
	}
	c.Assert(nonUniqueString.String, Equals, "0")
	c.Assert(indexNameString.String, Equals, "i1")
	c.Assert(seqInIndexString.String, Equals, "1")
	c.Assert(columnString.String, Equals, "c01")

	tidbServer1.Close()

	// TODO: remove time sleep
	time.Sleep(time.Second)

	tidbServer2, err := NewTiDBServer(NewOptions(c.MkDir()))
	// can create another tidb server after the tidb close.
	c.Assert(err, IsNil)
	tidbServer2.Close()

	// tidb server not exist after close
	_, err = GetTiDBServer()
	c.Assert(err, NotNil)
}
