package server

const userRecentSize = 1

type User struct {
	phoneNumber string
	recent      []*User
	blocked     []*User
}

func NewUser(phoneNumber string) *User {
	return &User{
		phoneNumber: phoneNumber,
		recent:      make([]*User, 0, userRecentSize),
		blocked:     make([]*User, 0),
	}
}

func (user *User) AddToRecents(other *User) {
	uLen := len(user.recent)
	uCap := cap(user.recent)
	if uLen < uCap {
		user.recent = user.recent[:uLen+1]
		user.recent[uLen] = other
	} else {
		copy(user.recent, user.recent[1:])
		user.recent[uLen-1] = other
	}
}
