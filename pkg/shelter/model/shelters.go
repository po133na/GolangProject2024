package model

import (
	"database/sql"
	"errors"
	"log"
)

type Shelter struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Coordinates string `json:"coordinates"`
}

type ShelterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var shelters = []Shelter{
	{
		Id:      "1",
		Title:   "Paws Haven",
		Address: "42 Bark Lane, Cityville, CA 12345",
	},
	{
		Id:      "2",
		Title:   "Whisker Sanctuary",
		Address: "15 Meow Street, Furrytown, TX 67890",
	},
	{
		Id:      "3",
		Title:   "Feathered Friends Shelter",
		Address: "78 Tweet Avenue, Birdsville, NY 54321",
	},
	{
		Id:      "4",
		Title:   "Scales and Tails Rescue",
		Address: "31 Slither Road, Reptileville, FL 98765",
	},
	{
		Id:      "5",
		Title:   "Horse Haven",
		Address: "55 Gallop Lane, Equinetown, CA 24680",
	},
}

func GetShelters() []Shelter {
	return shelters
}

// changed id to string check if goes wrong
func GetShelter(id string) (*Shelter, error) {
	for _, s := range shelters {
		if s.Id == id {
			return &s, nil
		}
	}
	return nil, errors.New("Shelter not found")
}
