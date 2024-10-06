package recipes

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	PREP_TIME_STR    string = "Preparation time:"
	COOKING_TIME_STR string = "Cooking time:"
	SERVES_STR       string = "Serves:"
	SOURCE_STR       string = "Source:"
)

func removeAndStrip(s string, removed string) string {
	return strings.TrimSpace(strings.Replace(s, removed, "", 1))
}

func createIngredient(ingredient_str string) Ingredient {
	f, s, _ := strings.Cut(ingredient_str, ":")
	return Ingredient{Name: strings.TrimSpace(f), Amount: strings.TrimSpace(s)}
}

// Super niave parser - does the job
func ParseMarkdownFile(fileData string) Recipe {

	scanner := bufio.NewScanner(strings.NewReader(fileData))

	// Iterate over each line in the file getting recipe values
	r := Recipe{}
	sub_heading := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "":
		case strings.HasPrefix(line, "# "):
			r.Title = removeAndStrip(line, "#")
		case strings.HasPrefix(line, PREP_TIME_STR):
			r.PrepTime = removeAndStrip(line, PREP_TIME_STR)
		case strings.HasPrefix(line, COOKING_TIME_STR):
			r.CookingTime = removeAndStrip(line, COOKING_TIME_STR)
		case strings.HasPrefix(line, SERVES_STR):
			r.Serves = removeAndStrip(line, SERVES_STR)
		case strings.HasPrefix(line, SOURCE_STR):
			r.Source = removeAndStrip(line, SOURCE_STR)
		case r.Serves != "" && r.Source == "":
			r.Description += line
		case line == "## Ingredients:":
			sub_heading = "ingredients"
		case line == "## Method:":
			sub_heading = "method"
		case line == "## Other notes:":
			sub_heading = "other notes"
		case sub_heading == "ingredients" && strings.HasPrefix(line, "-"):
			r.Ingredients = append(r.Ingredients, createIngredient(removeAndStrip(line, "-")))
		case sub_heading == "method":
			split := strings.SplitN(line, ".", 2)
			r.Method = append(r.Method, strings.TrimSpace(split[1]))
		case sub_heading == "other notes":
			r.OtherNotes += line
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		errMsg := fmt.Sprint("Error reading file:", err)
		panic(errMsg)
	}
	return r
}
