package main

import (
	_ "bytes"
	_ "context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/assign-scratch/arangodb"
	"github.com/assign-scratch/webServer/verification"
	"golang.org/x/crypto/bcrypt"
	_ "io"
	"io/ioutil"
	"net/http"
	_ "reflect"
	_ "strconv"
	_ "strings"
)

func init() {
	flag.Parse()
	arangodb.Connect()
}

type input struct {
	Username string
	Password string
}

type response struct {
	Message    string `json:"message"`
	Authorized bool   `json:"authorized"`
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

	var resp response
	input := input{}
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &input)
		if err != nil {
			panic(err)
		}

		userExists := arangodb.CheckIfUserExists("users", input.Username)
		if !userExists {
			resp.Authorized = false
			resp.Message = "incorrect username or password"
			js, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}

			_, err = w.Write(js)
			if err != nil {
				panic(err)
			}
			return
		}

		fmt.Println("user =>", input.Username)
		fmt.Println("pass =>", input.Password)

		hashPass, _ := verification.HashPassword(input.Password)
		fmt.Println("hash =>", hashPass)

		authUser := arangodb.VerifyUserPassword("users", input.Username, input.Password)
		if authUser {
			resp.Authorized = true
			resp.Message = "user authorized"
			js, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}

			_, err = w.Write(js)
			if err != nil {
				panic(err)
			}
		} else if !authUser {
			resp.Authorized = false
			resp.Message = "user not authorized"
			js, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}

			_, err = w.Write(js)
			if err != nil {
				panic(err)
			}
		}

	}
}

func sendData() {

}

func test() {
	// t := arangodb.CheckIfUserExists("users", "joshua")
	// z := arangodb.CheckIfUserExists("users", "phillip")
	// x := arangodb.CheckIfUserExists("users", "michael")
	// r := arangodb.CheckIfUserExists("users", "cracken rackin")
	// y := arangodb.CheckIfUserExists("users", " ")
	// fmt.Println("t =>", t)
	// fmt.Println("z =>", z)
	// fmt.Println("x =>", x)
	// fmt.Println("r =>", r)
	// fmt.Println("y =>", y)

	// arangodb.AddUserDoc("users", arangodb.User{Username: "William", Age: 29, Email: "williamg@hotmail.com", DOB: "12/25/1989", Town: "Tenafly", Admin: false})

	// arangodb.RemoveUserDoc("users", "William")

	// arangodb.UpdateUserDoc("users", "rodsan", "username", "joshua")

}

func main() {

	http.HandleFunc("/gate/", loginService)
	http.Handle("/assign-scratch/", http.StripPrefix("/assign-scratch/", http.FileServer(http.Dir("/Users/mektievp/assign-scratch/"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
