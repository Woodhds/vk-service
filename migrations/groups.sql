CREATE TABLE IF NOT EXISTS favorite_groups
(
    id   integer UNIQUE NOT NULL,
    name text,
    avatar text
)