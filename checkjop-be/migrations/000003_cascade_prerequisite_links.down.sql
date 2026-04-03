ALTER TABLE prerequisite_course_links DROP CONSTRAINT IF EXISTS fk_prerequisite_groups_prerequisite_courses;
ALTER TABLE prerequisite_course_links ADD CONSTRAINT fk_prerequisite_groups_prerequisite_courses
    FOREIGN KEY (group_id) REFERENCES prerequisite_groups(id);

ALTER TABLE prerequisite_course_links DROP CONSTRAINT IF EXISTS fk_prerequisite_course_links_prerequisite_course;
ALTER TABLE prerequisite_course_links ADD CONSTRAINT fk_prerequisite_course_links_prerequisite_course
    FOREIGN KEY (prerequisite_course_id) REFERENCES courses(id);
