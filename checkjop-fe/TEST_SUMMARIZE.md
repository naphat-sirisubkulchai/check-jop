# CheckJop — Test Plan Summary

## Unit Tests

### 1. App Store (`appStore.ts`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-ST-01 | `addCoursePlan` — add a new course to empty study plan | Course is added, study plan length = 1 |
| UT-ST-02 | `addCoursePlan` — add duplicate course (same code, year, semester) | Course is NOT added, returns existing plan |
| UT-ST-03 | `addCoursePlan` — add same course to different semester | Course is added (allowed) |
| UT-ST-04 | `removeCoursePlan` — remove existing course by code | Course removed from study plan |
| UT-ST-05 | `removeCoursePlan` — remove non-existent course code | Study plan unchanged, no error |
| UT-ST-06 | `editCoursePlan` — update grade of existing course | Grade updated correctly |
| UT-ST-07 | `editCoursePlan` — edit non-existent course code | Study plan unchanged |
| UT-ST-08 | `setStudyPlan` — set entire study plan array | Study plan replaced completely |
| UT-ST-09 | `clearStudyPlan` — clear study plan | Study plan becomes empty array |
| UT-ST-10 | `addExemption` — add new exemption code | Exemption added to list |
| UT-ST-11 | `addExemption` — add duplicate exemption code | No duplicate (uses Set) |
| UT-ST-12 | `removeExemption` — remove existing exemption | Exemption removed |
| UT-ST-13 | `clearExemptions` — clear all exemptions | Exemptions becomes empty |
| UT-ST-14 | `getCourseByCode` — find existing course | Returns correct course object |
| UT-ST-15 | `getCourseByCode` — find non-existent course | Returns undefined |
| UT-ST-16 | `setSelectedCurriculum` — set curriculum | Selected curriculum is updated |
| UT-ST-17 | `setSelectedCurriculum(null)` — clear selection | Selected curriculum becomes null |

### 2. Export/Import Utilities (`exportImport.ts`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-EI-01 | `exportStudyPlan` — export with valid data | Returns object with version, exportDate, studyPlan, metadata |
| UT-EI-02 | `exportStudyPlan` — metadata totalCredits calculation | Sum of all course credits is correct |
| UT-EI-03 | `exportStudyPlan` — metadata coursesCount | Equals studyPlan.length |
| UT-EI-04 | `exportStudyPlan` — with null curriculum | selectedCurriculum is null in output |
| UT-EI-05 | `validateImportData` — valid complete data | `isValid: true`, no errors |
| UT-EI-06 | `validateImportData` — null input | `isValid: false`, error "Invalid data format" |
| UT-EI-07 | `validateImportData` — missing version field | Error includes "Missing version field" |
| UT-EI-08 | `validateImportData` — studyPlan is not array | Error includes "studyPlan must be an array" |
| UT-EI-09 | `validateImportData` — plan entry missing course_code | Error includes "missing course_code" |
| UT-EI-10 | `validateImportData` — plan entry with non-number academicYear | Error includes "academicYear must be a number" |
| UT-EI-11 | `validateImportData` — plan entry with non-number semester | Error includes "semester must be a number" |
| UT-EI-12 | `validateImportData` — plan entry with non-number credits | Error includes "credits must be a number" |
| UT-EI-13 | `validateImportData` — exemptions is not array | Error includes "exemptions must be an array" |
| UT-EI-14 | `validateImportData` — selectedCurriculum is null | Warning "No curriculum selected" |
| UT-EI-15 | `validateImportData` — version mismatch (e.g. "2.0.0") | Warning about version compatibility |
| UT-EI-16 | `parseImportFile` — valid JSON string | Returns parsed data with `isValid: true` |
| UT-EI-17 | `parseImportFile` — invalid JSON string | Returns `isValid: false`, error "Failed to parse JSON" |
| UT-EI-18 | `saveStudyPlanToCookie` — data within 4KB limit | Returns true |
| UT-EI-19 | `saveStudyPlanToCookie` — data exceeds 4KB limit | Returns false |

### 3. Graph Utilities (`graph/utils/index.ts`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-GR-01 | `parseRelations` — empty string | Returns empty array |
| UT-GR-02 | `parseRelations` — single course code "2110101" | Returns `["2110101"]` |
| UT-GR-03 | `parseRelations` — comma-separated "2110101,2110102" | Returns `["2110101","2110102"]` |
| UT-GR-04 | `parseRelations` — OR relation "2110101 OR 2110102" | Returns `["2110101","2110102"]` |
| UT-GR-05 | `parseRelations` — with parentheses "(2110101,2110102)" | Returns codes without parentheses |
| UT-GR-06 | `parseRelations` — mixed OR and comma | Returns all unique codes |
| UT-GR-07 | `parseRelations` — duplicate codes | Returns deduplicated array |
| UT-GR-08 | `parseRelations` — numeric input (number type) | Converts to string and parses |
| UT-GR-09 | `generateGraphLayout` — nodes get correct positions | All nodes have x, y positions assigned |
| UT-GR-10 | `generateGraphLayout` — edges preserved after layout | Same number of edges in output |

### 4. General Utilities (`utils/index.ts`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-UT-01 | `uniqueCoursesByCode` — list with duplicate codes | Returns deduplicated list (last occurrence wins) |
| UT-UT-02 | `uniqueCoursesByCode` — empty array | Returns empty array |
| UT-UT-03 | `uniqueCoursesByCode` — all unique courses | Returns same list |

### 5. API Transform Functions (`courseApi.ts`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-AP-01 | `transformCurriculum` — snake_case to camelCase | All fields mapped correctly (name_th → nameTH, etc.) |
| UT-AP-02 | `transformCurriculum` — with categories array | Categories transformed with correct fields |
| UT-AP-03 | `transformCurriculum` — with courses array | Courses transformed with correct fields |
| UT-AP-04 | `transformCurriculum` — missing categories/courses | Returns empty arrays |
| UT-AP-05 | `getAllCourse` response mapping — handles name_en/nameEN variants | Maps all field name variants correctly |

### 6. CSV Adapter Functions (`admin/curriculum/create/page.tsx`)

| ID | Test Case | Expected Result |
|---|---|---|
| UT-CS-01 | `adaptCurriculumData` — valid CSV row | Returns object with nameTH, nameEN, year, minTotalCredits, isActive |
| UT-CS-02 | `adaptCurriculumData` — missing fields | Defaults to empty string / 0 |
| UT-CS-03 | `adaptCategoryData` — valid CSV row | Returns object with nameTH, nameEN, curriculumName, minCredits |
| UT-CS-04 | `adaptCategoryData` — handles typo field "cirrculumName" | Falls back correctly |
| UT-CS-05 | `adaptCourseData` — valid CSV row | Returns object with code, nameTH, nameEN, credits, prerequisites, corequisites |
| UT-CS-06 | `adaptCourseData` — missing optional fields | Defaults to empty strings |

### 7. Hook: `useCourseGraph`

| ID | Test Case | Expected Result |
|---|---|---|
| UT-HK-01 | `createNode` — creates node with correct position | Node has x, y based on index grid |
| UT-HK-02 | `createEdges` — course with prerequisites | Creates red-colored prerequisite edges |
| UT-HK-03 | `createEdges` — course with corequisites | Creates green-colored corequisite edges |
| UT-HK-04 | `createEdges` — prerequisite course not in dataset | Edge is NOT created (skip missing) |
| UT-HK-05 | `loadCourseGraph` — null courses input | Sets nodes/edges to empty arrays |

### 8. Hook: `useStudyPlan`

| ID | Test Case | Expected Result |
|---|---|---|
| UT-HK-06 | `handleCalculate` — no curriculum selected | Sets error "Please select a curriculum" |
| UT-HK-07 | `handleCalculate` — empty study plan | Sets error "Please add courses to your study plan" |
| UT-HK-08 | `handleCalculate` — calculates correct admission year | Uses Math.min of all academicYear values |
| UT-HK-09 | `handleCalculate` — transforms payload correctly | Maps academicYear → year, includes exemptions |
| UT-HK-10 | `handleCalculate` — API error | Sets error message, isLoading set to false |

---

## Integration Tests

### 1. Study Plan Workflow (Store + API)

| ID | Test Case | Expected Result |
|---|---|---|
| IT-SP-01 | Select curriculum → load courses → courses appear in store | Courses array populated with correct data |
| IT-SP-02 | Select curriculum → add course to plan → verify plan updated | Study plan contains the added course |
| IT-SP-03 | Add course → set grade → verify grade persisted | Course in plan has correct grade |
| IT-SP-04 | Add course → remove course → verify plan updated | Course no longer in study plan |
| IT-SP-05 | Add multiple courses → clear plan → verify empty | Study plan is empty |
| IT-SP-06 | Add courses → calculate graduation → result stored | Result object set in store |

### 2. Export/Import Flow

| ID | Test Case | Expected Result |
|---|---|---|
| IT-EX-01 | Add courses → export → import same file → verify data match | Imported study plan matches original |
| IT-EX-02 | Export → modify JSON → import invalid file | Shows validation errors |
| IT-EX-03 | Import file with different curriculum → courses reload | Courses fetched for imported curriculum |
| IT-EX-04 | Import file with warnings → display warnings to user | Warnings shown in UI status |

### 3. Cookie Persistence

| ID | Test Case | Expected Result |
|---|---|---|
| IT-CK-01 | Add courses → save to cookie → reload → load from cookie | Study plan restored correctly |
| IT-CK-02 | Save to cookie → clear cookie → load returns false | No data found |
| IT-CK-03 | Save large study plan (>4KB) to cookie | Returns false, data not saved |

### 4. Curriculum Management (Admin)

| ID | Test Case | Expected Result |
|---|---|---|
| IT-AD-01 | Load admin page → fetch all curricula → display in table | Table shows all curricula |
| IT-AD-02 | Search curricula by name | Filtered results shown |
| IT-AD-03 | Filter curricula by year | Only matching year shown |
| IT-AD-04 | Delete curriculum → confirm → removed from list | Curriculum removed from table |
| IT-AD-05 | Upload 3 CSV files → preview parsed data → submit | Curriculum created, redirect to admin |
| IT-AD-06 | Upload incomplete files (missing 1 CSV) → submit | Submit button disabled |

### 5. Graduation Check (End-to-End)

| ID | Test Case | Expected Result |
|---|---|---|
| IT-GC-01 | Select curriculum → add all required courses → calculate | `canGraduate: true` |
| IT-GC-02 | Select curriculum → add insufficient courses → calculate | `canGraduate: false`, shows missing courses |
| IT-GC-03 | Add course with missing prerequisites → calculate | Shows prerequisite violations |
| IT-GC-04 | Add courses missing category credits → calculate | Shows unsatisfied categories |
| IT-GC-05 | Calculate → navigate to /calculate → result displayed | Result page shows graduation result |

### 6. Course Graph Visualization

| ID | Test Case | Expected Result |
|---|---|---|
| IT-CG-01 | Load courses with prerequisites → render graph | Nodes and prerequisite edges displayed |
| IT-CG-02 | Load courses with corequisites → render graph | Corequisite edges (dashed green) displayed |
| IT-CG-03 | Switch curriculum → graph updates | New nodes/edges for new curriculum |

### 7. Course List & Search

| ID | Test Case | Expected Result |
|---|---|---|
| IT-CL-01 | Select curriculum → courses load → search by code | Filtered courses match query |
| IT-CL-02 | Search by English name | Matching courses shown |
| IT-CL-03 | Search with no results | "No courses found" message shown |
| IT-CL-04 | Switch curriculum → course list reloads | Courses change to new curriculum |

### 8. Semester Card Interaction

| ID | Test Case | Expected Result |
|---|---|---|
| IT-SC-01 | Drag course from list → drop on semester card | Course added to that semester |
| IT-SC-02 | Use QuickAdd combobox → select course → click Add | Course added to semester |
| IT-SC-03 | Add courses → verify credit total shown correctly | Credits sum matches courses |
| IT-SC-04 | Summer semester (sem 3) max credits = 12 | Max credits display shows 12 |

---

## Summary

| Type | Count |
|---|---|
| **Unit Tests** | 48 |
| **Integration Tests** | 28 |
| **Total** | **76** |
