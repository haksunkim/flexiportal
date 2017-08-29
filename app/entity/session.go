package entity

type Session struct {
	Id string
	Authenticated bool
	Unauthenticated bool
	User User
}
