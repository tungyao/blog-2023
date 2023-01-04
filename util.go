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

func page(request uc.Request) (int, int) {
	page, _ := strconv.Atoi(request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))

	return page - 1, limit
}
