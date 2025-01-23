package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/akanshrv/CacheGo/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}
type Server struct {
	ServerOpts

	followers map[net.Conn]struct{} //only key (address) matters not the value

	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return fmt.Errorf("listen error; %s", err)
	}
	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	//if not leader then dial leader
	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("connected with leader:", s.LeaderAddr)
			s.handleConn(conn)
		}()

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)

	if s.IsLeader {
		s.followers[conn] = struct{}{}
	}
	fmt.Println("connection made:", conn.RemoteAddr())

	for {

		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		go s.handleCommand(conn, buf[:n])
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {

	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command", err)
		conn.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("recieved command: %s\n ", msg.Cmd)

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCmd(conn, msg)
	case CMDGet:
		err = s.handleGetCmd(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command: ", err)
		conn.Write([]byte(err.Error()))

	}
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}
	_, err = conn.Write(val)

	return err
}

func (s *Server) handleSetCmd(con net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}
	go s.sendToFollowers(context.TODO(), msg)
	return nil
}

// if data is deleted( due to failure of the node ) it is sent to followers
func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Println("forwarding to follower")
		rawMsg := msg.ToBytes()
		fmt.Println("Forwarding rawmsg to follower", string(rawMsg))
		_, err := conn.Write(msg.ToBytes())
		if err != nil {
			log.Println("write to follower error: ", err)
			continue
		}
	}
	return nil
}
