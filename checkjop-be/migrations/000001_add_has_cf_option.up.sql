-- Add HasCFOption field to courses table
-- This field indicates whether a course allows C.F. (Consent of Faculty) exemption
-- based on whether C.F. appears in its prerequisites or corequisites

ALTER TABLE courses ADD COLUMN IF NOT EXISTS has_cf_option BOOLEAN DEFAULT FALSE;

-- Create index for faster queries on has_cf_option
CREATE INDEX IF NOT EXISTS idx_courses_has_cf_option ON courses(has_cf_option);

-- Note: To populate this field with existing data, you need to re-import the CSV files
-- The course_service.go will automatically set has_cf_option=true when it finds C.F. in prerequisites or corequisites
