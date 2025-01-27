package main

import (
	"context"
	"flag"
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

		for i := 0; i < 10; i++ {
			SendCommand(client)

		}
		client.Close()
		time.Sleep(time.Millisecond * 1)

	}()

	server := NewServer(opts, cache.NewCache())
	server.Start()
}

func SendCommand(c *client.Client) {

	_, err := c.Set(context.Background(), []byte("ak"), []byte("akanksh"), 0)
	if err != nil {
		log.Fatal(err)
	}

}
