# Task Completion Checklist

When completing a coding task, ensure:

## Backend (Go)
1. **Code Compiles**: Run `go build` or `make build`
2. **Dependencies Updated**: Run `go mod tidy` if new packages added
3. **Server Runs**: Test with `make run` or `make dev`
4. **API Endpoints Work**: Test new/modified endpoints manually

## Frontend (Vue/TypeScript)
1. **TypeScript Compiles**: Run `npm run build` to check for type errors
2. **Dev Server Works**: Run `npm run dev` and check for runtime errors
3. **UI Renders Correctly**: Visually verify components
4. **API Integration**: Ensure frontend correctly calls backend APIs

## General Checks
1. **No Console Errors**: Check browser console and server logs
2. **Code Follows Conventions**: Match existing code style
3. **No Hardcoded Values**: Use environment variables or config
4. **Security**: No exposed secrets, proper input validation

## Currently Missing (Consider Adding)
- Automated tests (unit/integration)
- Linting tools (golangci-lint, ESLint)
- Code formatting (gofmt, Prettier)
- Pre-commit hooks

Note: Since the project lacks configured linting/testing commands, manual verification is essential.