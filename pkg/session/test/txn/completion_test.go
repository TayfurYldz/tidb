// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package txn

import (
	"testing"

	"github.com/pingcap/tidb/pkg/testkit"
)

func TestTransactionCompletionChain(t *testing.T) {
	store := testkit.CreateMockStore(t)
	tk := testkit.NewTestKit(t, store)
	tk.MustExec("use test")

	t.Run("commit and chain", func(t *testing.T) {
		tk.MustExec("create table commit_chain (id int primary key)")
		tk.MustExec("begin")
		tk.MustExec("insert into commit_chain values (1)")
		tk.MustExec("commit and chain")

		tk.MustExec("insert into commit_chain values (2)")
		tk.MustExec("rollback")

		tk.MustQuery("select id from commit_chain order by id").
			Check(testkit.Rows("1"))
	})

	t.Run("rollback and chain", func(t *testing.T) {
		tk.MustExec("create table rollback_chain (id int primary key)")
		tk.MustExec("begin")
		tk.MustExec("insert into rollback_chain values (1)")
		tk.MustExec("rollback and chain")

		tk.MustExec("insert into rollback_chain values (2)")
		tk.MustExec("rollback")

		tk.MustQuery("select id from rollback_chain order by id").
			Check(testkit.Rows())
	})
}
