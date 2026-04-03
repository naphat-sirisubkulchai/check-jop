# Test Suite Documentation - CheckJop Backend

เอกสารนี้สรุปรายละเอียดของ test cases ทั้งหมดในโฟลเดอร์ `tests/` โดยแบ่งตามไฟล์และหน้าที่การทำงาน

---

## 📁 ไฟล์ทดสอบ

### 1. `mocks.go`
**คำอธิบาย**: Mock implementations สำหรับ Repository interfaces ใช้ในการทดสอบโดยไม่ต้องเชื่อมต่อกับฐานข้อมูลจริง

**Mock Repositories**:
- `MockCurriculumRepository` - จัดการข้อมูล Curriculum
- `MockCourseRepository` - จัดการข้อมูล Course และ Prerequisites/Corequisites
- `MockCategoryRepository` - จัดการข้อมูล Category

**Helper Functions**:
- `setupGraduationService()` - สร้าง GraduationService พร้อม mock repositories สำหรับใช้ในการทดสอบ

---

## 🧪 Test Cases

### 2. `course_year_test.go`
**คำอธิบาย**: ทดสอบ Course Versioning - การจัดการ prerequisites ที่แตกต่างกันตามปีของ course

#### Test Cases:
**TestCourseYearVersioning_DifferentVersions**
- **Input**:
  - นักศึกษาสมัครปี 2023 ลงวิชา 2301170 (version 2023 ไม่มี prerequisites)
  - นักศึกษาสมัครปี 2024 ลงวิชา 2301170 (version 2024 ต้องการ 2301172)
- **Expected Output**:
  - นักศึกษา 2023: ไม่มี violation ✅
  - นักศึกษา 2024: มี violation ระบุขาด 2301172 ❌
- **จุดประสงค์**: ยืนยันว่าระบบจัดการ course requirements ที่เปลี่ยนแปลงตามปีการศึกษาได้อย่างถูกต้อง

---

### 3. `free_electives_test.go`
**คำอธิบาย**: ทดสอบการจัดการวิชาเลือกเสรี (Free Electives) ทั้งที่มีและไม่มีในฐานข้อมูล

#### Test Cases:

**TestFreeElective_WithPrerequisites_ShouldCheck**
- **Input**: นักศึกษาลงวิชา 2301999 (Free Elective ที่มี prerequisite 2301170) โดยไม่มี 2301170
- **Expected Output**: violation ระบุขาด prerequisite 2301170 ❌
- **จุดประสงค์**: แม้เป็น Free Elective แต่ถ้ามี prerequisites ต้องตรวจสอบ

**TestFreeElective_NotInDB_ShouldSkip**
- **Input**: นักศึกษาลงวิชา 9999999 (Free Elective ไม่มีในฐานข้อมูล)
- **Expected Output**: ไม่มี violation (ข้ามการตรวจสอบ) ✅
- **จุดประสงค์**: วิชา Free Elective ที่ไม่มีในระบบไม่ถูกตรวจสอบ prerequisites

**TestFreeElective_NotInDB_CountsCredits**
- **Input**: นักศึกษาลงวิชา 9999999 (Free Elective ไม่มีในฐานข้อมูล) 3 หน่วยกิต
- **Expected Output**:
  - CanGraduate = true ✅
  - TotalCredits = 3
  - Category "Free Electives" IsSatisfied = true
- **จุดประสงค์**: วิชา Free Elective ที่ไม่มีในระบบยังคงนับหน่วยกิตได้ปกติ

---

### 4. `graduation_service_cf_test.go`
**คำอธิบาย**: ทดสอบ C.F. (Consent of Faculty) Conditions - การอนุญาตพิเศษให้ลงวิชาโดยไม่ต้องมี prerequisites

#### Test Cases:

**TestValidatePrerequisites_CFCondition_NoPermission**
- **Input**: ลงวิชา 2301199 (ต้องการ "2301170 OR C.F.") โดยไม่มี 2301170 และไม่มี C.F. permission
- **Expected Output**: violation ระบุขาด 2301170 ❌
- **จุดประสงค์**: กรณีไม่มี prerequisite และไม่มี C.F. permission

**TestValidatePrerequisites_CFCondition_WithPermission**
- **Input**: ลงวิชา 2301199 โดยไม่มี 2301170 แต่มี C.F. permission (ใน Exemptions)
- **Expected Output**: ไม่มี violation ✅
- **จุดประสงค์**: C.F. permission สามารถแทน prerequisite ได้

**TestValidatePrerequisites_CFCondition_WithPrerequisite**
- **Input**: ลงวิชา 2301199 โดยมี 2301170 แล้ว (ไม่ต้องการ C.F.)
- **Expected Output**: ไม่มี violation ✅
- **จุดประสงค์**: ถ้ามี prerequisite ครบไม่ต้องใช้ C.F. permission

---

### 5. `graduation_service_gpax_test.go`
**คำอธิบาย**: ทดสอบการคำนวณ GPAX และการจัดการเกรด F

#### GPAX Calculation Tests:

**TestCalculateGPAX_AllGraded**
- **Input**: วิชา A (3 หน่วย, เกรด A=4.0), B (3 หน่วย, เกรด B=3.0), C (2 หน่วย, เกรด C+=2.5)
- **Expected Output**: GPAX = 3.25 (26 คะแนน / 8 หน่วย)
- **จุดประสงค์**: คำนวณ GPAX จากวิชาที่มีเกรดตัวอักษรทั้งหมด

**TestCalculateGPAX_WithNonGraded**
- **Input**: วิชา A (3 หน่วย, เกรด A=4.0), B (3 หน่วย, เกรด W), C (3 หน่วย, เกรด S)
- **Expected Output**: GPAX = 4.0 (12 คะแนน / 3 หน่วย)
- **จุดประสงค์**: เกรด W และ S ไม่นับรวมในการคำนวณ GPAX

**TestCalculateGPAX_WithF**
- **Input**: วิชา A (3 หน่วย, เกรด A=4.0), B (3 หน่วย, เกรด F=0.0)
- **Expected Output**: GPAX = 2.0 (12 คะแนน / 6 หน่วย)
- **จุดประสงค์**: เกรด F นับเป็น 0.0 และนับรวมในการคำนวณ GPAX

#### F Grade in Prerequisites Tests:

**TestValidatePrerequisites_FGrade_Violation**
- **Input**: ลงวิชา PRE (เกรด F) แล้วลงวิชา MAIN (ต้องการ prerequisite PRE)
- **Expected Output**: violation ระบุขาด PRE ❌
- **จุดประสงค์**: เกรด F ไม่ถือว่าผ่าน prerequisite

**TestValidatePrerequisites_Corequisite_FGrade_Allowed**
- **Input**: ลงวิชา CO (เกรด F) และ MAIN (ต้องการ corequisite CO) ในเทอมเดียวกัน
- **Expected Output**: ไม่มี violation ✅
- **จุดประสงค์**: เกรด F ใน corequisite ไม่ถือว่าผิดเงื่อนไข (ลงพร้อมกัน)

---

### 6. `graduation_service_simple_test.go`
**คำอธิบาย**: ทดสอบพื้นฐานของ prerequisite และ corequisite validation

#### Test Cases:

**TestValidatePrerequisites_Simple_NoPrerequisites**
- **Input**: ลงวิชา 2301170 (ไม่มี prerequisites)
- **Expected Output**: ไม่มี violation ✅

**TestValidatePrerequisites_Simple_MissingPrerequisite**
- **Input**: ลงวิชา 2301180 (ต้องการ 2301170 OR 2301173) โดยไม่มีทั้งสองวิชา
- **Expected Output**: violation ระบุขาด 2301170 และ 2301173 ❌

**TestValidatePrerequisites_Simple_PrerequisiteSatisfied**
- **Input**: ลง 2301170 (เทอม 1) แล้วลง 2301180 (เทอม 2)
- **Expected Output**: ไม่มี violation ✅

**TestValidatePrerequisites_Simple_StrictTermRequirement**
- **Input**: ลง 2301220 และ 2301230 (ต้องการ 2301220) ในเทอมเดียวกัน
- **Expected Output**: violation - 2301230 มี prerequisite ถูกลงในเทอมเดียวกัน ❌

**TestValidatePrerequisites_Simple_TransitiveChain**
- **Input**: ลง 2301365 (ต้องการ 2301263 ต้องการ 2301260) โดยไม่มีทั้งสองวิชา
- **Expected Output**: violation ระบุขาด 2301263 และ 2301260 ❌

**TestValidatePrerequisites_Simple_Corequisites_WrongTerm**
- **Input**: ลง 2301172 (ต้องการ corequisite 2301170) คนละเทอม
- **Expected Output**: violation - corequisite ถูกลงผิดเทอม ❌

**TestValidatePrerequisites_Simple_Corequisites_SameTerm**
- **Input**: ลง 2301172 และ 2301170 (corequisite) ในเทอมเดียวกัน
- **Expected Output**: ไม่มี violation ✅

**TestValidateCreditLimits_Simple_RegularSemester**
- **Input**: ลงวิชารวม 23 หน่วยกิตในเทอมปกติ
- **Expected Output**: violation - เกิน 22 หน่วยกิต ❌

**TestValidateCreditLimits_Simple_SummerSemester**
- **Input**: ลงวิชารวม 12 หน่วยกิตในเทอมฤดูร้อน
- **Expected Output**: violation - เกิน 10 หน่วยกิต ❌

**TestValidatePrerequisites_Simple_CourseNotFound**
- **Input**: ลงวิชา INVALID ที่ไม่มีในระบบ
- **Expected Output**: ไม่มี violation (ข้ามวิชานี้) ✅

**TestValidatePrerequisites_Simple_PartialTransitiveChain**
- **Input**: ลง 2301365 โดยมี 2301260 แต่ไม่มี 2301263
- **Expected Output**: violation ระบุขาด 2301263 ❌

---

### 7. `graduation_service_test.go`
**คำอธิบาย**: Test suite ที่ครอบคลุมที่สุดสำหรับ GraduationService ครอบคลุมทุก edge cases

#### Basic Prerequisite Tests:

**TestValidatePrerequisites_NoPrerequisites**
- ทดสอบวิชาที่ไม่มี prerequisites (2301170)
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_BasicPrerequisite_2301180**
- ทดสอบ OR prerequisite (2301170 OR 2301173) กรณีขาด
- **Expected**: violation ระบุขาดทั้งสอง ❌

**TestValidatePrerequisites_PrerequisiteSatisfied_2301180**
- ทดสอบ OR prerequisite กรณีผ่าน
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_StrictTermRequirement_2301230**
- ทดสอบการลง prerequisite ในเทอมเดียวกัน (ต้องลงก่อน)
- **Expected**: violation - PrereqsTakenInWrongTerm ❌

#### Transitive Prerequisite Tests:

**TestValidatePrerequisites_TransitiveChain_DataStructures**
- ทดสอบ transitive chain: 2301365 → 2301263 → 2301260
- **Expected**: violation ระบุขาดทั้ง 2301263 และ 2301260 ❌

**TestValidatePrerequisites_ComplexPrerequisites_2301367**
- ทดสอบ partial chain (มี base แต่ขาด intermediate)
- **Expected**: violation ระบุขาด 2301375 ❌

#### Credit Limit Tests:

**TestValidateCreditLimits_RegularSemester**
- ทดสอบเทอมปกติ (23 หน่วยกิต > 22)
- **Expected**: violation ❌

**TestValidateCreditLimits_SummerSemester**
- ทดสอบเทอมฤดูร้อน (12 หน่วยกิต > 10)
- **Expected**: violation ❌

#### Corequisite Tests (Basic):

**TestValidatePrerequisites_Corequisites_2301172**
- ทดสอบ corequisite คนละเทอม
- **Expected**: violation - CoreqsTakenInWrongTerm ❌

**TestValidatePrerequisites_Corequisites_SameTerm_Valid**
- ทดสอบ corequisite เทอมเดียวกัน
- **Expected**: ไม่มี violation ✅

#### Category Requirements Tests:

**TestCheckCategoryRequirements**
- ทดสอบการนับหน่วยกิตตาม category
- **Expected**: แสดง earned credits, required credits, และ IsSatisfied status

**TestCheckGraduation_Complete**
- ทดสอบเงื่อนไขจบครบทุกอย่าง
- **Expected**: CanGraduate = true ✅

#### Error Handling Tests:

**TestValidatePrerequisites_CourseNotFound**
- ทดสอบวิชาที่ไม่มีในระบบ
- **Expected**: ไม่มี violation (ข้ามวิชา) ✅

**TestValidatePrerequisites_MultipleViolationTypes**
- ทดสอบหลาย violation types พร้อมกัน
- **Expected**: แสดง violations ทุกประเภท ❌

#### Comprehensive Corequisite Tests:

**TestValidatePrerequisites_Corequisites_MissingCorequisite**
- ทดสอบการขาด corequisite
- **Expected**: violation - MissingCoreqs ❌

**TestValidatePrerequisites_Corequisites_OrGroup_OneSatisfied**
- ทดสอบ OR corequisite (มีตัวใดตัวหนึ่ง)
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_Corequisites_OrGroup_NoneSatisfied**
- ทดสอบ OR corequisite (ไม่มีเลย)
- **Expected**: violation ระบุขาดทุกตัว ❌

**TestValidatePrerequisites_Corequisites_OrGroup_WrongTerm**
- ทดสอบ OR corequisite (มีแต่คนละเทอม)
- **Expected**: violation - CoreqsTakenInWrongTerm ❌

**TestValidatePrerequisites_Corequisites_AndGroup_AllSatisfied**
- ทดสอบ AND corequisite (ครบทั้งหมดเทอมเดียวกัน)
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_Corequisites_AndGroup_OneMissing**
- ทดสอบ AND corequisite (ขาดตัวใดตัวหนึ่ง)
- **Expected**: violation ระบุขาดตัวที่หาย ❌

**TestValidatePrerequisites_Corequisites_TransitiveCorequisites**
- ทดสอบ transitive corequisite chain (ครบทั้งหมด)
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_Corequisites_TransitiveMissing**
- ทดสอบ transitive corequisite chain (ขาดตัวใดตัวหนึ่ง)
- **Expected**: violations ทุกวิชาที่ไม่ครบ ❌

**TestValidatePrerequisites_Corequisites_MixedPrereqAndCoreq**
- ทดสอบวิชาที่มีทั้ง prerequisite และ corequisite
- **Expected**: ไม่มี violation เมื่อทำถูกต้อง ✅

#### Complex Prerequisite Tests:

**TestValidatePrerequisites_ComplexOrGroups_BothGroupsSatisfied**
- ทดสอบ (2301265 OR 2301274) AND (2301279 OR 2301369) - ครบทุก group
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_ComplexOrGroups_OnlyFirstGroupSatisfied**
- ทดสอบหลาย OR groups - ครบแค่ group แรก
- **Expected**: violation ระบุขาด group ที่สอง ❌

**TestValidatePrerequisites_ComplexOrGroups_NoGroupsSatisfied**
- ทดสอบหลาย OR groups - ไม่ครบเลย
- **Expected**: violation ระบุขาดทุกตัวจากทุก group ❌

**TestValidatePrerequisites_MixedOrAndRequirements_BothSatisfied**
- ทดสอบ OR group และ single requirement - ครบทั้งหมด
- **Expected**: ไม่มี violation ✅

**TestValidatePrerequisites_MixedOrAndRequirements_OrGroupMissing**
- ทดสอบ mixed requirements - ขาด OR group
- **Expected**: violation ระบุขาด OR group ❌

---

### 8. `cross_curriculum_test.go`
**คำอธิบาย**: ทดสอบการตรวจสอบ prerequisite ข้าม curriculum

#### Test Cases:

**TestCrossCurriculum_PrerequisiteCheck**
- **Input**:
  - Course 2301170 ถูก define ใน Curriculum A (Year 2023) มี prerequisite 2301172
  - นักศึกษาอยู่ใน Curriculum B ลงวิชา 2301170
- **Expected Output**: violation ระบุขาด 2301172 ❌
- **จุดประสงค์**: แม้นักศึกษาอยู่คนละ curriculum แต่ prerequisite ของ course ยังต้องตรวจสอบ

---

### 9. `csv_import_cross_curriculum_test.go`
**คำอธิบาย**: ทดสอบการ import CSV ที่มี cross-curriculum prerequisites

#### Test Cases:

**TestImportFromCSV_CrossCurriculumPrerequisite**
- **Input**:
  ```csv
  code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
  101,C1,C1,3,,,Cat1,Curr1,2023
  102,C2,C2,3,101,,Cat1,Curr2,2023
  ```
- **Expected Output**: import สำเร็จ (ระบบควรใช้ fallback เพื่อค้นหา prerequisite ข้าม curriculum) ✅
- **จุดประสงค์**: ทดสอบการจัดการ prerequisite ที่อยู่คนละ curriculum ในการ import CSV

---

### 10. `prerequisite_complex_test.go`
**คำอธิบาย**: ทดสอบการ parse complex prerequisite expressions จาก CSV

#### Test Cases:

**TestImportFromCSV_ComplexPrerequisites**
- **Input**:
  ```csv
  code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
  101,C1,C1,3,,,Cat1,Curr1,2023
  102,C2,C2,3,,,Cat1,Curr1,2023
  103,C3,C3,3,,,Cat1,Curr1,2023
  200,Target,Target,3,(101 AND 102) OR 103,,Cat1,Curr1,2023
  ```
- **Expected Output**:
  - Course 200 มี 2 prerequisite groups:
    - Group 1 (OR): [101, 103]
    - Group 2 (OR): [102, 103]
  - Import สำเร็จโดย convert "(101 AND 102) OR 103" เป็น CNF ได้ถูกต้อง ✅
- **จุดประสงค์**: ทดสอบการ parse และ convert boolean expressions เป็น prerequisite groups

---

### 11. `duplicate_key_debug_test.go`
**คำอธิบาย**: ทดสอบการจัดการ duplicate entries ใน CSV import

#### Test Cases:

**TestImportFromCSV_DuplicateKeyDebug**
- **Input**:
  ```csv
  code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
  101,C1,C1,3,,,Cat1,Curr1,2023
  101,C1,C1,3,,,Cat1,Curr1,2023
  ```
- **Expected Output**: BulkUpsert ถูกเรียกด้วย slice length = 1 (deduped) ✅
- **จุดประสงค์**: ระบบควร deduplicate entries ที่ซ้ำกันก่อน bulk insert

---

### 12. `reset_database_test.go`
**คำอธิบาย**: ทดสอบการ reset ฐานข้อมูล

#### Test Cases:

**TestResetDatabase**
- **Expected Output**:
  - เรียก DeleteAll ตามลำดับ: Course → Category → Curriculum
  - ทุก DeleteAll สำเร็จ ✅
- **จุดประสงค์**: ทดสอบการลบข้อมูลทั้งหมดในลำดับที่ถูกต้อง (reverse dependency order)

---

## 📊 สรุป Test Coverage

### ✅ Features ที่ครอบคลุม:

1. **Prerequisite Validation**
   - Basic prerequisites (AND, OR groups)
   - Transitive prerequisite chains
   - Cross-curriculum prerequisites
   - Complex boolean expressions
   - Strict term requirements (ต้องลงก่อน)
   - F grade handling (ไม่นับเป็นผ่าน prerequisite)

2. **Corequisite Validation**
   - Basic corequisites (ต้องลงเทอมเดียวกัน)
   - OR group corequisites
   - AND group corequisites
   - Transitive corequisite chains
   - Mixed prerequisites และ corequisites
   - F grade handling (ใน corequisite ไม่ถือว่าผิด)

3. **GPAX Calculation**
   - All graded courses
   - Non-graded courses (W, S)
   - F grade handling (นับเป็น 0.0)

4. **Credit Limits**
   - Regular semester limit (22 หน่วยกิต)
   - Summer semester limit (10 หน่วยกิต)

5. **Category Requirements**
   - Credit counting per category
   - Satisfaction checking
   - Manual credits support

6. **Free Electives**
   - Prerequisites checking (ถ้ามีในระบบ)
   - Credit counting (แม้ไม่มีในระบบ)
   - Skipping prerequisite check (ถ้าไม่มีในระบบ)

7. **C.F. (Consent of Faculty) Conditions**
   - Permission-based prerequisite bypass
   - Exemptions tracking

8. **Course Versioning**
   - Different prerequisites by year
   - Admission year tracking

9. **CSV Import**
   - Basic import
   - Cross-curriculum handling
   - Complex prerequisite parsing
   - Duplicate entry deduplication

10. **Database Operations**
    - Reset database in correct order

---

## 🎯 Test Statistics

- **Total Test Files**: 12
- **Total Test Functions**: 70+
- **Coverage Areas**:
  - Prerequisite Validation ✅
  - Corequisite Validation ✅
  - GPAX Calculation ✅
  - Credit Limits ✅
  - Category Requirements ✅
  - Free Electives ✅
  - C.F. Conditions ✅
  - Course Versioning ✅
  - CSV Import ✅
  - Error Handling ✅

---

## 📝 หมายเหตุ

### Prerequisite vs Corequisite
- **Prerequisite**: วิชาที่ต้องเรียนและผ่านก่อน (ต้องลงคนละเทอม)
- **Corequisite**: วิชาที่ต้องลงพร้อมกัน (ต้องลงในเทอมเดียวกัน)

### เกรดที่มีผลต่อการตรวจสอบ
- **เกรดปกติ (A, B+, B, C+, C, D+, D, F)**: นับเป็น prerequisite (ยกเว้น F)
- **เกรด F**: ไม่นับเป็นการผ่าน prerequisite
- **เกรด W (Withdraw)**: ไม่นับรวมใน GPAX
- **เกรด S (Satisfactory)**: ไม่นับรวมใน GPAX

### Transitive Chains
- ระบบตรวจสอบ transitive prerequisites และ corequisites อัตโนมัติ
- เช่น ถ้า A ต้องการ B และ B ต้องการ C ดังนั้น A ต้องการ C โดยอัตโนมัติ

---

## 🚀 การรัน Tests

```bash
# รัน test ทั้งหมด
go test ./tests/...

# รัน test เฉพาะไฟล์
go test ./tests/graduation_service_test.go ./tests/mocks.go

# รัน test พร้อม coverage
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

**Last Updated**: 2025-12-03
**Total Test Cases**: 70+
**Test Status**: ✅ All Passing
