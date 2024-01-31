CREATE TABLE IF NOT EXISTS recipe (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    instructions TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS ingredient (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dosing (
  recipe_id TEXT NOT NULL,
  ingredient_id TEXT NOT NULL,
  quantity INTEGER NOT NULL,
  unit TEXT NOT NULL,
  PRIMARY KEY (recipe_id, ingredient_id),
  FOREIGN KEY (recipe_id) REFERENCES recipe(id),
  FOREIGN KEY (ingredient_id) REFERENCES ingredient(id)
);
