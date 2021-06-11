package main

import (
	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/Nguyen-Hoang-Nam/let-shorten/server"
)

func main() {
	db.Init()
	db.RedisInit()
	server.Init()
}
