create table
    "refresh_token" (
        "id" int generated by default as identity primary key,
        "refresh_token" varchar(255) not null unique,
        "access_token" varchar(255) not null,
        "user_id" int not null,
        "created_at" timestamp DEFAULT 'now()' NOT NULL,
        "updated_at" timestamp DEFAULT 'now()' NOT NULL,
        "expired_at" timestamp not null,
        "deleted_at" timestamp
    );

alter table "refresh_token" add constraint "refresh_token_user_id_foreign" foreign key ("user_id") references "user" ("id") on delete set null;