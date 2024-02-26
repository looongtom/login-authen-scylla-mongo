package database

import (
	"github.com/gocql/gocql"
	"os"
)

func ConnectScylla() *gocql.ClusterConfig {
	auth := gocql.PasswordAuthenticator{
		Username: os.Getenv("UsernameScylla"),
		Password: os.Getenv("PasswordScylla"),
	}

	// Connect to the ScyllaDB cluster
	cluster := gocql.NewCluster(os.Getenv("UrlScylla"))
	cluster.Keyspace = os.Getenv("KeyspaceScylla")
	cluster.Authenticator = auth

	return cluster
}
