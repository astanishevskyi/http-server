package connectors

import (
	"context"
	"fmt"
	"github.com/astanishevskyi/http-server/internal/apiserver/models"
	"github.com/astanishevskyi/http-server/pkg/api"
	"google.golang.org/grpc"
	"io"
	"log"
)

type GrpcConnector struct {
	client api.UserClient
}

func NewGRPC(grpcAddr string) GrpcConnector {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("dgsdgdf")
		log.Fatal(err)
	}
	c := api.NewUserClient(conn)
	client := GrpcConnector{client: c}
	return client
}

func (c *GrpcConnector) GetUser(id uint64) (*api.UserObject, error) {
	return c.client.GetUser(context.Background(), &api.UserId{Id: uint32(id)})
}

func (c *GrpcConnector) GetUsers() ([]models.User, error) {
	grpcResp, err := c.client.GetUsers(context.Background(), &api.NoneObject{})
	if err != nil {
		return nil, err
	}

	userSlice := make([]models.User, 0)

	for {
		res, errRecv := grpcResp.Recv()
		if errRecv == io.EOF {
			break
		}
		if errRecv != nil {
			return nil, err
		}
		user := models.User{ID: res.GetId(), Age: uint8(res.GetAge()), Name: res.GetName(), Email: res.GetEmail()}
		userSlice = append(userSlice, user)
	}
	return userSlice, err
}

func (c *GrpcConnector) CreateUser(name, email string, age int32) (*api.UserObject, error) {
	return c.client.CreateUser(context.Background(), &api.NewUser{Name: name, Email: email, Age: age})
}

func (c *GrpcConnector) UpdateUser(id, age uint32, name, email string) (*api.UserObject, error) {
	return c.client.UpdateUser(context.Background(), &api.UserObject{Id: id, Name: name, Email: email, Age: age})
}

func (c *GrpcConnector) DeleteUser(id uint32) (*api.UserId, error) {
	return c.client.DeleteUser(context.Background(), &api.UserId{Id: id})
}
