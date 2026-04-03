# CheckJop Backend

A Go-based REST API backend for academic curriculum and graduation validation system built with Gin framework and PostgreSQL database.

## Features

- **Curriculum Management** - Create and manage academic curricula
- **Course Management** - Handle courses with prerequisites and corequisites
- **Graduation Validation** - Validate student progress and graduation eligibility
- **Prerequisite Checking** - Advanced prerequisite validation with OR/AND logic
- **Credit Limits** - Enforce credit limits per semester (22 regular, 10 summer)
- **Category Requirements** - Track progress across different course categories

## Tech Stack

- **Go 1.23**
- **Gin** - HTTP web framework
- **GORM** - ORM library with PostgreSQL
- **PostgreSQL** - Database
- **Docker** - Containerization
- **Testify** - Testing framework

## Project Structure

```
checkjop-be/
├── cmd/api/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   └── connection.go    # Database connection setup
│   ├── model/
│   │   ├── curriculum.go    # Curriculum data models
│   │   ├── course.go        # Course data models
│   │   ├── category.go      # Category data models
│   │   └── graduation.go    # Graduation validation models
│   ├── handler/
│   │   ├── curriculum_handler.go    # Curriculum HTTP handlers
│   │   ├── course_handler.go        # Course HTTP handlers
│   │   └── graduation_handler.go    # Graduation validation handlers
│   ├── repository/
│   │   ├── curriculum_repository.go # Curriculum data access
│   │   ├── course_repository.go     # Course data access
│   │   └── category_repository.go   # Category data access
│   ├── service/
│   │   └── graduation_service.go    # Graduation validation logic
│   └── routes/
│       └── routes.go        # API route definitions
├── tests/                   # Test files
│   ├── mocks.go            # Mock implementations
│   ├── graduation_service_test.go       # Comprehensive tests
│   └── graduation_service_simple_test.go # Basic tests
├── migrations/              # Database migrations
├── docker-compose.yml       # Docker services
├── Dockerfile              # Container build
├── Makefile                # Build commands
└── API_DOCUMENTATION.md     # API documentation
```

## Environment Variables

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=checkjop
DB_SSLMODE=disable
PORT=8080
```

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL
- Docker (optional)

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd checkjop-be
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/api/main.go
```

### Using Docker

1. Start services:
```bash
docker-compose up -d
```

2. Build and run:
```bash
make build
make run
```

## API Endpoints

All API endpoints are prefixed with `/api/v1`

### Curriculum Management
- `POST /api/v1/curricula/` - Create curriculum
- `GET /api/v1/curricula/` - Get all curricula
- `GET /api/v1/curricula/:id` - Get curriculum by ID
- `GET /api/v1/curricula/name/:name` - Get curriculum by name
- `GET /api/v1/curricula/year/:year` - Get active curricula by year
- `PUT /api/v1/curricula/:id` - Update curriculum
- `DELETE /api/v1/curricula/:id` - Delete curriculum
- `GET /api/v1/curricula/allwithout` - Get all curricula without categories and courses

### Category Management
- `POST /api/v1/categories/` - Create category
- `GET /api/v1/categories/` - Get all categories
- `GET /api/v1/categories/:id` - Get category by ID
- `GET /api/v1/categories/curriculum/:curriculum_id` - Get categories by curriculum
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Course Management  
- `POST /api/v1/courses/` - Create course
- `GET /api/v1/courses/` - Get all courses
- `GET /api/v1/courses/:id` - Get course by ID
- `GET /api/v1/courses/code/:code` - Get course by course code
- `GET /api/v1/courses/curriculum/:curriculum_id` - Get courses by curriculum
- `GET /api/v1/courses/test-relationships` - Test relationship functionality
- `PUT /api/v1/courses/:id` - Update course
- `DELETE /api/v1/courses/:id` - Delete course

### Graduation Validation
- `POST /api/v1/graduation/check` - Check overall graduation eligibility
- `POST /api/v1/graduation/check/name` - Check graduation by curriculum name
- `POST /api/v1/graduation/categories` - Check category credit requirements
- `POST /api/v1/graduation/prerequisites` - Validate prerequisite/corequisite rules
- `POST /api/v1/graduation/credit-limits` - Validate credit limits per term

### Set Default Management
- `POST /api/v1/set-defaults/` - Create set default
- `GET /api/v1/set-defaults/` - Get all set defaults
- `GET /api/v1/set-defaults/:id` - Get set default by ID
- `GET /api/v1/set-defaults/curriculum/name/:name` - Get set defaults by curriculum name
- `GET /api/v1/set-defaults/curriculum/:curriculum_id` - Get set defaults by curriculum ID
- `PUT /api/v1/set-defaults/:id` - Update set default
- `DELETE /api/v1/set-defaults/:id` - Delete set default

### CSV Import
- `POST /api/v1/import/curriculum-csv` - Import curricula from CSV
- `POST /api/v1/import/category-csv` - Import categories from CSV
- `POST /api/v1/import/course-csv` - Import courses from CSV (with upsert)
- `POST /api/v1/import/set-default-csv` - Import set defaults from CSV

## Testing

### Running Tests

Run all tests:
```bash
go test ./tests -v
```

Run simple tests only:
```bash
go test ./tests -run TestValidatePrerequisites_Simple -v
go test ./tests -run TestValidateCreditLimits_Simple -v
```

Run comprehensive tests only:
```bash
go test ./tests -run "^((?!Simple).)*$" -v
```

### Test Coverage

The project includes comprehensive test coverage for:
- ✅ Prerequisite validation with OR/AND logic
- ✅ Transitive prerequisite chains
- ✅ Corequisite validation
- ✅ Credit limit enforcement (22 regular, 10 summer)
- ✅ Category requirement checking
- ✅ Graduation eligibility validation

## Development

The application follows clean architecture principles with clear separation of concerns:

- **Handlers** - HTTP request/response handling
- **Services** - Business logic (graduation validation)
- **Repositories** - Data access layer
- **Models** - Data structures and entities
- **Routes** - API route definitions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request