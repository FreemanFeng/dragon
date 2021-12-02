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

package store

import (
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func dbPath(path string) string {
	prefix := GetPath(DBPath, DefaultDBPath)
	return strings.Join([]string{prefix, path}, SLASH)
}

func closeDB(dbs map[string]*leveldb.DB) {
	for _, db := range dbs {
		db.Close()
	}
}

func walkDB(db *leveldb.DB, k StoreType) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		k.BCh <- iter.Value()
		<-k.DCh
	}
	k.BCh <- []byte{}
	iter.Release()
}

func walkKeys(db *leveldb.DB, k StoreType) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		k.BCh <- iter.Key()
		<-k.DCh
	}
	k.BCh <- []byte{}
	iter.Release()
}

func emptyDB(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		err := db.Delete(iter.Key(), nil)
		if err != nil {
			Debug(ModuleDB, "Could not delete DB data for", iter.Key())
		}
	}
	iter.Release()
}

func readDB(db *leveldb.DB, k StoreType) {
	data, err := db.Get([]byte(k.Key), nil)
	if err != nil {
		data = []byte{}
	}
	k.BCh <- data
}

func Run() {
	ch := AnyChannel(ModuleDB)
	uniq := map[string]int{}
	dbs := map[string]*leveldb.DB{}
	defer closeDB(dbs)
	for {
		x := <-ch
		k := x.(StoreType)
		_, ok := uniq[k.Service]
		if !ok {
			path := dbPath(k.Service)
			db, err := leveldb.OpenFile(path, nil)
			if err != nil {
				Debug(ModuleDB, "Could not open DB", path)
				if k.Op == GetALLData || k.Op == Get {
					k.BCh <- []byte{}
				}
				continue
			}
			dbs[k.Service] = db
			uniq[k.Service] = 1
			Debug(ModuleDB, "Open DB", path)
		}
		switch k.Op {
		// 写数据
		case Set:
			err := dbs[k.Service].Put([]byte(k.Key), k.Data, nil)
			if err != nil {
				Debug(ModuleDB, "Could not write data into DB", err)
			}
		// 读数据
		case Get:
			go readDB(dbs[k.Service], k)
		// 读全部数据
		case GetALLData:
			go walkDB(dbs[k.Service], k)
		// 读全部数据
		case GetALLKeys:
			go walkKeys(dbs[k.Service], k)
		// 删除数据
		case Delete:
			err := dbs[k.Service].Delete([]byte(k.Key), nil)
			if err != nil {
				Debug(ModuleDB, "Could not delete DB data for key", k.Key)
			}
		// 清空数据
		case UNDEFINED:
			go emptyDB(dbs[k.Service])
		}
	}
}
