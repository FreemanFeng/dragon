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

import "strings"

func SearchLongestMatchKey(s []string, m map[string]int, b string) string {
	r := EMPTY
	n := len(s)
	if n == 1 {
		key := s[0]
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	for i := 1; i < n; i++ {
		key := strings.Join(s[:i], b)
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	return r
}

func SearchLongestMatchNode(s []string, m map[string]*NodeType, b string) string {
	r := EMPTY
	n := len(s)
	if n == 1 {
		key := s[0]
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	for i := 1; i < n; i++ {
		key := strings.Join(s[:i], b)
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	return r
}

func SearchLongestMatchSetting(s []string, m map[string][]string, b string) string {
	r := EMPTY
	n := len(s)
	if n == 1 {
		key := s[0]
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	for i := 1; i < n; i++ {
		key := strings.Join(s[:i], b)
		_, ok := m[key]
		if ok {
			r = key
		}
	}
	return r
}

func GetFlow(key string, c *ConfigType, st *TestSuite) *FlowType {
	if st == nil {
		k, ok := c.Flows[key]
		if !ok {
			return nil
		}
		return k
	}
	k, ok := st.Flows[key]
	if !ok {
		k, ok = c.Flows[key]
	}
	if !ok {
		return nil
	}
	return k
}

func IsFlow(key string, c *ConfigType, st *TestSuite) bool {
	if st == nil {
		_, ok := c.Flows[key]
		return ok
	}
	_, ok := st.Flows[key]
	if !ok {
		_, ok = c.Flows[key]
	}
	return ok
}

func IsGroup(key string, vt map[string]*VarType) bool {
	_, ok := vt[key]
	return ok
}
