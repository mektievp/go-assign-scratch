package arangodb

import (
	"context"
	"flag"
	driver "github.com/arangodb/go-driver"
	arangoHttp "github.com/arangodb/go-driver/http"
	"log"
	"time"
)

var (
	dbString = flag.String("dbhost", "http://localhost:8529", "arangodb host")
	dbUser   = flag.String("user", "root", "db username")
	dbPass   = flag.String("pass", "tacotime", "db password")
)

// This function is used to connect to arangodb. Typically goes into the init func
// of the file using it. Also call flag.Parse() in the init before Connect().
func Connect() {

	// Connect to endpoint: localhost:8529
	conn, err := arangoHttp.NewConnection(arangoHttp.ConnectionConfig{
		Endpoints: []string{*dbString},
	})
	if err != nil {
		log.Printf("%v: %v\n", time.Now(), err)
	}

	// Create a client and supply db user and db password.
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(*dbUser, *dbPass),
	})
	if err != nil {
		log.Printf("%v: %v\n", time.Now(), err)
	}

	// Create a context for server requests and connect to the "users" db.
	// If 404 is returned, we create "users" db. If another serious error
	// occurs, we return.
	CTX := context.Background()
	Db, err := client.Database(CTX, "users")
	if driver.IsNotFound(err) {
		Db, err = client.CreateDatabase(CTX, "users", &driver.CreateDatabaseOptions{})
		if err != nil {
			log.Printf("%v: %v\n", time.Now(), err)
		}
	}
	if err != nil {
		log.Printf("%v: %v\n", time.Now(), err)
		return
	}

	// Select "universal" collection. If 404 is returned, we create a "universal" collection.
	_, err = Db.Collection(CTX, "universal")
	if driver.IsNotFound(err) {
		options := &driver.CreateCollectionOptions{}
		_, err = Db.CreateCollection(CTX, "universal", options)
		if err != nil {
			log.Printf("%v: %v\n", time.Now(), err)
		}
	}
}
