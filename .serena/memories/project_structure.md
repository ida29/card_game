# Project Structure

## Root Directory
```
/
├── backend/          # Go backend server
├── frontend/         # Vue.js frontend application
├── .github/          # GitHub configuration
├── .serena/          # Serena MCP tool configuration
├── .claude/          # Claude AI configuration
├── README.md         # Project documentation
└── --context         # Context file
```

## Backend Structure
```
backend/
├── cmd/server/main.go     # Entry point
├── internal/              # Private application code
│   ├── database/          # Database connection
│   ├── models/            # Data models (Card, Deck, Game)
│   ├── services/          # Business logic layer
│   ├── handlers/          # HTTP handlers
│   ├── effects/           # Card effect implementations
│   ├── game/              # Game logic (events, targeting, etc.)
│   └── utils/             # Utility functions (card loaders)
├── data/                  # JSON data files
├── temp_card_images/      # Card image storage
├── Makefile              # Build commands
├── go.mod & go.sum       # Go dependencies
└── .air.toml             # Hot reload config
```

## Frontend Structure
```
frontend/
├── src/
│   ├── components/       # Vue components
│   │   └── game/        # Game-specific components
│   ├── views/           # Page components
│   ├── stores/          # Pinia state stores
│   ├── types/           # TypeScript type definitions
│   ├── router/          # Vue Router configuration
│   └── App.vue          # Root component
├── public/              # Static assets
├── package.json         # Node dependencies
├── vite.config.ts       # Vite configuration
├── tsconfig.json        # TypeScript configuration
└── tailwind.config.js   # TailwindCSS configuration
```

## Key Files
- API endpoints: `backend/cmd/server/main.go`
- Card/Deck models: `backend/internal/models/`
- Game components: `frontend/src/components/game/`
- State management: `frontend/src/stores/`