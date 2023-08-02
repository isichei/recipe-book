# A Simple Personal Recipe Book App

The idea of this is to have a Go App using htmx to search through recipes that I have stored in S3 (or as markdown files).

To run:

```
go run main.go
```

Note API also requires a `thumbnails/` folder for which has a thumbnail for each recipe with naming convention `<types.RecipeMetadata.Uid>.jpg`
