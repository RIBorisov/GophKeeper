begin transaction;

drop index if exists idx__login_is_unique;

drop table if exists users;

commit;