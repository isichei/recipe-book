package recipes

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseMarkdownFile(t *testing.T) {
	testFilePath := "testdata/example_recipe.md"
	dat, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Errorf("Failed to read test file: %s", testFilePath)
	}
	r := ParseMarkdownFile(string(dat))
	er := Recipe{
		Title:       "Cheese Sandwich",
		Description: "",
		PrepTime:    "2 mins",
		CookingTime: "blazingly fast",
		Serves:      "1",
		Source:      "https://github.com/isichei",
		Ingredients: []Ingredient{{"cheese", ""}, {"Bread", "2 slices"}, {"Some butter", ""}},
		Method:      []string{"Butter the bread", "Slice the cheese", "Put the cheese inbetween the bread"},
	}
	if r.Title != er.Title {
		t.Errorf("Got wrong title: `%s` exepected `%s`", r.Title, er.Title)
	}
	if r.PrepTime != er.PrepTime {
		t.Errorf("Got wrong prep time: `%s` expected: `%s`", r.PrepTime, er.PrepTime)
	}
	if r.CookingTime != er.CookingTime {
		t.Errorf("Got wrong cooking time: `%s` expected: `%s`", r.CookingTime, er.CookingTime)
	}
	if r.Serves != er.Serves {
		t.Errorf("Got wrong serves: `%s` expected: `%s`", r.Serves, er.Source)
	}
	if r.Source != er.Source {
		t.Errorf("Got wrong source: `%s` expected: `%s`", r.Source, er.Source)
	}
	if !(strings.HasPrefix(r.Description, "This is a ") && strings.HasSuffix(r.Description, "some more lines.")) {
		t.Errorf("Got wrong desc: `%s`", r.Description)
	}

	if len(r.Ingredients) != len(er.Ingredients) {
		t.Errorf("Size of ingredients wrong: got %d expected %d", len(r.Ingredients), len(er.Ingredients))
	} else {
		for i, expected := range er.Ingredients {
			testName := fmt.Sprintf("Check ingredient at %d", i)
			t.Run(testName, func(t *testing.T) {
				if r.Ingredients[i].Name != expected.Name || r.Ingredients[i].Amount != expected.Amount {
					t.Errorf("Got wrong ingedient: `%+v` expected `%+v`", r.Ingredients[i], expected)
				}
			})
		}
	}
	if len(r.Method) != len(er.Method) {
		t.Errorf("Size of Method wrong: got %d expected %d", len(r.Method), len(er.Method))
	} else {
		for i, expected := range er.Method {
			testName := fmt.Sprintf("Check method at %d", i)
			t.Run(testName, func(t *testing.T) {
				if r.Method[i] != expected {
					t.Errorf("Got wrong method: `%s` but expected `%s`", r.Method[i], expected)
				}
			})
		}
	}
}
