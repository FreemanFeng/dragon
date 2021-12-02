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
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func TestingPath(path string) string {
	prefix := GetPath(TPath, DefaultTPath)
	return strings.Join([]string{prefix, path, CasesFolder}, SLASH)
}

func ParserPath(top, project, parser string, params ...string) string {
	prefix := GetPath(TPath, DefaultTPath)
	s := []string{prefix, project, parser}
	s = append(s, params...)
	return strings.Join(s, string(os.PathSeparator))
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径文件夹是否存在
func IsDir(path string) bool {
	k, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return k.IsDir()
		}
		return false
	}
	return k.IsDir()
}

func IsFile(path string) bool {
	k, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return k.IsDir() == false
		}
		return false
	}
	return k.IsDir() == false
}

func IsNewer(new, old string) bool {
	sn, _ := os.Stat(new) //os.Stat获取文件信息
	so, _ := os.Stat(old) //os.Stat获取文件信息
	return sn.ModTime().Unix() > so.ModTime().Unix()
}

func GetBinPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return EMPTY
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return EMPTY
	}
	i := strings.LastIndex(path, string(os.PathSeparator))
	if i < 0 {
		return EMPTY
	}
	return path[0 : i+1]
}

func IsGoCode(p *PluginType) bool {
	s := []string{p.Path, string(os.PathSeparator), p.Name, ExtGo}
	h := strings.Join(s, EMPTY)
	return IsFile(h)
}
