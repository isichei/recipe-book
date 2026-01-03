package filesyncer

import (
	"maps"
	"testing"

	"github.com/isichei/recipe-book/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDbFileCache(t *testing.T) {
	// Create an in-memory sqlite db and run migrations
	db, err := database.CreateDbConnection(":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %s", err)
	}
	defer db.Close()

	database.RunDbMigrations(db)

	// Create DbFileCache with the db connection
	fc := &DbFileCache{dbEngine: db}

	// Test All() returns no data initially
	count := 0
	for range fc.All() {
		count++
	}
	assert.Equal(t, 0, count, "Expected no items in empty file cache")

	// Add some FileCacheData
	testData := map[string]fileCacheData{
		"recipe-a.md": {md5: "abc123", synced: false},
		"recipe-b.md": {md5: "def456", synced: true},
		"recipe-c.md": {md5: "ghi789", synced: false},
	}

	for filename, data := range testData {
		fc.Add(filename, data)
	}

	foundItems := maps.Collect(fc.All())
	assert.Equal(t, len(testData), len(foundItems), "Expected all added items to be returned")
	for filename, expectedData := range testData {
		actualData, exists := foundItems[filename]
		assert.True(t, exists, "Expected %s to exist in All() results", filename)
		assert.Equal(t, expectedData.md5, actualData.md5, "MD5 mismatch for %s", filename)
		assert.Equal(t, expectedData.synced, actualData.synced, "Synced mismatch for %s", filename)
	}

	// Call Get() for a specific item
	recipeA, exists := fc.Get("recipe-a.md")
	assert.True(t, exists, "Expected recipe-a.md to exist")
	assert.Equal(t, "abc123", recipeA.md5, "Expected MD5 to match for recipe-a.md")
	assert.Equal(t, false, recipeA.synced, "Expected synced to be false for recipe-a.md")

	// Change the FileCacheData and call Add() again (update)
	updatedData := fileCacheData{md5: "updated_md5", synced: true}
	fc.Add("recipe-a.md", updatedData)

	// Check that it was updated correctly
	recipeAUpdated, exists := fc.Get("recipe-a.md")
	assert.True(t, exists, "Expected recipe-a.md to still exist after update")
	assert.Equal(t, "updated_md5", recipeAUpdated.md5, "Expected MD5 to be updated")
	assert.Equal(t, true, recipeAUpdated.synced, "Expected synced to be updated to true")

	// Verify total count is still the same (update, not insert)
	countAfterUpdate := 0
	for range fc.All() {
		countAfterUpdate++
	}
	assert.Equal(t, len(testData), countAfterUpdate, "Expected same number of items after update")

	// Call Get() on a non-existing filename
	nonExistent, exists := fc.Get("does-not-exist.md")
	assert.False(t, exists, "Expected non-existent file to return exists=false")
	assert.Equal(t, "", nonExistent.md5, "Expected empty md5 for non-existent file")
	assert.Equal(t, false, nonExistent.synced, "Expected synced=false for non-existent file")
}
