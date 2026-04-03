# Pull Request

## 📝 Description  
Brief summary of changes and reason. Link related issues/PRs.

- **Related:** Closes/Connects #[issue number] (if any)  
- **Type:** 🐛 Bug fix / ✨ Feature / ♻️ Refactor / 📚 Docs / 🔧 Config / 🎨 UI/UX
- **Scope:** Frontend / Backend / Docker / CI/CD / Documentation

## 🔄 Changes Made
List key changes with clear descriptions:
- **Added:** New features or components
- **Changed:** Modified existing functionality  
- **Fixed:** Bug fixes or issues resolved
- **Removed:** Deprecated or unused code
- **Updated:** Dependencies, configurations, or documentation

## 🧪 Testing Instructions
Steps to verify the changes:

### Docker Testing (Recommended)
```bash
# Clone and test the PR branch
git checkout [branch-name]
make reset
# Visit http://localhost:3000
```

### Local Testing
```bash
yarn install
yarn lint
yarn build
yarn dev
```

### Specific Test Cases
1. Navigate to [specific page/feature]
2. Perform [specific action]
3. Verify [expected result]

## 📱 Screenshots/Demo
Attach screenshots, GIFs, or screen recordings to show:
- UI changes (before/after)
- New features in action
- Bug fixes demonstrated

## 🛠️ Technical Details
- **Framework:** Next.js changes (if any)
- **Styling:** Tailwind/shadcn/ui components used
- **Docker:** Dockerfile or build changes
- **Config:** Configuration file updates
- **Dependencies:** New packages added/updated

## ✅ Testing Checklist
- [ ] **Local development** - `make dev` works correctly
- [ ] **Docker build** - `make reset` builds and runs successfully  
- [ ] **Linting** - `make lint` passes without errors
- [ ] **Type checking** - TypeScript compiles without errors
- [ ] **Responsive design** - Works on mobile/tablet/desktop
- [ ] **Cross-browser** - Tested on Chrome/Firefox/Safari (if UI changes)
- [ ] **Performance** - No significant performance degradation

## 🔍 Code Quality Checklist
- [ ] **Code style** - Follows project conventions
- [ ] **Component structure** - Proper component organization
- [ ] **Accessibility** - ARIA labels and keyboard navigation (if applicable)
- [ ] **Error handling** - Proper error states and boundaries
- [ ] **Documentation** - README updated if needed
- [ ] **Clean code** - No console.logs or commented code left behind

## 🚀 Deployment Checklist
- [ ] **Environment variables** - No hardcoded values
- [ ] **Docker compatibility** - Works in containerized environment
- [ ] **Production ready** - Optimized for production build
- [ ] **Breaking changes** - None or properly documented

## 📋 Reviewer Notes
Additional context for reviewers:
- **Focus areas:** Specific areas that need careful review
- **Concerns:** Any potential issues or trade-offs
- **Future work:** Related tasks or follow-up items
- **Dependencies:** Other PRs or external dependencies

## 🔗 Related Links
- Design mockups: [link]
- Issue tracker: [link]
- Documentation: [link]
- Staging deployment: [link]

---

**Ready for review!** 🙏 
Please test both locally and with Docker to ensure everything works as expected.
