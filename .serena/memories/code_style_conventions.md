# Code Style and Conventions

## Go Backend
- **Package Structure**: Domain-driven design under `internal/`
  - `models/` - Data structures and domain models
  - `services/` - Business logic layer
  - `handlers/` - HTTP request handlers
  - `database/` - Database connection and utilities
  - `utils/` - Helper functions
  - `effects/` - Game effect implementations
  - `game/` - Game logic components
- **Naming**: 
  - Services use constructor pattern: `NewCardService()`
  - Exported types/functions use PascalCase
  - Method receivers use single letter: `(s *CardService)`
- **Error Handling**: Always return errors as second value
- **Database**: Use GORM for database operations

## Vue Frontend
- **Component Structure**: Single File Components (SFC) with `<template>`, `<script>`, `<style>`
- **TypeScript**: Strict typing with Vue 3 Composition API
- **Styling**: 
  - TailwindCSS utility classes
  - Scoped styles when needed
  - Responsive design with size classes
- **State Management**: Pinia stores for global state
- **API Calls**: Axios with centralized API configuration
- **Component Naming**: PascalCase for components (e.g., `GameCard.vue`)
- **Props**: Typed with TypeScript interfaces

## General Conventions
- No test files currently exist in the project
- Configuration through environment variables
- RESTful API design with JSON responses
- Git for version control