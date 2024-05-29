package database

import (
	"fmt"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	ids := make([]int, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		ids = append(ids, user.Id)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	latestID := 1
	if len(ids) > 0 {
		latestID = ids[len(ids)-1] + 1
	}

	fmt.Println(password)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		Id:       latestID,
		Email:    email,
		Password: string(hashBytes),
	}

	dbStructure.Users[newUser.Id] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return User{
		Id:    newUser.Id,
		Email: newUser.Email,
	}, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
