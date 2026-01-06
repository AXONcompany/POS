-- Initial schema based on current domain entities (User, Category, Product, Order, Table)

create table if not exists users (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  email text not null unique,
  password text not null
);

create table if not exists categories (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  category_name text not null
);

create table if not exists products (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  name text not null,
  price numeric(12,2) not null default 0,
  notes text not null default '',
  category_id bigint not null references categories(id)
);

create table if not exists orders (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  order_date timestamptz not null default now(),
  total numeric(12,2) not null default 0,
  client text not null default ''
);

create table if not exists tables (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  number int not null,
  capacity int not null,
  is_available boolean not null default true,
  occupied_at timestamptz null,
  released_at timestamptz null,
  order_id bigint null references orders(id)
);

create table if not exists order_products (
  order_id bigint not null references orders(id) on delete cascade,
  product_id bigint not null references products(id),
  primary key (order_id, product_id)
);

create index if not exists idx_products_category_id on products(category_id);
create index if not exists idx_tables_order_id on tables(order_id);
create index if not exists idx_order_products_product_id on order_products(product_id);
