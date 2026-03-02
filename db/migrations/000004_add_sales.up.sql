create table if not exists sales (
    id bigserial primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null,

    total decimal(10, 2) not null,
    payment_method varchar(50) not null,
    date timestamptz not null default now(),
    order_id integer not null references orders(id)
);