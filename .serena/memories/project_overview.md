# めめめのくらげ TCG Project Overview

## Project Purpose
A 1v1 trading card game (TCG) featuring characters ("friends") from "めめめのくらげ" (Mememe no Kurage). Players battle using decks of 50 cards with the goal of placing 7 cards in the opponent's negative energy area or depleting their deck.

## Tech Stack
- **Backend**: Go 1.21 with Gin web framework
  - Database: SQLite with GORM ORM
  - CORS support via gin-contrib/cors
  - Environment config via godotenv
- **Frontend**: Vue.js 3 with TypeScript
  - Vite build tool
  - Vue Router for navigation
  - Pinia for state management
  - Axios for API calls
  - TailwindCSS for styling
  - Phaser.js for game board (planned)

## Key Features
- Card collection system
- Deck building with validation rules
- Battle system with multiple game phases
- Three card types: ふれんど (Friend), サポート (Support), フィールド (Field)

## API Structure
RESTful API at `/api/v1/` with endpoints for:
- Card management (CRUD, search, filter by type/color)
- Deck management (CRUD, validation)
- Game state management (planned)