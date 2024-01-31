ALTER TABLE recipe DROP COLUMN class;
ALTER TABLE recipe ADD COLUMN when_to_eat text NOT NULL DEFAULT 'other'; -- Breakfast, Lunch, Dinner, Dessert, Snack, Other
