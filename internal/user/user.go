package user

import "time"

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Create(username, password string) (User, error) {
	newUser := User{
		Username: username,
	}

	return newUser, nil
}
