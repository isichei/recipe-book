package types

type RecipeMetadata struct {
	Uid         string
	Title       string
	Description string
}

type Recipe struct {
	Title       string
	Description string
	PrepTime    string
	CookingTime string
	Serves      string
	Ingredients []string
	Method      []string
	OtherNotes  string
	Source      string
}
