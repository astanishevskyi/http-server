package graphql

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.ID,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Email": &graphql.Field{
				Type: graphql.String,
			},
			"Age": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
