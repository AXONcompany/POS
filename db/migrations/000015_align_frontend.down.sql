ALTER TABLE categories  DROP COLUMN IF EXISTS color_class, DROP COLUMN IF EXISTS icon;
ALTER TABLE products    DROP COLUMN IF EXISTS description, DROP COLUMN IF EXISTS image_url, DROP COLUMN IF EXISTS category_id;
ALTER TABLE tables      DROP COLUMN IF EXISTS table_name, DROP COLUMN IF EXISTS x, DROP COLUMN IF EXISTS y,
                        DROP COLUMN IF EXISTS width, DROP COLUMN IF EXISTS height, DROP COLUMN IF EXISTS shape,
                        DROP COLUMN IF EXISTS rotation, DROP COLUMN IF EXISTS color, DROP COLUMN IF EXISTS floor,
                        DROP COLUMN IF EXISTS is_merged, DROP COLUMN IF EXISTS merged_from,
                        DROP COLUMN IF EXISTS guests, DROP COLUMN IF EXISTS assigned_waiter_id;
ALTER TABLE users       DROP COLUMN IF EXISTS pin_hash;
