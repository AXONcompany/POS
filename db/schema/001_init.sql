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

    FOREIGN KEY (table_id) REFERENCES tables(id_table) ON DELETE CASCADE,
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

create table if not exists products (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    product_name varchar(255) not null,
    sales_price decimal(10, 2) not null,
    is_active boolean not null
);

create table if not exists product_categories (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    product_id bigserial not null,
    category_id bigserial not null,

    foreign key (product_id) references products(id) on delete cascade,
    foreign key (category_id) references categories(id) on delete cascade

);

create table if not exists recipe (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    product_id bigint not null,
    ingredient_id bigint not null,

    quantity_required decimal(10, 4) not null,

    foreign key (product_id) references products(id) on delete cascade,
    foreign key (ingredient_id) references ingredients(id) on delete cascade

);

create table if not exists order_statuses (
  id serial primary key,
  name varchar(50) unique not null,
  description text
);

create table if not exists orders (
  id bigserial primary key, 
  restaurant_id integer not null references restaurants(id) on delete cascade,
  table_id bigint references tables(id_table) on delete set null,
  user_id integer not null references users(id), -- mesero
  status_id integer not null references order_statuses(id),
  total_amount decimal(10,2) not null default 0.00,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null
);

create table if not exists order_items (
  id bigserial primary key,
  order_id bigint not null references orders(id) on delete cascade,
  product_id bigint not null references products(id),
  quantity integer not null,
  unit_price decimal(10,2) not null,
  notes text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists sales (
  id bigserial primary key,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz null,

  total decimal(10, 2) not null,
  payment_method varchar(50) not null,
  date timestamptz not null default now(),
  order_id integer not null references orders(id) on delete cascade
);


