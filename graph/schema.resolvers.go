package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"gq_enrich/graph/model"
)

// DeletePerson is the resolver for the deletePerson field.
func (r *mutationResolver) DeletePerson(ctx context.Context, id int) (*model.Person, error) {
	var person = &model.Person{}

	queryDeletePersonByID := `DELETE FROM people
		WHERE id = $1`

	_, err := r.DB.ExecContext(
		context.Background(),
		queryDeletePersonByID,
		id,
	)

	if err != nil {
		return person, err
	}

	return person, nil
}

// ChangeSurname is the resolver for the changeSurname field.
func (r *mutationResolver) ChangeSurname(ctx context.Context, id int, surname string) (*model.Person, error) {
	var person = &model.Person{}

	queryChangeSurname := `UPDATE people
		SET surname = $2 WHERE id = $1`

	_, err := r.DB.ExecContext(
		context.Background(),
		queryChangeSurname,
		id,
		surname,
	)

	if err != nil {
		return person, err
	}

	return person, nil
}

// ChangeAge is the resolver for the changeAge field.
func (r *mutationResolver) ChangeAge(ctx context.Context, id int, age int) (*model.Person, error) {
	var person = &model.Person{}

	queryChangeAge := `UPDATE people
		SET age = $2 WHERE id = $1`

	_, err := r.DB.ExecContext(
		context.Background(),
		queryChangeAge,
		id,
		age,
	)

	if err != nil {
		return person, err
	}

	return person, nil
}

// Men is the resolver for the men field.
func (r *queryResolver) Men(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person

	querySelectMen := `SELECT *
		FROM people WHERE gender LIKE 'male'`

	err := r.DB.SelectContext(
		context.Background(),
		&persons,
		querySelectMen,
	)

	if err != nil {
		return persons, err
	}

	return persons, nil
}

// Women is the resolver for the women field.
func (r *queryResolver) Women(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person

	querySelectWomen := `SELECT *
		FROM people WHERE gender LIKE 'female'`

	err := r.DB.SelectContext(
		context.Background(),
		&persons,
		querySelectWomen,
	)

	if err != nil {
		return persons, err
	}

	return persons, nil
}

// People is the resolver for the people field.
func (r *queryResolver) People(ctx context.Context, name string) ([]*model.Person, error) {
	var persons []*model.Person

	queryGetPeopleByName := `SELECT *
        FROM people WHERE p_name ILIKE $1::text`

	err := r.DB.SelectContext(
		context.Background(),
		&persons,
		queryGetPeopleByName,
		name,
	)

	if err != nil {
		return persons, err
	}

	return persons, nil
}

// Pplage is the resolver for the pplage field.
func (r *queryResolver) Pplage(ctx context.Context, age int, less bool, desc bool) ([]*model.Person, error) {
	var persons []*model.Person

	queryGetPeopleByAgeLess := `SELECT * 
		FROM people WHERE age < $1
		ORDER BY (CASE WHEN $2 = true THEN age END) DESC,
				 (CASE WHEN $2 = false THEN age END) ASC`

	queryGetPeopleByAgeMore := `SELECT * 
	FROM people WHERE age > $1
		ORDER BY (CASE WHEN $2 = true THEN age END) DESC,
			 	 (CASE WHEN $2 = false THEN age END) ASC`

	var err error

	if less {
		err = r.DB.SelectContext(
			context.Background(),
			&persons,
			queryGetPeopleByAgeLess,
			age,
			desc,
		)
	} else {
		err = r.DB.SelectContext(
			context.Background(),
			&persons,
			queryGetPeopleByAgeMore,
			age,
			desc,
		)
	}

	if err != nil {
		return persons, err
	}

	return persons, nil
}

// Country is the resolver for the country field.
func (r *queryResolver) Country(ctx context.Context, name string) ([]*model.Country, error) {
	var countries []*model.Country

	queryGetCountriesByName := `SELECT rec_id, user_name, country_id, probability
		FROM nation WHERE user_name ILIKE $1::text`

	err := r.DB.SelectContext(
		context.Background(),
		&countries,
		queryGetCountriesByName,
		name,
	)

	if err != nil {
		return countries, err
	}

	return countries, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }