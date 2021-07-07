package graphql

import (
	"github.com/astanishevskyi/http-server/internal/apiserver/connectors"
	"github.com/graphql-go/graphql"
	"log"
)

func Schema(grpc connectors.GrpcConnector) (*graphql.Schema, error) {
	var fields = graphql.Fields{
		"user": &graphql.Field{
			Type:        UserType,
			Description: "Get user by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					user, err := grpc.GetUser(uint64(id))
					if err != nil {
						return nil, err
					}
					return user, nil
				}
				return nil, nil
			},
		},
		"users": &graphql.Field{
			Type:        graphql.NewList(UserType),
			Description: "Get user list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				users, err := grpc.GetUsers()
				if err != nil {
					return nil, err
				}
				return users, nil
			},
		},
	}
	var inputs = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addUser": &graphql.Field{
				Name:        "addUser",
				Description: "Create new user",
				Type:        UserType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, ok := p.Args["name"].(string)
					email, ok1 := p.Args["email"].(string)
					age, ok2 := p.Args["age"].(int)
					if ok {
						if ok1 {
							if ok2 {
								newUser, err := grpc.CreateUser(name, email, int32(age))
								if err != nil {
									return nil, err
								}
								return newUser, nil
							}
						}
					}
					return nil, nil
				},
			},
			"updateUser": &graphql.Field{
				Name:        "updateUser",
				Description: "Update user",
				Type:        UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					name, ok1 := p.Args["name"].(string)
					email, ok2 := p.Args["email"].(string)
					age, ok3 := p.Args["age"].(int)
					if ok {
						if ok1 {
							if ok2 {
								if ok3 {
									newUser, err := grpc.UpdateUser(uint32(id), uint32(age), name, email)
									if err != nil {
										return nil, err
									}
									return newUser, nil
								}
							}
						}
					}
					return nil, nil
				},
			},
			"deleteUser": &graphql.Field{
				Name:        "deleteUser",
				Description: "Delete user",
				Type:        UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						newUser, err := grpc.DeleteUser(uint32(id))
						if err != nil {
							return nil, err
						}
						return newUser, nil
					}
					return nil, nil
				},
			},
		},
	})
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: inputs}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return &schema, nil
}
