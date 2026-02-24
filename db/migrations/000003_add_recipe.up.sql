create table if not exists products (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    product_name varchar(255) not null,
    sales_price decimal(10, 2) not null default,
    is_active boolean not null,
);

create table if not exists categories (

    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    category_name varchar(255) not null,

    description text,
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

CREATE INDEX idx_recipe_product ON recipe_items(product_id);
CREATE INDEX idx_recipe_ingredient ON recipe_items(ingredient_id);