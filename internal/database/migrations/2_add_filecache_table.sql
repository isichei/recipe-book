CREATE TABLE IF NOT EXISTS file_cache (
    id TEXT PRIMARY KEY,  /* recipe id */
    md5 TEXT,
    deleted INTEGER,
    last_edited TEXT,
    synced INTEGER
);

UPDATE schema_version SET version = 2;
