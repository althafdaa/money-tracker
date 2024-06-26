create table
    "user_data" (
        "id" int generated by default as identity primary key,
        "name" varchar(255) not null,
        "email" varchar(255) not null unique,
        "hash" varchar(255) not null,
        "profile_picture_url" varchar(255),
        "created_at" timestamp DEFAULT 'now()' NOT NULL,
        "updated_at" timestamp DEFAULT 'now()' NOT NULL,
        "deleted_at" timestamp
    );

create index "user_email_index" on "user_data" ("email");