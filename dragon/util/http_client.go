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

package util

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func ConvertUrl(reqUrl string, headers http.Header) string {
	re := regexp.MustCompile(`(?P<scheme>^http.*?)//(?P<domain>\w.*?)/(?P<path>.*)`)
	match := re.FindStringSubmatch(reqUrl)
	groupNames := re.SubexpNames()
	result := make(map[string]string)
	for i, name := range groupNames {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	t := strings.Split(result["domain"], AT)
	if len(t) > 1 {
		result["domain"] = t[1]
		s := []string{BASIC, Base64Encode([]byte(t[0]))}
		headers[AUTHORIZATION] = []string{strings.Join(s, SPACE)}
	}
	s := []string{result["scheme"], "//", result["domain"], "/", result["path"]}
	return strings.Join(s, EMPTY)
}

func RequestHttp(reqUrl, method string, a, h, msg []byte) ([]byte, int, http.Header) {
	args := url.Values{}
	tm := ToJSON(a)
	if tm != nil {
		am := tm.(map[string]interface{})
		for k, v := range am {
			args[k] = []string{ToString(v)}
		}
	}
	if len(args) > 0 {
		s := []string{reqUrl, args.Encode()}
		reqUrl = strings.Join(s, QuestionMark)
	}
	headers := http.Header{}
	tm = ToJSON(h)
	if tm != nil {
		hm := tm.(map[string]interface{})
		for k, v := range hm {
			if b, ok := v.([]interface{}); ok {
				headers[k] = []string{}
				for _, x := range b {
					headers[k] = append(headers[k], ToString(x))
				}
			} else {
				headers[k] = []string{ToString(v)}
			}
		}
	}
	_, ok := headers[HeaderCT]
	if !ok && method == MethodPost {
		headers[HeaderCT] = []string{CTJson}
	}
	Debug(ModuleCommon, "请求参数", args)
	Debug(ModuleCommon, "请求头", headers)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	r := ConvertUrl(reqUrl, headers)
	req, err := http.NewRequest(method, r, bytes.NewReader(msg))
	if err != nil {
		Err(err)
		return []byte(""), http.StatusBadRequest, http.Header{}
	}
	defer req.Body.Close()

	for k, s := range headers {
		for _, v := range s {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		Err(err)
		return []byte(""), http.StatusInternalServerError, http.Header{}
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Err(err)
		return []byte(""), http.StatusInternalServerError, http.Header{}
	}
	Info(ModuleUtil, method, " ", r, " [", resp.Status, "]")
	Debug(ModuleUtil, "Request Header:", headers)
	Info(ModuleUtil, "Response Header:", resp.Header)
	Debug(ModuleUtil, "Response Body:", string(b))
	return b, resp.StatusCode, resp.Header
}
