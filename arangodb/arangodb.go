package arangodb

import (
	"context"
	"flag"
	"fmt"
	driver "github.com/arangodb/go-driver"
	arangoHttp "github.com/arangodb/go-driver/http"
	"github.com/assign-scratch/webServer/verification"
	_ "log"
	_ "reflect"
	_ "time"
)

var (
	client driver.Client
	CTX    context.Context
	Db     driver.Database

	dbString = flag.String("dbhost", "http://localhost:8529", "arangodb host")
	dbUser   = flag.String("user", "root", "db username")
	dbPass   = flag.String("pass", "tacotime", "db password")
)

type User struct {
	Username string `json:"username"`
	Town     string `json:"town"`
	Age      int64  `json:"age"`
	DOB      string `json:"dob"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

// This function is used to connect to arangodb. Typically goes into the init func
// of the file using it. Also call flag.Parse() in the init before Connect().
func Connect() {

	// Connect to endpoint: localhost:8529
	conn, err := arangoHttp.NewConnection(arangoHttp.ConnectionConfig{
		Endpoints: []string{*dbString},
	})
	if err != nil {
		panic(err)
	}

	// Create a client and supply db user and db password.
	client, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(*dbUser, *dbPass),
	})
	if err != nil {
		panic(err)
	}

	// Create a context for server requests and connect to the "users" db.
	// If 404 is returned, we create "users" db. If another serious error
	// occurs, we return.
	CTX = context.Background()
	Db, err = client.Database(CTX, "users")
	if driver.IsNotFound(err) {
		Db, err = client.CreateDatabase(CTX, "users", &driver.CreateDatabaseOptions{})
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
		return
	}

	// Select "users" collection. If 404 is returned, we create a "universal" collection.
	_, err = Db.Collection(CTX, "users")
	if driver.IsNotFound(err) {
		options := &driver.CreateCollectionOptions{}
		_, err = Db.CreateCollection(CTX, "users", options)
		if err != nil {
			panic(err)
		}
	}
}

// This gets passed the username from the client-side. If the user exists,
// returns true, otherwise false. Collection is also passed here.
func CheckIfUserExists(collection string, username string) bool {

	var doc User
	query := fmt.Sprintf("FOR doc IN %s FILTER doc.username == '%s' RETURN { username: doc.username, age: doc.age, email: doc.email, town: doc.town, DOB: doc.DOB, admin: doc.admin }", collection, username)
	cursor, err := Db.Query(CTX, query, nil)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for {
		_, err := cursor.ReadDocument(CTX, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
	}

	if doc.Username == "" {
		return false
	}

	return true
}

// This function adds a user to the database. Should be called after checking
// if the user already exists or not.
func AddUserDoc(collection string, user User) {

	col, err := Db.Collection(CTX, collection)
	if err != nil {
		panic(err)
	}
	meta, err := col.CreateDocument(CTX, user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created document w/ key '%s', revision '%s'\n", meta.Key, meta.Rev)

}

// This function removes a user from the database. Should be called after checking
// if the user already exists or not.
func RemoveUserDoc(collection string, username string) {

	var docKey []string
	query := fmt.Sprintf("FOR doc IN %s FILTER doc.username == '%s' RETURN { username: doc.username, age: doc.age, email: doc.email, town: doc.town, DOB: doc.DOB, admin: doc.admin }", collection, username)
	cursor, err := Db.Query(CTX, query, nil)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for {
		var doc User
		meta, err := cursor.ReadDocument(CTX, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}

		docKey = append(docKey, meta.Key)
	}

	col, err := Db.Collection(CTX, collection)
	if err != nil {
		panic(err)
	}
	for _, e := range docKey {
		_, err := col.RemoveDocument(CTX, e)
		if err != nil {
			panic(err)
		}
	}
}

// Updates a user in the database by supplying collection, username,
// the field you wish to update, and the update itself.
func UpdateUserDoc(collection string, username string, field string, update string) {

	var docKey []string
	query := fmt.Sprintf("FOR doc IN %s FILTER doc.username == '%s' RETURN doc", collection, username)
	cursor, err := Db.Query(CTX, query, nil)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for {
		var doc User
		meta, err := cursor.ReadDocument(CTX, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}

		docKey = append(docKey, meta.Key)
	}

	col, err := Db.Collection(CTX, collection)
	if err != nil {
		panic(err)
	}

	patch := map[string]interface{}{
		field: update,
	}

	_, err = col.UpdateDocument(CTX, docKey[0], patch)
	if err != nil {
		panic(err)
	}
}

func VerifyUserPassword(collection string, username string, password string) bool {

	var authorize bool
	query := fmt.Sprintf("FOR doc IN %s FILTER doc.username == '%s' RETURN doc.password", collection, username)
	cursor, err := Db.Query(CTX, query, nil)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for {
		var retrievedHash string
		_, err := cursor.ReadDocument(CTX, &retrievedHash)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}

		fmt.Println("verified =>", retrievedHash)
		authorize = verification.CheckPasswordHash(password, retrievedHash)

	}

	return authorize
}
