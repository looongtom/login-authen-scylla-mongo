package main

import (
	"fmt"
	"github.com/gocql/gocql"
)

func main() {
	// Set up the authentication credentials
	//auth := gocql.PasswordAuthenticator{
	//	Username: "scylla",
	//	Password: "YStO20J5bVCpXkQ",
	//}

	// Connect to the ScyllaDB cluster
	cluster := gocql.NewCluster("172.17.0.2") // Provide the IP addresses of your ScyllaDB nodes
	cluster.Port = 9042
	cluster.Keyspace = "mykeyspace" // Specify your keyspace
	//cluster.Authenticator = auth

	// Create a session to interact with the database
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Example query
	var result string
	if err := session.Query("SELECT * FROM mykeyspace.users").Scan(&result); err != nil {
		panic(err)
	}
	fmt.Println("ScyllaDB Version:", result)

	// Define your data model
	//history := models.HistoryConnection{
	//	Id:       gocql.TimeUUID(),
	//	Username: "test1",
	//	LoginAt:  time.Now(),
	//	LogoutAt: time.Now(),
	//}
	//
	//insertStmt, names := qb.Insert("history_connections").Columns("id", "login_at", "logout_at", "username").ToCql()
	//if err := gocqlx.Query(session.Query(insertStmt), names).BindStruct(history).ExecRelease(); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Data inserted successfully!")

	// Define the username to search for
	//username := "nam_em"

	// Define the query to retrieve history by username
	//query := "SELECT username, login_at, logout_at FROM history.history_connections WHERE username = ? ALLOW FILTERING"
	//
	//// Execute the query
	//iter := session.Query(query, username).Iter()
	//
	//// Iterate over the results
	//var history models.HistoryConnection
	//for iter.Scan(&history.Username, &history.LoginAt, &history.LogoutAt) {
	//	fmt.Printf("Username: %s, LoginAt: %v, LogoutAt: %v\n", history.Username, history.LoginAt, history.LogoutAt)
	//}
	//
	//if err := iter.Close(); err != nil {
	//	panic(err)
	//}

	//GET DATA BY USERNAME
	//query, names := qb.Select("history_connections").Where(qb.Eq("username")).AllowFiltering().ToCql()
	//// Execute the query
	//var history []models.HistoryConnection
	//if err := gocqlx.Query(session.Query(query, username), names).SelectRelease(&history); err != nil {
	//	panic(err)
	//}
	//
	//// Print the results
	//for _, h := range history {
	//	fmt.Printf("Username: %s, LoginAt: %v, LogoutAt: %v\n", h.Username, h.LoginAt, h.LogoutAt)
	//}

	//username := "nam_em"
	//
	//// Query the record by username
	////var history models.HistoryConnection
	//if (models.HistoryConnection{}) == history {
	//	fmt.Println("history is empty")
	//}
	//selectQuery, names := qb.Select("history_connections").Where(qb.Eq("username")).AllowFiltering().ToCql()
	////fmt.Println("history:", history)
	//if err := gocqlx.Query(session.Query(selectQuery, username), names).GetRelease(&history); err != nil {
	//	panic(err)
	//}
	//fmt.Println(history)
	//// Update the record
	//history.LogoutAt = time.Time{}
	//
	//// Define the update query
	//updateQuery, names := qb.Update("history_connections").
	//	Set("logout_at").Where(qb.Eq("id")).ToCql()
	//
	//// Execute the update query
	//if err := gocqlx.Query(session.Query(updateQuery, history.LogoutAt, history.Id), names).ExecRelease(); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Record updated successfully!")
}
