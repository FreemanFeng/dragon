package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func Sign(params ...interface{}) interface{} {
	var bs [][]byte
	if len(params) < 5 {
		return ""
	}
	i := params[0].(int64)
	ts := strconv.FormatInt(i, 10)
	i = params[1].(int64)
	nonce := strconv.FormatInt(i, 10)
	x := params[2].(int)
	id := strconv.Itoa(x)
	secret := params[3].(string)
	body := params[4].([]byte)
	s := []string{ts, nonce, id, secret}
	k := strings.Join(s, "")
	bs = append(bs, []byte(k), body)
	h := bytes.Join(bs, []byte(""))
	b := md5.Sum(h)
	return hex.EncodeToString(b[:])
}
