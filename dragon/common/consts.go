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

const (
	DETAIL  = "detail"
	DEBUG   = "debug"
	INFO    = "info"
	WARNING = "warning"
	ERROR   = "error"
)

const (
	PROJECT = "DRAGON"
)

const (
	ModuleProto   = "proto"
	ModuleControl = "control"
	ModuleCommon  = "common"
	ModuleDB      = "db"
	ModuleUtil    = "util"
	ModuleBuiltin = "builtin"
	ModulePlugin  = "plugin"
	ModuleParser  = "parser"
	ModuleRunner  = "runner"
	ModuleBuilder = "builder"
	ModuleRule    = "rule"
	ModuleReport  = "report"
)

const (
	MethodPost    = "POST"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
)

const (
	ArgJar = "-jar"
)

const (
	CmdPython = "python"
	CmdJava   = "java"
)

const (
	ParserGlobal        = "global"
	ParserCommon        = "common"
	ParserTemplates     = "templates"
	ParserPlugins       = "plugins"
	ParserTC            = "testcases"
	ParserRules         = "rules"
	ModuleConfig        = "config"
	ModuleControlFlow   = "control_flow"
	ModulePreCondition  = "pre_condition"
	ModulePostCondition = "post_condition"
	ModuleCheck         = "checkpoints"
)

const (
	SRC = "src"
	BIN = "bin"
	RPC = "rpc"
)

const (
	FuncProfiling = "profiling"
)

const (
	ProtoHTTP = "http"
	ProtoTCP  = "TCP"
	ProtoUDP  = "UDP"
)

const (
	ZERO             = "0"
	HYPHEN           = "-"
	UnderScope       = "_"
	DoubleUnderScope = "__"
	COMMA            = ","
	AND              = "&"
	DOT              = "."
	DoubleQuotes     = "\""
	DoubleQuotes2    = "\"\""
	EMPTY            = ""
	EOF              = "<!-- EOF -->"
	SPACE            = " "
	SLASH            = "/"
	PLUS             = "+"
	STAR             = "*"
	WILDCARD         = ".*"
	QuestionMark     = "?"
	COLON            = ":"
	AT               = "@"
	LeftBracket      = "["
	RightBracket     = "]"
	LeftParentheses  = "("
	RightParentheses = ")"
	LeftBrace        = "{"
	RightBrace       = "}"
	NUMBER           = "#"
	PARALLEL         = "||"
	VerticalBar      = "|"
	BACKSLASH        = "\\"
	NEWLINE          = "\n"
	OK               = "ok"
	BoundaryChars    = "[__0__]"
	NoChars          = " __NO_CHARS__ "
)

const (
	HeaderCT = "Content-Type"
	CTJson   = "application/json;charset=utf-8"
	HCookie  = "Cookie"
)

const (
	SpaceChar        = 32
	DoubleQuotesChar = 34
	SingleQuotesChar = 39
)

const (
	OSWindows = "windows"
)

const (
	PathInit = "/init"
	PathRun  = "/run"
)

// OpCode
const (
	OpEqual      = "="
	NotPresent   = "--"
	ExtraPresent = "++"
	NotEqual     = "!="
	MoreThan     = ">"
	LessThan     = "<"
	OpInput      = "<-"
	OpOr         = "|" // ????????????
	OpEOR        = "^" // ???????????????
	OpAnd        = "&" // ????????????
	OpNot        = "~" // ????????????
	OpConstruct  = "->"
	OpOn         = "on"
	SpaceOn      = " on "
	OR           = "or"
	SpaceOR      = " or "
	OpWith       = "with"
	SpaceWith    = " with "
	OpWithout    = "without"
	SpaceWithout = " without "
	OpMulti      = "*"
	OpDiv        = "/"
	OpSub        = "-"
	OpPlus       = "+"
	OpLP         = "("
	OpRP         = ")"
	OpMod        = "%"  // ??????
	OpPower      = "**" // ?????????
	OpLS         = "<<" // ??????
	OpRS         = ">>" // ??????
	OpColon      = ":"
	OpDot        = "."
	OpDotDot     = ".."
	OpDotDotDot  = "..."
)

// ?????????????????????
const (
	REQUEST  = "Request"
	REQUEST2 = "request"
	REQUEST3 = "REQUEST"
	CTX      = "x"
	CtxP     = "x."
)

// ??????????????????????????????
const (
	ActSend    = "send"
	AKSend     = "s"
	ActReceive = "receive"
	AKReceive  = "r"
	ActDB      = "db"
	ActKeyDB   = "db"
	ActCache   = "cache"
	AKCache    = "c"
	ActHttp    = "http"
	AKHttp     = "h"
	ActHttps   = "https"
	AKHttps    = "hs"
	ActMock    = "mock"
	AKMock     = "m"
	ActAndroid = "android"
	AKAndroid  = "a"
	ActWEB     = "web"
	ActKeyWEB  = "w"
	ActIOS     = "ios"
	ActKeyIOS  = "ios"
	ActROS     = "ros"
	ActKeyROS  = "ros"
	ActCOM     = "com"
	ActKeyCOM  = "com"
)

const (
	KeyArgs    = "a"
	Args       = "args"
	KeyBuffer  = "b"
	Buffer     = "buffer"
	KeyHeader  = "h"
	Header     = "header"
	KeyMainMsg = "msg"
	MainMsg    = "msg"
	KeyPath    = "p"
	Path       = "path"
	KeyResp    = "r"
	Resp       = "resp"
	KeyTopic   = "t"
	Topic      = "topic"
)

// ????????????
const (
	RAFile   = "file"
	RAFileZ  = "_file"
	RATotal  = "total"
	RATotalZ = "_total"
	RAResp   = "resp"
	RACode   = "code"
	RAHead   = "head"
)

const (
	CFUrl    = "URL"
	CFMethod = "METHOD"
	CFName   = "NAME"
)

const (
	SActions = "ACTIONS"
	SAttrs   = "ATTRS"
	SPorts   = "PORTS"
)

const (
	HttpOnReadySending = "http.OnReadySending"
	HttpOnSending      = "http.OnSending"
	HttpOnReceived     = "http.OnReceived"
	HttpOnError        = "http.OnError"
)

const (
	DefaultPort = 9899
	DefaultDNS  = "8.8.8.8:80"
)

const (
	StaticHosts = "HOSTS"
)

// Plugin Archived File
const (
	ZIP = "zip"
	GZ  = "gz"
	XZ  = "xz"
	TGZ = "tgz"
	TXZ = "txz"
	SZ  = "7z"
	TAR = "tar"
)

const (
	ExtJson = ".json"
	ExtTxt  = ".txt"
	ExtLog  = ".log"
	ExtCfg  = ".cfg"
	ExtIni  = ".ini"
	ExtMd   = ".md"
	ExtGo   = ".go"
	ExtJava = ".java"
	ExtBash = ".bash"
	ExtSh   = ".sh"
	ExtZsh  = ".zsh"
	ExtPy   = ".py"
	ExtJar  = ".jar"
	ExtJs   = ".js"
	ExtTs   = ".ts"
	ExtSo   = ".so"
	ExtDLL  = ".dll"
	ExtEXE  = ".exe"
)

const (
	SUCCESSFUL = 0
	FAILED     = 1
	UNDEFINED  = -1
	QUIT       = 1
	DONE       = 1
)

const (
	MsgSuccessful      = "success"
	MsgUnknownRule     = "unknown rule"
	MsgIncompleteRule  = "incomplete rule"
	MsgInvalid         = "invalid rule"
	MsgDuplicatedName  = "duplicated name"
	MsgUnknownFunction = "unknown function"
	MsgExecutionFailed = "execution failed"
	MsgSetupFailed     = "setup failed"
	MsgTeardownFailed  = "teardown failed"
)

const (
	MF  = "makefile"
	MF2 = "Makefile"
	MF3 = "MAKEFILE"
)

const (
	EXECUTING = 0
	PARSING   = 1
)

// ?????????DB?????????
const (
	Set        = 0
	Get        = 1
	Delete     = 2
	GetALLData = 3
	GetALLKeys = 4
)

const (
	PluginsPath        = "plugins_path"
	DBPath             = "db_path"
	TPath              = "testing"
	CasesFolder        = "testcases"
	DefaultPluginsPath = "plugins/dragon"
	DefaultPluginName  = "start.sh"
	DefaultDBPath      = "db/dragon"
	DefaultTPath       = "testing"
)

const (
	LOCALHOST = "localhost"
	PubDNS    = "publicDNS"
	PubIP     = "publicIP"
)

const (
	INIT     = 0
	MaxTry   = 1
	MaxRetry = 10
)

const (
	MarkdownType   = "*.md"
	MarkdownSuffix = ".md"
)

const (
	SETTING          = "# ??????"
	NODE             = "# ??????"
	MessageConstruct = "# ????????????"
	DataConstruct    = "# ????????????"
	PLUGIN           = "# ??????"
	ControlFlow      = "# ?????????"
	SETUP            = "# ????????????"
	TEARDOWN         = "# ????????????"
	CHECK            = "# ??????"
	CASE             = "# ??????"
	FLOW             = "# ??????"
)

const (
	SettingEN          = "# Setting"
	NodeEN             = "# Node"
	MessageConstructEN = "# Message Construct"
	DataConstructEN    = "# Data Construct"
	PluginEN           = "# Plugin"
	ControlFlowEN      = "# Control Flow"
	SetupEN            = "# Setup"
	TeardownEN         = "# Teardown"
	CheckEN            = "# Check"
	CaseEN             = "# Case"
	FlowEN             = "# Flow"
)

// ????????????
const (
	CaseVar = 0
	CFVar   = 1
	FlowVar = 2
)

const (
	MarkRegion    = "# "
	MarkRegion2   = "## "
	MarkRegion3   = "### "
	MarkClass     = "> "
	MarkClass2    = ">> "
	MarkComment   = "***"
	MarkComment2  = "//"
	MarkItem      = "* "
	MarkConstruct = "->"
	MarkExpect    = "=>"
)

const (
	MaxRange   = 100
	MaxQueue   = 6
	MaxBinRead = 100
)

const (
	DefaultTimeout = int64(0) // ????????????
)

const (
	// Random??????
	FuncNameUUID         = "UUID"
	FuncNameRandom       = "Random"
	FuncNameRandomString = "RandomString"
	FuncNameRandomDigit  = "RandomDigit"
	// Convert??????
	FuncNameFormatDate         = "FormatDate"
	FuncNameDATE               = "DATE"
	FuncNameFormatTime         = "FormatTime"
	FuncNameTIME               = "TIME"
	FuncNameFormatSeconds      = "FormatSeconds"
	FuncNameSEC                = "SEC"
	FuncNameFormatNowInSeconds = "FormatNowInSeconds"
	FuncNameNOW                = "NOW"
	FuncNameNowAddSeconds      = "NowAddSeconds"
	FuncNameNowMS              = "NowMS"
	FuncNameMD5                = "MD5"
	FuncNameSHA1               = "SHA1"
	FuncNameBase64             = "Base64"
	// Testing??????
	FuncNameRequest = "Request"
	// ??????
	FuncNameLog = "Log"
)

// ????????????Task??????
const (
	WaitTask = 0
	RunTask  = 1
)

// Random????????????
const (
	CHARS              = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DIGITS             = "0123456789"
	DefaultCharsLength = 10
)

// ????????????
const (
	ALL    = "all"
	RANDOM = "random"
)

const (
	FillRandom = 0
	FillAll    = 1
)

const (
	ScopeGroup   = "G"
	ScopeCase    = "C"
	ScopeSuite   = "S"
	ScopeAll     = "A"
	ScopeDiscard = "D"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

const (
	AuthorCN      = "??????"
	AuthorEN      = "author"
	MaintainerCN  = "?????????"
	MaintainerEN  = "maintainer"
	BugsCN        = "??????"
	BugsEN        = "bugs"
	FeaturesCN    = "??????"
	FeaturesEN    = "features"
	VersionCN     = "??????"
	VersionEN     = "version"
	TagsCN        = "??????"
	TagsEN        = "tags"
	LevelCN       = "??????"
	LevelEN       = "level"
	ModeCN        = "??????"
	ModeEN        = "mode"
	ProtoCN       = "??????"
	ProtoEN       = "proto"
	ClosedCN      = "??????"
	ClosedEN      = "closed"
	HostCN        = "??????"
	HostEN        = "host"
	PortCN        = "??????"
	PortEN        = "port"
	TimeoutCN     = "??????"
	TimeoutEN     = "T"
	ConcurrenceCN = "??????"
	ConcurrenceEN = "C"
	RoundsCN      = "??????"
	RoundsEN      = "R"
	FillCN        = "??????"
	FillEN        = "fill"
	NameCN        = "??????"
	NameEN        = "name"
	StartCN       = "??????"
	StartEN       = "start"
	BuildCN       = "??????"
	BuildEN       = "build"
	CallsCN       = "??????"
	CallsEN       = "calls"
	MutesCN       = "??????"
	MutesEN       = "mutes"
	HostsCN       = "??????"
	HostsEN       = "hosts"
	ProxyCN       = "??????"
	ProxyEN       = "proxy"
	PluginsCN     = "??????"
	PluginsEN     = "plugins"
	ArgsCN        = "??????"
	ArgsEN        = "args"
	MatchCN       = "??????"
	MatchEN       = "match"
	ScopeCN       = "??????"
	ScopeEN       = "scope"
	SetupsCN      = "??????"
	SetupsEN      = "setup"
	TeardownsCN   = "??????"
	TeardownsEN   = "teardown"
	ExtendCN      = "??????"
	ExtendEN      = "E"
	UserCN        = "??????"
	UserEN        = "user"
	PassCN        = "??????"
	PassEN        = "pass"
	PathCN        = "??????"
	PathEN        = "path"
	MethodCN      = "??????"
	MethodEN      = "method"
	TunnelCN      = "??????"
	TunnelEN      = "tunnel"
	TunUserCN     = "????????????"
	TunUserEN     = "sshuser"
	TunPassCN     = "????????????"
	TunPassEN     = "sshpass"
	RandomCN      = "??????"
	RandomEN      = "random"
	LocalCN       = "??????"
	LocalEN       = "local"
)

const (
	STR       = "STR"
	BYTES     = "BYTES"
	INT       = "INT"
	UINT      = "UINT"
	FLOAT     = "FLOAT"
	BOOL      = "BOOL"
	LIST      = "LIST"
	ListBytes = "ListBytes"
	MAP       = "MAP"
	MapBytes  = "MapBytes"
	NULL      = "NULL"
)
