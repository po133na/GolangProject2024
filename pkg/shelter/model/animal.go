package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Animal struct {
	ID             string `json:"id"`
	Kind_Of_Animal string `json:"kind_of_animal"`
	Kind_Of_Breed  string `json:"kind_of_breed"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	Description    string `json:"description"`
}

type AnimalModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (a AnimalModel) Insert(animal *Animal) error {
	// check for ID needed here if error
	query := `
		INSERT INTO Animals (Kind_Of_Animal, Kind_Of_Breed, Name, Age, Description) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, Name
		`
	// check if its animal of Animals in case of error
	args := []interface{}{animal.Kind_Of_Animal, animal.Kind_Of_Breed, animal.Name, animal.Age, animal.Description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&animal.ID, &animal.Name)
}

func (a AnimalModel) Get(id int) (*Animal, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT id, Kind_Of_Animal, Kind_Of_Breed, Name, Age, Description
		FROM Animals
		WHERE ID = $1
		`
	var animal Animal
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// again animal or Animals?
	row := a.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&animal.ID, &animal.Kind_Of_Animal, &animal.Kind_Of_Breed, &animal.Name, &animal.Age, &animal.Description)
	if err != nil {
		return nil, err
	}
	return &animal, nil
}

func (a AnimalModel) GetSort(Kind_Of_Breed, Kind_of_Animal string, filters Filters) ([]*Animal, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, Kind_Of_Animal, Kind_Of_Breed, Name, Age, Description
		FROM Animals
		WHERE (LOWER(Kind_Of_Breed) = LOWER($1) OR $1 = '')
		AND (LOWER(Kind_of_Animal) = LOWER($2) OR $2 = '')
		--AND (Age >= $2 OR $2 = 0)
		--AND (Age <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{Kind_Of_Breed, Kind_of_Animal, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			a.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var animals []*Animal
	for rows.Next() {
		var animal Animal
		err := rows.Scan(&totalRecords, &animal.ID, &animal.Kind_Of_Animal, &animal.Kind_Of_Breed, &animal.Name, &animal.Age, &animal.Description)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		animals = append(animals, &animal)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return animals, metadata, nil
}

func (a AnimalModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM Animals
		WHERE ID = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.DB.ExecContext(ctx, query, id)
	return err
}

func (a AnimalModel) Update(animal *Animal) error {
	// Update a specific animal in the database.
	query := `
        UPDATE Animals
        SET Kind_Of_Animal = $2, Kind_Of_Breed = $3, Name = $4, Age = $5, Description = $6
        WHERE ID = $1
        RETURNING ID
        `
	args := []interface{}{animal.ID, animal.Kind_Of_Animal, animal.Kind_Of_Breed, animal.Name, animal.Age, animal.Description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&animal.ID)
}
