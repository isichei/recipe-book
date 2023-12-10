package recipes

import (
	"bufio"
	"fmt"
	"github.com/isichei/recipe-book/types"
	"os"
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

// Super niave parser - does the job
func ParseMarkdownFile(filepath string) types.Recipe {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic("ARGHHHGHGHH!")
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file getting recipe values
	r := types.Recipe{}
	sub_heading := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fmt.Println("===")
		fmt.Printf("Line: %s\n", line)
		switch {
		case line == "":
			fmt.Println("case empty")
		case strings.HasPrefix(line, "# "):
			r.Title = removeAndStrip(line, "#")
			fmt.Println("case title")
		case strings.HasPrefix(line, PREP_TIME_STR):
			r.PrepTime = removeAndStrip(line, PREP_TIME_STR)
			fmt.Println("case prep")
		case strings.HasPrefix(line, COOKING_TIME_STR):
			r.CookingTime = removeAndStrip(line, COOKING_TIME_STR)
			fmt.Println("case cooking")
		case strings.HasPrefix(line, SERVES_STR):
			r.Serves = removeAndStrip(line, SERVES_STR)
			fmt.Println("case serves")
		case strings.HasPrefix(line, SOURCE_STR):
			r.Source = removeAndStrip(line, SOURCE_STR)
			fmt.Println("case source")
		case r.Serves != "" && r.Source == "":
			r.Description += line
			fmt.Println("case description")
		case line == "## Ingredients:":
			sub_heading = "ingredients"
			fmt.Println("case set I")
		case line == "## Method:":
			sub_heading = "method"
			fmt.Println("case set M")
		case line == "## Other notes:":
			sub_heading = "other notes"
			fmt.Println("case set ON")
		case sub_heading == "ingredients" && strings.HasPrefix(line, "-"):
			r.Ingredients = append(r.Ingredients, removeAndStrip(line, "-"))
			fmt.Println("case ingredients")
		case sub_heading == "method":
			split := strings.SplitN(line, ".", 2)
			fmt.Println(split)
			r.Method = append(r.Method, strings.TrimSpace(split[1]))
			fmt.Println("case method")
		case sub_heading == "other notes":
			r.OtherNotes += line
			fmt.Println("case other notes")
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		panic("!")
	}
	return r
}

