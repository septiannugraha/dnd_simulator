# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a D&D (Dungeons & Dragons) platform MVP that enables multiplayer tabletop RPG sessions with AI-powered Dungeon Master capabilities. The platform supports campaign management, character creation, real-time multiplayer sessions, and AI narrative generation.

## Technology Stack (Planned)

- **Backend**: Go with Gin framework
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Real-time Communication**: WebSockets
- **AI Integration**: OpenAI/Claude API
- **Frontend**: Not specified in roadmap

## Core Architecture

### Backend Services
- **Authentication System**: JWT-based user registration/login
- **Campaign Management**: DM-controlled campaigns with invitation system
- **Character Management**: D&D 5e compliant character creation with auto-calculated stats
- **Session Management**: Real-time game sessions with turn order tracking
- **WebSocket Hub**: Real-time message routing and state synchronization
- **AI Service**: Context-aware narrative generation using campaign and character data

### Key Data Models
```go
// Core entities from roadmap
User, Character, Campaign, GameSession

// AI Context Structure
type AIContext struct {
    Campaign      Campaign
    CurrentScene  string
    Characters    []Character
    RecentEvents  []GameEvent
    PlayerAction  PlayerAction
}
```

### API Structure (Planned)
```
Authentication:
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/me

Campaigns:
POST /api/campaigns
GET  /api/campaigns
GET  /api/campaigns/:id
PUT  /api/campaigns/:id
POST /api/campaigns/:id/invite
POST /api/campaigns/:id/join

Characters:
POST /api/characters
GET  /api/characters
GET  /api/characters/:id
PUT  /api/characters/:id
POST /api/characters/:id/assign

Sessions:
POST /api/sessions
GET  /api/sessions/:id
POST /api/sessions/:id/join
POST /api/sessions/:id/start
PUT  /api/sessions/:id/turn
```

### WebSocket Message Types
```
- player_join_session
- player_leave_session
- chat_message
- character_update
- dice_roll
```

## Development Phases

The project follows a 9-week MVP roadmap:
1. **Foundation** (Weeks 1-2): Core infrastructure and authentication
2. **Campaign Foundation** (Weeks 2-3): Campaign management system
3. **Character Management** (Weeks 3-4): D&D 5e character creation
4. **Real-Time Communication** (Weeks 4-5): WebSocket implementation
5. **Game Session Management** (Weeks 5-6): Structured sessions with turn management
6. **AI DM Integration** (Weeks 6-8): AI-powered narrative responses
7. **MVP Polish** (Weeks 8-9): Testing and bug fixes

## Commands

### Development
```bash
# Start development environment (MongoDB only)
make dev-up

# Run app locally with dev database
make dev-run

# Stop development environment
make dev-down

# View database via Mongo Express UI
# http://localhost:8081 (after dev-up)
```

### Production
```bash
# Start full production environment
make prod-up

# Stop production environment
make prod-down

# View production logs
make prod-logs
```

### Database
```bash
# Connect to MongoDB shell
make db-shell

# Clean up all Docker containers and volumes
make clean
```

### Building
```bash
# Build Go application
make build

# Build Docker image
make docker-build

# Run tests
make test
```

## Key Requirements

### Performance Targets
- Support 6+ concurrent players per session
- AI response time under 10 seconds
- Session persistence through disconnections

### Core Features for MVP
- Campaign management with DM permissions
- D&D 5e character creation with auto-calculations
- Real-time multiplayer sessions
- AI Dungeon Master with campaign context
- Text chat and dice rolling
- Turn order management

### Out of Scope for MVP
- 3D Virtual Tabletop
- Voice/Video Chat
- Advanced inventory management
- Complex combat automation
- Mobile apps (PWA sufficient)