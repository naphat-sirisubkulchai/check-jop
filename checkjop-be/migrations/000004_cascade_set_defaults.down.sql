ALTER TABLE set_defaults DROP CONSTRAINT IF EXISTS fk_set_defaults_course;
ALTER TABLE set_defaults ADD CONSTRAINT fk_set_defaults_course
    FOREIGN KEY (course_id) REFERENCES courses(id);

ALTER TABLE set_defaults DROP CONSTRAINT IF EXISTS fk_set_defaults_curriculum;
ALTER TABLE set_defaults ADD CONSTRAINT fk_set_defaults_curriculum
    FOREIGN KEY (curriculum_id) REFERENCES curriculums(id);
