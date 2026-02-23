package service

import "fmt"

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
	return []UserDTO{}

}

func T1() {
	fmt.Println("HELLO")
}
