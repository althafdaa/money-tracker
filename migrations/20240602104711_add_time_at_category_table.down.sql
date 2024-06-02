-- Drop columns from "category" table
ALTER TABLE category
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;

-- Drop columns from "subcategory" table
ALTER TABLE subcategory
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;