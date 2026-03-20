ALTER TABLE students ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id);
ALTER TABLE students ADD CONSTRAINT students_user_id_unique UNIQUE (user_id);
