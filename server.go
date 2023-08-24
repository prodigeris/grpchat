package main

import (
	"github.com/prodigeris/grpchat/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var servers []proto.Chat_StreamServer

func serve() {
	newServer()
}

type server struct {
	proto.UnimplementedChatServer
}

func (s server) Stream(srv proto.Chat_StreamServer) error {
	log.Println("start new server")
	servers = append(servers, srv)
	ctx := srv.Context()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		// continue if number reveived from stream
		// less than max
		if req.Message == "Arnas" {
			continue
		}
		resp := proto.MessageResponse{Message: req.Message, From: req.From}
		for _, serv := range servers {
			if srv == serv {
				continue
			}
			if err := serv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("send %s", resp.Message)
		}
	}
}

func newServer() {
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterChatServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
