package main

import (
	"flag"

	"github.com/akanshrv/CacheGo/cache"
)

func main() {
	// conn, err := net.Dial("tcp", ":3000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = conn.Write([]byte("SET Foo Bar 4000000"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// select {}
	// return
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

	server := NewServer(opts, cache.NewCache())
	server.Start()
}
