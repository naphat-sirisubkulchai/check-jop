# CheckJop - Study Plan Management System

## 📌 Project Overview

### Description
**CheckJop** (ชื่อโปรเจกต์: `senior-plan-checking`) เป็นระบบจัดการแผนการเรียนสำหรับนักศึกษาจุฬาลงกรณ์มหาวิทยาลัย ช่วยให้นักศึกษาสามารถวางแผนการเรียน ตรวจสอบสถานะการจบการศึกษา และติดตามความคืบหน้าของหน่วยกิตได้อย่างมีประสิทธิภาพ

### Purpose & Goals
- **Planning**: วางแผนการเรียนตามภาคเรียนและปีการศึกษา
- **Validation**: ตรวจสอบ prerequisites, corequisites และข้อกำหนดต่างๆ
- **Tracking**: ติดตามหน่วยกิตและความคืบหน้าแบบ real-time
- **Graduation**: ตรวจสอบคุณสมบัติการจบการศึกษา
- **Visualization**: แสดงผลแผนการเรียนในรูปแบบกราฟและตาราง

### Target Users
- นักศึกษาจุฬาลงกรณ์มหาวิทยาลัยที่ต้องการวางแผนการเรียน
- อาจารย์ที่ปรึกษาที่ต้องการช่วยเหลือนักศึกษา
- เจ้าหน้าที่ทะเบียนที่ต้องการตรวจสอบข้อมูลนักศึกษา

---

## 🛠️ Tech Stack

### Frontend Framework
- **Next.js** `15.4.1` - React framework with App Router
- **React** `19.1.0` - UI library
- **TypeScript** `^5` - Type safety

### UI Libraries
- **shadcn/ui** `^2.10.0` - Component library built on Radix UI
- **Radix UI** - Accessible component primitives
  - `@radix-ui/react-dialog`
  - `@radix-ui/react-select`
  - `@radix-ui/react-tabs`
  - `@radix-ui/react-progress`
  - And more...
- **Tailwind CSS** `^4` - Utility-first CSS framework
- **Lucide React** `^0.525.0` - Icon library
- **next-themes** `^0.4.6` - Theme management

### State Management
- **Zustand** `^5.0.7` - Lightweight state management
- **js-cookie** `^3.0.5` - Cookie management for persistence

### Visualization Tools
- **@xyflow/react** `^12.8.4` - Flow/graph visualization
- **@dagrejs/dagre** `^1.1.5` - Graph layout algorithm

### HTTP & Data
- **Axios** `^1.11.0` - HTTP client
- **PapaParse** `^5.5.3` - CSV parsing

### Other Dependencies
- **sonner** `^2.0.7` - Toast notifications
- **class-variance-authority** `^0.7.1` - CSS variant management
- **clsx** `^2.1.1` - Conditional className utility
- **tailwind-merge** `^3.3.1` - Merge Tailwind classes

---

## 📁 Project Structure

```
checkjop-fe/
├── src/
│   ├── app/                      # Next.js App Router pages
│   │   ├── page.tsx             # Landing page (/)
│   │   ├── setup/               # Setup flow (/setup)
│   │   ├── home/                # Main planning interface (/home)
│   │   ├── calculate/           # Graduation results (/calculate)
│   │   ├── admin/               # Admin panel (/admin)
│   │   │   └── curriculum/      # Curriculum management
│   │   ├── portal/              # Portal page
│   │   └── testing/             # Testing page
│   │
│   ├── api/                     # API clients & services
│   │   ├── apiClient.ts         # Axios instance configuration
│   │   ├── courseApi.ts         # Course-related API calls
│   │   └── gradApi.ts           # Graduation checking API calls
│   │
│   ├── components/              # Reusable components
│   │   ├── ui/                  # shadcn/ui components
│   │   ├── Header.tsx           # App header
│   │   └── StudyPlanInitializer.tsx
│   │
│   ├── graph/                   # Graph visualization components
│   │   ├── components/          # Graph-specific components
│   │   ├── hooks/               # Graph-related hooks
│   │   └── utils/               # Graph utilities
│   │
│   ├── store/                   # State management
│   │   └── appStore.ts          # Zustand store (unified app state)
│   │
│   ├── types/                   # TypeScript type definitions
│   │   └── index.ts             # Core types (Course, Plan, Curriculum, etc.)
│   │
│   └── utils/                   # Utility functions
│
├── docs/                        # Documentation
│   ├── API_USAGE_GUIDE.md      # Backend API documentation
│   ├── PROJECT_OVERVIEW.md     # This file
│   └── studyplan-complete.json # Sample study plan data
│
├── public/                      # Static assets
├── package.json                 # Dependencies & scripts
├── tsconfig.json               # TypeScript configuration
├── tailwind.config.ts          # Tailwind CSS configuration
└── next.config.js              # Next.js configuration
```

### Key Directories

- **`src/app/`**: Next.js pages using App Router structure
- **`src/api/`**: API integration layer
- **`src/components/`**: Reusable UI components
- **`src/graph/`**: Flow diagram visualization logic
- **`src/store/`**: Global state management with Zustand
- **`src/types/`**: TypeScript type definitions
- **`src/utils/`**: Helper functions and utilities
- **`docs/`**: Project documentation

---

## ✨ Features

### 1. **Admin Authentication**
- Simple login for admin access (username/password from environment variables)
- Protected admin routes for curriculum management
- No authentication required for student-facing features

### 2. **Setup Flow** (`/setup`)
- Select curriculum from available programs
- Map academic years (Year 1-4) to actual years (พ.ศ.)
- Auto-populate subsequent years when Year 1 is selected
- Validation before proceeding
- Reset confirmation dialog if changing settings with existing plan

### 3. **Study Plan Management** (`/home`)
- Add/remove/edit courses in study plan
- Organize courses by:
  - Academic Year (ปีการศึกษา)
  - Year of Study (ชั้นปี 1-4)
  - Semester (1, 2, 3)
- Track credits per course and total credits
- Category tracking (หมวดวิชา)
- Manual course entry for free electives

### 4. **Graduation Checking** (`/calculate`)
- **GPAX Calculation**: Automatic GPA calculation
- **Credit Validation**: Check if total credits meet requirements
- **Category Requirements**: Verify each category's credit requirements
- **Prerequisite Validation**: Detect missing or incorrectly sequenced prerequisites
- **Corequisite Validation**: Check courses that must be taken together
- **Credit Limit Violations**: Detect over-enrollment (22 for normal semester, 10 for summer)
- **Missing Courses**: Identify required courses not yet taken
- **C.F. Exemptions**: Support for Consent of Faculty exemptions

### 5. **Data Persistence**
- **Cookie Storage**: Auto-save study plan to cookies
- **Export/Import**: Download study plan as JSON
- **Auto-load**: Restore previous session on return

### 6. **Admin Features** (`/admin`)
- Curriculum management
- CSV import for courses
- Bulk data operations
- Preview before import

### 7. **Visualization**
- Flow diagram of course dependencies
- Prerequisites/corequisites graph
- Interactive course network

---

## 🔌 API Integration

### Backend Base URL
```typescript
const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";
```

### Key Endpoints

#### 1. **Curriculum APIs**
```typescript
GET  /curricula/                      // Get all curricula (without details)
GET  /curricula/name/:name            // Get curriculum by name (with courses & categories)
```

#### 2. **Course APIs**
```typescript
GET  /courses/curriculum/:id          // Get courses by curriculum ID
GET  /courses/code/:code/cf-option    // Check C.F. exemption option
POST /import/course-csv-with-year     // Import course CSV with year
```

#### 3. **Graduation Checking APIs**
```typescript
POST /graduation/check/name           // Check graduation eligibility
```

### API Client Setup
Located at `src/api/apiClient.ts`:
- Axios instance with base URL and 10s timeout
- Request interceptor for content-type handling
- Supports both JSON and FormData

### Data Models

#### Request Model (Graduation Check)
```typescript
{
  name_th: string;              // Curriculum name (Thai)
  admission_year: number;       // Year of admission
  courses: Array<{
    course_code: string;
    year: number;
    semester: number;
    grade?: string;
    credits: number;
    category_name?: string;
  }>;
  manual_credits?: {            // Manual credit input
    [category_name: string]: number;
  };
  exemptions?: string[];        // C.F. exemptions
}
```

#### Response Model (Graduation Result)
```typescript
{
  can_graduate: boolean;
  gpax: number;
  total_credits: number;
  required_credits: number;
  category_results: CategoryResult[];
  missing_courses: string[];
  prerequisite_violations: PrerequisiteViolation[];
  credit_limit_violations: CreditLimitViolation[];
}
```

For detailed API documentation, see [`docs/API_USAGE_GUIDE.md`](./API_USAGE_GUIDE.md).

---

## 💻 Development

### Prerequisites
- **Node.js** 20.x or higher
- **npm** or **yarn** or **bun**
- **Backend API** running on `http://localhost:8080` (or configured endpoint)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd checkjop-fe
   ```

2. **Install dependencies**
   ```bash
   npm install
   # or
   yarn install
   # or
   bun install
   ```

3. **Set up environment variables**
   Create a `.env.local` file (see `.env.example`):
   ```env
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
   NEXT_PUBLIC_ADMIN_USERNAME=your_admin_username
   NEXT_PUBLIC_ADMIN_PASSWORD=your_admin_password
   ```

### Running Locally

```bash
# Development mode (with hot reload)
npm run dev

# Build for production
npm run build

# Start production server
npm run start

# Run linter
npm run lint
```

The app will be available at `http://localhost:3000`

### Available Scripts

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run start` | Start production server |
| `npm run lint` | Run ESLint |

---

## 🚀 Deployment

### Build Process

1. **Environment Variables**
   - Set `NEXT_PUBLIC_API_BASE_URL` to production backend URL
   - Example: `https://api.checkjop.com/api/v1`

2. **Build the application**
   ```bash
   npm run build
   ```

3. **Start production server**
   ```bash
   npm start
   ```

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `NEXT_PUBLIC_API_BASE_URL` | Backend API base URL | No | `http://localhost:8080/api/v1` |
| `NEXT_PUBLIC_ADMIN_USERNAME` | Admin login username | Yes (for admin access) | - |
| `NEXT_PUBLIC_ADMIN_PASSWORD` | Admin login password | Yes (for admin access) | - |

### Deployment Platforms

This Next.js application can be deployed to:
- **Vercel** (Recommended - zero configuration)
- **Netlify**
- **Docker** containers
- Any **Node.js** hosting service

#### Vercel Deployment
1. Connect your GitHub repository to Vercel
2. Set environment variables in Vercel dashboard
3. Deploy automatically on push to main branch

---

## 📝 Notes & Conventions

### Coding Conventions

#### File Naming
- **Components**: PascalCase (e.g., `Header.tsx`, `CourseCard.tsx`)
- **Pages**: lowercase with hyphen (e.g., `page.tsx`, `setup/page.tsx`)
- **Utils**: camelCase (e.g., `exportImport.ts`)
- **Types**: camelCase (e.g., `index.ts`)

#### Component Structure
```typescript
// Imports
import { ... } from '...';

// Types (if needed)
type Props = { ... };

// Component
export default function ComponentName({ props }: Props) {
  // Hooks
  // State
  // Effects
  // Handlers
  // Render
}
```

#### State Management
- Use **Zustand** for global state (curriculum, study plan, etc.)
- Use **React useState** for local component state
- Use **cookies** for persistence

### Best Practices

1. **Type Safety**: Always define TypeScript types for props and data
2. **Error Handling**: Wrap API calls in try-catch blocks
3. **Loading States**: Show loading indicators during async operations
4. **User Feedback**: Use toast notifications (sonner) for success/error messages
5. **Accessibility**: Use semantic HTML and ARIA labels
6. **Responsive Design**: Ensure mobile-friendly layouts

### Code Organization

- Keep components small and focused (single responsibility)
- Extract reusable logic into custom hooks
- Use barrel exports for cleaner imports
- Group related files in feature folders

### Styling

- Use **Tailwind CSS** utility classes
- Follow **mobile-first** approach
- Use **shadcn/ui** components when possible
- Custom colors defined in `tailwind.config.ts`:
  - `chula-active`: Primary pink color
  - `chula-soft`: Light pink background

---

## Known Issues & Limitations

### Current Limitations

1. **Admin Authentication**: Simple username/password authentication (not production-ready)
2. **Graph Performance**: Large curricula (>100 courses) may cause slow rendering
3. **Browser Support**: Optimized for modern browsers (Chrome, Firefox, Safari, Edge)
4. **Mobile UX**: Some admin features are desktop-only

### Future Improvements

- [ ] Enhanced authentication system (JWT, OAuth, etc.)
- [ ] Offline support with service workers
- [ ] Advanced filtering and search in course selection
- [ ] Undo/redo functionality for study plan changes
- [ ] Multi-language support (Thai/English toggle)
- [ ] PDF export of study plan
- [ ] Share study plan via URL
- [ ] Course recommendation based on prerequisites

---

## 🔗 Related Documentation

- [API Usage Guide](./API_USAGE_GUIDE.md) - Backend API documentation
- [Sample Study Plan](./studyplan-complete.json) - Example study plan data
- [Next.js Documentation](https://nextjs.org/docs)
- [shadcn/ui Documentation](https://ui.shadcn.com)
- [Zustand Documentation](https://zustand-demo.pmnd.rs)

---

## 📄 License

© 2024 CheckJop. All rights reserved.

---

**Last Updated**: 2024-03-25
**Version**: 0.1.0
**Maintainers**: CheckJop Development Team
