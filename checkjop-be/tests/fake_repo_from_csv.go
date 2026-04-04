package tests

// FakeRepoFromCSV builds in-memory repositories by parsing the real CSV data files.
// It replicates the same prerequisite parsing logic as courseService so that unit tests
// run against realistic course graphs without needing a database connection.
//
// Usage:
//
//	repos, err := LoadRealData()
//	graduationService := service.NewGraduationService(repos.Curriculum, repos.Course, repos.Category)

import (
	"checkjop-be/internal/model"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ─── Entry point ────────────────────────────────────────────────────────────

// RealDataRepos holds the three in-memory repositories built from real CSV data.
type RealDataRepos struct {
	Curriculum *FakeCurriculumRepo
	Course     *FakeCourseRepo
	Category   *FakeCategoryRepo
}

// LoadRealData parses the version10 CSV files and returns populated in-memory repos.
// Path is relative to the project root; adjust if tests are run from a different CWD.
func LoadRealData() (*RealDataRepos, error) {
	return loadFromCSVFiles(false)
}

// LoadRealDataWithTestYear is like LoadRealData but also loads the test_year CSV
// (which contains modified prereq rules for 2567, e.g. 2301234 has "2301220 OR C.F.").
func LoadRealDataWithTestYear() (*RealDataRepos, error) {
	return loadFromCSVFiles(true)
}

func loadFromCSVFiles(includeTestYear bool) (*RealDataRepos, error) {
	base := "../csv_data/version10"

	currRepo := newFakeCurriculumRepo()
	catRepo := newFakeCategoryRepo()
	courseRepo := newFakeCourseRepo()

	// 1. Load curricula
	if err := loadCurricula(currRepo, base+"/example_checkjop - curriculum.csv"); err != nil {
		return nil, fmt.Errorf("curricula: %w", err)
	}

	// 2. Load categories
	if err := loadCategories(catRepo, currRepo, base+"/example_checkjop - catagory.csv"); err != nil {
		return nil, fmt.Errorf("categories: %w", err)
	}

	// 3. Load courses for each year (year derived from filename)
	csvFiles := []struct {
		path string
		year int
	}{
		{base + "/example_checkjop - course_Present_2566-7.csv", 2566},
		{base + "/example_checkjop - course_Present_2567-6.csv", 2567},
		{base + "/example_checkjop - course_Present_2568-6.csv", 2568},
	}

	// test_year overrides 2567 with different prereq rules for testing year-versioning
	if includeTestYear {
		csvFiles[1] = struct {
			path string
			year int
		}{base + "/version10_for_test_year/example_checkjop - course_Present_2567_for_test_year-3.csv", 2567}
	}

	for _, f := range csvFiles {
		if err := loadCourses(courseRepo, catRepo, currRepo, f.path, f.year); err != nil {
			return nil, fmt.Errorf("courses year %d: %w", f.year, err)
		}
	}

	// 4. Wire prerequisite/corequisite groups
	for _, f := range csvFiles {
		if err := wireRelationships(courseRepo, currRepo, f.path, f.year); err != nil {
			return nil, fmt.Errorf("relationships year %d: %w", f.year, err)
		}
	}

	return &RealDataRepos{
		Curriculum: currRepo,
		Course:     courseRepo,
		Category:   catRepo,
	}, nil
}

// ─── CSV loaders ─────────────────────────────────────────────────────────────

func loadCurricula(repo *FakeCurriculumRepo, path string) error {
	records, err := readCSV(path)
	if err != nil {
		return err
	}
	// header: curriculumNameEN,curriculumNameTH,Year,minTotalCredit,isActive
	for i, r := range records[1:] {
		if len(r) < 5 {
			return fmt.Errorf("row %d: too few columns", i+2)
		}
		year, _ := strconv.Atoi(strings.TrimSpace(r[2]))
		min, _ := strconv.Atoi(strings.TrimSpace(r[3]))
		c := &model.Curriculum{
			ID:              uuid.New(),
			NameEN:          strings.TrimSpace(r[0]),
			NameTH:          strings.TrimSpace(r[1]),
			Year:            year,
			MinTotalCredits: min,
			IsActive:        strings.TrimSpace(r[4]) == "TRUE",
		}
		repo.store[c.NameTH] = c
		repo.store[c.NameEN] = c
	}
	return nil
}

func loadCategories(catRepo *FakeCategoryRepo, currRepo *FakeCurriculumRepo, path string) error {
	records, err := readCSV(path)
	if err != nil {
		return err
	}
	// header: categoryNameEn,categoryNameTH,cirrculumName,minCredit,Year
	for i, r := range records[1:] {
		if len(r) < 4 {
			return fmt.Errorf("row %d: too few columns", i+2)
		}
		nameEN := strings.TrimSpace(r[0])
		nameTH := strings.TrimSpace(r[1])
		minCred, _ := strconv.Atoi(strings.TrimSpace(r[3]))
		currNames := strings.Split(strings.TrimSpace(r[2]), ",")

		for _, cn := range currNames {
			cn = strings.TrimSpace(cn)
			if cn == "" {
				continue
			}
			curr, ok := currRepo.store[cn]
			if !ok {
				continue // skip unknown curriculum
			}
			cat := &model.Category{
				ID:           uuid.New(),
				CurriculumID: curr.ID,
				NameEN:       nameEN,
				NameTH:       nameTH,
				MinCredits:   minCred,
			}
			key := catKey(curr.ID, nameEN)
			if _, exists := catRepo.byKey[key]; !exists {
				catRepo.byKey[key] = cat
				catRepo.byCurriculum[curr.ID] = append(catRepo.byCurriculum[curr.ID], *cat)
			}
			keyTH := catKey(curr.ID, nameTH)
			if _, exists := catRepo.byKey[keyTH]; !exists {
				catRepo.byKey[keyTH] = cat
			}
		}
	}
	return nil
}

func loadCourses(courseRepo *FakeCourseRepo, catRepo *FakeCategoryRepo, currRepo *FakeCurriculumRepo, path string, year int) error {
	records, err := readCSV(path)
	if err != nil {
		return err
	}
	// header: code,courseNameEN,courseNameTH,credit,prerequisites,corequisites,category,curriculum
	for i, r := range records[1:] {
		if len(r) < 8 {
			return fmt.Errorf("row %d: too few columns", i+2)
		}
		code := strings.TrimSpace(r[0])
		nameEN := strings.TrimSpace(r[1])
		nameTH := strings.TrimSpace(r[2])
		credits, _ := strconv.Atoi(strings.TrimSpace(r[3]))
		catName := strings.TrimSpace(r[6])
		currNames := strings.Split(strings.TrimSpace(r[7]), ",")

		for _, cn := range currNames {
			cn = strings.TrimSpace(cn)
			if cn == "" {
				continue
			}
			curr, ok := currRepo.store[cn]
			if !ok {
				continue
			}
			// resolve category (try EN then TH)
			cat := catRepo.byKey[catKey(curr.ID, catName)]
			if cat == nil {
				continue
			}
			c := &model.Course{
				ID:                 uuid.New(),
				Code:               code,
				NameEN:             nameEN,
				NameTH:             nameTH,
				Credits:            credits,
				CurriculumID:       curr.ID,
				CategoryID:         cat.ID,
				Year:               year,
				PrerequisiteGroups: []model.PrerequisiteGroup{},
				CorequisiteGroups:  []model.PrerequisiteGroup{},
			}
			key := courseKey(code, curr.ID, year)
			if _, exists := courseRepo.byKey[key]; !exists {
				courseRepo.byKey[key] = c
				courseRepo.byCodeYear[codeYearKey(code, year)] = c
			}
		}
	}
	return nil
}

func wireRelationships(courseRepo *FakeCourseRepo, currRepo *FakeCurriculumRepo, path string, year int) error {
	records, err := readCSV(path)
	if err != nil {
		return err
	}
	for _, r := range records[1:] {
		if len(r) < 8 {
			continue
		}
		code := strings.TrimSpace(r[0])
		preStr := strings.TrimSpace(r[4])
		coStr := strings.TrimSpace(r[5])
		currNames := strings.Split(strings.TrimSpace(r[7]), ",")

		if preStr == "" && coStr == "" {
			continue
		}

		for _, cn := range currNames {
			cn = strings.TrimSpace(cn)
			if cn == "" {
				continue
			}
			curr, ok := currRepo.store[cn]
			if !ok {
				continue
			}
			course, ok := courseRepo.byKey[courseKey(code, curr.ID, year)]
			if !ok {
				continue
			}

			hasCF := false

			if preStr != "" {
				groups := parsePrereqString(preStr)
				for _, g := range groups {
					if g.HasCFCondition {
						hasCF = true
					}
					var links []model.PrerequisiteCourseLink
					for _, pc := range g.CourseCodes {
						// look up in same curriculum first, then any curriculum
						prereq := courseRepo.byKey[courseKey(pc, curr.ID, year)]
						if prereq == nil {
							prereq = courseRepo.byCodeYear[codeYearKey(pc, year)]
						}
						if prereq == nil {
							continue
						}
						links = append(links, model.PrerequisiteCourseLink{
							ID:                   uuid.New(),
							PrerequisiteCourseID: prereq.ID,
							PrerequisiteCourse:   *prereq,
						})
					}
					course.PrerequisiteGroups = append(course.PrerequisiteGroups, model.PrerequisiteGroup{
						ID:                  uuid.New(),
						IsOrGroup:           g.IsOrGroup,
						HasCFCondition:      g.HasCFCondition,
						PrerequisiteCourses: links,
					})
				}
			}

			if coStr != "" {
				groups := parsePrereqString(coStr)
				for _, g := range groups {
					if g.HasCFCondition {
						hasCF = true
					}
					var links []model.PrerequisiteCourseLink
					for _, pc := range g.CourseCodes {
						coreq := courseRepo.byKey[courseKey(pc, curr.ID, year)]
						if coreq == nil {
							coreq = courseRepo.byCodeYear[codeYearKey(pc, year)]
						}
						if coreq == nil {
							continue
						}
						links = append(links, model.PrerequisiteCourseLink{
							ID:                   uuid.New(),
							PrerequisiteCourseID: coreq.ID,
							PrerequisiteCourse:   *coreq,
						})
					}
					course.CorequisiteGroups = append(course.CorequisiteGroups, model.PrerequisiteGroup{
						ID:                  uuid.New(),
						IsOrGroup:           g.IsOrGroup,
						HasCFCondition:      g.HasCFCondition,
						PrerequisiteCourses: links,
					})
				}
			}

			if hasCF {
				course.HasCFOption = true
			}
		}
	}
	return nil
}

// ─── Prerequisite string parser (mirrors courseService logic) ────────────────

type parsedGroup struct {
	CourseCodes    []string
	IsOrGroup      bool
	HasCFCondition bool
}

func parsePrereqString(s string) []parsedGroup {
	if s == "" {
		return nil
	}
	re := regexp.MustCompile(`(?i)\s+AND\s+`)
	s = re.ReplaceAllString(s, ",")
	cnf := parseToCNF(s)
	var groups []parsedGroup
	for _, codes := range cnf {
		g := parsedGroup{IsOrGroup: len(codes) > 1}
		for _, c := range codes {
			c = strings.TrimSpace(c)
			if strings.ToUpper(c) == "C.F." {
				g.HasCFCondition = true
			} else if c != "" {
				g.CourseCodes = append(g.CourseCodes, c)
			}
		}
		if len(g.CourseCodes) > 0 || g.HasCFCondition {
			groups = append(groups, g)
		}
	}
	return groups
}

func parseToCNF(str string) [][]string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	parts := splitOutside(str, ',')
	if len(parts) > 1 {
		var result [][]string
		for _, p := range parts {
			result = append(result, parseToCNF(p)...)
		}
		return result
	}
	orParts := splitByOR(str)
	if len(orParts) > 1 {
		var cnfs [][][]string
		for _, p := range orParts {
			cnfs = append(cnfs, parseToCNF(p))
		}
		return mergeCNFOr(cnfs)
	}
	if strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")") && isEnclosed(str) {
		return parseToCNF(str[1 : len(str)-1])
	}
	return [][]string{{str}}
}

func mergeCNFOr(cnfs [][][]string) [][]string {
	if len(cnfs) == 0 {
		return nil
	}
	result := cnfs[0]
	for i := 1; i < len(cnfs); i++ {
		var next [][]string
		for _, g1 := range result {
			for _, g2 := range cnfs[i] {
				merged := make([]string, len(g1)+len(g2))
				copy(merged, g1)
				copy(merged[len(g1):], g2)
				next = append(next, merged)
			}
		}
		result = next
	}
	return result
}

func splitOutside(str string, delim rune) []string {
	var parts []string
	var cur strings.Builder
	depth := 0
	for _, ch := range str {
		switch ch {
		case '(':
			depth++
			cur.WriteRune(ch)
		case ')':
			depth--
			cur.WriteRune(ch)
		case delim:
			if depth == 0 {
				parts = append(parts, cur.String())
				cur.Reset()
			} else {
				cur.WriteRune(ch)
			}
		default:
			cur.WriteRune(ch)
		}
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func splitByOR(str string) []string {
	var parts []string
	var cur strings.Builder
	depth := 0
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		switch ch {
		case '(':
			depth++
			cur.WriteRune(ch)
		case ')':
			depth--
			cur.WriteRune(ch)
		default:
			cur.WriteRune(ch)
			if depth == 0 {
				buf := cur.String()
				if len(buf) >= 4 && strings.EqualFold(buf[len(buf)-4:], " OR ") {
					parts = append(parts, buf[:len(buf)-4])
					cur.Reset()
				}
			}
		}
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func isEnclosed(str string) bool {
	if !strings.HasPrefix(str, "(") || !strings.HasSuffix(str, ")") {
		return false
	}
	depth := 0
	for i, ch := range str {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		}
		if depth == 0 && i < len(str)-1 {
			return false
		}
	}
	return depth == 0
}

// ─── In-memory repository implementations ───────────────────────────────────

// key helpers
func courseKey(code string, curriculumID uuid.UUID, year int) string {
	return fmt.Sprintf("%s|%s|%d", code, curriculumID, year)
}
func codeYearKey(code string, year int) string { return fmt.Sprintf("%s|%d", code, year) }
func catKey(curriculumID uuid.UUID, name string) string {
	return fmt.Sprintf("%s|%s", curriculumID, name)
}

// ── FakeCurriculumRepo ───────────────────────────────────────────────────────

type FakeCurriculumRepo struct {
	store map[string]*model.Curriculum // keyed by nameTH and nameEN
}

func newFakeCurriculumRepo() *FakeCurriculumRepo {
	return &FakeCurriculumRepo{store: make(map[string]*model.Curriculum)}
}

func (r *FakeCurriculumRepo) GetByID(id uuid.UUID) (*model.Curriculum, error) {
	for _, c := range r.store {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("curriculum %s not found", id)
}
func (r *FakeCurriculumRepo) GetByName(name string) (*model.Curriculum, error) {
	if c, ok := r.store[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("curriculum '%s' not found", name)
}
func (r *FakeCurriculumRepo) GetAll() ([]model.Curriculum, error) {
	var out []model.Curriculum
	seen := map[uuid.UUID]bool{}
	for _, c := range r.store {
		if !seen[c.ID] {
			out = append(out, *c)
			seen[c.ID] = true
		}
	}
	return out, nil
}
func (r *FakeCurriculumRepo) Create(c *model.Curriculum) error                    { return nil }
func (r *FakeCurriculumRepo) Update(c *model.Curriculum) error                    { return nil }
func (r *FakeCurriculumRepo) Delete(id uuid.UUID) error                           { return nil }
func (r *FakeCurriculumRepo) DeleteAll() error                                    { return nil }
func (r *FakeCurriculumRepo) Upsert(c *model.Curriculum) error                    { return nil }
func (r *FakeCurriculumRepo) BulkUpsert(cs []model.Curriculum) error              { return nil }
func (r *FakeCurriculumRepo) GetActiveByYear(year int) ([]model.Curriculum, error) { return nil, nil }
func (r *FakeCurriculumRepo) GetAllWithOutCatAndCourse() ([]model.Curriculum, error) {
	return r.GetAll()
}

// ── FakeCategoryRepo ─────────────────────────────────────────────────────────

type FakeCategoryRepo struct {
	byKey        map[string]*model.Category
	byCurriculum map[uuid.UUID][]model.Category
}

func newFakeCategoryRepo() *FakeCategoryRepo {
	return &FakeCategoryRepo{
		byKey:        make(map[string]*model.Category),
		byCurriculum: make(map[uuid.UUID][]model.Category),
	}
}

func (r *FakeCategoryRepo) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error) {
	return r.byCurriculum[curriculumID], nil
}
func (r *FakeCategoryRepo) GetByID(id uuid.UUID) (*model.Category, error) {
	for _, c := range r.byKey {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("category %s not found", id)
}
func (r *FakeCategoryRepo) GetByName(name string) (*model.Category, error)   { return nil, nil }
func (r *FakeCategoryRepo) GetAll() ([]model.Category, error)                { return nil, nil }
func (r *FakeCategoryRepo) Create(c *model.Category) error                   { return nil }
func (r *FakeCategoryRepo) Update(c *model.Category) error                   { return nil }
func (r *FakeCategoryRepo) Delete(id uuid.UUID) error                        { return nil }
func (r *FakeCategoryRepo) DeleteAll() error                                 { return nil }
func (r *FakeCategoryRepo) Upsert(c *model.Category) error                   { return nil }
func (r *FakeCategoryRepo) BulkUpsert(cats []model.Category) error           { return nil }

// ── FakeCourseRepo ───────────────────────────────────────────────────────────

type FakeCourseRepo struct {
	byKey      map[string]*model.Course // courseKey(code, currID, year)
	byCodeYear map[string]*model.Course // codeYearKey(code, year) — first match wins
}

func newFakeCourseRepo() *FakeCourseRepo {
	return &FakeCourseRepo{
		byKey:      make(map[string]*model.Course),
		byCodeYear: make(map[string]*model.Course),
	}
}

func (r *FakeCourseRepo) GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error) {
	if c, ok := r.byKey[courseKey(code, curriculumID, year)]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("course %s not found in curriculum %s year %d", code, curriculumID, year)
}

func (r *FakeCourseRepo) CatalogYearExists(curriculumID uuid.UUID, year int) bool {
	for k := range r.byKey {
		parts := strings.Split(k, "|")
		if len(parts) == 3 && parts[1] == curriculumID.String() && parts[2] == strconv.Itoa(year) {
			return true
		}
	}
	return false
}

func (r *FakeCourseRepo) GetLatestAvailableCatalogYear(curriculumID uuid.UUID, maxYear int) (int, bool) {
	best := 0
	for k := range r.byKey {
		parts := strings.Split(k, "|")
		if len(parts) == 3 && parts[1] == curriculumID.String() {
			y, _ := strconv.Atoi(parts[2])
			if y <= maxYear && y > best {
				best = y
			}
		}
	}
	if best == 0 {
		return 0, false
	}
	return best, true
}

func (r *FakeCourseRepo) CourseHasCFOptionInAnyCatalogYear(code string, curriculumID uuid.UUID) bool {
	for k, c := range r.byKey {
		parts := strings.Split(k, "|")
		if len(parts) == 3 && parts[0] == code && parts[1] == curriculumID.String() {
			if c.HasCFOption {
				return true
			}
		}
	}
	return false
}

func (r *FakeCourseRepo) ExistsByCodeAndCurriculumID(code string, curriculumID uuid.UUID) bool {
	for k := range r.byKey {
		parts := strings.Split(k, "|")
		if len(parts) == 3 && parts[0] == code && parts[1] == curriculumID.String() {
			return true
		}
	}
	return false
}

func (r *FakeCourseRepo) GetByCodeAndYear(code string, year int) (*model.Course, error) {
	if c, ok := r.byCodeYear[codeYearKey(code, year)]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("course %s year %d not found", code, year)
}

// stubs for unused interface methods
func (r *FakeCourseRepo) Create(c *model.Course) error                                      { return nil }
func (r *FakeCourseRepo) GetByID(id uuid.UUID) (*model.Course, error)                       { return nil, nil }
func (r *FakeCourseRepo) GetByName(n string) (*model.Course, error)                         { return nil, nil }
func (r *FakeCourseRepo) GetByCurriculumID(id uuid.UUID) ([]model.Course, error)            { return nil, nil }
func (r *FakeCourseRepo) GetByCategoryID(id uuid.UUID) ([]model.Course, error)              { return nil, nil }
func (r *FakeCourseRepo) GetAll() ([]model.Course, error)                                   { return nil, nil }
func (r *FakeCourseRepo) Update(c *model.Course) error                                      { return nil }
func (r *FakeCourseRepo) Delete(id uuid.UUID) error                                         { return nil }
func (r *FakeCourseRepo) CreateFromCSV(courses []model.Course) error                        { return nil }
func (r *FakeCourseRepo) Upsert(c *model.Course) error                                      { return nil }
func (r *FakeCourseRepo) BulkUpsert(courses []model.Course) error                           { return nil }
func (r *FakeCourseRepo) SetPrerequisites(id uuid.UUID, ids []uuid.UUID) error              { return nil }
func (r *FakeCourseRepo) SetCorequisites(id uuid.UUID, ids []uuid.UUID) error               { return nil }
func (r *FakeCourseRepo) SetPrerequisiteGroups(id uuid.UUID, g []model.PrerequisiteGroup) error { return nil }
func (r *FakeCourseRepo) SetCorequisiteGroups(id uuid.UUID, g []model.PrerequisiteGroup) error  { return nil }
func (r *FakeCourseRepo) DeleteAll() error                                                  { return nil }
func (r *FakeCourseRepo) DeleteByYear(year int) error                                       { return nil }

// ─── CSV utility ─────────────────────────────────────────────────────────────

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return csv.NewReader(f).ReadAll()
}
