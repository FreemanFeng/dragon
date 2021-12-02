package common

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"strings"
)

func Base64Decode(data string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		Log("decode error:", err)
		return []byte{}
	}
	return decoded
}

func Base64Encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

func MD5String(data ...string) string {
	b := MD5Sum(data...)
	return hex.EncodeToString(b[:])
}

func MD5Sum(data ...string) [16]byte {
	h := append([]string{}, data...)
	k := strings.Join(h, "")
	return md5.Sum([]byte(k))
}

func SHA1Sum(data ...string) []byte {
	h := sha1.New()
	for _, k := range data {
		h.Write([]byte(k))
	}
	return h.Sum(nil)
}

func GobEncode(v interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		Err(ModuleCommon, "failed in gob encode", err)
		return []byte{}
	}
	return buf.Bytes()
}

func GobDecode(p []byte, v interface{}) error {
	var buf bytes.Buffer
	buf.Write(p)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(v)
	if err != nil {
		Err(ModuleCommon, "failed in gob decode", err)
	}
	return err
}
