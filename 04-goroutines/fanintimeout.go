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

func query(userID string, ch chan<- User) { // use context to check if can send to channel
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users/" + userID)
	if err != nil {
		log.Printf("[ERROR] could not fetch userID %s: %v", userID, err)
	}

	defer resp.Body.Close()
	user := User{}
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &user)

	ch <- user
}

func main() {
	userIDs := []string{"1", "2", "3", "4", "5", "42"}
	// start_main OMIT

	ch := make(chan User, len(userIDs))

	start := time.Now()
	for _, id := range userIDs {
		go query(id, ch) // HL
	}

	timeout := time.After(190 * time.Millisecond) // HL
	users := collect(len(userIDs), ch, timeout)   // HL

	elapsed := time.Now().Sub(start)
	fmt.Printf("total elapsed time: %dms to fetch %d users\n",
		elapsed/time.Millisecond, len(users))

	fmt.Printf("%f%% success rate\n", percent(users, userIDs))
	// end_main OMIT
}

// start_collect OMIT
func collect(size int, ch <-chan User, stopCh <-chan time.Time) []User {
	var users []User

	for i := 0; i < size; i++ {
		select {
		case <-stopCh: // if we get a signal here, we must stop // HL
			fmt.Printf("too long, timing out! fetched %d out of %d users\n", // HL
				len(users), size) // HL
			return users // HL
		case u := <-ch:
			users = append(users, u)
		}
	}

	return users
}

// end_collect OMIT

func percent(users []User, userIDs []string) float64 {
	return (float64(len(users)) / float64(len(userIDs))) * 100
}
