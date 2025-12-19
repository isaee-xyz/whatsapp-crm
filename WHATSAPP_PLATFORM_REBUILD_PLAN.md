# Whatomate - WhatsApp Business Platform

## Executive Summary

**Whatomate** is a modern WhatsApp Business Platform rebuilt from the ground up, replacing the three Frappe WhatsApp applications (`frappe_whatsapp`, `whatsapp_chat`, `frappe_whatsapp_chatbot`). Built with **Go (Fastglue)** for the backend and **Vue.js (shadcn-vue)** for the frontend.

---

## Table of Contents

1. [Current System Analysis](#1-current-system-analysis)
2. [Architecture Overview](#2-architecture-overview)
3. [Technology Stack](#3-technology-stack)
4. [Database Schema](#4-database-schema)
5. [Backend Implementation Plan (Go)](#5-backend-implementation-plan-go)
6. [Frontend Implementation Plan (Vue.js)](#6-frontend-implementation-plan-vuejs)
7. [API Design](#7-api-design)
8. [Real-time Communication](#8-real-time-communication)
9. [AI Integration](#9-ai-integration)
10. [Development Phases](#10-development-phases)
11. [Deployment Strategy](#11-deployment-strategy)
12. [Migration Strategy](#12-migration-strategy)

---

## 1. Current System Analysis

### 1.1 Application Overview

| App | Purpose | Key Features |
|-----|---------|--------------|
| **frappe_whatsapp** | Core WhatsApp Business API integration | Message sending/receiving, Templates, Flows, Bulk messaging, Notifications, Multi-account |
| **whatsapp_chat** | Real-time chat UI | Agent chat interface, Contact management, Socket.IO notifications |
| **frappe_whatsapp_chatbot** | Automated bot responses | Keyword matching, Conversation flows, AI responses, Agent transfer |

### 1.2 Current Data Models

**Core Entities (from frappe_whatsapp):**
- WhatsApp Account (credentials, phone_id, business_id)
- WhatsApp Message (all incoming/outgoing messages)
- WhatsApp Templates (Meta-approved templates)
- WhatsApp Flow (interactive forms)
- WhatsApp Notification (automation rules)
- Bulk WhatsApp Message (batch sending)
- WhatsApp Profiles (contact profiles)

**Chat Entities (from whatsapp_chat):**
- WhatsApp Contact (chat room abstraction)

**Chatbot Entities (from frappe_whatsapp_chatbot):**
- WhatsApp Chatbot (settings)
- WhatsApp Chatbot Flow (conversation flows)
- WhatsApp Flow Step (flow steps)
- WhatsApp Keyword Reply (keyword rules)
- WhatsApp Chatbot Session (active conversations)
- WhatsApp AI Context (AI knowledge base)
- WhatsApp Agent Transfer (handoff tracking)

### 1.3 Current Integration Points

```
Meta WhatsApp Cloud API
        ↓↑
frappe_whatsapp (webhook, message sending)
        ↓↑
whatsapp_chat (UI layer) ←→ frappe_whatsapp_chatbot (automation)
        ↓
Frappe Framework (ORM, permissions, scheduler)
```

---

## 2. Architecture Overview

### 2.1 Proposed Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Vue.js)                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────────┐│
│  │Dashboard │ │Chat UI   │ │Templates │ │Chatbot Config        ││
│  │& Reports │ │(Realtime)│ │& Flows   │ │(Keywords/Flows/AI)   ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
                              │
                    REST API + WebSocket
                              │
┌─────────────────────────────────────────────────────────────────┐
│                     Backend (Go + Fastglue)                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                      API Gateway                          │   │
│  │  (Auth, Rate Limiting, Request Validation, Logging)      │   │
│  └──────────────────────────────────────────────────────────┘   │
│                              │                                   │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────────┐   │
│  │ Messaging │ │ Templates │ │  Chatbot  │ │  Bulk Jobs    │   │
│  │  Service  │ │  Service  │ │  Service  │ │   Service     │   │
│  └───────────┘ └───────────┘ └───────────┘ └───────────────┘   │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────────┐   │
│  │  Flows    │ │   AI      │ │ Webhook   │ │ Notification  │   │
│  │  Service  │ │  Service  │ │  Handler  │ │   Service     │   │
│  └───────────┘ └───────────┘ └───────────┘ └───────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────┴───────┐   ┌─────────┴─────────┐   ┌──────┴──────┐
│  PostgreSQL   │   │      Redis        │   │  AI APIs    │
│  (Primary DB) │   │ (Cache/PubSub/Q)  │   │(OpenAI etc) │
└───────────────┘   └───────────────────┘   └─────────────┘
                              │
                    ┌─────────┴─────────┐
                    │  Meta WhatsApp    │
                    │    Cloud API      │
                    └───────────────────┘
```

### 2.2 Key Architectural Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Backend Framework | Fastglue | Built on Fasthttp, high performance, clean API, great for REST APIs |
| Database | PostgreSQL | JSONB support, strong consistency, excellent Go drivers |
| Cache/Queue | Redis | PubSub for real-time, job queue for bulk messaging |
| Real-time | WebSocket (Gorilla/Melody) | Native Go support, efficient memory usage |
| AI Integration | Direct API calls | Flexibility, no heavy SDK dependencies |
| Authentication | JWT + API Keys | Stateless, scalable |
| File Storage | S3-compatible | Scalable media handling |
| UI Components | shadcn-vue | Beautiful, accessible, customizable components with Radix Vue primitives |

---

## 3. Technology Stack

### 3.1 Backend (Go + Fastglue)

```
Core:
├── Go 1.21+
├── zerodha/fastglue (HTTP framework built on fasthttp)
├── GORM v2 (ORM with PostgreSQL)
├── sqlx (for complex queries)
└── validator/v10 (request validation)

Fastglue Features:
├── Built on valyala/fasthttp (10x faster than net/http)
├── Clean handler interface with context
├── Built-in request binding and validation
├── Middleware support (auth, logging, recovery)
├── Static file serving
└── Template rendering

Real-time:
├── gorilla/websocket
├── melody (WebSocket framework)
└── Redis PubSub

Background Jobs:
├── asynq (Redis-based job queue)
└── robfig/cron (scheduled tasks)

AI Integration:
├── sashabaranov/go-openai
├── anthropics/anthropic-sdk-go
└── google/generative-ai-go

Utilities:
├── knadh/koanf (configuration - works great with fastglue)
├── zerodha/logf (structured logging)
├── ozzo-validation (business validation)
└── golang-jwt/jwt (authentication)
```

**Fastglue Example:**
```go
package main

import (
    "github.com/zerodha/fastglue"
    "github.com/valyala/fasthttp"
)

type App struct {
    g *fastglue.Fastglue
    // services, repos, etc.
}

func main() {
    g := fastglue.New()
    app := &App{g: g}

    // Middleware
    g.Use(app.authMiddleware)
    g.Use(app.loggingMiddleware)

    // Routes
    g.GET("/api/v1/contacts", app.handleGetContacts)
    g.POST("/api/v1/messages", app.handleSendMessage)
    g.GET("/api/v1/webhook", app.handleWebhookVerify)
    g.POST("/api/v1/webhook", app.handleWebhookMessage)

    // Group routes
    chatbot := g.Group("/api/v1/chatbot")
    chatbot.GET("/settings", app.handleGetSettings)
    chatbot.PUT("/settings", app.handleUpdateSettings)
    chatbot.GET("/keywords", app.handleGetKeywords)

    fasthttp.ListenAndServe(":8080", g.Handler())
}

func (a *App) handleGetContacts(r *fastglue.Request) error {
    var req GetContactsRequest
    if err := r.Decode(&req); err != nil {
        return r.SendErrorEnvelope(fasthttp.StatusBadRequest, err.Error(), nil)
    }

    contacts, err := a.contactService.List(r.Context(), req)
    if err != nil {
        return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, err.Error(), nil)
    }

    return r.SendEnvelope(contacts)
}
```

### 3.2 Frontend (Vue.js + shadcn-vue)

```
Core:
├── Vue 3 (Composition API)
├── TypeScript
├── Vite (build tool)
└── Pinia (state management)

UI (shadcn-vue):
├── shadcn-vue (beautiful, accessible components)
├── Radix Vue (headless UI primitives)
├── TailwindCSS (styling)
├── tailwindcss-animate (animations)
├── class-variance-authority (component variants)
├── clsx + tailwind-merge (class utilities)
└── lucide-vue-next (icons)

Real-time:
├── Native WebSocket API
└── @vueuse/core (composables)

HTTP:
├── axios (API calls)
└── @tanstack/vue-query (data fetching/caching)

Forms:
├── vee-validate (form validation)
├── @vee-validate/zod (schema integration)
└── zod (schema validation - works with shadcn forms)

Additional:
├── vue-router (routing)
├── vue-i18n (internationalization)
├── @vueuse/motion (animations)
└── vue-sonner (toast notifications)
```

**shadcn-vue Components Used:**
```
├── Button, Input, Textarea, Select
├── Card, Dialog, Sheet, Dropdown Menu
├── Table, Data Table (with TanStack Table)
├── Tabs, Accordion, Collapsible
├── Avatar, Badge, Tooltip
├── Form, Label, Checkbox, Radio Group
├── Command (for search/command palette)
├── Toast (via vue-sonner)
├── Skeleton (loading states)
├── Separator, ScrollArea
└── Alert, AlertDialog
```

**shadcn-vue Setup:**
```bash
# Initialize shadcn-vue in your Vue project
npx shadcn-vue@latest init

# Add components as needed
npx shadcn-vue@latest add button
npx shadcn-vue@latest add card
npx shadcn-vue@latest add dialog
npx shadcn-vue@latest add form
npx shadcn-vue@latest add input
npx shadcn-vue@latest add table
npx shadcn-vue@latest add dropdown-menu
npx shadcn-vue@latest add command
npx shadcn-vue@latest add avatar
npx shadcn-vue@latest add badge
npx shadcn-vue@latest add toast
```

### 3.3 Infrastructure

```
Database:
├── PostgreSQL 15+
└── Redis 7+

Storage:
├── MinIO (S3-compatible)
└── Local filesystem (dev)

Deployment:
├── Docker + Docker Compose
├── Kubernetes (production)
└── Nginx (reverse proxy)

Monitoring:
├── Prometheus + Grafana
├── Sentry (error tracking)
└── ELK Stack (logging)
```

---

## 4. Database Schema

### 4.1 Core Tables

```sql
-- Organizations (multi-tenant support)
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'agent', -- admin, manager, agent
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- WhatsApp Accounts
CREATE TABLE whatsapp_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    phone_id VARCHAR(100) NOT NULL,
    business_id VARCHAR(100) NOT NULL,
    access_token TEXT NOT NULL, -- encrypted
    webhook_verify_token VARCHAR(255),
    api_version VARCHAR(20) DEFAULT 'v18.0',
    is_default_incoming BOOLEAN DEFAULT false,
    is_default_outgoing BOOLEAN DEFAULT false,
    auto_read_receipt BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(organization_id, phone_id)
);

-- Contacts (WhatsApp Profiles)
CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    phone_number VARCHAR(20) NOT NULL,
    profile_name VARCHAR(255),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    assigned_user_id UUID REFERENCES users(id),
    last_message_at TIMESTAMPTZ,
    last_message_preview TEXT,
    is_read BOOLEAN DEFAULT true,
    tags JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(organization_id, phone_number)
);

-- Messages
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    contact_id UUID REFERENCES contacts(id),
    whatsapp_message_id VARCHAR(255), -- Meta's message ID
    conversation_id VARCHAR(255),
    direction VARCHAR(10) NOT NULL, -- incoming, outgoing
    message_type VARCHAR(20) NOT NULL, -- text, image, video, audio, document, template, interactive, flow, reaction, location, contact
    content TEXT,
    media_url TEXT,
    media_mime_type VARCHAR(100),
    media_filename VARCHAR(255),
    template_name VARCHAR(255),
    template_params JSONB,
    interactive_data JSONB, -- buttons, lists, flows
    flow_response JSONB,
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, delivered, read, failed
    error_message TEXT,
    is_reply BOOLEAN DEFAULT false,
    reply_to_message_id UUID REFERENCES messages(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Templates
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    meta_template_id VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    language VARCHAR(10) NOT NULL,
    category VARCHAR(50), -- MARKETING, UTILITY, AUTHENTICATION
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    header_type VARCHAR(20), -- TEXT, IMAGE, DOCUMENT, VIDEO
    header_content TEXT,
    body_content TEXT NOT NULL,
    footer_content TEXT,
    buttons JSONB DEFAULT '[]',
    sample_values JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(organization_id, name, language)
);

-- WhatsApp Flows (Meta's interactive forms)
CREATE TABLE whatsapp_flows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    meta_flow_id VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, PUBLISHED, DEPRECATED, BLOCKED
    category VARCHAR(50),
    json_version VARCHAR(10) DEFAULT '6.0',
    flow_json JSONB,
    screens JSONB DEFAULT '[]',
    preview_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Notifications (Automation Rules)
CREATE TABLE notification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    trigger_type VARCHAR(50) NOT NULL, -- webhook, scheduler, api
    trigger_config JSONB NOT NULL,
    template_id UUID REFERENCES templates(id),
    field_mappings JSONB DEFAULT '{}',
    conditions JSONB DEFAULT '{}',
    attachment_config JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Bulk Messages
CREATE TABLE bulk_message_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    name VARCHAR(255) NOT NULL,
    template_id UUID REFERENCES templates(id),
    status VARCHAR(20) DEFAULT 'draft', -- draft, queued, processing, completed, failed
    total_recipients INT DEFAULT 0,
    sent_count INT DEFAULT 0,
    delivered_count INT DEFAULT 0,
    failed_count INT DEFAULT 0,
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bulk_message_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID REFERENCES bulk_message_campaigns(id) ON DELETE CASCADE,
    phone_number VARCHAR(20) NOT NULL,
    recipient_name VARCHAR(255),
    template_params JSONB DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'pending',
    message_id UUID REFERENCES messages(id),
    error_message TEXT,
    sent_at TIMESTAMPTZ
);
```

### 4.2 Chatbot Tables

```sql
-- Chatbot Settings (per organization)
CREATE TABLE chatbot_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID UNIQUE REFERENCES organizations(id),
    is_enabled BOOLEAN DEFAULT false,
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    process_all_accounts BOOLEAN DEFAULT false,
    default_response TEXT,
    business_hours_enabled BOOLEAN DEFAULT false,
    business_hours JSONB DEFAULT '[]', -- [{day, enabled, start_time, end_time}]
    out_of_hours_message TEXT,
    ai_enabled BOOLEAN DEFAULT false,
    ai_provider VARCHAR(20), -- openai, anthropic, google
    ai_api_key TEXT, -- encrypted
    ai_model VARCHAR(100),
    ai_max_tokens INT DEFAULT 500,
    ai_temperature DECIMAL(3,2) DEFAULT 0.7,
    ai_system_prompt TEXT,
    ai_include_history BOOLEAN DEFAULT true,
    ai_history_limit INT DEFAULT 4,
    session_timeout_minutes INT DEFAULT 30,
    excluded_numbers JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Keyword Rules
CREATE TABLE keyword_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 10,
    keywords TEXT[] NOT NULL,
    match_type VARCHAR(20) DEFAULT 'contains', -- exact, contains, starts_with, regex
    case_sensitive BOOLEAN DEFAULT false,
    response_type VARCHAR(20) NOT NULL, -- text, template, media, flow, script
    response_content JSONB NOT NULL, -- varies by type
    conditions TEXT, -- additional Python/JS conditions
    active_from TIMESTAMPTZ,
    active_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Chatbot Flows (conversation flows)
CREATE TABLE chatbot_flows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    description TEXT,
    trigger_keywords TEXT[],
    trigger_button_id VARCHAR(100),
    initial_message TEXT,
    initial_message_type VARCHAR(20) DEFAULT 'text',
    initial_template_id UUID REFERENCES templates(id),
    completion_message TEXT,
    on_complete_action VARCHAR(20), -- none, webhook, create_record
    completion_config JSONB, -- webhook_url, doctype, field_mapping, etc.
    timeout_message TEXT,
    cancel_keywords TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Flow Steps
CREATE TABLE chatbot_flow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID REFERENCES chatbot_flows(id) ON DELETE CASCADE,
    step_name VARCHAR(100) NOT NULL,
    step_order INT NOT NULL,
    message TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- text, template, script
    template_id UUID REFERENCES templates(id),
    input_type VARCHAR(20), -- none, text, number, email, phone, date, select, button, whatsapp_flow
    input_config JSONB, -- buttons, options, validation, whatsapp_flow_id, etc.
    validation_regex VARCHAR(255),
    validation_error TEXT,
    store_as VARCHAR(100),
    next_step VARCHAR(100),
    conditional_next JSONB, -- {"option1": "step_a", "default": "step_b"}
    skip_condition TEXT,
    retry_on_invalid BOOLEAN DEFAULT true,
    max_retries INT DEFAULT 3,
    UNIQUE(flow_id, step_name)
);

-- Chatbot Sessions
CREATE TABLE chatbot_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    contact_id UUID REFERENCES contacts(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, completed, cancelled, timeout
    current_flow_id UUID REFERENCES chatbot_flows(id),
    current_step VARCHAR(100),
    step_retries INT DEFAULT 0,
    session_data JSONB DEFAULT '{}',
    started_at TIMESTAMPTZ DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Session Messages (history)
CREATE TABLE chatbot_session_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES chatbot_sessions(id) ON DELETE CASCADE,
    direction VARCHAR(10) NOT NULL, -- incoming, outgoing
    message TEXT,
    step_name VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- AI Context
CREATE TABLE ai_contexts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 10,
    context_type VARCHAR(20) NOT NULL, -- static, query
    trigger_keywords TEXT[],
    static_content TEXT,
    query_config JSONB, -- table, fields, filters, max_results
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Agent Transfers
CREATE TABLE agent_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    contact_id UUID REFERENCES contacts(id),
    whatsapp_account_id UUID REFERENCES whatsapp_accounts(id),
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, resumed
    agent_id UUID REFERENCES users(id),
    notes TEXT,
    transferred_at TIMESTAMPTZ DEFAULT NOW(),
    resumed_at TIMESTAMPTZ,
    resumed_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_messages_contact ON messages(contact_id, created_at DESC);
CREATE INDEX idx_messages_whatsapp_id ON messages(whatsapp_message_id);
CREATE INDEX idx_messages_conversation ON messages(conversation_id);
CREATE INDEX idx_contacts_phone ON contacts(organization_id, phone_number);
CREATE INDEX idx_contacts_assigned ON contacts(assigned_user_id, is_read);
CREATE INDEX idx_sessions_phone ON chatbot_sessions(organization_id, phone_number, status);
CREATE INDEX idx_keyword_rules_priority ON keyword_rules(organization_id, is_enabled, priority DESC);
CREATE INDEX idx_agent_transfers_active ON agent_transfers(organization_id, phone_number, status);
```

---

## 5. Backend Implementation Plan (Go)

### 5.1 Project Structure

```
whatsapp-platform/
├── cmd/
│   ├── server/
│   │   └── main.go              # HTTP server entry point
│   ├── worker/
│   │   └── main.go              # Background worker entry point
│   └── migrate/
│       └── main.go              # Database migrations
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── database/
│   │   ├── postgres.go          # PostgreSQL connection
│   │   └── redis.go             # Redis connection
│   ├── models/
│   │   ├── organization.go
│   │   ├── user.go
│   │   ├── whatsapp_account.go
│   │   ├── contact.go
│   │   ├── message.go
│   │   ├── template.go
│   │   ├── flow.go
│   │   ├── notification.go
│   │   ├── bulk_campaign.go
│   │   ├── chatbot_settings.go
│   │   ├── keyword_rule.go
│   │   ├── chatbot_flow.go
│   │   ├── chatbot_session.go
│   │   ├── ai_context.go
│   │   └── agent_transfer.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── message_repo.go
│   │   ├── contact_repo.go
│   │   ├── template_repo.go
│   │   ├── session_repo.go
│   │   └── ...
│   ├── services/
│   │   ├── messaging/
│   │   │   ├── service.go       # Message sending/receiving
│   │   │   ├── sender.go        # WhatsApp API client
│   │   │   └── media.go         # Media handling
│   │   ├── templates/
│   │   │   ├── service.go       # Template CRUD
│   │   │   └── sync.go          # Meta template sync
│   │   ├── flows/
│   │   │   ├── service.go       # WhatsApp Flow management
│   │   │   └── builder.go       # Flow JSON builder
│   │   ├── bulk/
│   │   │   ├── service.go       # Bulk campaign management
│   │   │   └── processor.go     # Background processing
│   │   ├── notifications/
│   │   │   ├── service.go       # Notification rules
│   │   │   └── trigger.go       # Event triggers
│   │   ├── chatbot/
│   │   │   ├── processor.go     # Message processing pipeline
│   │   │   ├── keyword.go       # Keyword matching
│   │   │   ├── flow_engine.go   # Conversation flow engine
│   │   │   ├── session.go       # Session management
│   │   │   └── response.go      # Response building
│   │   ├── ai/
│   │   │   ├── service.go       # AI response generation
│   │   │   ├── openai.go        # OpenAI provider
│   │   │   ├── anthropic.go     # Anthropic provider
│   │   │   ├── google.go        # Google AI provider
│   │   │   └── context.go       # Context building
│   │   └── webhook/
│   │       ├── handler.go       # Webhook processing
│   │       └── parser.go        # Payload parsing
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── accounts.go
│   │   ├── contacts.go
│   │   ├── messages.go
│   │   ├── templates.go
│   │   ├── flows.go
│   │   ├── bulk.go
│   │   ├── chatbot.go
│   │   ├── webhook.go
│   │   └── websocket.go
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── org.go               # Organization context
│   │   ├── ratelimit.go         # Rate limiting
│   │   └── logging.go           # Request logging
│   ├── websocket/
│   │   ├── hub.go               # WebSocket hub
│   │   ├── client.go            # Client connection
│   │   └── events.go            # Event types
│   └── workers/
│       ├── bulk_sender.go       # Bulk message worker
│       ├── session_cleanup.go   # Session expiry worker
│       └── template_sync.go     # Template sync worker
├── pkg/
│   ├── whatsapp/
│   │   ├── client.go            # Meta API client
│   │   ├── types.go             # API types
│   │   └── errors.go            # Error handling
│   ├── encryption/
│   │   └── crypto.go            # Token encryption
│   └── validator/
│       └── validator.go         # Custom validators
├── migrations/
│   └── *.sql                    # SQL migrations
├── api/
│   └── openapi.yaml             # OpenAPI specification
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

### 5.2 Core Services Implementation

#### 5.2.1 WhatsApp API Client

```go
// pkg/whatsapp/client.go
package whatsapp

type Client struct {
    httpClient *http.Client
    baseURL    string
    token      string
    phoneID    string
    version    string
}

func NewClient(token, phoneID string, opts ...Option) *Client

// Message sending
func (c *Client) SendTextMessage(ctx context.Context, to, text string) (*SendResponse, error)
func (c *Client) SendTemplateMessage(ctx context.Context, to string, template TemplatePayload) (*SendResponse, error)
func (c *Client) SendMediaMessage(ctx context.Context, to string, media MediaPayload) (*SendResponse, error)
func (c *Client) SendInteractiveMessage(ctx context.Context, to string, interactive InteractivePayload) (*SendResponse, error)
func (c *Client) SendFlowMessage(ctx context.Context, to string, flow FlowPayload) (*SendResponse, error)
func (c *Client) MarkAsRead(ctx context.Context, messageID string) error

// Template management
func (c *Client) GetTemplates(ctx context.Context, businessID string) ([]Template, error)
func (c *Client) CreateTemplate(ctx context.Context, businessID string, template CreateTemplateRequest) (*Template, error)
func (c *Client) DeleteTemplate(ctx context.Context, businessID, name string) error

// Flow management
func (c *Client) CreateFlow(ctx context.Context, businessID string, flow CreateFlowRequest) (*Flow, error)
func (c *Client) UploadFlowJSON(ctx context.Context, flowID string, json []byte) error
func (c *Client) PublishFlow(ctx context.Context, flowID string) error
func (c *Client) DeprecateFlow(ctx context.Context, flowID string) error
func (c *Client) DeleteFlow(ctx context.Context, flowID string) error

// Media
func (c *Client) DownloadMedia(ctx context.Context, mediaID string) ([]byte, string, error)
func (c *Client) UploadMedia(ctx context.Context, phoneID string, media io.Reader, mimeType string) (string, error)
```

#### 5.2.2 Message Processing Service

```go
// internal/services/messaging/service.go
package messaging

type Service struct {
    repo       repository.MessageRepository
    contactRepo repository.ContactRepository
    waClient   *whatsapp.Client
    wsHub      *websocket.Hub
    redis      *redis.Client
}

func (s *Service) SendMessage(ctx context.Context, req SendMessageRequest) (*models.Message, error)
func (s *Service) ProcessIncoming(ctx context.Context, payload webhook.MessagePayload) error
func (s *Service) UpdateStatus(ctx context.Context, messageID, status string) error
func (s *Service) GetConversation(ctx context.Context, contactID uuid.UUID, pagination Pagination) ([]models.Message, error)
```

#### 5.2.3 Chatbot Processor

```go
// internal/services/chatbot/processor.go
package chatbot

type Processor struct {
    settings       *models.ChatbotSettings
    keywordMatcher *KeywordMatcher
    flowEngine     *FlowEngine
    sessionMgr     *SessionManager
    aiService      *ai.Service
    messenger      *messaging.Service
}

func (p *Processor) Process(ctx context.Context, message *models.Message) error {
    // 1. Check if chatbot is enabled
    // 2. Check excluded numbers
    // 3. Check agent transfer status
    // 4. Check business hours
    // 5. Check active session -> process flow
    // 6. Check keyword rules -> send response
    // 7. Check flow triggers -> start flow
    // 8. Try AI response (if enabled)
    // 9. Send default response
}
```

#### 5.2.4 Flow Engine

```go
// internal/services/chatbot/flow_engine.go
package chatbot

type FlowEngine struct {
    flowRepo    repository.ChatbotFlowRepository
    sessionRepo repository.SessionRepository
    messenger   *messaging.Service
}

func (e *FlowEngine) CheckTrigger(ctx context.Context, message, buttonPayload string) (*models.ChatbotFlow, error)
func (e *FlowEngine) StartFlow(ctx context.Context, session *models.ChatbotSession, flow *models.ChatbotFlow) error
func (e *FlowEngine) ProcessInput(ctx context.Context, session *models.ChatbotSession, input string, buttonPayload string) error
func (e *FlowEngine) ValidateInput(step *models.ChatbotFlowStep, input string) (bool, string)
func (e *FlowEngine) GetNextStep(step *models.ChatbotFlowStep, input string) string
func (e *FlowEngine) CompleteFlow(ctx context.Context, session *models.ChatbotSession) error
func (e *FlowEngine) ExecuteCompletionAction(ctx context.Context, flow *models.ChatbotFlow, data map[string]interface{}) error
```

#### 5.2.5 AI Service

```go
// internal/services/ai/service.go
package ai

type Service struct {
    settings    *models.ChatbotSettings
    contextRepo repository.AIContextRepository
    providers   map[string]Provider
}

type Provider interface {
    GenerateResponse(ctx context.Context, req GenerateRequest) (string, error)
}

func (s *Service) GenerateResponse(ctx context.Context, message string, history []Message) (string, error) {
    // 1. Build context from AI Context documents
    // 2. Format conversation history
    // 3. Call appropriate provider
    // 4. Return response
}

func (s *Service) BuildContext(ctx context.Context, message string, phoneNumber string) (string, error)
```

### 5.3 Background Workers

```go
// internal/workers/bulk_sender.go
package workers

type BulkSenderWorker struct {
    repo      repository.BulkCampaignRepository
    messenger *messaging.Service
}

func (w *BulkSenderWorker) Process(ctx context.Context, task *asynq.Task) error {
    // Process bulk message batch
    // Update progress
    // Handle failures with retry
}

// internal/workers/session_cleanup.go
type SessionCleanupWorker struct {
    repo      repository.SessionRepository
    messenger *messaging.Service
}

func (w *SessionCleanupWorker) Run(ctx context.Context) {
    // Find expired sessions
    // Send timeout messages
    // Update status to timeout
}
```

---

## 6. Frontend Implementation Plan (Vue.js)

### 6.1 Project Structure

```
frontend/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/
│   │   └── index.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── contacts.ts
│   │   ├── messages.ts
│   │   ├── templates.ts
│   │   ├── chatbot.ts
│   │   └── websocket.ts
│   ├── composables/
│   │   ├── useWebSocket.ts
│   │   ├── useNotifications.ts
│   │   ├── useInfiniteScroll.ts
│   │   └── useMediaUpload.ts
│   ├── services/
│   │   ├── api.ts              # Axios instance
│   │   ├── auth.ts
│   │   ├── accounts.ts
│   │   ├── contacts.ts
│   │   ├── messages.ts
│   │   ├── templates.ts
│   │   ├── flows.ts
│   │   ├── bulk.ts
│   │   └── chatbot.ts
│   ├── components/
│   │   ├── common/
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppSidebar.vue
│   │   │   ├── LoadingSpinner.vue
│   │   │   ├── EmptyState.vue
│   │   │   └── ConfirmDialog.vue
│   │   ├── chat/
│   │   │   ├── ChatLayout.vue
│   │   │   ├── ContactList.vue
│   │   │   ├── ContactItem.vue
│   │   │   ├── ChatWindow.vue
│   │   │   ├── MessageList.vue
│   │   │   ├── MessageItem.vue
│   │   │   ├── MessageInput.vue
│   │   │   ├── MediaPreview.vue
│   │   │   └── TypingIndicator.vue
│   │   ├── templates/
│   │   │   ├── TemplateList.vue
│   │   │   ├── TemplateForm.vue
│   │   │   ├── TemplatePreview.vue
│   │   │   └── ButtonEditor.vue
│   │   ├── flows/
│   │   │   ├── FlowList.vue
│   │   │   ├── FlowBuilder.vue
│   │   │   ├── ScreenEditor.vue
│   │   │   └── FieldEditor.vue
│   │   ├── bulk/
│   │   │   ├── CampaignList.vue
│   │   │   ├── CampaignForm.vue
│   │   │   ├── RecipientImport.vue
│   │   │   └── CampaignProgress.vue
│   │   ├── chatbot/
│   │   │   ├── ChatbotSettings.vue
│   │   │   ├── KeywordRuleList.vue
│   │   │   ├── KeywordRuleForm.vue
│   │   │   ├── ChatbotFlowList.vue
│   │   │   ├── ChatbotFlowBuilder.vue
│   │   │   ├── FlowStepEditor.vue
│   │   │   ├── AIContextList.vue
│   │   │   ├── AIContextForm.vue
│   │   │   └── AgentTransfers.vue
│   │   └── dashboard/
│   │       ├── StatsCards.vue
│   │       ├── MessageChart.vue
│   │       └── RecentActivity.vue
│   ├── views/
│   │   ├── auth/
│   │   │   ├── LoginView.vue
│   │   │   └── RegisterView.vue
│   │   ├── dashboard/
│   │   │   └── DashboardView.vue
│   │   ├── chat/
│   │   │   └── ChatView.vue
│   │   ├── templates/
│   │   │   ├── TemplatesView.vue
│   │   │   └── TemplateEditView.vue
│   │   ├── flows/
│   │   │   ├── FlowsView.vue
│   │   │   └── FlowEditView.vue
│   │   ├── bulk/
│   │   │   ├── CampaignsView.vue
│   │   │   └── CampaignEditView.vue
│   │   ├── chatbot/
│   │   │   ├── ChatbotView.vue
│   │   │   ├── KeywordsView.vue
│   │   │   ├── ChatbotFlowsView.vue
│   │   │   └── AIContextsView.vue
│   │   └── settings/
│   │       ├── SettingsView.vue
│   │       └── AccountsView.vue
│   ├── types/
│   │   ├── api.ts
│   │   ├── models.ts
│   │   └── websocket.ts
│   └── utils/
│       ├── formatters.ts
│       ├── validators.ts
│       └── constants.ts
├── public/
│   └── sounds/
│       ├── notification.mp3
│       └── send.mp3
├── index.html
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

### 6.2 Key Components

#### 6.2.1 Chat Interface (with shadcn-vue)

```vue
<!-- src/components/chat/ChatWindow.vue -->
<template>
  <div class="flex flex-col h-full bg-background">
    <!-- Header -->
    <div class="flex items-center p-4 border-b">
      <Avatar class="h-10 w-10">
        <AvatarImage :src="contact.avatarUrl" />
        <AvatarFallback>{{ getInitials(contact.profileName) }}</AvatarFallback>
      </Avatar>
      <div class="ml-3 flex-1">
        <h3 class="font-semibold text-foreground">
          {{ contact.profileName || contact.phoneNumber }}
        </h3>
        <span class="text-sm text-muted-foreground">
          {{ formatPhoneNumber(contact.phoneNumber) }}
        </span>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon">
            <MoreVertical class="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem @click="transferToAgent">
            <UserPlus class="mr-2 h-4 w-4" />
            Transfer to Agent
          </DropdownMenuItem>
          <DropdownMenuItem @click="viewContact">
            <User class="mr-2 h-4 w-4" />
            View Contact
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem class="text-destructive" @click="deleteConversation">
            <Trash class="mr-2 h-4 w-4" />
            Delete Conversation
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <!-- Messages -->
    <ScrollArea ref="messagesContainer" class="flex-1 p-4">
      <div v-if="loading" class="space-y-4">
        <Skeleton v-for="i in 5" :key="i" class="h-16 w-3/4" :class="i % 2 === 0 ? 'ml-auto' : ''" />
      </div>
      <MessageList v-else :messages="messages" @load-more="loadMoreMessages" />
      <div v-if="isTyping" class="flex items-center gap-2 text-muted-foreground">
        <Loader2 class="h-4 w-4 animate-spin" />
        <span class="text-sm">Typing...</span>
      </div>
    </ScrollArea>

    <!-- Input -->
    <div class="border-t p-4">
      <div class="flex items-end gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <Paperclip class="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem @click="attachImage">
              <Image class="mr-2 h-4 w-4" /> Image
            </DropdownMenuItem>
            <DropdownMenuItem @click="attachDocument">
              <FileText class="mr-2 h-4 w-4" /> Document
            </DropdownMenuItem>
            <DropdownMenuItem @click="showTemplates">
              <LayoutTemplate class="mr-2 h-4 w-4" /> Template
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <Textarea
          v-model="messageText"
          placeholder="Type a message..."
          class="min-h-[40px] max-h-[120px] resize-none"
          @keydown.enter.exact.prevent="sendMessage"
        />

        <Button @click="sendMessage" :disabled="!messageText.trim() || sending">
          <Loader2 v-if="sending" class="h-4 w-4 animate-spin" />
          <Send v-else class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMessages } from '@/stores/messages'
import { useWebSocket } from '@/composables/useWebSocket'
import {
  Avatar, AvatarImage, AvatarFallback
} from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuItem, DropdownMenuSeparator
} from '@/components/ui/dropdown-menu'
import {
  MoreVertical, Send, Paperclip, Image, FileText,
  LayoutTemplate, Loader2, UserPlus, User, Trash
} from 'lucide-vue-next'

const props = defineProps<{ contactId: string }>()

const messagesStore = useMessages()
const { subscribe, unsubscribe } = useWebSocket()

const messages = computed(() => messagesStore.getByContact(props.contactId))
const loading = computed(() => messagesStore.loading)
const sending = ref(false)
const messageText = ref('')
const isTyping = ref(false)

onMounted(() => {
  messagesStore.fetchMessages(props.contactId)
  subscribe(`contact:${props.contactId}`, handleNewMessage)
})

onUnmounted(() => {
  unsubscribe(`contact:${props.contactId}`, handleNewMessage)
})

const handleNewMessage = (message: Message) => {
  messagesStore.addMessage(message)
  scrollToBottom()
  playNotificationSound()
}

const sendMessage = async () => {
  if (!messageText.value.trim() || sending.value) return
  sending.value = true
  try {
    await messagesStore.sendMessage(props.contactId, messageText.value)
    messageText.value = ''
  } finally {
    sending.value = false
  }
}
</script>
```

#### 6.2.2 Chatbot Flow Builder (with shadcn-vue)

```vue
<!-- src/components/chatbot/ChatbotFlowBuilder.vue -->
<template>
  <div class="space-y-6">
    <!-- Flow Settings -->
    <Card>
      <CardHeader>
        <CardTitle>Flow Settings</CardTitle>
        <CardDescription>Configure the basic settings for your chatbot flow</CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="saveFlow" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="name">Flow Name</Label>
              <Input id="name" v-model="flow.name" placeholder="e.g., Customer Support" />
            </div>
            <div class="space-y-2">
              <Label for="keywords">Trigger Keywords</Label>
              <Input id="keywords" v-model="flow.triggerKeywords" placeholder="help, support, assist" />
            </div>
          </div>
          <div class="space-y-2">
            <Label for="initial">Initial Message</Label>
            <Textarea id="initial" v-model="flow.initialMessage" placeholder="Welcome! How can I help you today?" />
          </div>
          <div class="space-y-2">
            <Label for="completion">Completion Message</Label>
            <Textarea id="completion" v-model="flow.completionMessage" placeholder="Thank you! Your request has been submitted." />
          </div>
        </form>
      </CardContent>
    </Card>

    <!-- Steps -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Flow Steps</CardTitle>
            <CardDescription>Define the conversation steps</CardDescription>
          </div>
          <Button @click="addStep">
            <Plus class="mr-2 h-4 w-4" />
            Add Step
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="flow.steps.length === 0" class="text-center py-8 text-muted-foreground">
          <MessageSquare class="mx-auto h-12 w-12 mb-4 opacity-50" />
          <p>No steps yet. Add your first step to get started.</p>
        </div>

        <Draggable v-else v-model="flow.steps" item-key="stepName" handle=".drag-handle" class="space-y-4">
          <template #item="{ element, index }">
            <Card class="border-l-4 border-l-primary">
              <CardHeader class="pb-2">
                <div class="flex items-center gap-2">
                  <GripVertical class="drag-handle h-5 w-5 cursor-move text-muted-foreground" />
                  <Badge variant="outline">Step {{ index + 1 }}</Badge>
                  <Input
                    v-model="element.stepName"
                    class="max-w-[200px] h-8"
                    placeholder="Step name"
                  />
                  <div class="ml-auto">
                    <Button variant="ghost" size="icon" @click="deleteStep(index)">
                      <Trash class="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <FlowStepEditor
                  :step="element"
                  :all-steps="flow.steps"
                  @update="updateStep(index, $event)"
                />
              </CardContent>
            </Card>
          </template>
        </Draggable>
      </CardContent>
    </Card>

    <!-- Actions -->
    <div class="flex justify-end gap-2">
      <Button variant="outline" @click="previewFlow">
        <Eye class="mr-2 h-4 w-4" />
        Preview
      </Button>
      <Button @click="saveFlow" :disabled="saving">
        <Loader2 v-if="saving" class="mr-2 h-4 w-4 animate-spin" />
        <Save v-else class="mr-2 h-4 w-4" />
        Save Flow
      </Button>
    </div>

    <!-- Preview Dialog -->
    <Dialog v-model:open="showPreview">
      <DialogContent class="max-w-md">
        <DialogHeader>
          <DialogTitle>Flow Preview</DialogTitle>
          <DialogDescription>Test your flow conversation</DialogDescription>
        </DialogHeader>
        <FlowPreview :flow="flow" />
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Draggable from 'vuedraggable'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Plus, Trash, GripVertical, Eye, Save, Loader2, MessageSquare } from 'lucide-vue-next'
import FlowStepEditor from './FlowStepEditor.vue'
import FlowPreview from './FlowPreview.vue'

const flow = ref({
  name: '',
  triggerKeywords: '',
  initialMessage: '',
  completionMessage: '',
  steps: []
})

const saving = ref(false)
const showPreview = ref(false)

const addStep = () => {
  flow.value.steps.push({
    stepName: `step_${flow.value.steps.length + 1}`,
    message: '',
    inputType: 'text',
    storeAs: '',
    nextStep: ''
  })
}

const deleteStep = (index: number) => {
  flow.value.steps.splice(index, 1)
}

const updateStep = (index: number, updates: any) => {
  flow.value.steps[index] = { ...flow.value.steps[index], ...updates }
}

const previewFlow = () => {
  showPreview.value = true
}

const saveFlow = async () => {
  saving.value = true
  // API call to save flow
  saving.value = false
}
</script>
```

### 6.3 State Management (Pinia)

```typescript
// src/stores/messages.ts
import { defineStore } from 'pinia'
import { messagesApi } from '@/services/messages'

interface MessagesState {
  byContact: Record<string, Message[]>
  loading: boolean
  sending: boolean
}

export const useMessages = defineStore('messages', {
  state: (): MessagesState => ({
    byContact: {},
    loading: false,
    sending: false
  }),

  getters: {
    getByContact: (state) => (contactId: string) => {
      return state.byContact[contactId] || []
    }
  },

  actions: {
    async fetchMessages(contactId: string, pagination?: Pagination) {
      this.loading = true
      try {
        const messages = await messagesApi.getByContact(contactId, pagination)
        this.byContact[contactId] = messages
      } finally {
        this.loading = false
      }
    },

    async sendMessage(contactId: string, content: string, type: string = 'text') {
      this.sending = true
      try {
        const message = await messagesApi.send({ contactId, content, type })
        this.addMessage(message)
        return message
      } finally {
        this.sending = false
      }
    },

    addMessage(message: Message) {
      const contactId = message.contactId
      if (!this.byContact[contactId]) {
        this.byContact[contactId] = []
      }
      // Avoid duplicates
      if (!this.byContact[contactId].find(m => m.id === message.id)) {
        this.byContact[contactId].push(message)
      }
    }
  }
})
```

### 6.4 WebSocket Integration

```typescript
// src/composables/useWebSocket.ts
import { ref, onUnmounted } from 'vue'
import { io, Socket } from 'socket.io-client'
import { useAuth } from '@/stores/auth'

const socket = ref<Socket | null>(null)
const connected = ref(false)
const subscriptions = new Map<string, Set<Function>>()

export function useWebSocket() {
  const auth = useAuth()

  const connect = () => {
    if (socket.value?.connected) return

    socket.value = io(import.meta.env.VITE_WS_URL, {
      auth: { token: auth.token },
      transports: ['websocket']
    })

    socket.value.on('connect', () => {
      connected.value = true
    })

    socket.value.on('disconnect', () => {
      connected.value = false
    })

    socket.value.on('message', (data: { event: string; payload: any }) => {
      const handlers = subscriptions.get(data.event)
      handlers?.forEach(handler => handler(data.payload))
    })
  }

  const subscribe = (event: string, handler: Function) => {
    if (!subscriptions.has(event)) {
      subscriptions.set(event, new Set())
    }
    subscriptions.get(event)!.add(handler)
  }

  const unsubscribe = (event: string, handler: Function) => {
    subscriptions.get(event)?.delete(handler)
  }

  const emit = (event: string, payload: any) => {
    socket.value?.emit(event, payload)
  }

  return { connect, subscribe, unsubscribe, emit, connected }
}
```

---

## 7. API Design

### 7.1 RESTful Endpoints

```yaml
# Authentication
POST   /api/v1/auth/login
POST   /api/v1/auth/register
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout

# WhatsApp Accounts
GET    /api/v1/accounts
POST   /api/v1/accounts
GET    /api/v1/accounts/:id
PUT    /api/v1/accounts/:id
DELETE /api/v1/accounts/:id

# Contacts
GET    /api/v1/contacts
POST   /api/v1/contacts
GET    /api/v1/contacts/:id
PUT    /api/v1/contacts/:id
DELETE /api/v1/contacts/:id
PUT    /api/v1/contacts/:id/assign

# Messages
GET    /api/v1/contacts/:id/messages
POST   /api/v1/messages
POST   /api/v1/messages/template
POST   /api/v1/messages/media
PUT    /api/v1/messages/:id/read

# Templates
GET    /api/v1/templates
POST   /api/v1/templates
GET    /api/v1/templates/:id
PUT    /api/v1/templates/:id
DELETE /api/v1/templates/:id
POST   /api/v1/templates/sync

# WhatsApp Flows
GET    /api/v1/flows
POST   /api/v1/flows
GET    /api/v1/flows/:id
PUT    /api/v1/flows/:id
DELETE /api/v1/flows/:id
POST   /api/v1/flows/:id/publish
POST   /api/v1/flows/:id/deprecate

# Bulk Campaigns
GET    /api/v1/campaigns
POST   /api/v1/campaigns
GET    /api/v1/campaigns/:id
PUT    /api/v1/campaigns/:id
DELETE /api/v1/campaigns/:id
POST   /api/v1/campaigns/:id/start
GET    /api/v1/campaigns/:id/progress
POST   /api/v1/campaigns/:id/recipients/import

# Chatbot Settings
GET    /api/v1/chatbot/settings
PUT    /api/v1/chatbot/settings

# Keyword Rules
GET    /api/v1/chatbot/keywords
POST   /api/v1/chatbot/keywords
GET    /api/v1/chatbot/keywords/:id
PUT    /api/v1/chatbot/keywords/:id
DELETE /api/v1/chatbot/keywords/:id

# Chatbot Flows
GET    /api/v1/chatbot/flows
POST   /api/v1/chatbot/flows
GET    /api/v1/chatbot/flows/:id
PUT    /api/v1/chatbot/flows/:id
DELETE /api/v1/chatbot/flows/:id

# AI Contexts
GET    /api/v1/chatbot/ai-contexts
POST   /api/v1/chatbot/ai-contexts
GET    /api/v1/chatbot/ai-contexts/:id
PUT    /api/v1/chatbot/ai-contexts/:id
DELETE /api/v1/chatbot/ai-contexts/:id

# Agent Transfers
GET    /api/v1/chatbot/transfers
POST   /api/v1/chatbot/transfers
PUT    /api/v1/chatbot/transfers/:id/resume

# Sessions (for debugging/admin)
GET    /api/v1/chatbot/sessions
GET    /api/v1/chatbot/sessions/:id

# Webhook (Meta callback)
GET    /api/v1/webhook
POST   /api/v1/webhook

# Dashboard/Analytics
GET    /api/v1/analytics/overview
GET    /api/v1/analytics/messages
GET    /api/v1/analytics/chatbot
```

### 7.2 WebSocket Events

```typescript
// Client -> Server
interface ClientEvents {
  'subscribe': { room: string }
  'unsubscribe': { room: string }
  'typing': { contactId: string, isTyping: boolean }
  'read': { contactId: string }
}

// Server -> Client
interface ServerEvents {
  'message:new': Message
  'message:status': { messageId: string, status: string }
  'contact:updated': Contact
  'contact:new': Contact
  'session:started': ChatbotSession
  'session:completed': ChatbotSession
  'transfer:new': AgentTransfer
  'transfer:resumed': AgentTransfer
}
```

---

## 8. Real-time Communication

### 8.1 WebSocket Hub (Go)

```go
// internal/websocket/hub.go
package websocket

type Hub struct {
    clients    map[*Client]bool
    rooms      map[string]map[*Client]bool
    broadcast  chan Event
    register   chan *Client
    unregister chan *Client
    redis      *redis.Client
}

func (h *Hub) Run() {
    // Subscribe to Redis PubSub for multi-instance support
    pubsub := h.redis.Subscribe(ctx, "ws:broadcast")
    go h.handleRedisMessages(pubsub)

    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            h.removeClient(client)
        case event := <-h.broadcast:
            h.broadcastToRoom(event)
        }
    }
}

func (h *Hub) BroadcastToUser(userID string, event Event) {
    // Publish to Redis for multi-instance
    h.redis.Publish(ctx, "ws:broadcast", event)
}

func (h *Hub) BroadcastToRoom(room string, event Event) {
    clients := h.rooms[room]
    for client := range clients {
        client.send <- event
    }
}
```

### 8.2 Event Publishing

```go
// When a new message arrives
func (s *MessagingService) ProcessIncoming(ctx context.Context, payload WebhookPayload) error {
    // ... save message ...

    // Publish to WebSocket
    s.wsHub.BroadcastToRoom(fmt.Sprintf("contact:%s", contact.ID), Event{
        Type:    "message:new",
        Payload: message,
    })

    // Also notify the assigned user
    if contact.AssignedUserID != nil {
        s.wsHub.BroadcastToUser(*contact.AssignedUserID, Event{
            Type:    "contact:updated",
            Payload: contact,
        })
    }
}
```

---

## 9. AI Integration

### 9.1 Provider Interface

```go
// internal/services/ai/provider.go
package ai

type Provider interface {
    GenerateResponse(ctx context.Context, req GenerateRequest) (string, error)
    Name() string
}

type GenerateRequest struct {
    SystemPrompt string
    Messages     []Message
    MaxTokens    int
    Temperature  float64
}

type Message struct {
    Role    string // system, user, assistant
    Content string
}
```

### 9.2 OpenAI Provider

```go
// internal/services/ai/openai.go
package ai

type OpenAIProvider struct {
    client *openai.Client
    model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
    return &OpenAIProvider{
        client: openai.NewClient(apiKey),
        model:  model,
    }
}

func (p *OpenAIProvider) GenerateResponse(ctx context.Context, req GenerateRequest) (string, error) {
    messages := make([]openai.ChatCompletionMessage, len(req.Messages)+1)
    messages[0] = openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleSystem,
        Content: req.SystemPrompt,
    }

    for i, msg := range req.Messages {
        messages[i+1] = openai.ChatCompletionMessage{
            Role:    msg.Role,
            Content: msg.Content,
        }
    }

    resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model:       p.model,
        Messages:    messages,
        MaxTokens:   req.MaxTokens,
        Temperature: float32(req.Temperature),
    })

    if err != nil {
        return "", fmt.Errorf("openai error: %w", err)
    }

    return resp.Choices[0].Message.Content, nil
}
```

### 9.3 Context Building

```go
// internal/services/ai/context.go
package ai

func (s *Service) BuildContext(ctx context.Context, message, phoneNumber string) (string, error) {
    contexts, err := s.contextRepo.FindEnabled(ctx, s.orgID)
    if err != nil {
        return "", err
    }

    var builder strings.Builder
    builder.WriteString("Context information:\n\n")

    for _, aiCtx := range contexts {
        // Check trigger keywords
        if len(aiCtx.TriggerKeywords) > 0 && !containsAny(message, aiCtx.TriggerKeywords) {
            continue
        }

        switch aiCtx.ContextType {
        case "static":
            builder.WriteString(aiCtx.StaticContent)
            builder.WriteString("\n\n")
        case "query":
            data, err := s.executeQuery(ctx, aiCtx.QueryConfig, phoneNumber)
            if err != nil {
                continue
            }
            builder.WriteString(fmt.Sprintf("%s: %s\n\n", aiCtx.Name, formatJSON(data)))
        }
    }

    return builder.String(), nil
}
```

---

## 10. Development Phases

### Phase 1: Foundation (4-6 weeks)
**Goal:** Core infrastructure and basic messaging

- [ ] Project setup (Go + Vue.js scaffolding)
- [ ] Database schema and migrations
- [ ] Authentication system (JWT)
- [ ] WhatsApp Account management
- [ ] Webhook handler for incoming messages
- [ ] Basic message sending (text only)
- [ ] Contact management
- [ ] Basic chat UI (contact list + message thread)
- [ ] WebSocket infrastructure
- [ ] Real-time message updates

**Deliverables:**
- Working API server
- Basic frontend with login and chat
- Can send/receive text messages

### Phase 2: Templates & Media (3-4 weeks)
**Goal:** Full message type support

- [ ] Template management (CRUD)
- [ ] Template sync from Meta
- [ ] Template message sending
- [ ] Media upload/download
- [ ] Image/video/audio/document messages
- [ ] Interactive messages (buttons/lists)
- [ ] Message status tracking
- [ ] Read receipts

**Deliverables:**
- Template management UI
- Media message support
- Interactive message support

### Phase 3: WhatsApp Flows (3-4 weeks)
**Goal:** Meta's interactive forms

- [ ] Flow management (CRUD)
- [ ] Flow JSON builder
- [ ] Flow publishing lifecycle
- [ ] Flow message sending
- [ ] Flow response handling
- [ ] Flow preview/testing

**Deliverables:**
- Flow builder UI
- Flow integration in chat

### Phase 4: Bulk Messaging (2-3 weeks)
**Goal:** Mass message campaigns

- [ ] Campaign management
- [ ] Recipient import (CSV/manual)
- [ ] Background job processing
- [ ] Progress tracking
- [ ] Retry failed messages
- [ ] Campaign analytics

**Deliverables:**
- Bulk campaign UI
- Background worker
- Progress dashboard

### Phase 5: Chatbot - Keywords & Flows (4-5 weeks)
**Goal:** Automated responses

- [ ] Chatbot settings
- [ ] Keyword rule engine
- [ ] Keyword matching (exact/contains/regex)
- [ ] Chatbot conversation flows
- [ ] Flow step builder
- [ ] Input validation
- [ ] Session management
- [ ] Conditional routing
- [ ] Completion actions (webhook/create record)

**Deliverables:**
- Chatbot configuration UI
- Keyword rules working
- Conversation flows working

### Phase 6: AI Integration (2-3 weeks)
**Goal:** AI-powered responses

- [ ] AI provider integration (OpenAI, Anthropic, Google)
- [ ] Context management
- [ ] Conversation history
- [ ] Dynamic context queries
- [ ] AI response fallback

**Deliverables:**
- AI configuration UI
- AI responses working

### Phase 7: Agent Transfer & Advanced Features (2-3 weeks)
**Goal:** Human handoff and polish

- [ ] Agent transfer system
- [ ] Business hours
- [ ] Excluded numbers
- [ ] User assignment
- [ ] Transfer notifications
- [ ] Dashboard & analytics

**Deliverables:**
- Agent transfer UI
- Analytics dashboard

### Phase 8: Testing & Optimization (2-3 weeks)
**Goal:** Production readiness

- [ ] Unit tests (Go)
- [ ] Integration tests
- [ ] E2E tests (Playwright/Cypress)
- [ ] Performance optimization
- [ ] Security audit
- [ ] Documentation
- [ ] Docker production config

**Deliverables:**
- Test suite
- Documentation
- Production-ready deployment

---

## 11. Deployment Strategy

### 11.1 Docker Compose (Development/Staging)

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: docker/Dockerfile.api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/whatsapp
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - db
      - redis

  worker:
    build:
      context: .
      dockerfile: docker/Dockerfile.worker
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/whatsapp
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - api

  db:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=whatsapp

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - api
      - frontend

volumes:
  postgres_data:
  redis_data:
```

### 11.2 Kubernetes (Production)

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whatsapp-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: whatsapp-api
  template:
    metadata:
      labels:
        app: whatsapp-api
    spec:
      containers:
      - name: api
        image: whatsapp-platform/api:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: whatsapp-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

---

## 12. Migration Strategy

### 12.1 Data Migration

```go
// cmd/migrate/main.go
package main

// Migration from Frappe/MariaDB to PostgreSQL
func migrateData() {
    // 1. Export from Frappe
    exportWhatsAppAccounts()
    exportMessages()
    exportTemplates()
    exportFlows()
    exportContacts()
    exportChatbotSettings()
    exportKeywordRules()
    exportChatbotFlows()

    // 2. Transform data
    transformPhoneNumbers() // Normalize format
    transformMediaPaths()   // Update file references
    transformJSONFields()   // Convert Frappe JSON to PostgreSQL JSONB

    // 3. Import to PostgreSQL
    importOrganizations()
    importUsers()
    importAccounts()
    importContacts()
    importMessages()
    importTemplates()
    importFlows()
    importChatbotConfig()
}
```

### 12.2 Parallel Running

1. **Week 1-2:** Deploy new system alongside Frappe
2. **Week 3-4:** Webhook forwarding (Frappe receives, forwards to new system)
3. **Week 5-6:** Switch webhook to new system
4. **Week 7-8:** Decommission Frappe app

### 12.3 Rollback Plan

- Keep Frappe system running for 30 days after migration
- Maintain webhook forwarding capability
- Database backups before each migration step
- Feature flags for gradual rollout

---

## Appendix A: Estimated Effort

| Phase | Duration | Backend | Frontend |
|-------|----------|---------|----------|
| Phase 1: Foundation | 4-6 weeks | 70% | 30% |
| Phase 2: Templates & Media | 3-4 weeks | 60% | 40% |
| Phase 3: WhatsApp Flows | 3-4 weeks | 50% | 50% |
| Phase 4: Bulk Messaging | 2-3 weeks | 60% | 40% |
| Phase 5: Chatbot Keywords & Flows | 4-5 weeks | 70% | 30% |
| Phase 6: AI Integration | 2-3 weeks | 80% | 20% |
| Phase 7: Agent Transfer | 2-3 weeks | 50% | 50% |
| Phase 8: Testing & Polish | 2-3 weeks | 50% | 50% |
| **Total** | **22-31 weeks** | | |

## Appendix B: Technology Choices

| Component | Primary Choice | Alternative |
|-----------|---------------|-------------|
| Go Framework | **Fastglue** (zerodha/fastglue) | Fiber, Echo, Gin, Chi |
| ORM | GORM | sqlx, sqlc, Ent |
| Frontend Framework | Vue 3 | React, Svelte |
| UI Library | **shadcn-vue** | Radix Vue, PrimeVue, Vuetify |
| CSS Framework | TailwindCSS | (required for shadcn-vue) |
| State Management | Pinia | Vuex |
| Job Queue | Asynq | Faktory, Machinery |
| WebSocket | Gorilla | Melody, Centrifugo |
| Configuration | knadh/koanf | viper |
| Logging | zerodha/logf | zap, zerolog |

### Why Fastglue?

**Fastglue** (github.com/zerodha/fastglue) is a minimal, opinionated HTTP framework built on top of `fasthttp`:

1. **Performance**: Built on `fasthttp` which is 10x faster than `net/http`
2. **Clean API**: Simple handler interface with request context
3. **Built-in features**: Request binding, validation, JSON responses
4. **Battle-tested**: Used in production at Zerodha (India's largest stock broker)
5. **Minimal footprint**: Small codebase, easy to understand and extend

### Why shadcn-vue?

**shadcn-vue** (shadcn-vue.com) provides beautifully designed, accessible components:

1. **Copy-paste components**: Components live in your project, fully customizable
2. **Radix Vue primitives**: Built on accessible, unstyled primitives
3. **TailwindCSS**: Uses Tailwind for styling, easy to customize
4. **No runtime dependency**: Components are copied to your project
5. **TypeScript**: Full TypeScript support
6. **Dark mode**: Built-in dark mode support
7. **Active community**: Growing ecosystem with excellent documentation

## Appendix C: Security Considerations

1. **API Keys & Tokens**
   - Encrypt at rest (AES-256)
   - Never log sensitive data
   - Rotate keys periodically

2. **Authentication**
   - JWT with short expiry (15 min)
   - Refresh tokens (7 days)
   - Rate limiting on auth endpoints

3. **Webhook Security**
   - Verify Meta signature
   - IP allowlisting (optional)
   - Idempotency for duplicate events

4. **Data Protection**
   - GDPR compliance
   - Message retention policies
   - Data export capability

---

*Document Version: 1.0*
*Created: December 2024*
*Last Updated: December 2024*
