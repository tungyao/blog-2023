package main

import (
	"database/sql"
	"sync"
)

// Cache 简易的cache功能
type Cache struct {
	mutex sync.RWMutex
	data  map[string]*Post
}

func NewCache() *Cache {
	return &Cache{
		mutex: sync.RWMutex{},
		data:  make(map[string]*Post),
	}
}

func (n *Cache) Update(name string, data *Post) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.data[name] = data
}
func (n *Cache) Get(name string) *Post {
	n.mutex.RLock()
	data := n.data[name]
	if data == nil {
		n.mutex.RUnlock()
		n.mutex.Lock()
		row := Db.QueryRow("select * from post where name=?", name)
		dt := &Post{}
		err := row.Scan(&dt.Name, &dt.Data, &dt.CreateTime, &dt.UpdateTIme, &dt.Permission, &dt.IsDelete)
		if err != sql.ErrNoRows {
			n.data[name] = dt
		}
		n.mutex.Unlock()
		return dt
	}
	n.mutex.RUnlock()
	return data
}

// 从数据库中读取缓存

func ReadFromDb() {
	Db.Query("select * from ")
}
