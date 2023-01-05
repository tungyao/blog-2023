package main

import (
	"encoding/json"
	uc "github.com/tungyao/ultimate-cedar"
	"io"
	"strconv"
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
