# CheckJop Backend - System Overview

## 🏗️ 1. System Architecture

### Clean Architecture Overview

```mermaid
graph TB
    subgraph External["🌐 External Layer"]
        Client[Client/API Requests]
    end
    
    subgraph Interface["🎯 Interface Layer"]
        Handler[Handlers<br/>HTTP Controllers]
        Router[Routes & Middleware]
    end
    
    subgraph UseCase["💼 Use Case Layer"]
        Service[Services<br/>Business Logic]
    end
    
    subgraph Data["💾 Data Layer"]
        Repository[Repositories<br/>Data Access]
    end
    
    subgraph Framework["🔧 Framework & Drivers"]
        DB[(PostgreSQL)]
        Gin[Gin Framework]
        GORM[GORM ORM]
    end
    
    Client -->|HTTP| Router
    Router --> Handler
    Handler -->|Dependency| Service
    Service -->|Dependency| Repository
    Repository --> GORM
    GORM --> DB
    Router -.->|Uses| Gin
    
    style External fill:#e1f5ff
    style Interface fill:#ffe1e1
    style UseCase fill:#e1ffe1
    style Data fill:#f0e1ff
    style Framework fill:#ffd700
```

**Dependency Rule:** แต่ละ layer ขึ้นกับ layer ด้านในเท่านั้น (Inner layers ไม่รู้จัก outer layers)

### Detailed Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Client[Client/Frontend]
    end
    
    subgraph "API Layer"
        Router[Gin Router<br/>routes/routes.go]
        MW[Middleware<br/>CORS, Logger, RateLimit, Error]
    end
    
    subgraph "Handler Layer"
        CH[Category Handler]
        CRH[Course Handler]
        CUH[Curriculum Handler]
        GH[Graduation Handler]
        SDH[SetDefault Handler]
        CSVH[CSV Import Handler]
    end
    
    subgraph "Service Layer"
        CS[Category Service]
        CRS[Course Service]
        CUS[Curriculum Service]
        GS[Graduation Service]
        SDS[SetDefault Service]
    end
    
    subgraph "Repository Layer"
        CR[Category Repository]
        CRR[Course Repository]
        CUR[Curriculum Repository]
        SDR[SetDefault Repository]
    end
    
    subgraph "Database Layer"
        DB[(PostgreSQL<br/>GORM)]
    end
    
    Client -->|HTTP Request| Router
    Router --> MW
    MW --> CH & CRH & CUH & GH & SDH & CSVH
    
    CH --> CS
    CRH --> CRS
    CUH --> CUS
    GH --> GS
    SDH --> SDS
    CSVH --> CRS & CUS
    
    CS --> CR
    CRS --> CRR
    CUS --> CUR
    GS --> CRR & CUR & SDR
    SDS --> SDR
    
    CR --> DB
    CRR --> DB
    CUR --> DB
    SDR --> DB
    
    style Client fill:#e1f5ff
    style Router fill:#fff4e1
    style MW fill:#fff4e1
    style CH fill:#ffe1e1
    style CRH fill:#ffe1e1
    style CUH fill:#ffe1e1
    style GH fill:#ffe1e1
    style SDH fill:#ffe1e1
    style CSVH fill:#ffe1e1
    style CS fill:#e1ffe1
    style CRS fill:#e1ffe1
    style CUS fill:#e1ffe1
    style GS fill:#e1ffe1
    style SDS fill:#e1ffe1
    style CR fill:#f0e1ff
    style CRR fill:#f0e1ff
    style CUR fill:#f0e1ff
    style SDR fill:#f0e1ff
    style DB fill:#ffd700
```

## Data Model Relationships

```mermaid
erDiagram
    CURRICULUM ||--o{ CATEGORY : contains
    CURRICULUM ||--o{ COURSE : contains
    CURRICULUM ||--o{ SET_DEFAULT : has
    CATEGORY ||--o{ COURSE : categorizes
    COURSE ||--o{ PREREQUISITE_GROUP : has
    COURSE ||--o{ SET_DEFAULT : included_in
    PREREQUISITE_GROUP ||--o{ PREREQUISITE_COURSE_LINK : contains
    PREREQUISITE_COURSE_LINK }o--|| COURSE : references

    CURRICULUM {
        uuid id PK
        string name_th
        string name_en
        int year
        int min_total_credits
        bool is_active
    }

    CATEGORY {
        uuid id PK
        uuid curriculum_id FK
        string name_th
        string name_en
        int min_credits
        int sort_order
    }

    COURSE {
        uuid id PK
        uuid curriculum_id FK
        uuid category_id FK
        string code
        string name_th
        string name_en
        int credits
        string description
    }

    PREREQUISITE_GROUP {
        uuid id PK
        uuid course_id FK
        string group_type
        bool is_or_group
    }

    PREREQUISITE_COURSE_LINK {
        uuid id PK
        uuid group_id FK
        uuid prerequisite_course_id FK
    }

    SET_DEFAULT {
        uuid id PK
        uuid curriculum_id FK
        uuid course_id FK
        int year
        int semester
    }
```

## API Flow - Graduation Check Example

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router
    participant GH as Graduation Handler
    participant GS as Graduation Service
    participant CR as Curriculum Repo
    participant COR as Course Repo
    participant CAR as Category Repo
    participant DB as Database

    C->>R: POST /api/v1/graduation/check
    R->>GH: CheckGraduation(progress)
    GH->>GS: CheckGraduation(progress)
    
    GS->>CR: GetByID(curriculum_id)
    CR->>DB: SELECT * FROM curricula
    DB-->>CR: Curriculum data
    CR-->>GS: Curriculum
    
    GS->>COR: GetByCurriculumID(curriculum_id)
    COR->>DB: SELECT * FROM courses
    DB-->>COR: Course list
    COR-->>GS: Courses with prerequisites
    
    GS->>CAR: GetByCurriculumID(curriculum_id)
    CAR->>DB: SELECT * FROM categories
    DB-->>CAR: Category list
    CAR-->>GS: Categories
    
    Note over GS: Validate Prerequisites
    Note over GS: Check Credit Limits
    Note over GS: Check Category Requirements
    
    GS-->>GH: GraduationCheckResult
    GH-->>R: JSON Response
    R-->>C: 200 OK with result
```

## Business Logic Flow - Prerequisite Validation

```mermaid
flowchart TD
    Start[Start Validation] --> GetCourses[Get Student's Completed Courses]
    GetCourses --> Loop[For Each Completed Course]
    
    Loop --> GetPrereqGroups[Get Prerequisite Groups]
    GetPrereqGroups --> HasPrereq{Has Prerequisites?}
    
    HasPrereq -->|No| NextCourse[Continue to Next Course]
    HasPrereq -->|Yes| CheckGroups[Check Each Group]
    
    CheckGroups --> IsORGroup{Is OR Group?}
    
    IsORGroup -->|Yes| CheckOR[Check if ANY course satisfied]
    IsORGroup -->|No| CheckAND[Check if ALL courses satisfied]
    
    CheckOR --> ORSatisfied{Any Satisfied?}
    ORSatisfied -->|No| AddViolation[Add Violation]
    ORSatisfied -->|Yes| CheckTerm[Check Term Requirement]
    
    CheckAND --> ANDSatisfied{All Satisfied?}
    ANDSatisfied -->|No| AddMissing[Add Missing Prerequisites]
    ANDSatisfied -->|Yes| CheckTerm
    
    CheckTerm --> SameTerm{Taken in Same Term?}
    SameTerm -->|Yes| AddTermViolation[Add Term Violation]
    SameTerm -->|No| ValidTransitive[Check Transitive Prerequisites]
    
    ValidTransitive --> NextCourse
    AddViolation --> NextCourse
    AddMissing --> NextCourse
    AddTermViolation --> NextCourse
    
    NextCourse --> MoreCourses{More Courses?}
    MoreCourses -->|Yes| Loop
    MoreCourses -->|No| Return[Return Violations List]
    
    Return --> End[End]
```

## Component Responsibilities

```mermaid
mindmap
  root((CheckJop System))
    API Layer
      Routing
      Middleware
        CORS
        Rate Limiting
        Logging
        Error Handling
    Handlers
      Request Validation
      Response Formatting
      HTTP Status Codes
    Services
      Business Logic
        Prerequisite Validation
        Credit Calculation
        Graduation Checking
      CSV Import Processing
      Data Transformation
    Repositories
      Database Queries
      Data Persistence
      Transaction Management
    Models
      Data Structures
        Curriculum
        Course
        Category
        Prerequisites
        Graduation Status
```

## Deployment Architecture

```mermaid
graph LR
    subgraph "Development"
        Dev[Developer] --> Git[Git Repository]
    end

    subgraph "CI/CD Pipeline"
        Git --> GHA[GitHub Actions]
        GHA --> Test[Run Tests]
        GHA --> Lint[Code Quality Check]
        GHA --> Security[Security Scan]
        Test --> Build[Build Binary]
        Lint --> Build
        Security --> Build
        Build --> Docker[Build Docker Image]
    end

    subgraph "Container Registry"
        Docker --> Hub[Docker Hub]
    end

    subgraph "Production Environment"
        Hub --> Deploy[Container Deployment]
        Deploy --> App[Go Application]
        App --> PG[(PostgreSQL)]
        LB[Load Balancer] --> App
    end

    subgraph "Monitoring"
        App --> Logs[Logs]
        App --> Metrics[Metrics]
    end

    Users[End Users] --> LB
```

## CSV Import Flow

```mermaid
flowchart TD
    Start[Upload CSV File] --> Parse[Parse CSV]
    Parse --> ValidateFormat{Valid Format?}
    
    ValidateFormat -->|No| Error[Return Error]
    ValidateFormat -->|Yes| CheckType{File Type?}
    
    CheckType -->|Curriculum| ParseCurriculum[Parse Curriculum Data]
    CheckType -->|Category| ParseCategory[Parse Category Data]
    CheckType -->|Course| ParseCourse[Parse Course Data]
    CheckType -->|SetDefault| ParseSetDefault[Parse SetDefault Data]
    
    ParseCurriculum --> ValidateCurriculum[Validate Data]
    ParseCategory --> ValidateCategory[Validate Data + Find Curriculum]
    ParseCourse --> ValidateCourse[Validate Data + Find Category]
    ParseSetDefault --> ValidateSetDefault[Validate Data + Find Course]
    
    ValidateCurriculum --> Transaction1[Start Transaction]
    ValidateCategory --> Transaction2[Start Transaction]
    ValidateCourse --> Transaction3[Start Transaction]
    ValidateSetDefault --> Transaction4[Start Transaction]
    
    Transaction1 --> Delete1[Delete Existing by Key]
    Transaction2 --> Delete2[Delete Existing by name_th + name_en + curriculum_id]
    Transaction3 --> Delete3[Delete Existing by code + curriculum_id]
    Transaction4 --> Delete4[Delete Existing by curriculum_id + course_id]
    
    Delete1 --> Insert1[Insert All New Records]
    Delete2 --> Insert2[Insert All New Records]
    Delete3 --> ProcessRelations[Process Prerequisites/Corequisites]
    Delete4 --> Insert4[Insert All New Records]
    
    ProcessRelations --> Insert3[Insert All New Records]
    
    Insert1 --> Commit1[Commit Transaction]
    Insert2 --> Commit2[Commit Transaction]
    Insert3 --> Commit3[Commit Transaction]
    Insert4 --> Commit4[Commit Transaction]
    
    Commit1 --> Success[Return Success]
    Commit2 --> Success
    Commit3 --> Success
    Commit4 --> Success
    
    Error --> End[End]
    Success --> End
    
    Delete1 -.Rollback on Error.-> Rollback[Rollback Transaction]
    Delete2 -.Rollback on Error.-> Rollback
    Delete3 -.Rollback on Error.-> Rollback
    Delete4 -.Rollback on Error.-> Rollback
    Insert1 -.Rollback on Error.-> Rollback
    Insert2 -.Rollback on Error.-> Rollback
    Insert3 -.Rollback on Error.-> Rollback
    Insert4 -.Rollback on Error.-> Rollback
    Rollback --> Error
```

**หมายเหตุ:** ใช้ Delete-then-Insert pattern ภายใน Transaction (ไม่ใช่ drop table)
- Category: ลบที่มี name_th, name_en, curriculum_id เหมือนกัน แล้ว insert ใหม่
- Course: ลบที่มี code, curriculum_id เหมือนกัน แล้ว insert ใหม่พร้อม relationships
- SetDefault: ลบที่มี curriculum_id, course_id เหมือนกัน แล้ว insert ใหม่
- มี rollback เมื่อเกิด error ระหว่าง transaction

## Key Features

### 1. Curriculum Management
- สร้าง/แก้ไข/ลบหลักสูตร
- Import หลักสูตรจาก CSV
- Query หลักสูตรตามปี/ชื่อ

### 2. Course Management
- จัดการรายวิชาในหลักสูตร
- กำหนด Prerequisites (OR/AND groups)
- กำหนด Corequisites
- Support Transitive Prerequisites

### 3. Graduation Checking
- ตรวจสอบเงื่อนไขการจบ
- Validate Prerequisites/Corequisites
- Check Credit Limits (ปกติ 22, ฤดูร้อน 10)
- Check Category Requirements

### 4. Set Default
- กำหนดแผนการเรียนมาตรฐาน
- ระบุปี/เทอมที่แนะนำ

### 5. CSV Import
- Import ข้อมูลจาก CSV files
- Bulk upsert operations
- Transaction support

## Technology Stack

```mermaid
graph TB
    subgraph "Backend"
        Go[Go 1.21+]
        Gin[Gin Framework]
        GORM[GORM ORM]
    end

    subgraph "Database"
        PostgreSQL[PostgreSQL]
        UUID[UUID Support]
    end

    subgraph "DevOps"
        Docker[Docker]
        GHA[GitHub Actions]
        Make[Makefile]
    end

    subgraph "Code Quality"
        GolangCI[golangci-lint]
        Gosec[Gosec]
        CodeQL[CodeQL]
    end

    Go --> Gin
    Gin --> GORM
    GORM --> PostgreSQL
    
    Docker -.-> Go
    GHA -.-> Docker
    Make -.-> Go
    
    GolangCI -.-> Go
    Gosec -.-> Go
    CodeQL -.-> Go
```
