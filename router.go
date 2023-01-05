package main

import (
	"crypto/sha1"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	uc "github.com/tungyao/ultimate-cedar"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	Name       string `json:"name"`
	Data       string `json:"data,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
	Permission int    `json:"permission,omitempty"`
	IsDelete   int    `json:"is_delete,omitempty"`
}

type IndexPageData struct {
	Count int
	Data  []*Post
}

//go:embed static
var indexHtml embed.FS

func Index(writer uc.ResponseWriter, request uc.Request) {
	page, limit := pagination(request)
	pageCache := Spruce.Get([]byte(fmt.Sprintf("%s %d %d", "index", page, limit)))
	if pageCache == nil {
		out := make([]*Post, 0)
		rows, _ := Db.Query("select * from post where is_delete=0 limit ? offset ?", limit, limit*page)
		for rows.Next() {
			p := &Post{}
			err := rows.Scan(&p.Name, &p.Data, &p.CreateTime, &p.UpdateTime, &p.Permission, &p.IsDelete)
			if err == nil {
				out = append(out, p)
			}
		}
		var count int
		row := Db.QueryRow("select count(*) from post where is_delete=0")
		row.Scan(&count)
		c := &IndexPageData{
			Count: count,
			Data:  out,
		}
		pageCache = c
		if request.URL.Query().Get("plat") == "api" {

			data, _ := json.Marshal(out)
			writer.Data(data).Send()
			Spruce.Set([]byte(fmt.Sprintf("%s %d %d", "index", page, limit)), out, time.Now().Unix()+3600)
			return
		}
	}
	if request.URL.Query().Get("plat") == "api" {
		writer.Data(pageCache).Send()
		return
	}
	// 渲染静态界面
	t, err := template.ParseFS(indexHtml, "static/index.html")
	log.Println(err)
	err = t.Execute(writer, pageCache)
	log.Println(err)
}

func OnlyOne(writer uc.ResponseWriter, request uc.Request) {
	data := Spruce.Get([]byte(request.Data.Get("name")))
	if data == nil {
		writer.Data(`{}`).Status(404).Send()
		return
	}
	writer.Data(data).Send()
}

func MgPostGet(writer uc.ResponseWriter, request uc.Request) {
	page, limit := pagination(request)
	out := make([]*Post, 0)
	rows, _ := Db.Query("select * from post limit where is_delete=0 ? offset ?", limit, limit*page)
	for rows.Next() {
		p := &Post{}
		err := rows.Scan(&p.Name, &p.Data, &p.CreateTime, &p.UpdateTime, &p.Permission, &p.IsDelete)
		if err == nil {
			out = append(out, p)
		}
	}
	data, _ := json.Marshal(out)
	writer.Data(data).Send()
}

func MgPostAdd(writer uc.ResponseWriter, request uc.Request) {
	pt := &Post{}
	if err := IoRead(request, pt); err != nil {
		writer.Data(err).Send()
		return
	}
	res, _ := Db.Exec("insert into post (`name`,`data`,`create_time`,`update_time`,`permission`,`is_delete`) values (?,?,?,?,?,?)",
		pt.Name, pt.Data, Now(), Now(), 1, 0,
	)
	id, _ := res.LastInsertId()
	Caches.Update(pt.Name, pt)
	writer.Data(fmt.Sprintf(`{"id":%d}`, id)).Send()
}
func MgPostUpdate(writer uc.ResponseWriter, request uc.Request) {
	pt := &Post{}
	if err := IoRead(request, pt); err != nil {
		writer.Data(err).Send()
		return
	}
	res, _ := Db.Exec("update post set name=?,data=?,update_time=?,permission=?", pt.Name, pt.Data, Now(), pt.Permission)
	id, _ := res.RowsAffected()
	preCache := Caches.Get(pt.Name)
	preCache.Name = pt.Name
	preCache.Data = pt.Data
	preCache.UpdateTime = Now()
	preCache.Permission = pt.Permission
	Caches.Update(pt.Name, preCache)
	writer.Data(fmt.Sprintf(`{"id":%d}`, id)).Send()
}
func MgPostDelete(writer uc.ResponseWriter, request uc.Request) {
	name := request.URL.Query().Get("name")
	res, _ := Db.Exec("update postt set is_delete=1 where name=?", name)
	id, _ := res.RowsAffected()
	Caches.Delete(name)
	writer.Data(fmt.Sprintf(`{"id":%d}`, id)).Send()
}
func FileUpload(writer uc.ResponseWriter, request uc.Request) {
	fv, header, _ := request.FormFile("file")
	fs, err := os.OpenFile(FilePathPrefix+header.Filename, os.O_CREATE|os.O_WRONLY, 775)
	if err != nil {
		log.Println(err)
		writer.Data(err).Send()
		return
	}
	rd := io.TeeReader(fv, fs)
	preData, err := io.ReadAll(rd)
	if err != nil {
		log.Println(err)
		writer.Data(err).Send()
		return
	}

	_, err = io.Copy(fs, fv)
	if err != nil {
		log.Println(err)
		writer.Data(err).Send()
		return
	}
	fs.Close()
	fv.Close()
	lastExt := strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]
	newName := Sha1(preData) + "." + lastExt
	os.Rename(FilePathPrefix+header.Filename, FilePathPrefix+newName)
	writer.Data(fmt.Sprintf(`{"url":"%s"}`, newName)).Send()
}
func FileUpdate(writer uc.ResponseWriter, request uc.Request) {

}

type User struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

var user = make(map[string]*User)

// Login 信息保存在内存中
func Login(writer uc.ResponseWriter, request uc.Request) {
	u := User{}
	if err := IoRead(request, &u); err != nil {
		http.Error(writer, err.Error(), 403)
		return
	}
	row := Db.QueryRow("select count(*) from user where name=? and pass=?", u.Name, u.Pass)
	yes := 0
	row.Scan(&yes)
	if yes == 0 {
		http.Error(writer, "check failed", 401)
		return
	}
	sh := sha1.New()
	sh.Write([]byte(time.Now().String() + u.Pass + u.Name))
	token := fmt.Sprintf("%x", sh.Sum(nil))
	user[token] = &u
	writer.Data(fmt.Sprintf(`{"token":"%s"}`, token)).Send()
}
