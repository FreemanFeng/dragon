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
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

func Decompress(dst, src string) (int, []string) {
	t := ZipType(src)
	switch t {
	case ZIP:
		return Unzip(dst, src)
	case TGZ:
		return GunzipTar(dst, src)
	}
	return FAILED, []string{}
}

func Unzip(dst, src string) (int, []string) {
	var files []string
	r, err := zip.OpenReader(src)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED, files
	}
	defer r.Close()
	rm := map[io.ReadCloser]int{}
	defer closeReadCloser(rm)
	for _, f := range r.Reader.File {
		Debug(ModuleCommon, "Extracting file", f.Name, "...")
		files = append(files, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(Join(SLASH, dst, f.Name), 0755)
			if err != nil {
				Err(ModuleCommon, err)
				return FAILED, files
			}
			continue
		}
		rc, err := f.Open()
		if err != nil {
			Err(ModuleCommon, err)
			return FAILED, files
		}
		rm[rc] = 1
		c := extractFile(dst, f.Name, f.Mode(), rc)
		if c != SUCCESSFUL {
			return c, files
		}
	}
	os.Remove(src)
	return SUCCESSFUL, files
}

func GunzipTar(dst, src string) (int, []string) {
	var files []string
	fr, err := os.Open(src)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED, files
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED, files
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			Err(ModuleCommon, err)
			return FAILED, files
		}
		Debug(ModuleCommon, "Extracting file", h.Name, "...")
		files = append(files, h.Name)
		if h.FileInfo().IsDir() {
			err := os.MkdirAll(Join(SLASH, dst, h.Name), 0755)
			if err != nil {
				Err(ModuleCommon, err)
				return FAILED, files
			}
			continue
		}
		c := extractFile(dst, h.Name, os.FileMode(h.Mode), tr)
		if c != SUCCESSFUL {
			return c, files
		}
	}
	os.Remove(src)
	return SUCCESSFUL, files
}

func extractFile(dst, name string, mode os.FileMode, rc io.Reader) int {
	filename := Join(SLASH, dst, name)
	p := GetDir(filename)
	err := os.MkdirAll(p, 0755)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED
	}
	//w, err := os.Create(filename)
	w, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	if err != nil {
		Err(ModuleCommon, err)
		return FAILED
	}
	return SUCCESSFUL
}

func GetDir(path string) string {
	p := strings.Split(path, SLASH)
	n := len(p) - 1
	return strings.Join(p[:n], SLASH)
}

func closeReadCloser(rm map[io.ReadCloser]int) {
	for r := range rm {
		r.Close()
	}
}
