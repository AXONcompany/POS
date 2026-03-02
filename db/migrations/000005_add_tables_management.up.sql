CREATE TABLE IF NOT EXISTS tables (
  id_table bigserial PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz NULL,

  table_number integer NOT NULL UNIQUE,
  capacity integer NOT NULL,
  status varchar(16) NOT NULL DEFAULT 'LIBRE', 
  arrival_time timestamptz NULL
);

CREATE TABLE IF NOT EXISTS table_waitress (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    deleted_at timestamptz NULL,

    -- FOREIGN KEYS
    table_id bigint NOT NULL,
    waitress_id bigint NOT NULL,

    -- Ya corregido: apunta a id_table
    CONSTRAINT fk_table FOREIGN KEY (table_id) REFERENCES tables(id_table) ON DELETE CASCADE,
    CONSTRAINT fk_waitress FOREIGN KEY (waitress_id) REFERENCES waitress(id_user) ON DELETE CASCADE
);