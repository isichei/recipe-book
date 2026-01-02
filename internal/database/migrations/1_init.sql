-- Going to use CREATE TABLE IF NOT EXISTS so that I can just run through each sql file (can't see I am going to have many)
CREATE TABLE IF NOT EXISTS recipes (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	prep_time TEXT,
	cooking_time TEXT,
	serves TEXT,
	other_notes TEXT,
	source TEXT
    );

CREATE TABLE IF NOT EXISTS ingredients (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	recipe_id TEXT,
	name TEXT NOT NULL,
	amount TEXT NULL,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS methods (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	recipe_id TEXT,
	step TEXT NOT NULL,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER
);

INSERT INTO schema_version (version) VALUES (1);
