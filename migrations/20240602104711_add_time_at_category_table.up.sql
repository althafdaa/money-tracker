ALTER TABLE "category"
ADD COLUMN "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
ADD COLUMN "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
ADD COLUMN "deleted_at" TIMESTAMP;

alter table "subcategory"
add column "created_at" timestamp default current_timestamp not null,
add column "updated_at" timestamp default current_timestamp not null,
add column "deleted_at" timestamp;