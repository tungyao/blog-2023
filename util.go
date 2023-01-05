package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	uc "github.com/tungyao/ultimate-cedar"
	"io"
	"strconv"
	"time"
)

func IoRead(request uc.Request, pt any) error {
	ft, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	request.Body.Close()
	err = json.Unmarshal(ft, pt)
	if err != nil {
		return err
	}
	return nil
}

func pagination(request uc.Request) (int, int) {
	p := request.URL.Query().Get("page")
	if p == "" {
		p = "1"
	}
	l := request.URL.Query().Get("limit")
	if l == "" {
		l = "10"
	}
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)

	return page - 1, limit
}
func Now() int64 {
	return time.Now().Unix()
}
func Sha1(data []byte) string {
	sh := sha1.New()
	sh.Write(data)
	return fmt.Sprintf("%x", sh.Sum(nil))
}
