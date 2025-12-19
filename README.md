# Whatomate

A modern WhatsApp Business Platform built with Go (Fastglue) and Vue.js (shadcn-vue).

## Features

- **Multi-tenant Architecture**: Support multiple organizations with isolated data
- **WhatsApp Cloud API Integration**: Connect with Meta's WhatsApp Business API
- **Real-time Chat**: Live messaging with WebSocket support
- **Template Management**: Create and manage message templates
- **Bulk Messaging**: Send campaigns to multiple contacts
- **Chatbot Automation**:
  - Keyword-based auto-replies
  - Conversation flows with branching logic
  - AI-powered responses (OpenAI, Anthropic, Google)
  - Agent transfer support
- **Analytics Dashboard**: Track messages, engagement, and campaign performance

## Tech Stack

### Backend
- **Go 1.21+** with [Fastglue](https://github.com/zerodha/fastglue) (fasthttp-based HTTP framework)
- **PostgreSQL** for data storage with GORM v2
- **Redis** for caching and pub/sub
- **JWT** for authentication

### Frontend
- **Vue 3** with Composition API and TypeScript
- **Vite** for build tooling
- **shadcn-vue** / Radix Vue for UI components
- **TailwindCSS** for styling
- **Pinia** for state management
- **Vue Query** for server state

## Project Structure

```
whatomate/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connections
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   └── services/         # Business logic
├── docker/               # Docker configuration
├── frontend/             # Vue.js frontend
│   ├── src/
│   │   ├── components/   # UI components
│   │   ├── views/        # Page views
│   │   ├── stores/       # Pinia stores
│   │   ├── services/     # API services
│   │   └── lib/          # Utilities
│   └── ...
├── config.example.toml   # Example configuration
├── Makefile              # Build commands
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 6+
- Docker & Docker Compose (optional)

### Development Setup

1. **Clone and configure**:
   ```bash
   cd whatomate
   cp config.example.toml config.toml
   # Edit config.toml with your settings
   ```

2. **Start backend**:
   ```bash
   # Install dependencies
   go mod download

   # Run database migrations
   make migrate

   # Start the server
   make run
   ```

3. **Start frontend**:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### Docker Setup

```bash
cd docker
docker compose up -d
```

Access the application at `http://localhost:3000`

## Configuration

Copy `config.example.toml` to `config.toml` and update the values:

```toml
[app]
name = "Whatomate"
environment = "development"
debug = true

[server]
host = "0.0.0.0"
port = 8080

[database]
host = "localhost"
port = 5432
user = "whatomate"
password = "your-password"
name = "whatomate"
ssl_mode = "disable"

[redis]
host = "localhost"
port = 6379
password = ""
db = 0

[jwt]
secret = "your-jwt-secret"
access_token_expiry = "15m"
refresh_token_expiry = "7d"

[whatsapp]
api_version = "v18.0"
webhook_verify_token = "your-webhook-verify-token"

[ai]
openai_api_key = ""
anthropic_api_key = ""
google_api_key = ""
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh token
- `POST /api/auth/logout` - Logout

### WhatsApp Accounts
- `GET /api/accounts` - List accounts
- `POST /api/accounts` - Create account
- `PUT /api/accounts/:id` - Update account
- `DELETE /api/accounts/:id` - Delete account

### Contacts
- `GET /api/contacts` - List contacts
- `POST /api/contacts` - Create contact
- `GET /api/contacts/:id/messages` - Get messages
- `POST /api/contacts/:id/messages` - Send message

### Templates
- `GET /api/templates` - List templates
- `POST /api/templates/sync` - Sync from Meta

### Chatbot
- `GET /api/chatbot/settings` - Get settings
- `PUT /api/chatbot/settings` - Update settings
- `GET /api/chatbot/keywords` - List keyword rules
- `GET /api/chatbot/flows` - List flows
- `GET /api/chatbot/ai-contexts` - List AI contexts

### Webhooks
- `GET /api/webhook` - Webhook verification
- `POST /api/webhook` - Receive messages

## WhatsApp Setup

1. Create a Meta Developer account at [developers.facebook.com](https://developers.facebook.com)
2. Create a new app and add the WhatsApp product
3. Get your Phone Number ID and Business Account ID
4. Generate a permanent access token
5. Configure the webhook URL to point to `/api/webhook`
6. Set the webhook verify token in your configuration

## License

MIT License
