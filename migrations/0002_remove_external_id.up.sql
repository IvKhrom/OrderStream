-- Historical migration: in earlier versions we had external_id column.
-- Now we keep external identifier inside payload->>'id', so we ensure external_id is absent.

ALTER TABLE IF EXISTS orders DROP COLUMN IF EXISTS external_id;

-- Historical migration: in earlier versions we had external_id column.
-- Now we keep external identifier inside payload->>'id', so we ensure external_id is absent.

ALTER TABLE IF EXISTS orders DROP COLUMN IF EXISTS external_id;

-- Remove external_id column and its index; we keep payload->>'id' as external identifier
DROP INDEX IF EXISTS idx_orders_external_user;

-- Drop column from parent partitioned table; this will remove it from partitions as well
ALTER TABLE IF EXISTS orders DROP COLUMN IF EXISTS external_id;

-- sanity: ensure no leftover indexes referencing external_id
DROP INDEX IF EXISTS idx_orders_external_user;
