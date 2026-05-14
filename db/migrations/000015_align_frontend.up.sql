-- Users: PIN login for waiters
ALTER TABLE users ADD COLUMN pin_hash VARCHAR(255);

-- Tables: canvas layout data for the frontend editor
ALTER TABLE tables
  ADD COLUMN table_name  VARCHAR(100),
  ADD COLUMN x           INTEGER      NOT NULL DEFAULT 100,
  ADD COLUMN y           INTEGER      NOT NULL DEFAULT 100,
  ADD COLUMN width       INTEGER      NOT NULL DEFAULT 110,
  ADD COLUMN height      INTEGER      NOT NULL DEFAULT 110,
  ADD COLUMN shape       VARCHAR(20)  NOT NULL DEFAULT 'square',
  ADD COLUMN rotation    INTEGER      NOT NULL DEFAULT 0,
  ADD COLUMN color       VARCHAR(50),
  ADD COLUMN floor       INTEGER      NOT NULL DEFAULT 1,
  ADD COLUMN is_merged   BOOLEAN      NOT NULL DEFAULT false,
  ADD COLUMN merged_from JSONB,
  ADD COLUMN guests      INTEGER      NOT NULL DEFAULT 0,
  ADD COLUMN assigned_waiter_id INTEGER REFERENCES users(id);

-- Products: description, image, and direct category link
ALTER TABLE products
  ADD COLUMN description TEXT,
  ADD COLUMN image_url   TEXT,
  ADD COLUMN category_id BIGINT REFERENCES categories(id);

-- Categories: UI metadata (icon + color for the frontend)
ALTER TABLE categories
  ADD COLUMN color_class VARCHAR(50),
  ADD COLUMN icon        VARCHAR(50);
