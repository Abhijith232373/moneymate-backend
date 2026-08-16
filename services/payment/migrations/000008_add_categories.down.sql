-- 0000XX_add_categories.down.sql
ALTER TABLE payment.transactions DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS payment.categories;

