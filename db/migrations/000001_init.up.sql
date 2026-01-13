-- Initial schema based on current domain entities (User, Category, Product, Order, Table)

create table if not exists users (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  email text not null unique,
  password text not null
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