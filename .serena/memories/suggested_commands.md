# Development Commands

## Backend (Go)
```bash
cd backend

# Development with hot reload
make dev          # or: air (if installed)

# Run without hot reload
make run          # or: go run cmd/server/main.go

# Build binary
make build        # outputs to bin/server

# Install dependencies
go mod tidy

# Clean temporary files
make clean
```

## Frontend (Vue.js)
```bash
cd frontend

# Development server
npm run dev       # Runs Vite dev server

# Build for production
npm run build     # TypeScript check + Vite build

# Preview production build
npm run preview

# Install dependencies
npm install
```

## System Utilities (macOS/Darwin)
- `git` - Version control
- `ls` - List files/directories
- `cd` - Change directory
- `rg` (ripgrep) - Fast content search
- `find` - Find files by name/pattern

## Testing & Linting
Currently no test suite or linting commands are configured in the project. Consider adding:
- Go: `go test`, `golangci-lint`
- Vue: `npm run lint`, `npm run test`