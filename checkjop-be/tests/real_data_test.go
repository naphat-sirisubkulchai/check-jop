package tests

// real_data_test.go — integration-style unit tests that use the actual course
// data from csv_data/version10 via FakeRepoFromCSV instead of hand-crafted mocks.
//
// Each test maps to a sensitive_case JSON file. The study-plan items come
// directly from those files; the course graph comes from the real CSV data.

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// shared repos — loaded once per package test run (lazy init via sync.Once would
// also work; for simplicity we reload in each test since tests are parallelisable)
func mustLoadRealData(t *testing.T) *RealDataRepos {
	t.Helper()
	repos, err := LoadRealData()
	require.NoError(t, err, "failed to load real CSV data")
	return repos
}

func realGraduationService(repos *RealDataRepos) service.GraduationService {
	return service.NewGraduationService(repos.Curriculum, repos.Course, repos.Category)
}

// ─── 1. Prereq: 2301108 ต้องการ 2301107 แต่ลงพร้อมกันในเทอมเดียวกัน ──────────
// sensitive_case/1.Prereg.json
// ผลลัพธ์ที่คาดหวัง: violation — 2301108 ถูกลงพร้อม prerequisite 2301107 ในเทอมเดียวกัน

func TestRealData_Case1_Prereq_WrongTerm(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301108", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301107", Year: 2566, Semester: 2, Credits: 3},
			{CourseCode: "2301234", Year: 2566, Semester: 2, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// 2301108 requires 2301107 as prerequisite — 2301108 sem1, 2301107 sem2 → violation
	// 2301234 requires 2301220 — not taken → violation
	// Both should appear
	violationCodes := make(map[string]bool)
	for _, v := range violations {
		violationCodes[v.CourseCode] = true
	}
	assert.True(t, violationCodes["2301108"], "expected violation for 2301108")
	assert.True(t, violationCodes["2301234"], "expected violation for 2301234")
}

// ─── 2. Coreq missing: 2303108 ต้องการ coreq 2303107 แต่ไม่ได้ลง ────────────
// sensitive_case/2.Coreq-2.json
// ผลลัพธ์ที่คาดหวัง: violation — 2303108 ขาด coreq 2303107

func TestRealData_Case2_Coreq_Missing(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2303108", Year: 2566, Semester: 2, Credits: 1},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Equal(t, "2303108", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2303107")
}

// ─── 3. Coreq same term: 2301107 + 2301108 ลงพร้อมกัน ───────────────────────
// sensitive_case/3.Coreq_Same.json — 2301108 prereq = 2301107, ลง sem เดียวกัน
// ผลลัพธ์ที่คาดหวัง: violation — 2301108 ต้องการ prereq ลงก่อน ไม่ใช่พร้อมกัน

func TestRealData_Case3_Prereq_SameTerm(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301107", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301108", Year: 2566, Semester: 1, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// 2301108 has 2301107 as prerequisite → must be taken BEFORE not same term
	require.Len(t, violations, 1)
	assert.Equal(t, "2301108", violations[0].CourseCode)
	assert.Contains(t, violations[0].PrereqsTakenInWrongTerm, "2301107")
}

// ─── 4. Coreq same term valid: 2303107 + 2303108 (coreq) ลงพร้อมกัน ─────────
// sensitive_case/4.Coreq_Same.json
// ผลลัพธ์ที่คาดหวัง: ไม่มี violation

func TestRealData_Case4_Coreq_SameTerm_Valid(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2303107", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2303108", Year: 2566, Semester: 1, Credits: 1},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ─── 5. Prereq wrong term: 2301108 ลงก่อน 2301107 ────────────────────────────
// sensitive_case/5.Prereq_WrongTerm.json
// 2301108 sem1, 2301107 sem2 → 2301108 ลงก่อน prereq → violation

func TestRealData_Case5_Prereq_TakenBefore(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301108", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301107", Year: 2566, Semester: 2, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// 2301108 sem1, 2301107 sem2 → 2301107 is in completedMap but taken AFTER → PrereqsTakenInWrongTerm
	require.Len(t, violations, 1)
	assert.Equal(t, "2301108", violations[0].CourseCode)
	assert.Contains(t, violations[0].PrereqsTakenInWrongTerm, "2301107")
}

// ─── 6. Coreq wrong term: 2303108 ลง sem1 แต่ coreq 2303107 ลง sem2 ──────────
// sensitive_case/6.Coreq-WrongTerm.json
// ผลลัพธ์ที่คาดหวัง: violation — coreq ต้องลงเทอมเดียวกัน

func TestRealData_Case6_Coreq_WrongTerm(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2303107", Year: 2566, Semester: 2, Credits: 3},
			{CourseCode: "2303108", Year: 2566, Semester: 1, Credits: 1},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Equal(t, "2303108", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2303107")
}

// ─── 7. Coreq → F → Coreq again ──────────────────────────────────────────────
// sensitive_case/7.Co->F->Co.json
// 2303107 sem1 F, 2303108 sem1 F, 2303108 sem2 B (retake without coreq)
// ผลลัพธ์ที่คาดหวัง: violation ในเทอม 2 — 2303108 ขาด coreq 2303107

func TestRealData_Case7_Coreq_F_Retake(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกโท-สหกิจ-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2303107", Year: 2566, Semester: 1, Credits: 3, Grade: "F"},
			{CourseCode: "2303108", Year: 2566, Semester: 1, Credits: 1, Grade: "F"},
			{CourseCode: "2303108", Year: 2566, Semester: 2, Credits: 1, Grade: "B"},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// sem1: 2303108 with coreq 2303107 same term → valid (both F)
	// sem2: 2303108 retaken but 2303107 not retaken → missing coreq
	found := false
	for _, v := range violations {
		if v.CourseCode == "2303108" {
			found = true
		}
	}
	assert.True(t, found, "expected violation for 2303108 retake without coreq")
}

// ─── 8a. OR C.F. — no permission: 2301234 ต้องการ 2301220 OR C.F. ──────────
// sensitive_case/8.OR-C.F.-1json.json
// ใช้ test_year CSV ที่ 2301234 ปี 2567 มี prereq = "2301220 OR C.F."
// นักศึกษาลง 2301234 ปี 2567 โดยไม่มี 2301220 และไม่มี exemption → violation

func TestRealData_Case8a_ORCF_NoPermission(t *testing.T) {
	repos, err := LoadRealDataWithTestYear()
	require.NoError(t, err)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301234", Year: 2567, Semester: 1, Credits: 3},
		},
		Exemptions: []string{},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Equal(t, "2301234", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301220")
}

// ─── 8b. OR C.F. — with permission ──────────────────────────────────────────
// sensitive_case/8.OR-C.F.-2.json
// นักศึกษาลง 2301234 ปี 2567 ไม่มี 2301220 แต่มี exemption 2301234 → ผ่าน

func TestRealData_Case8b_ORCF_WithPermission(t *testing.T) {
	repos, err := LoadRealDataWithTestYear()
	require.NoError(t, err)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301234", Year: 2567, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"2301234"},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ─── 13a. Year fallback: 2301286 ปี 2567 มี prereq (test_year CSV) ───────────
// sensitive_case/13.YearWithMock-2567.json
// ปี 2567 2301286 ต้องการ prereq แต่ version10_for_test_year ไม่ใช่ default CSV
// ทดสอบว่าเมื่อลง 2301286 ปี 2567 โดยไม่มีอะไร → ใช้ catalog 2567 ถ้ามี

func TestRealData_Case13a_Year2567_NoPrereq(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301286", Year: 2567, Semester: 1, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// In the standard 2567 CSV, 2301286 has no prerequisite → no violation
	assert.Empty(t, violations)
}

// ─── 13b. Year fallback: 2301286 ปี 2566 ไม่มี prereq ───────────────────────
// sensitive_case/13.YearWithoutMock-2566.json
// ผลลัพธ์ที่คาดหวัง: ไม่มี violation (2566 catalog ไม่มี prereq)

func TestRealData_Case13b_Year2566_NoPrereq(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301286", Year: 2566, Semester: 1, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ─── 14a. Only C.F. — no permission: 2301290 ต้องการ C.F. เท่านั้น ──────────
// sensitive_case/14.Only C.F.-NotPass.json

func TestRealData_Case14a_OnlyCF_NoPermission(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301290", Year: 2566, Semester: 1, Credits: 1},
		},
		Exemptions: []string{},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Equal(t, "2301290", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "C.F.")
}

// ─── 14b. Only C.F. — with permission ────────────────────────────────────────
// sensitive_case/14.Only C.F.-Pass.json

func TestRealData_Case14b_OnlyCF_WithPermission(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301290", Year: 2566, Semester: 1, Credits: 1},
		},
		Exemptions: []string{"2301290"},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ─── 11. CS elective: 2301350 ต้องการ prereq 2301170 OR 2301173 ───────────────
// sensitive_case/11.CS.json
// ลง 2301350 โดยไม่มี prereq → violation

func TestRealData_Case11_CS_Elective_MissingPrereq(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301350", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301250", Year: 2566, Semester: 1, Credits: 3},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	// 2301350 requires 2301170 OR 2301173 — neither present → violation
	found := false
	for _, v := range violations {
		if v.CourseCode == "2301350" {
			found = true
			assert.True(t, len(v.MissingPrereqs) > 0)
		}
	}
	assert.True(t, found, "expected violation for 2301350")
	// 2301250 has no prereq → no violation for it
	for _, v := range violations {
		assert.NotEqual(t, "2301250", v.CourseCode)
	}
}

// ─── 2301260 prerequisite: (2301170 AND 2301172) OR 2301173 ──────────────────
// ทดสอบ complex prereq expression จาก CSV จริง

func TestRealData_2301260_ComplexPrereq_Satisfied_Via173(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2566, Semester: 1, Credits: 4},
			{CourseCode: "2301260", Year: 2566, Semester: 2, Credits: 4},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestRealData_2301260_ComplexPrereq_Satisfied_Via170And172(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2566, Semester: 1, Credits: 1},
			{CourseCode: "2301260", Year: 2566, Semester: 2, Credits: 4},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestRealData_2301260_ComplexPrereq_Missing(t *testing.T) {
	repos := mustLoadRealData(t)
	svc := realGraduationService(repos)

	curr, err := repos.Curriculum.GetByName("เอกเดี่ยว-ฝึกงาน-2566")
	require.NoError(t, err)

	progress := &model.StudentProgress{
		CurriculumID:  curr.ID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			// only 2301170, missing 2301172 AND missing 2301173
			{CourseCode: "2301170", Year: 2566, Semester: 1, Credits: 3},
			{CourseCode: "2301260", Year: 2566, Semester: 2, Credits: 4},
		},
	}

	violations, err := svc.ValidatePrerequisites(progress)
	assert.NoError(t, err)
	found := false
	for _, v := range violations {
		if v.CourseCode == "2301260" {
			found = true
		}
	}
	assert.True(t, found, "expected violation for 2301260 when only 2301170 taken without 2301172")
}
