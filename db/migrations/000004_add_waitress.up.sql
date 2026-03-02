
create table if not exists waitress(
  id_user bigint primary key,
  FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);
 