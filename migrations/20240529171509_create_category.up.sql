create table
    "category" (
        "id" serial primary key,
        "name" varchar(255) not null unique,
        "slug" varchar(255) not null unique
    );