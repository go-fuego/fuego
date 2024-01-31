ALTER TABLE recipe DROP COLUMN when_to_eat; 
ALTER TABLE recipe ADD COLUMN class text NOT NULL DEFAULT 'other';
