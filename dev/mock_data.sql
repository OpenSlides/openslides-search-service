-- This script should only be used for an database filled with base_data.sql.

BEGIN;

INSERT INTO gender_t (name) VALUES ('power wizard');

INSERT INTO motion_t (sequential_number, title, state_id) VALUES (0, "test", 0);

COMMIT;
