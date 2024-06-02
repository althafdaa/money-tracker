ALTER TABLE "subcategory"
ALTER COLUMN "name" TYPE varchar(255),
ALTER COLUMN "name"
SET
    NOT NULL,
    ADD CONSTRAINT "unique_name" UNIQUE ("name");