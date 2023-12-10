package recipes

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/isichei/recipe-book/types"
)

func TestParseMarkdownFile(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	r := ParseMarkdownFile("./example_recipe.md")
	er := types.Recipe{
		Title:       "Cheese Sandwich",
		Description: "",
		PrepTime:    "2 mins",
		CookingTime: "blazingly fast",
		Serves:      "1",
		Source:      "https://github.com/isichei",
		Ingredients: []string{"cheese", "2 slices of bread", "Some butter"},
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
				if r.Ingredients[i] != expected {
					t.Errorf("Got wrong ingedient: `%s` expected `%s`", r.Ingredients[i], expected)
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
