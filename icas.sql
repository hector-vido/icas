-- htpasswd -nbBC 12 "" "mypassword" | tr -d ':\n'
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL DEFAULT "", email TEXT NOT NULL DEFAULT "", password TEXT NOT NULL DEFAULT "");
CREATE TABLE animals (id INTEGER PRIMARY KEY, name TEXT NOT NULL DEFAULT "");
CREATE TABLE animals_coordinates (animal_id INTEGER NOT NULL, latitude REAL NOT NULL, longitude REAL NOT NULL, registered_at TEXT NOT NULL DEFAULT "");
