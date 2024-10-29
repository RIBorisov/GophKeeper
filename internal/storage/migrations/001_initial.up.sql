begin transaction;

create table if not exists users(
    id uuid primary key default gen_random_uuid(),
    login varchar(200),
    password bytea
);

create unique index if not exists idx__login_is_unique ON users (login);

commit;