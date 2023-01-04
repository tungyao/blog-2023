package main

import uc "github.com/tungyao/ultimate-cedar"

type Post struct {
	Name       string `json:"name"`
	Data       string `json:"data"`
	CreateTime int    `json:"create_time"`
	UpdateTIme int    `json:"update_t_ime"`
	Permission int    `json:"permission"`
	IsDelete   int    `json:"is_delete"`
}

func Index(writer uc.ResponseWriter, request uc.Request) {
	//data := Caches.Get(request.Data.Get("name"))

}
