package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/akanshrv/CacheGo/cache"
	"github.com/akanshrv/CacheGo/client"
)

func main() {

	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen addres of the leader")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   *leaderAddr == "",
		LeaderAddr: *leaderAddr,
	}
	go func() {
		time.Sleep(time.Second * 2)
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		err = client.Set(context.Background(), []byte("foo"), []byte("bar"), 0)
		if err != nil {
			log.Fatal(err)
		}

		value, err := client.Get(context.Background(), []byte("foo"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(value))

		client.Close()

	}()

	server := NewServer(opts, cache.NewCache())
	server.Start()
}
