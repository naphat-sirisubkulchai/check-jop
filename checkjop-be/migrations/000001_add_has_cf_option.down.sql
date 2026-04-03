DROP INDEX IF EXISTS idx_courses_has_cf_option;
ALTER TABLE courses DROP COLUMN IF EXISTS has_cf_option;
