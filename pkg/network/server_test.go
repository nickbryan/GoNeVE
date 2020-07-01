package network

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/nickbryan/GoNeVE/pkg/network/grpc/domain"

	"github.com/nickbryan/GoNeVE/pkg/engine"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/impl"
	"github.com/nickbryan/GoNeVE/pkg/network/grpc/service"

	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	t.Run("returns a pointer to a new Server", func(t *testing.T) {
		got := reflect.TypeOf(NewServer(":"))
		want := reflect.TypeOf(&Server{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

	t.Run("sets the Addr on the server", func(t *testing.T) {
		addr := ":7777"
		srv := NewServer(addr)
		got := srv.Addr
		want := addr
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

func TestServer(t *testing.T) {
	t.Run("can register grpc services", func(t *testing.T) {
		l := bufconn.Listen(1024 * 1024)

		srv := NewServer("bufnet")
		srv.RegisterServices(func(server *grpc.Server) {
			service.RegisterPlayerServer(server, &impl.Player{EntityManager: engine.NewEntityManager()})
		})
		go func() {
			if err := srv.Serve(l); err != nil {
				t.Fatalf("unable to start Server via Serve: %v", err)
			}
		}()

		conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (conn net.Conn, e error) {
			return l.Dial()
		}), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("unable to dial bufnet: %v", err)
		}
		c := service.NewPlayerClient(conn)
		resp, err := c.Login(context.Background(), &domain.LoginRequest{})
		if err != nil {
			t.Errorf("unable to login: %v", err)
		}
		if resp.Status != domain.LoginResponse_OK {
			t.Errorf("got %d expected %d", resp.Status, domain.LoginResponse_OK)
		}
	})

	//t.Run("allows a grpc service to Queue an event on the network", func(t *testing.T) {
	//
	//})
}

func TestServer_ListenAndServe(t *testing.T) {
	t.Run("panics if no address is set on the server", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("server did not panic when address not set")
			}
		}()

		srv := &Server{}
		err := srv.ListenAndServe()
		if err != nil {
			t.Fatalf("error returned from ListenAndServe: %v", err)
		}
	})

	t.Run("starts listening on the specified address", func(t *testing.T) {
		addr := ":7777"
		srv := &Server{Addr: addr}

		go func() {
			if err := srv.ListenAndServe(); err != nil {
				t.Fatalf("unable to ListenAndServe: %v", err)
			}
		}()

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatalf("unable to connect to server: %v", err)
		}

		err = conn.Close()
		if err != nil {
			t.Fatalf("unable to close connection: %v", err)
		}
	})
}
