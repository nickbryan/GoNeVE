package network

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/nickbryan/GoNeVE/pkg/engine"

	"github.com/nickbryan/GoNeVE/pkg/network/grpc/impl"

	"google.golang.org/grpc/test/bufconn"

	"github.com/nickbryan/GoNeVE/pkg/network/grpc/domain"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/service"
	"google.golang.org/grpc"
)

func TestNewClient(t *testing.T) {
	t.Run("returns a pointer to a new Client", func(t *testing.T) {
		got := reflect.TypeOf(NewClient(":"))
		want := reflect.TypeOf(&Client{})
		if got != want {
			t.Errorf("got %v expected %T", got, want)
		}
	})
}

func TestClient(t *testing.T) {
	t.Run("can handle grpc usage", func(t *testing.T) {
		c := NewClient("bufnet")

		l := bufconn.Listen(1024 * 1024)
		srv := grpc.NewServer()
		service.RegisterPlayerServer(srv, &impl.Player{EntityManager: engine.NewEntityManager()})
		go func() {
			if err := srv.Serve(l); err != nil {
				t.Fatalf("unable to start server: %v", err)
			}
		}()

		err := c.Dial(grpc.WithContextDialer(func(context.Context, string) (conn net.Conn, e error) {
			return l.Dial()
		}))
		if err != nil {
			t.Fatalf("client Dial failed: %v", err)
		}

		var resp *domain.LoginResponse
		c.Handle(HandlerFunc(func(cc *grpc.ClientConn) {
			pc := service.NewPlayerClient(cc)

			resp, err = pc.Login(context.Background(), &domain.LoginRequest{})

			if err != nil {
				t.Fatalf("error logging in to server: %v", err)
			}

			if resp.Status != domain.LoginResponse_OK {
				t.Errorf("unable to login to server: %s", resp.Message)
			}
		}))
		if resp == nil {
			t.Errorf("handler was not called: response is nil")
		}
	})
}
