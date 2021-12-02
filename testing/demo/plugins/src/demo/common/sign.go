package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Sign(ts, nonce, id, secret string, body []byte) []byte {
	var bs [][]byte
	s := []string{ts, nonce, id, secret}
	h := strings.Join(s, "")
	fmt.Println("签名前，h", h, "body:", string(body))
	bs = append(bs, []byte(h), body)
	hb := bytes.Join(bs, []byte(""))
	b := md5.Sum(hb)
	fmt.Println("签名为", hex.EncodeToString(b[:]))
	return []byte(hex.EncodeToString(b[:]))
}
