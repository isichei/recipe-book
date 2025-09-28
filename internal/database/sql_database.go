package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/isichei/recipe-book/internal/recipes"
	_ "github.com/mattn/go-sqlite3"
)

type SqlDatabase struct {
	dbEngine *sql.DB
}

// Search the SQL lite DB for matching terms in the given text
func (db SqlDatabase) SearchRecipes(text string) []recipes.RecipeMetadata {
	var recipeUid, title, desc, likeSearch string

	if strings.TrimSpace(text) == "" {
		return db.getSomeRecipes(9)
	}
	foundRecipes := make(Set)
	recipeMetadata := []recipes.RecipeMetadata{}
	searchTerms := strings.Split(strings.ToLower(text), " ")

	query := `SELECT id, title, description
	FROM recipes
	WHERE title LIKE ? OR description LIKE ?;
	`

	for i, searchTerm := range searchTerms {
		likeSearch = "%" + searchTerm + "%"
		rows, err := db.dbEngine.Query(query, likeSearch, likeSearch)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&recipeUid, &title, &desc); err != nil {
				log.Fatal(err)
			}
			if !foundRecipes[recipeUid] {
				recipeMetadata = append(recipeMetadata, recipes.RecipeMetadata{Uid: recipeUid, Title: title, Description: desc})
				foundRecipes.Add(recipeUid)
			}
		}
		if i > 10 {
			break // who searches that much
		}
	}
	return recipeMetadata
}

func (db SqlDatabase) getSomeRecipes(n int) []recipes.RecipeMetadata {
	var recipeUid, title, desc string
	recipeMetadata := []recipes.RecipeMetadata{}

	query := `SELECT id, title, description
	FROM recipes
	limit ?;
	`
	rows, err := db.dbEngine.Query(query, n)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&recipeUid, &title, &desc)
		recipeMetadata = append(recipeMetadata, recipes.RecipeMetadata{Uid: recipeUid, Title: title, Description: desc})
	}
	return recipeMetadata
}

func (db SqlDatabase) GetRecipeMetadata(recipeUid string) recipes.RecipeMetadata {
	var title, desc string

	query := `SELECT title, description
	FROM recipes
	WHERE id = ?;
	`

	err := db.dbEngine.QueryRow(query, recipeUid).Scan(&title, &desc)
	if err != nil {
		log.Printf("%s. No recipe data found for %s\n", err, recipeUid)
		return recipes.RecipeMetadata{}
	}
	return recipes.RecipeMetadata{Uid: recipeUid, Title: title, Description: desc}
}

func (db SqlDatabase) GetRecipe(recipeUid string) recipes.Recipe {
	var title, description, prepTime, cookingTime, serves, otherNotes, source string

	// Get main part of recipe
	query := `SELECT title, description, prep_time, cooking_time, serves, other_notes, source
	FROM recipes
	WHERE id = ?;
	`
	err := db.dbEngine.QueryRow(query, recipeUid).Scan(&title, &description, &prepTime, &cookingTime, &serves, &otherNotes, &source)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("No recipe data found for %s\n", recipeUid)
			return recipes.Recipe{}
		}
		log.Fatal(err)
	}

	ingredients, err := db.getIngredients(recipeUid)
	if err != nil {
		log.Fatal(fmt.Sprintf("Recipe %s is missing ingredients or could not get ingredients from db - %s", recipeUid, err))
	}

	method, err := db.getMethod(recipeUid)
	if err != nil {
		log.Fatal(fmt.Sprintf("Recipe %s is missing method or could not get method from db - %s", recipeUid, err))
	}

	r := recipes.Recipe{
		Title:       title,
		Description: description,
		PrepTime:    prepTime,
		CookingTime: cookingTime,
		Serves:      serves,
		Ingredients: ingredients,
		Method:      method,
		OtherNotes:  otherNotes,
		Source:      source,
	}
	return r
}

// Get Ingredients from DB
func (db SqlDatabase) getIngredients(recipeUid string) ([]recipes.Ingredient, error) {
	var name, amount string

	ingredients := []recipes.Ingredient{}
	// Get ingredients
	query := `SELECT name, amount
	FROM ingredients
	WHERE recipe_id = ?
	ORDER by id;
	`
	rows, err := db.dbEngine.Query(query, recipeUid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name, &amount)
		ingredients = append(ingredients, recipes.Ingredient{Name: name, Amount: amount})
	}
	return ingredients, nil
}

// Get Methods from DB
func (db SqlDatabase) getMethod(recipeUid string) ([]string, error) {
	var step string
	method := []string{}

	query := `SELECT step
	FROM methods
	WHERE recipe_id = ?
	ORDER by id;
	`
	rows, err := db.dbEngine.Query(query, recipeUid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&step)
		method = append(method, step)
	}
	return method, nil
}

// Initialize connection and return the DB
func NewSqlDatabase(dataSourceName string) (*SqlDatabase, error) {
	// Open up the DB and init
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to the DB")
	}

	// Create main recipe table
	query := `CREATE TABLE IF NOT EXISTS recipes (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	prep_time TEXT,
	cooking_time TEXT,
	serves TEXT,
	other_notes TEXT,
	source TEXT
    );`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Create the ingredients table
	query = `CREATE TABLE IF NOT EXISTS ingredients (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	recipe_id TEXT,
	name TEXT NOT NULL,
	amount TEXT NULL,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Create the methods table
	query = `CREATE TABLE IF NOT EXISTS methods (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	recipe_id TEXT,
	step TEXT NOT NULL,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	query = "SELECT count(*) as n from recipes"
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total number of recipes in %s: %d\n", dataSourceName, count)
	return &SqlDatabase{dbEngine: db}, nil
}

// Close method to be deferred
func (s *SqlDatabase) Close() error {
	return s.dbEngine.Close()
}

// Write data to the SQL database
func (db *SqlDatabase) AddRecipe(rUid string, r recipes.Recipe) error {
	query := `
	INSERT INTO recipes (id, title, description, prep_time, cooking_time, serves, other_notes, source)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	result, err := db.dbEngine.Exec(query, rUid, r.Title, r.Description, r.PrepTime, r.CookingTime, r.Serves, r.OtherNotes, r.Source)
	if err != nil {
		log.Println(err)
		return err
	}

	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); rowsAffected != 1 || err != nil {
		log.Fatal(fmt.Sprintf("%d %s", rowsAffected, err))
	}

	// Add ingredients
	query = `INSERT INTO ingredients (recipe_id, name, amount)
	VALUES (?, ?, ?);
	`
	for _, ingredient := range r.Ingredients {
		_, err := db.dbEngine.Exec(query, rUid, ingredient.Name, ingredient.Amount)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Add methods
	query = `INSERT INTO methods (recipe_id, step)
	VALUES (?, ?);
	`
	for _, step := range r.Method {
		_, err := db.dbEngine.Exec(query, rUid, step)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (db *SqlDatabase) AddFiles(dir string) {
	fileGetter := LocalMarkdownFileGetter{dir}

	for _, file := range fileGetter.files() {
		uid, fullRecipe := fileGetter.getRecipeFromFilePath(file)
		fmt.Printf("File %s added...\n", file)
		db.AddRecipe(uid, fullRecipe)
	}
}
