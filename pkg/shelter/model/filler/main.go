package filler

import (
	model "github.com/po133na/go-mid/pkg/shelter/model"
)

func PopulateDatabase(models model.Models) error {
	for _, animal := range animals {
		models.Animals.Insert(&animal)
	}
	// TODO: Implement restaurants population
	// TODO: Implement the relationship between shelters and animals
	return nil
}

var animals = []model.Animal{
	{ID: "1", Name: "Dumbo", Age: "15", Description: "A large mammal with a trunk and tusks"},
	{ID: "2", Name: "Simba", Age: "8", Description: "A large carnivorous feline, king of the jungle"},
	{ID: "3", Name: "Flipper", Age: "10", Description: "An intelligent aquatic mammal"},
	{ID: "4", Name: "Longneck", Age: "12", Description: "A tall mammal with a long neck"},
	{ID: "5", Name: "Pingu", Age: "6", Description: "A flightless bird that lives in cold climates"},
	{ID: "6", Name: "Joey", Age: "7", Description: "A marsupial with powerful hind legs"},
	{ID: "7", Name: "Whitey", Age: "5", Description: "A large bear adapted to living in cold climates"},
	{ID: "8", Name: "Swift", Age: "4", Description: "The fastest land animal, capable of high speeds"},
	{ID: "9", Name: "Polly", Age: "20", Description: "A colorful bird known for its ability to mimic speech"},
	{ID: "10", Name: "Moby", Age: "25", Description: "A large marine mammal, among the largest animals"},
}
