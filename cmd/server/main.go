package main

import (
	"log"

	"github.com/nickbryan/GoNeVE/pkg/network/grpc/service"

	"github.com/nickbryan/GoNeVE/pkg/engine"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/impl"

	"google.golang.org/grpc"

	"github.com/nickbryan/GoNeVE/pkg/network"
)

func main() {
	srv := network.NewServer(":7777")

	srv.RegisterServices(func(server *grpc.Server) {
		service.RegisterPlayerServer(server, &impl.Player{EntityManager: engine.NewEntityManager()})
	})

	log.Fatal(srv.ListenAndServe())
}
