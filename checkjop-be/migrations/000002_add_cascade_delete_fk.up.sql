-- Update foreign key constraints to use ON DELETE CASCADE
-- This ensures child records are automatically deleted when parent records are removed

-- courses.category_id -> categories.id
ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_categories_courses;
ALTER TABLE courses ADD CONSTRAINT fk_categories_courses
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

-- courses.curriculum_id -> curriculums.id
ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_curriculums_courses;
ALTER TABLE courses ADD CONSTRAINT fk_curriculums_courses
    FOREIGN KEY (curriculum_id) REFERENCES curriculums(id) ON DELETE CASCADE;

-- categories.curriculum_id -> curriculums.id
ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_curriculums_categories;
ALTER TABLE categories ADD CONSTRAINT fk_curriculums_categories
    FOREIGN KEY (curriculum_id) REFERENCES curriculums(id) ON DELETE CASCADE;
