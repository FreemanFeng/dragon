//   Copyright 2019 Freeman Feng<freeman@nuxim.cn>
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package report

import (
	"strconv"
	"time"

	. "github.com/nuxim/dragon/dragon/common"
)

func Run() {
	ch := AnyChannel(ModuleReport)
	for {
		x := <-ch
		r := x.(TaskRequest)
		go report(r)
		Debug(ModuleReport, r)
	}
}

func report(r TaskRequest) {
	read := 0                   // 消费计数器
	write := 0                  // 生产计数器
	max := 0                    // 最大尝试次数
	ch := AnyChannel(r.Task)    // 测试执行结果通道
	ich := IntChannel(r.Task)   // 请求通道
	dch := AnyChannel(ModuleDB) // 数据库通道
	rch := BytesChannel(r.Task) // 应答通道
	for {
		select {
		case x := <-ch:
			k := x.(TestReport)
			dch <- StoreType{Service: r.Task, Op: Set, Key: strconv.Itoa(write), Data: GobEncode(k)}
			write++
			Debug(ModuleReport, r.Task, k)
		case x := <-ich:
			if x == INIT {
				read = 0
			}

			bch := make(chan []byte)
			var v []byte
			for {
				if read == UNDEFINED {
					rch <- []byte(EOF)
					continue
				}
				for max < MaxTry {
					dch <- StoreType{Service: r.Task, Op: Get, Key: strconv.Itoa(read), BCh: bch}
					v = <-bch
					if len(v) == 0 {
						max++
						Debug(ModuleReport, r.Task, "no data!")
						time.Sleep(time.Duration(100) * time.Microsecond)
					} else {
						max = 0
						break
					}
				}
				if max >= MaxTry {
					rch <- []byte(EOF)
					continue
				}
				var q TestReport
				err := GobDecode(v, &q)
				if err != nil {
					Debug(ModuleReport, r.Task, "error decoding", err)
					rch <- []byte{}
					break
				}
				rch <- formatTCReport(read, q)
				read++
				if q.Line == UNDEFINED {
					read = UNDEFINED
				}
			}
		}
	}
}

func formatTCReport(id int, q TestReport) []byte {
	key := JoinKeys(VerticalBar, EMPTY, id, q.Project, q.Suite, q.Line, q.Title, q.Passed, q.Reason, EMPTY)
	return []byte(key)
}
