package main

import (
	_ "context"
	"encoding/json"
	"fmt"
	"github.com/assign-scratch/arangodb"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	_ "reflect"
	_ "strconv"
	_ "strings"
)

type user struct {
	Username string
	Password string
}

func init() {
	arangodb.Connect()
	// conn, err := arangoHttp.NewConnection(arangoHttp.ConnectionConfig{
	// 	Endpoints: []string{"http://localhost:8529"},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// client, err := driver.NewClient(driver.ClientConfig{
	// 	Connection:     conn,
	// 	Authentication: driver.BasicAuthentication("root", ""),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("client =>", client)

	// ctx := context.Background()
	// db, err := client.Database(ctx, "users")
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// col, err := db.Collection(ctx, "users")
	// if err != nil {
	// 	panic(err)
	// }

	// joshPassword, _ := hashPassword("dragi")
	// doc := user{
	// 	Username: "joshua",
	// 	Password: joshPassword,
	// }

	// ctx = context.Background()
	// meta, err := col.CreateDocument(ctx, doc)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Created document with key '%s', revision '%s'\n", meta.Key, meta.Rev)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func loginService(w http.ResponseWriter, r *http.Request) {
	user := user{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &user)
		fmt.Println("err =>", err)
		if err != nil {
			panic(err)
		}

		fmt.Println("user =>", user.Username)
		fmt.Println("pass =>", user.Password)

		hashPass, _ := hashPassword(user.Password)
		fmt.Println("hash =>", hashPass)
	}
}

func main() {
	http.HandleFunc("/gate/", loginService)
	http.Handle("/assign-scratch/", http.StripPrefix("/assign-scratch/", http.FileServer(http.Dir("/Users/mektievp/assign-scratch"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
