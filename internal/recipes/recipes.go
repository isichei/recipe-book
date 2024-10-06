package recipes

import (
	"errors"
	"fmt"
)

type RecipeMetadata struct {
	Uid         string
	Title       string
	Description string
}

type Ingredient struct {
	Name   string
	Amount string
}
type Recipe struct {
	Title       string
	Description string
	PrepTime    string
	CookingTime string
	Serves      string
	Ingredients []Ingredient
	Method      []string
	OtherNotes  string
	Source      string
}

func (r Recipe) Validate() error {
	if r.Title == "" {
		return errors.New("Title is required")
	}
	if r.Description == "" {
		return errors.New("Description is required")
	}
	if r.PrepTime == "" {
		return errors.New("PrepTime is required")
	}
	if r.CookingTime == "" {
		return errors.New("CookingTime is required")
	}
	if r.Serves == "" {
		return errors.New("Serves is required")
	}
	if len(r.Ingredients) == 0 {
		return errors.New("At least one ingredient is required")
	}
	for i, ingredient := range r.Ingredients {
		if ingredient.Name == "" || ingredient.Amount == "" {
			return fmt.Errorf("Ingredient at index %d is incomplete", i)
		}
	}
	if len(r.Method) == 0 {
		return errors.New("At least one step in the method is required")
	}
	if r.OtherNotes == "" {
		return errors.New("Other notes is required")
	}
	if r.Source == "" {
		return errors.New("Source is required")
	}

	return nil
}
