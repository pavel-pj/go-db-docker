package main

import (
	"strings"
	"sync"
)

// User represents a domain model.
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Age       int
}

// UserDTO represents a public view model.
type UserDTO struct {
	ID       int
	FullName string
	Email    string
	AgeGroup string
}

const (
	AgeGroupMinor  = "minor"
	AgeGroupAdult  = "adult"
	AgeGroupSenior = "senior"
)

func MapUsers(users []User) []UserDTO {
	var wg sync.WaitGroup

	var userDTO []UserDTO

	for _, user := range users {

		wg.Add(1)

		go func(u User, wg *sync.WaitGroup) {
			defer wg.Done()
			firstName := strings.TrimSpace(u.FirstName)
			lastName := strings.TrimSpace(u.LastName)
			FullName := firstName + " " + lastName
			Email := strings.ToLower(u.Email)
			var AgeGroup string
			switch {
			case u.Age < 18:
				AgeGroup = "minor"
			case u.Age >= 18 && u.Age < 65:
				AgeGroup = "adult"
			default:
				AgeGroup = "senior"
			}

			userDto = append(userDto, UserDTO{
				ID:       u.ID,
				FullName: FullName,
				Email:    Email,
				AgeGroup: AgeGroup,
			})

		}(user, &wg)

	}

	wg.Wait()
	return userDTO

}

func main() {

}
