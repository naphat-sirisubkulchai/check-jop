# คู่มือการใช้ API สำหรับ CheckJop Backend

## CSV Import API

### 1. Import Course CSV with Year (Version 4)

สำหรับ CSV ที่ไม่มีคอลัมน์ Year (รูปแบบ version 4) ให้ใช้ endpoint นี้พร้อมระบุปีการศึกษา

#### Endpoint
```
POST /api/v1/import/course-csv-with-year
```

#### Request Format
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file): CSV file to import
  - `year` (int): ปีการศึกษา (เช่น 2566, 2567, 2568)

#### CSV Format (Version 4)
```csv
code,courseNameEN,courseNameTH,credit,prerequisites,corequisites,category,curriculum
2301170,Computer and Programming,คอมพิวเตอร์และการโปรแกรม,3,,,Major Courses,เอกเดี่ยว-2561
2301173,Programming,การเขียนโปรแกรม,4,,,Major Courses,"เอกเดี่ยว-ฝึกงาน-2566,เอกเดี่ยว-สหกิจ-2566"
```

**หมายเหตุ**: ไฟล์ CSV รูปแบบนี้ไม่มีคอลัมน์ `Year` - ระบบจะใช้ค่า `year` จาก parameter แทน

#### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/import/course-csv-with-year" \
  -F "file=@course_Present_2566.csv" \
  -F "year=2566"
```

#### Success Response
```json
{
  "message": "Course CSV imported successfully for year 2566"
}
```

#### Error Response
```json
{
  "error": "Year parameter is required"
}
```

```json
{
  "error": "found 3 error(s) during course relationship setup:\n- row 5: prerequisite course 'CS101' (year 2566) not found in curriculum 'Computer Science' for course 'CS201'\n- row 8: corequisite course 'MATH102' (year 2566) not found in curriculum 'Computer Science' for course 'CS202'\n- row 12: course 'CS301' (year 2566) not found in curriculum 'Engineering'"
}
```

#### พฤติกรรมของ API
1. **Hard Delete by Year**: ลบ courses ทั้งหมดของปีที่ระบุออกก่อน (hard delete รวมถึง relationships)
2. **Import New Data**: นำเข้า courses ใหม่พร้อมเติม year ให้ทุก record
3. **Error Collection**: รวบรวม errors ทั้งหมดเกี่ยวกับ prerequisites/corequisites ที่หาไม่เจอและแสดงพร้อมกัน

#### ตัวอย่างการใช้งาน

##### Import ข้อมูลปี 2566
```bash
curl -X POST "http://localhost:8080/api/v1/import/course-csv-with-year" \
  -F "file=@csv_data/version4/example_checkjop - course_Present_2566.csv" \
  -F "year=2566"
```

##### Import ข้อมูลปี 2567
```bash
curl -X POST "http://localhost:8080/api/v1/import/course-csv-with-year" \
  -F "file=@csv_data/version4/example_checkjop - course_Present_2567.csv" \
  -F "year=2567"
```

##### Import ข้อมูลปี 2568
```bash
curl -X POST "http://localhost:8080/api/v1/import/course-csv-with-year" \
  -F "file=@csv_data/version4/example_checkjop - course_Present_2568.csv" \
  -F "year=2568"
```

---

### 2. Import Course CSV (Version 3)

สำหรับ CSV ที่มีคอลัมน์ Year อยู่แล้ว (รูปแบบ version 3)

#### Endpoint
```
POST /api/v1/import/course-csv
```

#### CSV Format (Version 3)
```csv
code,courseNameEN,courseNameTH,credit,prerequisites,corequisites,category,curriculum,Year
2301170,Computer and Programming,คอมพิวเตอร์และการโปรแกรม,3,,,Major Courses,เอกเดี่ยว-2561,2568
2301173,Programming,การเขียนโปรแกรม,4,,,Major Courses,"เอกเดี่ยว-ฝึกงาน-2566,เอกเดี่ยว-สหกิจ-2566",2568
```

**หมายเหตุ**: ไฟล์ CSV รูปแบบนี้มีคอลัมน์ `Year` ในตำแหน่งสุดท้าย และจะลบข้อมูลเก่าทั้งหมด (DeleteAll) ก่อน import

#### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/import/course-csv" \
  -F "file=@csv_data/version3/example_checkjop - course_Present.csv"
```

#### Success Response
```json
{
  "message": "Course CSV imported successfully"
}
```

---

### ความแตกต่างระหว่าง Version 3 และ Version 4

| Feature | Version 3 (`/course-csv`) | Version 4 (`/course-csv-with-year`) |
|---------|---------------------------|--------------------------------------|
| CSV Format | มีคอลัมน์ `Year` | ไม่มีคอลัมน์ `Year` |
| Year Parameter | ไม่ต้องส่ง | **ต้องส่ง** year parameter |
| Delete Strategy | DeleteAll (ลบทั้งหมด) | DeleteByYear (ลบเฉพาะปีที่ระบุ) |
| Use Case | Import ครั้งแรกหรือ reset ทั้งหมด | Import/Update เฉพาะปีที่ระบุ |

---

# คู่มือการใช้ API `/api/v1/graduation/check/name`

## ภาพรวม
API นี้ใช้สำหรับตรวจสอบสถานะการจบการศึกษาของนักศึกษา โดยระบุชื่อหลักสูตร (ภาษาไทย) แทนการใช้ UUID

## Endpoint
```
POST /api/v1/graduation/check/name
```

## Request Body Structure

```json
{
  "name_th": "string",
  "admission_year": number,
  "courses": [
    {
      "course_code": "string",
      "year": number,
      "semester": number,
      "grade": "string (optional)",
      "credits": number,
      "category_name": "string (optional)"
    }
  ],
  "manual_credits": {
    "category_name": number
  },
  "exemptions": ["course_code"]
}
```

## Response Structure

```json
{
  "can_graduate": boolean,
  "gpax": number,
  "total_credits": number,
  "required_credits": number,
  "category_results": [
    {
      "category_name": "string",
      "earned_credits": number,
      "required_credits": number,
      "is_satisfied": boolean
    }
  ],
  "missing_courses": ["string"],
  "prerequisite_violations": [
    {
      "course_code": "string",
      "missing_prereqs": ["string"],
      "prereqs_taken_in_wrong_term": ["string"],
      "taken_in_wrong_term": boolean,
      "missing_coreqs": ["string"],
      "coreqs_taken_in_wrong_term": ["string"]
    }
  ],
  "credit_limit_violations": [
    {
      "year": number,
      "semester": number,
      "credits": number,
      "max_credits": number
    }
  ]
}
```

---

# ตัวอย่างการใช้งาน (Test Cases)

> **หมายเหตุ**: ตัวอย่างทั้งหมดใช้ข้อมูลจาก `csv_data/version3` และหลักสูตร **เอกเดี่ยว-ฝึกงาน-2566** (admission_year: 2566)
>
> วิชาหลักที่ใช้ในตัวอย่าง:
> - `2301173` — การเขียนโปรแกรม (4 หน่วยกิต, ไม่มี prerequisite)
> - `2301260` — เทคนิคการทำโปรแกรม (4 หน่วยกิต, prerequisite: `(2301170 AND 2301172) OR 2301173`)
> - `2301263` — โครงสร้างข้อมูลและขั้นตอนวิธีหลักมูล (4 หน่วยกิต, prerequisite: `2301260`)
> - `2301362` — Computer Network Design (3 หน่วยกิต, corequisite: `2301279 OR 2301369`)
> - `2301279` — Introduction to Computer Network (3 หน่วยกิต, ไม่มี prerequisite)
> - `2301172` — ปฏิบัติการคอมพิวเตอร์และการโปรแกรม (1 หน่วยกิต, corequisite: `2301170`, อยู่ในหลักสูตร `เอกเดี่ยว-2561` เท่านั้น)
> - `2301170` — คอมพิวเตอร์และการโปรแกรม (3 หน่วยกิต, ไม่มี prerequisite, อยู่ในหลักสูตร `เอกเดี่ยว-2561` เท่านั้น)
> - `2301290` — โครงงานขนาดเล็กทางวิทยาการคอมพิวเตอร์ (1 หน่วยกิต, prerequisite: `C.F.`)

---

## 1. กรณีพื้นฐาน - วิชาไม่มี Prerequisites

### Request
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "can_graduate": false,
  "gpax": 0,
  "total_credits": 4,
  "required_credits": 136,
  "category_results": [...],
  "missing_courses": [...],
  "prerequisite_violations": [],
  "credit_limit_violations": []
}
```

**อธิบาย**: วิชา 2301173 (การเขียนโปรแกรม) ไม่มี prerequisite สามารถลงได้โดยไม่มีการละเมิด

---

## 2. การละเมิด Prerequisites - วิชาต้องใช้เงื่อนไข `(A AND B) OR C`

### Request (ผิด - ไม่มี prerequisite)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "can_graduate": false,
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}
```

**อธิบาย**: วิชา 2301260 (เทคนิคการทำโปรแกรม) ต้องการ `(2301170 AND 2301172) OR 2301173` แต่ไม่ได้ลงสักวิชา ระบบแจ้งรายชื่อทั้งหมดที่เป็นไปได้

### Request (ถูก - ใช้ 2301173 เป็น prerequisite)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: ลง 2301173 ใน Semester 1 แล้วลง 2301260 ใน Semester 2 (ถูกต้อง — ผ่านเงื่อนไข OR ด้วย 2301173)

### Request (ถูก - ใช้ 2301170 + 2301172 เป็น prerequisite)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301170",
      "year": 2024,
      "semester": 1,
      "credits": 3
    },
    {
      "course_code": "2301172",
      "year": 2024,
      "semester": 1,
      "credits": 1
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: ลง 2301170 และ 2301172 ใน Semester 1 แล้วลง 2301260 ใน Semester 2 (ถูกต้อง — ผ่านเงื่อนไข AND ด้วย 2301170 AND 2301172)

**หมายเหตุ**: แม้ว่า 2301170 และ 2301172 จะอยู่ในหลักสูตร `เอกเดี่ยว-2561` แต่ระบบตรวจสอบว่าวิชาเหล่านี้อยู่ใน `completedCourses` (ลงทะเบียนแล้ว) และเรียนก่อน 2301260 ซึ่งผ่านเงื่อนไข prerequisites — กฎ corequisite ของ 2301172 เองจะไม่ถูกบังคับกับนักศึกษาในหลักสูตรอื่น

---

## 3. การละเมิด Strict Term Requirement

### Request (ผิด - ลง prerequisite พร้อม main course ในเทอมเดียวกัน)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301172"],
      "prereqs_taken_in_wrong_term": ["2301173"],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}
```

**อธิบาย**: วิชา 2301260 ต้องการ 2301173 เป็น prerequisite แต่ลงพร้อมกันในเทอมเดียวกัน (ผิด — ต้องลง prerequisites ก่อนอย่างน้อย 1 เทอม)

- `prereqs_taken_in_wrong_term: ["2301173"]` — มี 2301173 แต่ลงพร้อมกัน (ไม่นับ)
- `missing_prereqs: ["2301170", "2301172"]` — ยังขาด 2301170 AND 2301172 (อีกเส้นทางหนึ่ง)

---

## 4. Transitive Prerequisites (โซ่ของ Prerequisites)

### Request (ผิด - ขาด prerequisites ในโซ่ทั้งหมด)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301263",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301263",
      "missing_prereqs": ["2301260", "2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}
```

**อธิบาย**:
- 2301263 (โครงสร้างข้อมูลและขั้นตอนวิธีหลักมูล) ต้องการ 2301260
- 2301260 (เทคนิคการทำโปรแกรม) ต้องการ `(2301170 AND 2301172) OR 2301173`
- ระบบตรวจสอบ transitive prerequisites และแจ้งว่าขาดทั้ง 2301260 และ prerequisites ของมัน

### Request (ผิด - มีบางส่วนของโซ่)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    },
    {
      "course_code": "2301263",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301263",
      "missing_prereqs": ["2301260"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}
```

**อธิบาย**: มี 2301173 แล้ว (ซึ่งเป็น prerequisite ของ 2301260) แต่ยังขาด 2301260 ซึ่งเป็น direct prerequisite ของ 2301263

---

## 5. Corequisites (วิชาที่ต้องลงพร้อมกัน)

วิชา **2301362** (Computer Network Design) มี corequisite: `2301279 OR 2301369` (ต้องลง 2301279 หรือ 2301369 พร้อมกันในเทอมเดียวกัน)

### Request (ผิด - ไม่ได้ลง corequisite เลย)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301362",
      "year": 2024,
      "semester": 1,
      "credits": 3
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301362",
      "missing_prereqs": [],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": ["2301279", "2301369"],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}
```

**อธิบาย**: 2301362 ต้องการ 2301279 หรือ 2301369 เป็น corequisite แต่ไม่ได้ลงเลย จึงแจ้งว่าขาดทั้งสองวิชา

### Request (ผิด - ลง corequisite คนละเทอม)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301279",
      "year": 2024,
      "semester": 1,
      "credits": 3
    },
    {
      "course_code": "2301362",
      "year": 2024,
      "semester": 2,
      "credits": 3
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301362",
      "missing_prereqs": [],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": true,
      "missing_coreqs": ["2301369"],
      "coreqs_taken_in_wrong_term": ["2301279"]
    }
  ]
}
```

**อธิบาย**: ลง 2301279 คนละเทอมกับ 2301362 จึงไม่นับเป็น corequisite (ต้องลงพร้อมกัน) และยังขาด 2301369 ในอีกเส้นทางของ OR

### Request (ถูก - ลงพร้อมกัน)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301279",
      "year": 2024,
      "semester": 1,
      "credits": 3
    },
    {
      "course_code": "2301362",
      "year": 2024,
      "semester": 1,
      "credits": 3
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

### Request (ถูก - ลงพร้อมกันแม้ได้เกรด F)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301279",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "F"
    },
    {
      "course_code": "2301362",
      "year": 2024,
      "semester": 1,
      "credits": 3
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: แม้ว่า 2301279 จะได้เกรด F แต่เนื่องจากเป็น corequisite (ลงพร้อมกัน) จึงยังนับว่าผ่านเงื่อนไข corequisite (ต้องลงพร้อมกันเท่านั้น ไม่สนใจเกรด)

**หมายเหตุ**: นี่เป็นข้อแตกต่างสำคัญระหว่าง prerequisite และ corequisite
- **Prerequisite + เกรด F**: ถือว่าไม่ผ่าน ต้องเรียนใหม่
- **Corequisite + เกรด F**: ถือว่าผ่าน เพราะเงื่อนไขคือต้อง "ลงพร้อมกัน" ไม่ได้กำหนดว่าต้องได้เกรดผ่าน

---

## 6. Credit Limits (จำกัดหน่วยกิตต่อเทอม)

### Request (ผิด - เทอมปกติเกิน 22 หน่วยกิต)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 15
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 1,
      "credits": 8
    }
  ]
}
```

### Response
```json
{
  "credit_limit_violations": [
    {
      "year": 2024,
      "semester": 1,
      "credits": 23,
      "max_credits": 22
    }
  ]
}
```

**อธิบาย**: เทอมปกติ (1, 2) จำกัดที่ 22 หน่วยกิต

### Request (ผิด - เทอมฤดูร้อนเกิน 10 หน่วยกิต)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 3,
      "credits": 8
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 3,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "credit_limit_violations": [
    {
      "year": 2024,
      "semester": 3,
      "credits": 12,
      "max_credits": 10
    }
  ]
}
```

**อธิบาย**: เทอมฤดูร้อน (semester 3) จำกัดที่ 10 หน่วยกิต

---

## 7. Course Year Versioning (วิชาเวอร์ชันต่างปี)

### Request
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

**อธิบาย**: ระบบจะใช้เวอร์ชันของวิชาตาม `admission_year` ของนักศึกษา
- นักศึกษาที่เข้าปี 2566 จะใช้วิชาเวอร์ชัน 2566
- นักศึกษาที่เข้าปี 2567 จะใช้วิชาเวอร์ชัน 2567
- วิชาแต่ละปีอาจมี prerequisites แตกต่างกัน

---

## 8. Complex Prerequisites (เงื่อนไข Prerequisites `(A AND B) OR C`)

วิชา **2301260** (เทคนิคการทำโปรแกรม) มี prerequisite: `(2301170 AND 2301172) OR 2301173`

แปลว่า: ต้องเรียน **(2301170 และ 2301172 ทั้งคู่)** หรือ **เรียน 2301173** อย่างใดอย่างหนึ่ง

### Request (ผิด - ไม่มีทั้ง 3 วิชา)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 1,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301173", "2301172"]
    }
  ]
}
```

**อธิบาย**: ไม่มีทั้ง 3 วิชา จึงแจ้งว่าขาดทุกวิชาที่เป็นไปได้

### Request (ถูก - มี 2301173)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: มี 2301173 ผ่านเงื่อนไข OR group ฝั่ง C

### Request (ถูก - มี 2301170 และ 2301172)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301170",
      "year": 2024,
      "semester": 1,
      "credits": 3
    },
    {
      "course_code": "2301172",
      "year": 2024,
      "semester": 1,
      "credits": 1
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: มี 2301170 AND 2301172 ทั้งคู่ในรายการ completedCourses และเรียนก่อน 2301260 ผ่านเงื่อนไข AND group ฝั่ง (A AND B)

---

## 9. Cross Curriculum Prerequisites (Prerequisites ข้ามหลักสูตร)

ระบบตรวจสอบ prerequisites และ corequisites โดยใช้ทั้ง **curriculum ID** และ **ปีการรับเข้า** ดังนั้นวิชาที่อยู่ในหลักสูตรอื่นจะไม่ถูกบังคับกฎจากหลักสูตรนั้น

วิชา **2301172** (ปฏิบัติการคอมพิวเตอร์และการโปรแกรม) อยู่ในหลักสูตร `เอกเดี่ยว-2561` เท่านั้น ไม่อยู่ใน `เอกเดี่ยว-ฝึกงาน-2566`

### Request (นักศึกษา เอกเดี่ยว-ฝึกงาน-2566 ลง 2301172 ที่เป็นวิชาจากหลักสูตรอื่น)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301170",
      "year": 2024,
      "semester": 1,
      "credits": 3
    },
    {
      "course_code": "2301172",
      "year": 2024,
      "semester": 2,
      "credits": 1
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**:
- 2301172 อยู่ในหลักสูตร `เอกเดี่ยว-2561` ไม่ใช่ `เอกเดี่ยว-ฝึกงาน-2566`
- ระบบค้นหาวิชาด้วย curriculum ID + year ของนักศึกษา → ไม่พบ 2301172 ในหลักสูตรนี้
- จึงไม่มีการตรวจ corequisite สำหรับวิชาที่ไม่อยู่ในหลักสูตรของนักศึกษา → ไม่มี violation
- **พฤติกรรมนี้เจตนา**: วิชาจากหลักสูตรอื่นไม่ควรบังคับกฎ prerequisites/corequisites ของหลักสูตรนั้น

---

## 10. Free Electives (วิชาเลือกเสรี)

### Request (วิชาเลือกเสรีที่ไม่มีใน DB)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "9999999",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "category_name": "วิชาเสรี"
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [],
  "category_results": [
    {
      "category_name": "วิชาเสรี",
      "earned_credits": 3,
      "required_credits": 6,
      "is_satisfied": false
    }
  ]
}
```

**อธิบาย**: วิชาที่ไม่มีใน DB แต่ระบุเป็น "วิชาเสรี" จะนับหน่วยกิตได้ และข้ามการตรวจสอบ prerequisite

**หมายเหตุ**: `category_name` สามารถใช้ได้ทั้งภาษาไทย (`"วิชาเสรี"`) และภาษาอังกฤษ (`"Free Electives"`)

---

## 11. C.F. Permission (Exemptions)

วิชา **2301290** (โครงงานขนาดเล็กทางวิทยาการคอมพิวเตอร์) มี prerequisite เป็น `C.F.` (Consent of Faculty) เพียงอย่างเดียว

### Request (ไม่มีสิทธิ์ C.F.)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301290",
      "year": 2024,
      "semester": 1,
      "credits": 1
    }
  ],
  "exemptions": []
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**หมายเหตุ**: เมื่อ prerequisite เป็น `C.F.` เพียงอย่างเดียว (ไม่มีวิชาอื่นร่วม) ระบบจะไม่แจ้ง violation ไม่ว่าจะมีสิทธิ์ C.F. หรือไม่ก็ตาม เนื่องจากไม่มี course code จริงที่ต้องตรวจสอบ

กรณีที่ `exemptions` มีผลจริงคือเมื่อ prerequisite เป็นแบบ `"courseCode OR C.F."` — หากนักศึกษามี courseCode แล้ว จะผ่านโดยอัตโนมัติ แต่ถ้าไม่มี courseCode จะต้องระบุ course นั้นใน `exemptions` เพื่อบอกว่ามีสิทธิ์ C.F.

### Request (ใช้ exemptions กับ course ที่มีเงื่อนไข OR C.F.)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301290",
      "year": 2024,
      "semester": 1,
      "credits": 1
    }
  ],
  "exemptions": ["2301290"]
}
```

### Response
```json
{
  "prerequisite_violations": []
}
```

**อธิบาย**: เมื่อระบุ course code ใน `exemptions` หมายความว่านักศึกษาได้รับอนุมัติ C.F. สำหรับวิชานั้น จึงสามารถลงได้

---

## 12. GPAX Calculation

### Request (เกรดปกติ)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "A",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "A"
    },
    {
      "course_code": "B",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "B"
    },
    {
      "course_code": "C",
      "year": 2024,
      "semester": 1,
      "credits": 2,
      "grade": "C+"
    }
  ]
}
```

### Response
```json
{
  "gpax": 3.25,
  "total_credits": 8
}
```

**อธิบาย**:
- A (4.0 × 3) + B (3.0 × 3) + C+ (2.5 × 2) = 26 คะแนน
- 26 ÷ 8 หน่วยกิต = 3.25

### Request (มีเกรด W, S ที่ไม่นับ GPAX)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "A",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "A"
    },
    {
      "course_code": "B",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "W"
    },
    {
      "course_code": "C",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "S"
    }
  ]
}
```

### Response
```json
{
  "gpax": 4,
  "total_credits": 9
}
```

**อธิบาย**: เกรด W (Withdraw) และ S (Satisfactory) ไม่นับใน GPAX (GPAX = 4.0 จากวิชา A เพียงวิชาเดียว) แต่ `total_credits` นับหน่วยกิตทุกวิชาที่ลงทะเบียน

### Request (มีเกรด F)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "A",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "A"
    },
    {
      "course_code": "B",
      "year": 2024,
      "semester": 1,
      "credits": 3,
      "grade": "F"
    }
  ]
}
```

### Response
```json
{
  "gpax": 2,
  "total_credits": 6
}
```

**อธิบาย**:
- A (4.0 × 3) + F (0.0 × 3) = 12 คะแนน
- 12 ÷ 6 หน่วยกิต = 2.0
- เกรด F นับใน GPAX (0.0 คะแนน)

### Request (เกรด F ไม่นับเป็น prerequisite)
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4,
      "grade": "F"
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4,
      "grade": "A"
    }
  ]
}
```

### Response
```json
{
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ],
  "gpax": 2,
  "total_credits": 8
}
```

**อธิบาย**: แม้จะลง 2301173 แล้ว แต่ได้เกรด F จึงถือว่ายังไม่ผ่าน prerequisite ของ 2301260

---

## 13. Manual Credits (ระบุหน่วยกิตเอง)

### Request
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 12
    }
  ],
  "manual_credits": {
    "วิชาพื้นฐานวิทยาศาสตร์": 12
  }
}
```

### Response
```json
{
  "category_results": [
    {
      "category_name": "วิชาพื้นฐานวิทยาศาสตร์",
      "earned_credits": 12,
      "required_credits": 12,
      "is_satisfied": true
    }
  ]
}
```

**อธิบาย**: สามารถเพิ่มหน่วยกิตในหมวดเอง (เช่น transfer credit) ได้ผ่าน `manual_credits` โดยระบุชื่อหมวดตรงตาม category ในระบบ

**หมวดวิชาที่มีในหลักสูตรเอกเดี่ยว-ฝึกงาน-2566**:
- `วิชาแกน` (14 หน่วยกิต)
- `วิชาเฉพาะด้าน` (39 หน่วยกิต)
- `วิชาพื้นฐานวิทยาศาสตร์` (12 หน่วยกิต)
- `วิชาศึกษาทั่วไป-สหกิจ` (3 หน่วยกิต)
- `วิชาศึกษาทั่วไป-มนุษย์` (3 หน่วยกิต)
- `วิชาศึกษาทั่วไป-สังคม` (3 หน่วยกิต)
- `วิชาศึกษาทั่วไป-วิทย์` (3 หน่วยกิต)
- `กลุ่มวิชาภาษา` (12 หน่วยกิต)
- `วิชาศึกษาทั่วไปกลุ่มพิเศษ` (6 หน่วยกิต)
- `วิชาบังคับเลือก` (3 หน่วยกิต)
- `วิชาเลือก` (24 หน่วยกิต)
- `วิชาประสบการณ์ภาคสนาม` (8 หน่วยกิต)
- `วิชาเสรี` (6 หน่วยกิต)

---

## 14. Category Requirements (ตรวจสอบหมวดวิชา)

### Request
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 4
    },
    {
      "course_code": "2301180",
      "year": 2024,
      "semester": 2,
      "credits": 2
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 4
    }
  ]
}
```

### Response
```json
{
  "category_results": [
    {
      "category_name": "วิชาเฉพาะด้าน",
      "earned_credits": 10,
      "required_credits": 39,
      "is_satisfied": false
    }
  ]
}
```

**อธิบาย**: แสดงสถานะของแต่ละหมวดว่าได้หน่วยกิตเท่าไหร่และต้องการเท่าไหร่ (เฉพาะหมวดที่มีหน่วยกิตสะสม)

---

## 15. Complete Graduation Check (ตรวจสอบจบการศึกษาแบบสมบูรณ์)

### Request
```json
{
  "name_th": "เอกเดี่ยว-ฝึกงาน-2566",
  "admission_year": 2566,
  "courses": [
    {
      "course_code": "2301173",
      "year": 2024,
      "semester": 1,
      "credits": 15,
      "grade": "A"
    },
    {
      "course_code": "2301260",
      "year": 2024,
      "semester": 2,
      "credits": 15,
      "grade": "B+"
    },
    {
      "course_code": "2301263",
      "year": 2025,
      "semester": 1,
      "credits": 22,
      "grade": "A"
    },
    {
      "course_code": "2301263",
      "year": 2025,
      "semester": 2,
      "credits": 22,
      "grade": "B"
    },
    {
      "course_code": "2301263",
      "year": 2026,
      "semester": 1,
      "credits": 22,
      "grade": "A"
    },
    {
      "course_code": "2301263",
      "year": 2026,
      "semester": 2,
      "credits": 22,
      "grade": "A"
    },
    {
      "course_code": "2301263",
      "year": 2027,
      "semester": 1,
      "credits": 2,
      "grade": "B+"
    }
  ],
  "manual_credits": {
    "วิชาพื้นฐานวิทยาศาสตร์": 12,
    "วิชาแกน": 14,
    "กลุ่มวิชาภาษา": 12
  }
}
```

**อธิบาย**: ตัวอย่างการส่ง courses หลายเทอมพร้อม manual_credits เพื่อตรวจสอบแบบครบทุกด้าน:
- หน่วยกิตรวมจาก courses
- แต่ละหมวดวิชา
- GPAX
- การละเมิด prerequisites/credit limits

---

## สรุปกฎสำคัญ

### Prerequisites
1. **OR Group**: ต้องมีอย่างน้อย 1 วิชาในกลุ่ม
2. **AND Group**: ต้องมีครบทุกวิชาในกลุ่ม
3. **Strict Term**: prerequisite ต้องเรียนก่อนอย่างน้อย 1 เทอม (ไม่สามารถลงพร้อมกันได้)
4. **Transitive**: ระบบตรวจสอบโซ่ของ prerequisites ทั้งหมด
5. **Grade F**: ถือว่าไม่ผ่าน prerequisite
6. **C.F. Permission**: ใช้ `exemptions` เพื่อยกเว้น prerequisite กรณีเงื่อนไข "courseCode OR C.F."

### Corequisites
1. ต้องลงในเทอมเดียวกัน
2. Grade F ยังนับเป็นผ่าน corequisite (เพราะลงพร้อมกัน)

### Credit Limits
1. **เทอมปกติ (1, 2)**: สูงสุด 22 หน่วยกิต
2. **เทอมฤดูร้อน (3)**: สูงสุด 10 หน่วยกิต

### GPAX
1. เกรดที่นับ: A, B+, B, C+, C, D+, D, F
2. เกรดที่ไม่นับใน GPAX: W (Withdraw), S (Satisfactory), U (Unsatisfactory)
3. เกรด F นับเป็น 0.0 คะแนน
4. `total_credits` นับหน่วยกิตทุกวิชาที่ลงทะเบียน (รวม W และ S)

### Free Electives
1. ถ้ามี prerequisite ต้องตรวจสอบ
2. ถ้าไม่มีใน DB ข้ามการตรวจสอบ prerequisite แต่นับหน่วยกิต

### Course Year Versioning
1. ใช้เวอร์ชันวิชาตาม `admission_year` ของนักศึกษา
2. ระบบใช้ `GetByCodeAndCurriculumIDAndYear` ซึ่งผูกกับทั้ง curriculum ID **และ** ปีของวิชานั้นๆ
3. วิชาที่ไม่อยู่ในหลักสูตรของนักศึกษา จะไม่มีการตรวจ prerequisites/corequisites (ถือว่าไม่มีกฎ)

---

## Error Handling

### Bad Request (400)
```json
{
  "error": "invalid request body"
}
```

### Internal Server Error (500)
```json
{
  "error": "curriculum not found"
}
```

---

## หมายเหตุ
- ระบบตรวจสอบ prerequisites แบบ transitive (recursive)
- การนับ credits สำหรับแต่ละหมวดจะนับเฉพาะวิชาที่อยู่ในหมวดนั้น
- `manual_credits` ใช้สำหรับกรณี transfer credit หรือหน่วยกิตพิเศษ
- `exemptions` ใช้สำหรับกรณีได้รับอนุมัติ C.F. (Consent of Faculty)

---

## C.F. (Consent of Faculty) API

### 1. Check C.F. Option for a Course

ตรวจสอบว่าวิชาแต่ละวิชามีสิทธิ์ใช้ C.F. exemption หรือไม่

#### Endpoint
```
GET /api/v1/courses/code/:code/cf-option
```

#### Query Parameters
- `curriculum_id` (required) - UUID ของหลักสูตร
- `year` (required) - ปีการศึกษา (เช่น 2566, 2567, 2568)

#### Example Request
```bash
curl "http://localhost:8080/api/v1/courses/code/2301290/cf-option?curriculum_id=763dd07d-84d4-4c2a-b9c1-995d526ecc4b&year=2566"
```

#### Success Response (Course WITH C.F. option)
```json
{
  "course_code": "2301290",
  "course_name_en": "Computer Science Mini Project",
  "course_name_th": "โครงงานขนาดเล็กทางวิทยาการคอมพิวเตอร์",
  "has_cf_option": true,
  "prerequisite_groups": [
    {
      "is_or_group": false,
      "has_cf_condition": true,
      "course_codes": []
    }
  ],
  "corequisite_groups": [],
  "message": "This course allows C.F. exemption"
}
```

#### Success Response (Course WITHOUT C.F. option)
```json
{
  "course_code": "2301173",
  "course_name_en": "Programming",
  "course_name_th": "การเขียนโปรแกรม",
  "has_cf_option": false,
  "prerequisite_groups": [],
  "corequisite_groups": [],
  "message": "This course does NOT allow C.F. exemption"
}
```

#### Error Response
```json
{
  "error": "Course not found"
}
```

### 2. C.F. Exemption Validation

เมื่อส่ง `exemptions` array มาใน prerequisite validation API ระบบจะตรวจสอบว่า:
1. ✅ วิชาที่ใส่ใน exemptions มีอยู่ในหลักสูตร
2. ✅ วิชานั้นมี `has_cf_option = true` (มี C.F. ใน prerequisites หรือ corequisites)

#### Example: Valid C.F. Exemption

**Request**
```bash
curl -X POST http://localhost:8080/api/v1/graduation/prerequisites \
  -H "Content-Type: application/json" \
  -d '{
    "curriculum_id": "763dd07d-84d4-4c2a-b9c1-995d526ecc4b",
    "admission_year": 2566,
    "courses": [
      {
        "course_code": "2301290",
        "year": 2566,
        "semester": 1,
        "credits": 1,
        "grade": "A"
      }
    ],
    "exemptions": ["2301290"]
  }'
```

**Response**
```json
{
  "violations": []
}
```

**อธิบาย**: วิชา 2301290 มี C.F. ใน prerequisite ดังนั้นสามารถใช้ C.F. exemption ได้

#### Example: Invalid C.F. Exemption

**Request**
```bash
curl -X POST http://localhost:8080/api/v1/graduation/prerequisites \
  -H "Content-Type: application/json" \
  -d '{
    "curriculum_id": "763dd07d-84d4-4c2a-b9c1-995d526ecc4b",
    "admission_year": 2566,
    "courses": [
      {
        "course_code": "2301173",
        "year": 2566,
        "semester": 1,
        "credits": 4,
        "grade": "A"
      }
    ],
    "exemptions": ["2301173"]
  }'
```

**Response**
```json
{
  "error": "course '2301173' does not allow C.F. exemption (no C.F. option in prerequisites or corequisites)"
}
```

**อธิบาย**: วิชา 2301173 ไม่มี C.F. ใน prerequisite หรือ corequisite ดังนั้นไม่สามารถใช้ C.F. exemption ได้

### 3. How C.F. Option is Determined

ระบบจะตั้งค่า `has_cf_option = true` อัตโนมัติเมื่อ import CSV ถ้า:
- พบ "C.F." ใน column `prerequisites` ของวิชา **หรือ**
- พบ "C.F." ใน column `corequisites` ของวิชา

ตัวอย่าง CSV:
```csv
code,courseNameEN,courseNameTH,credit,prerequisites,corequisites,category,curriculum
2301290,Computer Science Mini Project,โครงงานขนาดเล็กทางวิทยาการคอมพิวเตอร์,1,C.F.,,Elective Courses(CS),เอกเดี่ยว-ฝึกงาน-2566
```

→ วิชานี้จะได้ `has_cf_option = true`

### 4. Use Cases

#### Use Case 1: ตรวจสอบว่าวิชานี้สามารถขอ C.F. ได้หรือไม่
```bash
# ขั้นตอนที่ 1: หา curriculum_id
curl http://localhost:8080/api/v1/curricula/ | jq '.[] | {id, name_th}'

# ขั้นตอนที่ 2: ตรวจสอบ C.F. option
curl "http://localhost:8080/api/v1/courses/code/2301290/cf-option?curriculum_id=YOUR_CURRICULUM_ID&year=2566"
```

#### Use Case 2: ตรวจสอบว่าวิชาใดบ้างในหลักสูตรที่มี C.F. option
```bash
curl "http://localhost:8080/api/v1/courses/curriculum/YOUR_CURRICULUM_ID" | \
  jq '.[] | select(.has_cf_option == true) | {code, name_th, has_cf_option}'
```

#### Use Case 3: Validate C.F. exemption ก่อนบันทึกข้อมูล
```bash
# ระบบจะ return error ทันทีถ้าพยายามใส่ C.F. ให้วิชาที่ไม่อนุญาต
curl -X POST http://localhost:8080/api/v1/graduation/prerequisites \
  -H "Content-Type: application/json" \
  -d '{
    "curriculum_id": "YOUR_CURRICULUM_ID",
    "admission_year": 2566,
    "courses": [...],
    "exemptions": ["COURSE_CODE"]
  }'
```

### 5. Error Messages

| Error | Meaning | Solution |
|-------|---------|----------|
| `course 'XXX' does not allow C.F. exemption` | วิชานี้ไม่มี C.F. ใน prerequisite/corequisite | ลบวิชานี้ออกจาก exemptions |
| `exemption course 'XXX' not found in curriculum` | วิชาไม่มีในหลักสูตร | ตรวจสอบ course code และ curriculum_id |
| `Course not found` | ไม่พบวิชาในระบบ | ตรวจสอบ course_code, curriculum_id, year |
| `curriculum_id is required` | ขาด query parameter | เพิ่ม `?curriculum_id=xxx` |
| `year is required` | ขาด query parameter | เพิ่ม `&year=2566` |

### 6. Database Schema

ฟิลด์ใหม่ที่เพิ่มเข้ามา:

```sql
-- Column: has_cf_option
-- Type: BOOLEAN
-- Default: FALSE
-- Description: Indicates whether this course allows C.F. exemption
-- Auto-populated during CSV import when "C.F." is found in prerequisites or corequisites
```

---
