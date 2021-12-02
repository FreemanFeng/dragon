package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func CreateFile(name string, data []byte, paths ...string) string {
	h := strings.Join(paths, SLASH)
	if _, err := os.Stat(h); os.IsNotExist(err) {
		os.MkdirAll(h, os.ModePerm)
		Debug(ModuleUtil, "Created Path", h)
	}
	s := []string{h, name}
	path := strings.Join(s, SLASH)
	err := ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		Debug(ModuleUtil, "Could not Write File to", path, err)
		return EMPTY
	}
	Debug(ModuleUtil, "Write File", path, "successfully!")
	return path
}

func GetDirList(path string) ([]string, error) {
	var list []string
	err := filepath.Walk(path,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				list = append(list, path)
				return nil
			}

			return nil
		})
	return list, err
}

func GetFiles(path, wildcard string) []string {
	var files []string
	dirLists, err := GetDirList(path)
	if err != nil {
		Err(ModuleUtil, err)
		return []string{}
	}
	for _, v := range dirLists {
		s := []string{v, wildcard}
		t, err := filepath.Glob(strings.Join(s, string(os.PathSeparator)))
		if err != nil {
			Err(ModuleUtil, err)
			continue
		}
		files = append(files, t...)
		Debug(ModuleUtil, files)
	}
	return files
}
