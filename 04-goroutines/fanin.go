package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Address  Address `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Company  Company `json:"company"`
}
type Geo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     Geo    `json:"geo"`
}
type Company struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
	Bs          string `json:"bs"`
}

// start_fetch OMIT

func fetch(userID string, ch chan<- User) { // HL_ret
	start := time.Now()

	resp, err := http.Get("https://jsonplaceholder.typicode.com/users/" + userID) // HL_core
	if err != nil {
		log.Printf("[ERROR] could not fetch userID %s: %v", userID, err)
	}

	defer resp.Body.Close()              // HL_core
	user := User{}                       // HL_core
	body, _ := ioutil.ReadAll(resp.Body) // HL_core

	_ = json.Unmarshal(body, &user) // HL_core

	elapsed := time.Now().Sub(start)
	fmt.Printf("elapsed %dms to fetch user %s\n", elapsed/time.Millisecond, userID)

	ch <- user // HL_ret
}

// end_fetch OMIT

// start_main OMIT
func main() {
	// start_main_fetch OMIT
	userIDs := []string{"1", "2", "3", "4", "5", "42"}

	// No goroutine will block when sending data.
	// If it's necessary or not depends on the specific usecase.
	ch := make(chan User, len(userIDs))

	start := time.Now()
	for _, id := range userIDs {
		go fetch(id, ch) // HL
	}
	// end_main_fetch OMIT

	// start_main_collect OMIT
	var users []User

	for i := 0; i < len(userIDs); i++ { // we know how manny messages we need to collect
		users = append(users, <-ch) // HL
	}
	// end_main_collect OMIT

	elapsed := time.Now().Sub(start)

	fmt.Printf("total elapsed time: %d to fetch %d users\n",
		elapsed/time.Millisecond, len(users))
}

// end_main OMIT
