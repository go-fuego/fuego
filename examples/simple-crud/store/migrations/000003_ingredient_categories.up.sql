ALTER TABLE ingredient ADD COLUMN category text NOT NULL DEFAULT 'other';
ALTER TABLE ingredient ADD COLUMN available_all_year boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_jan boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_feb boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_mar boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_apr boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_may boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_jun boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_jul boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_aug boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_sep boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_oct boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_nov boolean NOT NULL DEFAULT false;
ALTER TABLE ingredient ADD COLUMN available_dec boolean NOT NULL DEFAULT false;



