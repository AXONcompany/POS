create if not exists tables(
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