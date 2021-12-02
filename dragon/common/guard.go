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

import "time"

func WaitForQuit(sch, done, stop, quit chan int) {
	signal := UNDEFINED
	for {
		select {
		case <-sch:
			// 如果信号为QUIT，需要分别通知TCP/UDP/TLS层以及上层协议（i.e. HTTP）
			quit <- signal
			if signal == QUIT {
				Detail(ModuleCommon, "signal quit")
				stop <- signal
				return
			}
		case <-done:
			Detail(ModuleCommon, "quit service")
			signal = QUIT
		}
	}
}

func Finish(done chan int) {
	done <- DONE
}

func WaitTimeout(duration int, ch chan int, timeout chan bool) {
	if duration == 0 {
		return
	}
	for {
		<-ch
		time.Sleep(time.Duration(duration) * time.Second)
		timeout <- true
	}
}
