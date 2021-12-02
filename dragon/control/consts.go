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

package control

const (
	PathConfig      = "/config/:key/:value"
	PathTestProject = "/test/:project"
	PathStopProject = "/stop/:method/:project"
	PathTestReport  = "/report/:task"
	PathCallPlugin  = "/plugin/:plugin/:function"
	PathWSTask      = "/ws/:task"
	PathStopShow    = "/stopShow"
)

const (
	ParamPath        = "path"
	ParamProject     = "project"
	ParamSuite       = "s"
	ParamGroup       = "g"
	ParamCase        = "c"
	ParamTags        = "t"
	ParamSmoking     = "smoking"
	ParamRandom      = "random"
	ParamNoCheck     = "nocheck"
	ParamDebug       = "d"
	ParamHosts       = "h"
	ParamExtraHosts  = "h_"
	ParamConcurrence = "p"
	ParamTotal       = "n"
	ParamDuration    = "d"
	ParamTask        = "task"
	ParamReport      = "/report"
	ParamPlugin      = "plugin"
	ParamFunction    = "function"
	ParamKey         = "key"
	ParamValue       = "value"
)

const (
	LOG      = "log"
	WS       = "ws"
	TOTAL    = "total"
	INTERVAL = "interval"
)

const (
	CTTEXT = "text/html; charset=UTF-8"
	CTJSON = "application/json; charset=utf-8"
)

const (
	WsInitTask  = "INIT TASK"
	WsBeginTask = "BEGIN TASK"
	WsWaitTask  = "WAIT TASK"
)
