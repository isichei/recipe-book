package types

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
