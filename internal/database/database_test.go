package database

import (
	"github.com/isichei/recipe-book/internal/recipes"
	"testing"
)

func TestSet(t *testing.T) {
	s := make(Set)

	s.Add("something")
	if !s["something"] {
		t.Error("Was expecting string 'something' to be in Set")
	}

	if s["something else"] {
		t.Error("Was expecting string 'something else' to not be in Set")
	}
}

func TestSetCache(t *testing.T) {
	sc := make(SearchCache)
	sc.Add("something", "a.file")
	sc.Add("something", "b.file")
	set := sc["something"]
	if !set["b.file"] {
		t.Error("Was expecting string 'a.file' to have been added to 'something'")
	}
	if !set["b.file"] {
		t.Error("Was expecting string 'b.file' to have been added to 'something'")
	}
	if set["c.file"] {
		t.Error("Was expecting string 'c.file' to not have been added to 'something'")
	}
	if len(sc["something else"]) > 0 {
		t.Error("Was expecting 'something else' to return an empty set")
	}
}

func TestSqlDatabase(t *testing.T) {
	db, err := NewSqlDatabase(":memory:")
	if err != nil {
		t.Errorf("No error should have occurred. Error: %s\n", err)
	}

	testRecipe := recipes.Recipe{
		Title:       "Cheese Sandwich",
		Description: "",
		PrepTime:    "2 mins",
		CookingTime: "blazingly fast",
		Serves:      "1",
		Source:      "https://github.com/isichei",
		Ingredients: []recipes.Ingredient{{"cheese", ""}, {"Bread", "2 slices"}, {"Some butter", ""}},
		Method:      []string{"Butter the bread", "Slice the cheese", "Put the cheese inbetween the bread"},
	}
	// Test AddRecipe
	Uid := "cheese-sandwich"
	db.AddRecipe(Uid, testRecipe)

	// Test getSomeRecipes
	if lenGetSomeRecipes := len(db.getSomeRecipes(3)); lenGetSomeRecipes != 1 {
		t.Errorf("Expected 1 recipe to be returned from getSomeRecipes. Got %d", lenGetSomeRecipes)
	}

	// Test GetRecipeMetadata
	recipeMetadata := db.GetRecipeMetadata(Uid)
	if recipeMetadata.Uid != Uid {
		t.Errorf("Expected recipeMetadata to have the Uid: `%s` instead it has `%s`", Uid, recipeMetadata.Uid)
	}
	if recipeMetadata.Title != testRecipe.Title {
		t.Errorf("Expected recipeMetadata to have the Title: `%s` instead it has `%s`", testRecipe.Title, recipeMetadata.Title)
	}

	// Test GetRecipe
	outRecipe := db.GetRecipe(Uid)
	if outRecipe.Title != testRecipe.Title {
		t.Errorf("Expected GetRecipe to return the added recipe but the returned recipe title was %s", outRecipe.Title)
	}

	// Test getIngredients
	outIngredients, err := db.getIngredients(Uid)
	if err != nil {
		t.Errorf("Error on getIngredients %s", err)
	}
	if outIngredients[0] != testRecipe.Ingredients[0] || outIngredients[1] != testRecipe.Ingredients[1] || outIngredients[2] != testRecipe.Ingredients[2] {
		t.Errorf("Got wrong ingredients or order expected %v got %v", testRecipe.Ingredients, outIngredients)
	}

	// Test getMethod
	outMethod, err := db.getMethod(Uid)
	if outMethod[0][0] != 'B' || outMethod[1][0] != 'S' || outMethod[2][0] != 'P' {
		t.Errorf("Got wrong method or order expected %v got %v", testRecipe.Method, outMethod)
	}
}
