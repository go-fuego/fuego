ALTER TABLE recipe ADD COLUMN category text NOT NULL DEFAULT 'other'; -- Breakfast, Lunch, Dinner, Dessert, Snack, Other
ALTER TABLE recipe ADD COLUMN class text NOT NULL DEFAULT 'other'; -- Pasta, Soup, Salad, Sandwich, Other
ALTER TABLE recipe ADD COLUMN published boolean NOT NULL DEFAULT false;
ALTER TABLE recipe ADD COLUMN created_by text NOT NULL DEFAULT 'admin';
ALTER TABLE recipe ADD COLUMN calories integer NOT NULL DEFAULT 0;
ALTER TABLE recipe ADD COLUMN cost integer NOT NULL DEFAULT 0;
ALTER TABLE recipe ADD COLUMN prep_time integer NOT NULL DEFAULT 0;
ALTER TABLE recipe ADD COLUMN cook_time integer NOT NULL DEFAULT 0;
ALTER TABLE recipe ADD COLUMN servings integer NOT NULL DEFAULT 0;
ALTER TABLE recipe ADD COLUMN image_url text NOT NULL DEFAULT '';
ALTER TABLE recipe ADD COLUMN disclaimer text NOT NULL DEFAULT '';
