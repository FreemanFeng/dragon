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

package common

import (
	"os/exec"
	"sync"
)

type ControlInfo struct {
	Port     int    // 端口
	Proto    string // 协议 tcp/http/websocket/mqtt etc.
	Service  string // 服务
	Status   string // 服务状态
	Scenario string // 场景: ok/tr/rr/pr/tw/rw/pw
}

type LogConfig struct {
	LogLevel   string         //日志级别
	LogALL     int            // 记录全部日志开关
	LogModules map[string]int // 打日志模块
}

type StoreType struct {
	Service string // 服务名称
	Key     string // 缓存Key
	Data    []byte // 缓存数据
	Op      int    // 操作码，-1：删除服务数据，0：写数据，1：读数据
	BCh     chan []byte
	DCh     chan int
}

type DataInfo struct {
	Key  string
	Data []byte
}

type TaskRequest struct {
	IsMocking    bool     // 是否模拟服务
	IsHijacking  bool     // 是否劫持服务
	IsExtraHosts bool     // 不覆盖已配置的HOSTS
	IsStateless  bool     // 是否无状态，用例可以随意分配
	IsSmoking    bool     // 是否冒烟测试，每个用例只测试一份数据
	IsRandom     bool     // 是否随机测试，每个用例只随机测试一份数据
	IsNoCheck    bool     // 是否不校验，一般用于A/B测试或稳定性测试
	Task         string   // 任务，唯一标识
	Project      string   // 项目名称
	Suite        []string // 测试用例集名称，一般为用例文件名称，去掉后缀名
	Group        []string // 测试用例组，一个用例集包含多个用例组，每个用例组包含多个用例
	Case         []string // 测试用例
	Concurrence  int      // 并发数
	Duration     int64    // 请求总时长
	Tags         []string // 标签列表
	Hosts        []string // 主机列表(IP:PORT)
	Quit         chan int // 结束管道
}

type Testing struct {
	Project     string                   // 项目名称
	IsMocking   bool                     // 是否模拟服务
	IsHijacking bool                     // 是否劫持服务
	Config      ConfigType               // 全局配置
	Templates   map[string]*TemplateType // 全局模板
	ControlFlow map[string]*ControlType  // 全局控制流
	Setup       map[string]*ControlType  // 全局前置条件
	Teardown    map[string]*ControlType  // 全局后置条件
	Checks      map[string]*ControlType  // 全局校验
	Suites      map[string]*TestSuite    // 所有用例集
	Funcs       map[string]*PluginType   // 插件函数
	Bins        map[string]string        //执行函数映射表，Key为函数名，Value为可执行文件路径
}

type ConfigType struct {
	Actions    map[string]string         // 动作
	Attributes map[string]string         // 属性
	Settings   map[string][]string       // 设置
	Meta       map[string]int            // 值
	Plugins    map[string]*PluginType    // 插件
	Messages   map[string]*ConstructType // 消息构造
	Data       map[string]*ConstructType // 数据构造
	Nodes      map[string]*NodeType      // 节点映射表，Key为节点ID
	Flows      map[string]*FlowType      // 流程映射表，Key为流程ID
}

type NodeType struct {
	NID     string              // 服务标识
	Args    []string            // 输入参数
	Title   string              // 标题
	Desc    string              // 描述
	Hosts   []string            // 主机列表，元素值格式ip:port
	User    string              // 用户名
	Pass    string              // 密码
	Tunnel  string              // 通道，一般用于鉴权
	TunUser string              // 通道用户
	TunPass string              // 通道密码
	Path    string              // 路径，用于http协议
	Method  string              // 方法，用于http协议
	Proto   string              // 协议, i.e. http/mqtt/ws/pb
	Plugins []string            // 插件列表
	Extends map[string][]string // 扩展属性列表
	Steps   []StepType          // 控制流列表
}

type TestSuite struct {
	SID         string                  // 用例集ID
	ControlFlow map[string]*ControlType // 控制流
	Setup       map[string]*ControlType // 前置条件
	Teardown    map[string]*ControlType // 后置条件
	Checks      map[string]*ControlType // 校验
	Flows       map[string]*FlowType    // 流程映射表，Key为流程ID
	Groups      map[string]*GroupType   // 用例组
}

type GroupType struct {
	Line    int                  // 用例所在行号
	Timeout int64                // 超时时间
	GID     string               // 组ID
	Title   string               // 组标题
	Desc    string               // 描述
	Proto   string               // 协议 http/mqtt/ws/pb，默认 http
	CIDs    []string             // 控制流ID
	PreIDs  []string             // 前置条件ID
	PostIDs []string             // 后置条件ID
	Plugins []string             // 插件列表
	Extends map[string][]string  // 扩展属性列表
	Args    []OpType             // 参数值，一般是定义前置条件或后置条件的参数值
	Matches []OpType             // 匹配条件，一般用于模拟服务匹配请求特征
	Cases   map[string]*CaseType // 用例

}

type CaseType struct {
	Line        int                 // 用例所在行号
	Timeout     int64               // 超时时间
	Rounds      int64               // 跑的轮次
	Concurrence int                 // 并发数
	Fill        int                 // 对应值 random:0/all:1, 默认random
	Random      int                 // 随机次数
	Local       int                 // 上下文的作用域，1：本地，0：全局，默认0
	TID         string              // 用例ID
	Title       string              // 用例标题
	Desc        string              // 用例描述
	Author      string              // 作者
	Maintainer  string              // 维护者
	Version     string              // 版本
	Level       string              // 等级 A/B/C/D/E 自定义
	Proto       string              // 协议 http/mqtt/ws/pb，默认 http
	Bugs        []string            // 问题列表
	Features    []string            // 需求列表
	Tags        []string            // 标签
	Extends     map[string][]string // 扩展属性列表
	Plugins     []string            // 插件列表
	FIDs        []string            // 流程列表
	Args        []OpType            // 参数值，一般是定义前置条件或后置条件的参数值
	Matches     []OpType            // 匹配条件，一般用于模拟服务匹配请求特征
	Vars        map[string]*VarType // 变量
	Steps       []StepType          // 用例步骤列表
}

type FlowType struct {
	Line    int                 // 用例所在行号
	Random  int                 // 随机次数
	Local   int                 // 上下文的作用域，1：本地，0：全局，默认0
	FID     string              // 流程ID
	Title   string              // 用例标题
	Desc    string              // 用例描述
	Proto   string              // 协议 http/mqtt/ws/pb，默认 http
	Tags    []string            // 标签
	Extends map[string][]string // 扩展属性列表
	Plugins []string            // 插件列表
	FIDs    []string            // 流程列表
	Vars    map[string]*VarType // 变量
	Steps   []StepType          // 用例步骤列表
}

type PluginType struct {
	PID     string                 // 插件ID
	Title   string                 // 插件标题
	Desc    string                 // 插件描述
	Mode    string                 // 插件模式
	Name    string                 // 插件名称
	Build   string                 // 构建，一般用于定制构建
	Start   string                 // 启动，一般用于定制启动参数
	Path    string                 // 插件路径
	Proto   string                 // 协议
	Port    string                 // 端口
	Closed  int                    // 关闭
	IsSRC   bool                   // 是否源码
	IsRPC   bool                   // 是否RPC调用，用于非Go语言的源码插件
	Calls   []string               // 调用函数列表，作为过滤器，剔除掉不需要的函数，若为空，则不过滤，加载所有函数，bin模式不能为空，否则不能映射为可执行文件路径
	Mutes   []string               // 屏蔽函数列表，src模式的插件需要屏蔽的函数列表，bin模式不需要
	Extends map[string][]string    // 扩展属性列表
	Funcs   map[string]interface{} // 插件函数
	Bins    map[string]string      // 可执行文件调用函数映射表，Key为函数名，Value为可执行文件
	Cmd     *exec.Cmd              // 运行命令
}

type VarType struct {
	Key        string            // 变量名
	Node       string            // 节点
	IsTop      bool              // 是否顶级变量
	IsTemplate bool              // 是否模板变量
	IsOpt      bool              // 是否操作变量
	IsGroup    bool              // 是否模板组变量
	IsFlow     bool              // 是否是流程变量
	IsNode     bool              // 是否是节点变量
	Op         string            // 操作符
	Templates  []TemplateVarType // 模板列表(包括规则、节点、上下文、消息)
	Value      []interface{}     // 变量值
}

type TemplateVarType struct {
	ID      string // 标识
	Key     string // 值
	Node    string // 节点标识
	IsVar   bool   // 是否变量
	IsGroup bool   // 是否模板组
	IsFlow  bool   // 是否流程
	IsNode  bool   // 是否节点
	IsCtx   bool   // 是否上下文
	IsWait  bool   // 是否等待
	Count   int    // 倍数
	Targets []int  // 目标列表，用作消息模板的并发id列表
}

type TemplateType struct {
	Type    int               // 模板类型, 0:消息模板(json格式), 1:二进制, 2:文本文件
	Keys    []string          // 标识符列表
	Actions []string          // 动作
	Name    string            // 名称
	Files   map[string]string // 模板列表，Key为顶级目录名称, Value为文件路径
}

type GroupData struct {
	Host  string                // ip:port
	SID   string                // 用例集ID
	GID   string                // 组ID
	Scope string                // 范围
	Cases map[string][]*VarData // 用例数据,Key为用例ID，多个并发多份数据
	CFs   map[string]*VarData   // 控制流数据,Key为控制流ID，用于前置条件/后置条件
	FDs   map[string]*VarData   // 流程数据,Key为流程ID
}

// VarData 运行时用例数据
type VarData struct {
	Total int                     // 运行时实例总数
	ID    string                  // ID
	Title string                  // 标题
	Host  string                  // ip:port
	Funcs map[string]*PluginType  // 插件函数映射,Key为插件函数名称
	Bins  map[string]string       // 可执行函数映射,Key为函数，值为可执行文件路径
	Data  map[string][]*FileType  // 数据集,Key为变量名
	Tags  map[string]*OpType      // 标签替换,Key为标签
	Vars  map[string]interface{}  // 变量映射,Key为变量
	Ctx   map[string]interface{}  // 上下文映射,Key为变量，用于流程
	Msgs  map[string]interface{}  // 消息映射,Key为消息变量
	VPos  map[string]int          // 变量值下标，用于变量遍历取值
	Keys  map[string][]FormatType // 字段映射,Key为变量
	Tops  map[string][]string     // 标记消息存在
}

type FieldType struct {
	Key   string       // 字段名
	Value []FormatType // 值
}

type DepType struct {
	IsFuture bool         //是否未来赋值
	IsKey    bool         // 是否Key
	Var      string       // 对应变量
	Key      string       // 字段名
	Value    []FormatType // 数据
}

type FileType struct {
	IsVar   bool     // 是否变量
	IsFlow  bool     // 是否流程
	IsNode  bool     // 是否节点
	IsJSON  bool     // 是否JSON数据
	IsMap   bool     // JSON格式MAP结构
	IsList  bool     // JSON格式List结构
	IsBin   bool     // 是否二进制数据
	Top     string   // 顶级目录
	Name    string   // 文件名
	Path    string   // 路径
	Tags    []string // 动态标签列表
	Content []byte   // 内容
}

type TestReport struct {
	Stage   int      // 阶段, 0:解析规则配置阶段，1:执行用例阶段
	Code    int      // 状态码
	Line    int      // 用例所在行号
	Passed  int      // 是否通过，1：通过，0：不通过
	File    string   // 规则配置文件
	Title   string   // 用例标题
	Reason  string   // 用例失败原因
	Project string   // 项目名称
	Suite   string   // 用例集名称
	Group   string   // 用例组名称
	Tags    []string // 标签列表
}

type StepType struct {
	Line     int         // 行号
	No       int         // 序号
	Rule     string      // 规则原文
	IsExpect bool        // 是否校验步骤
	Setting  SettingType // 赋值
	Expects  []OpType    // 期待行为
}

type SettingType struct {
	IsList     bool     // 是否列表赋值，仅用于操作赋值
	Items      int      // 列表元素总数，仅用于操作赋值
	IsDep      bool     // 是否依赖
	IsTag      bool     // 是否标签，用于文本替换
	IsField    bool     // 是否字段，用于JSON替换
	IsInt      bool     // 能转成整型就设置为true，用""强制指定为字符串，值为false
	MustString bool     // 指定是字符串类型，用""强制指定为字符串，值为true，默认false，非指定
	Ops        []OpType // 操作赋值
}

type ControlType struct {
	CID     string              // 控制流ID
	Args    map[string]int      // 输入参数
	Title   string              // 标题
	Desc    string              // 描述
	Scope   string              // 范围，目前只针对前置条件和后置条件
	Plugins []string            // 插件列表
	Extends map[string][]string // 扩展属性列表
	Vars    map[string]*VarType // 变量
	Flows   []StepType          // 控制流列表
}

type OpType struct {
	CIDs     []int    // 并发请求ID列表
	Key      string   // 字段
	Op       string   // 操作码
	Value    []string // 值列表
	IsList   bool     // 变量是否列表类型
	IsMap    bool     // 变量是否字典类型
	IsStr    bool     // 值类型是字符串
	Continue bool     // 操作待续：true，有后续操作：false，无后续操作
}

type ConstructType struct {
	IsInt    bool         // 值是否整型
	IsList   bool         // 是否列表
	ID       string       // 名称
	Op       string       // 构造操作符：->或者=
	Formats  []FormatType // 构造格式
	Keys     []FormatType // 字段路径、数据变量
	Settings []OpType     // 赋值动作
}

type FormatType struct {
	IsInt      bool   // 能转成整型就设置为true，用""强制指定为字符串，值为false
	IsFloat    bool   // 能转成浮点数就设置为true，用""强制指定为字符串，值为false
	MustString bool   // 指定是字符串类型，用""强制指定为字符串，值为true，默认false，非指定
	Content    string // 内容
}

type GlobalConfig struct {
	Interval int
	Paths    sync.Map // 路径配置表
	IPs      sync.Map // IP列表
	Ports    sync.Map // Ports列表
	Calls    sync.Map // 调用映射
	Any      sync.Map // 任意数据
	Plugins  sync.Map // 插件
	Tasks    sync.Map
	Projects sync.Map
	Init     sync.Map
	Status   sync.Map
}
