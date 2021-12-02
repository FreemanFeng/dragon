package goplugin

import "sync"

type Config struct {
	Project string   // 项目名称
	Calls   sync.Map // 函数映射
}
