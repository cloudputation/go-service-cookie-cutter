# Agent Prompt Template

## Setup Instructions
```
Read GUIDE.md before starting. This Go boilerplate has:
- Config (HCL), Logger (multi-output), Stats (Prometheus), CLI (Cobra), API (HTTP)
- Rule: Extend existing packages, don't replace core files
- Follow: .factory/meta/SOFTWARE-ENGINEERING-GUIDELINES.md
```

## Prompt Template
```
Build [APP_NAME] using this Go boilerplate.

PURPOSE: [One sentence description]

FEATURES:
1. [Feature 1]
2. [Feature 2]

REQUIREMENTS:
- Use existing packages (config, logger, stats, cli, api)
- Add proper error handling, logging, metrics
- Extend, don't modify core boilerplate

DELIVERABLES:
- Working application
- Updated config (if needed)
- New CLI commands/API endpoints
- Tests
```

## Example: ChatGPT Integration
```
Build ChatGPT integration service using this Go boilerplate.

PURPOSE: HTTP API and CLI for ChatGPT with conversation management.

FEATURES:
1. CLI commands: chat (interactive), ask (single Q&A), history
2. API endpoints: POST /v1/chat/completions, GET /v1/chat/conversations
3. Conversation persistence
4. Streaming responses
5. Rate limiting

REQUIREMENTS:
- Extend config package for OpenAI settings (api_key, model, temperature)
- Use existing logger for all operations
- Add metrics: chat_requests, api_errors, response_time
- Follow existing handler patterns (method check, metrics, logging)

IMPLEMENTATION:
1. Add OpenAI config block to packages/config/config.go
2. Create packages/openai/ for API client
3. Add CLI commands in packages/cli/cli.go
4. Add API handlers in packages/api/v1/chat.go
5. Register endpoints in packages/api/server.go

DELIVERABLES:
- Complete ChatGPT service
- Updated config.hcl with openai block
- Tests for handlers and API client
```

## Key Points
- **Always read GUIDE.md first**
- **Extend, don't replace** existing code
- **Follow patterns**: error handling, logging, metrics
- **See `.factory/examples/simple-todo-api.md`** for complete example
