package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/plitn/wb_school_l0/reader"
	"github.com/plitn/wb_school_l0/server"
	"github.com/plitn/wb_school_l0/storage"
	"log"
)

func main() {
	c := storage.NewCache()
	err := c.Init()
	if err != nil {
		log.Println("cache init error")
		return
	}
	fmt.Println("cache inited")
	np := reader.NewNatsReader(c)
	fmt.Println("nats created")
	err = np.Init()
	fmt.Println("nats inited")

	s := server.NewServer(c)
	s.Init()
	fmt.Println("server inited")

	for {
		if err != nil {
			return
		}
	}

}
