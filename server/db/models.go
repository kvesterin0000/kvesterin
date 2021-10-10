package db

type User struct {
	Id       int
	Email    string
	Username string
	Password string
}

type Release struct {
	Id          int
	Cover       string
	Name        string
	Authors     string
	Status      string
	ReleaseDate string
}

type Track struct {
	Id      int
	Name    string
	Authors string
}
