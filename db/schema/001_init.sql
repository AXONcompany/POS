-- Schema snapshot

create table if not exists users (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  email text not null unique,
  password text not null
);

create table if not exists waitress(
  id_user bigint primary key,
  FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);

create table if not exists tables(
  id_table bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  table_number integer not null unique,
  capacity integer not null,
  status varchar(16) not null,
  arrival_time timestamptz
);

create table if not exists table_waitress (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    --FOREIGN KEYS
    table_id bigint not null,
    waitress_id bigint not null,

    FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE,
    FOREIGN KEY (waitress_id) REFERENCES waitress(id_user) ON DELETE CASCADE
);


create table if not exists categories (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  category_name text not null
);

create table if not exists ingredients (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  ingredient_name varchar(124) not null,
  unit_of_measure varchar(8) not null,
  ingredient_type varchar(24) not null,
  stock bigint not null default 0
);

