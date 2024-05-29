create table
    "subcategory" (
        "id" serial primary key,
        "name" varchar(255) not null unique,
        "slug" varchar(255) not null unique,
        "category_id" integer not null,
        "user_id" integer not null
    );

alter table "subcategory" add constraint "subcategory_category_id_foreign" foreign key ("category_id") references "category" ("id") on delete cascade;

alter table "subcategory" add constraint "subcategory_user_id_foreign" foreign key ("user_id") references "user_data" ("id") on delete cascade;