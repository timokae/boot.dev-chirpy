package database

import (
	"errors"
	"fmt"
	"sort"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

var ErrNotExist = errors.New("resource does not exist")

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// // CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	ids := make([]int, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		ids = append(ids, chirp.Id)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	latestID := 1
	if len(ids) > 0 {
		latestID = ids[len(ids)-1] + 1
	}

	newChirp := Chirp{
		Id:       latestID,
		Body:     body,
		AuthorId: authorId,
	}
	fmt.Println(newChirp)

	dbStructure.Chirps[newChirp.Id] = newChirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}
