ALTER TABLE students DROP CONSTRAINT IF EXISTS students_user_id_unique;
ALTER TABLE students DROP COLUMN IF EXISTS user_id;
