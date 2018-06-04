package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

func randByte(cnt int) (res []byte, err error) {
	var f *os.File
	f, err = os.Open("/dev/urandom")
	if nil == err {
		defer f.Close()
		res = make([]byte, cnt)
		var tmp []byte
		tmp = res
		sz := 0
		n := 0
		for sz < cnt {
			n, err = f.Read(tmp)
			fmt.Println(cnt, n, sz)
			if nil == err {
				tmp = tmp[n:]
				sz += n
			} else {
				break
			}
		}
	}
	return
}

func comper(s1, s2 []byte) (r bool) {
	le := 0
	if len(s1) < len(s2) {
		le = len(s1)
	} else {
		le = len(s2)
	}
	var diff byte
	i := 0
	for i < le {
		diff |= (s1[i] ^ s2[i])
		i++
	}
	return 0 == diff
}

func genSlotAndHash(password string) (slot, pa []byte, err error) {
	passwd, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return
	}
	slot, err = randByte(256)
	if err != nil {
		return
	}
	h := sha256.New()
	h.Write(slot[:56])
	h.Write(passwd)
	h.Write(slot[56:])
	pa = h.Sum(nil)
	return
}

func calReqHash(slot []byte, password string) (pa []byte, err error) {
	passwd, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return
	}
	h := sha256.New()
	h.Write(slot[:56])
	h.Write(passwd)
	h.Write(slot[56:])
	pa = h.Sum(nil)
	return
}
