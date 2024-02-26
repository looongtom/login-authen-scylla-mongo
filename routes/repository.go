package routes

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"go.mongodb.org/mongo-driver/bson"
	"login-user/auth"
	db "login-user/database"
	"login-user/models"
	"login-user/utils"
	"net/http"
	"os"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if govalidator.IsNull(username) || govalidator.IsNull(email) || govalidator.IsNull(password) {
		utils.JSON(w, 400, "Data can not empty")
		return
	}

	if !govalidator.IsEmail(email) {
		utils.JSON(w, 400, "Email is invalid")
		return
	}

	collection := db.ConnectUser(os.Getenv("UserCollectionMongo"))
	var result bson.M
	errFindUsername := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	errFindEmail := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)

	if errFindUsername == nil || errFindEmail == nil {
		utils.JSON(w, 409, "User does exists")
		return
	}

	password, err := models.Hash(password)

	if err != nil {
		utils.JSON(w, 500, "Register has failed")
		return
	}

	newUser := bson.M{"username": username, "email": email, "password": password}

	_, errs := collection.InsertOne(context.TODO(), newUser)

	if errs != nil {
		utils.JSON(w, 500, "Register has failed")
		return
	}

	utils.JSON(w, 201, "Register Succesfully")
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if govalidator.IsNull(username) || govalidator.IsNull(password) {
		utils.JSON(w, 400, "Data can not empty")
		return
	}
	username = models.Santize(username)
	password = models.Santize(password)

	collection := db.ConnectUser(os.Getenv("UserCollectionMongo"))
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)

	fmt.Println("user login:")
	fmt.Println(result)

	if err != nil {
		utils.JSON(w, 400, "Username or Password incorrect")
		return
	}

	hashedPassword := fmt.Sprintf("%v", result["password"])
	err = models.CheckPasswordHash(hashedPassword, password)

	if err != nil {
		utils.JSON(w, 401, "Username or Password incorrect")
		return
	}

	token, errCreate := auth.Create(username)
	if errCreate != nil {
		fmt.Println("errCreate")
		fmt.Println(errCreate)

		fmt.Println("token")
		fmt.Println(token)
		utils.JSON(w, 500, "Internal Server Error")
		return
	}

	collection = db.ConnectUser(os.Getenv("TokenCollectionMongo"))
	newToken := bson.M{"token": token, "user": username, "created_at": time.Now()}
	_, errs := collection.InsertOne(context.TODO(), newToken)

	if errs != nil {
		utils.JSON(w, 500, "Error create token")
		return
	}

	//collection = db.ConnectUser(history_collection)
	//
	//filter := bson.M{"username": username}
	//update := bson.M{"$set": bson.M{"logout_at": "", "login_at": time.Now().Format("2006-01-02 15:04:05")}}
	//updateOptions := options.Update().SetUpsert(true)
	//updateResult, err := collection.UpdateOne(context.TODO(), filter, update, updateOptions)
	//if err != nil {
	//	fmt.Println("Error updating or inserting document:", err)
	//	return
	//}
	//if updateResult.MatchedCount == 0 {
	//	fmt.Println("No matching document found. A new document has been inserted.")
	//} else {
	//	fmt.Println("Matched and updated document.")
	//}

	// Connect to the ScyllaDB cluster
	cluster := db.ConnectScylla()

	session, err := cluster.CreateSession()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	defer session.Close()

	var history models.HistoryConnection
	selectQuery, getNames := qb.Select("history_connections").Where(qb.Eq("username")).AllowFiltering().ToCql()
	if err := gocqlx.Query(session.Query(selectQuery, username), getNames).GetRelease(&history); err != nil {
		fmt.Println("Error getting history:", err)
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	switch history {
	case models.HistoryConnection{}:
		fmt.Println("inserting new history")
		newHistory := models.HistoryConnection{
			Id:       gocql.TimeUUID(),
			Username: username,
			LoginAt:  time.Now().Local(),
			LogoutAt: time.Time{},
		}
		insertStmt, names := qb.Insert("history_connections").Columns("id", "login_at", "logout_at", "username").ToCql()
		if err := gocqlx.Query(session.Query(insertStmt), names).BindStruct(newHistory).ExecRelease(); err != nil {
			utils.ERROR(w, http.StatusInternalServerError, err)
		}

	default:
		fmt.Println("updating history:", history)
		history.LogoutAt = time.Time{}
		history.LoginAt = time.Now().Local()

		updateQuery, names := qb.Update("history_connections").
			Set("login_at", "logout_at").Where(qb.Eq("id")).ToCql()

		// Execute the update query
		if err := gocqlx.Query(session.Query(updateQuery, history.LoginAt, history.LogoutAt, history.Id), names).ExecRelease(); err != nil {
			utils.ERROR(w, http.StatusInternalServerError, err)
		}

	}

	utils.JSON(w, 200, token)
}

func GetProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	userName, err := auth.GetSubjectFromToken(tokenString)
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err = auth.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}
	fmt.Fprint(w, "Welcome "+userName)
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	userName, err := auth.GetSubjectFromToken(tokenString)
	tokenString = tokenString[len("Bearer "):]

	// Connect to the token collection
	collection := db.ConnectUser(os.Getenv("TokenBlackListCollectionMongo"))

	_, err = collection.InsertOne(context.TODO(), bson.M{
		"token":          tokenString,
		"blacklisted_at": time.Now(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error during logout")
		return
	}

	//collection = db.ConnectUser(history_collection)
	//filter := bson.M{"username": userName}
	//update := bson.M{"$set": bson.M{"logout_at": time.Now().Format("2006-01-02 15:04:05")}}
	//updateResult := collection.FindOneAndUpdate(context.TODO(), filter, update)
	//
	//if updateResult.Err() != nil {
	//	fmt.Println("Error updating logout_at field:", updateResult.Err())
	//	return
	//}

	// Connect to the ScyllaDB cluster
	cluster := db.ConnectScylla()

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var history models.HistoryConnection
	selectQuery, getNames := qb.Select("history_connections").Where(qb.Eq("username")).AllowFiltering().ToCql()
	if err := gocqlx.Query(session.Query(selectQuery, userName), getNames).GetRelease(&history); err != nil {
		fmt.Println("Error getting history:", err)
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println("updating history:", history)
	history.LogoutAt = time.Now().Local()

	updateQuery, names := qb.Update("history_connections").
		Set("logout_at").Where(qb.Eq("id")).ToCql()

	// Execute the update query
	if err := gocqlx.Query(session.Query(updateQuery, history.LogoutAt, history.Id), names).ExecRelease(); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}

	// Send a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logout successful")
}
