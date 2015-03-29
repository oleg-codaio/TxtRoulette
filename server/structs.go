package server

const userRecentSize = 5

type User struct {
	phoneNumber string
	recent      []*User
	blocked     []*User
}

func New(phoneNumber string) *User {
	return &User{
		phoneNumber: phoneNumber,
		recent:      make([]*User, userRecentSize),
		blocked:     make([]*User, 0),
	}
}
