begin transaction;

create table if not exists metadata(
    id uuid DEFAULT NULL,
    user_id uuid,
    kind VARCHAR(200) CHECK (kind IN ('binary', 'card', 'credentials', 'text')),
    metadata VARCHAR(512),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

commit;