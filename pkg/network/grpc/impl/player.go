package impl

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/nickbryan/GoNeVE/pkg/engine"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/domain"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/service"
)

type Player struct {
	EntityManager *engine.EntityManager
}

func (p *Player) Login(ctx context.Context, r *domain.LoginRequest) (*domain.LoginResponse, error) {
	e := p.EntityManager.Create()

	return &domain.LoginResponse{
		Status: domain.LoginResponse_OK,
		Id:     e.ID,
	}, nil
}

func (p *Player) SyncMovement(stream service.Player_SyncMovementServer) error {
	for {
		mc, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// TODO: error hand;e
		e := p.EntityManager.Get(mc.Id)
		log.Println(e)
	}
}

func (p *Player) Handle(cc *grpc.ClientConn) {
	c := service.NewPlayerClient(cc)

	resp, err := c.Login(context.Background(), &domain.LoginRequest{})
	if err != nil {
		log.Fatalf("unable to login: %v", err)
	}
	if resp.Status != domain.LoginResponse_OK {
		log.Fatalf("got %d expected %d", resp.Status, domain.LoginResponse_OK)
	}

	log.Printf("logged in: %s", resp.Id)

}
