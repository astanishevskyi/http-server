package connectors

import (
	"github.com/astanishevskyi/http-server/pkg/api"
	"google.golang.org/grpc"
	"log"
)

func NewGRPC(grpcAddr string) api.UserClient {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := api.NewUserClient(conn)
	return c
}
