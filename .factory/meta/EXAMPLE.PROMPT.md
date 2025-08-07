# Example Prompt: Building Applications with Go Boilerplate

## 📋 Prerequisites

**CRITICAL**: Before proceeding, you MUST thoroughly read and understand the `guide.md` and `README.md` file in this repository. These files contain essential information about the boilerplate architecture, package structure, development patterns, and critical rules for extending the codebase.

## 🎯 Prompt Template

Use this structured prompt when working with a coding agent to build applications using this Go boilerplate:

---

### **Initial Setup Instructions**

```
You are a senior Go developer tasked with creating a new application using our existing Go boilerplate. 

MANDATORY FIRST STEP: Read and fully understand the `guide.md` and `README.md` file in this repository. These files contain critical information about:
- Project architecture and package structure
- Existing components (config, logger, stats, CLI, API, bootstrap)
- Development patterns and extension points
- Rules for extending vs. replacing code

Follow the software engineering guidelines in `.factory/meta/software-engineering-guidelines.md`.
```

### **Application Specification**

```
APPLICATION PURPOSE: [One clear sentence describing what the app does]

FEATURE REQUIREMENTS:
1. [Specific feature 1]
2. [Specific feature 2]
3. [Additional features as needed]

TECHNICAL REQUIREMENTS:
- Use existing boilerplate packages (config, logger, stats, CLI, API)
- Follow established error handling and logging patterns
- Implement proper metrics collection
- Maintain modular architecture
```

### **Implementation Guidelines**

```
DEVELOPMENT APPROACH:
1. Analyze existing packages and identify reusable components
2. Plan extension points without modifying core boilerplate files
3. Implement features by extending existing modules
4. Add proper error handling, logging, and metrics to all new code
5. Create comprehensive tests for new functionality

DELIVERABLES:
- Complete, working application
- Updated configuration schema (if needed)
- New CLI commands (if applicable)
- API endpoints (if applicable)
- Documentation updates
- Test coverage for new features
```

---

## 🚀 Example: OpenAI ChatGPT Integration Application

### **Complete Prompt for ChatGPT Integration App**

```
You are a senior Go developer creating a ChatGPT integration application using our existing Go boilerplate.

MANDATORY: First, read and understand the `guide.md` and `README.md` files completely. Pay special attention to:
- Package structure and existing components
- Development patterns for extending functionality
- Critical rules about leveraging vs. replacing existing code

APPLICATION PURPOSE: 
Create a service that provides ChatGPT integration with both CLI and HTTP API interfaces, including conversation management and response streaming.

FEATURE REQUIREMENTS:
1. CLI command to send messages to ChatGPT and receive responses
2. HTTP API endpoint for ChatGPT conversations
3. Conversation history management and persistence
4. Streaming response support
5. Configurable OpenAI API settings (model, temperature, max tokens)
6. Rate limiting and error handling for API calls
7. Metrics collection for API usage and response times

TECHNICAL REQUIREMENTS:
- Extend existing config package to include OpenAI settings
- Use existing logger for all operations
- Implement metrics using existing stats package
- Add new CLI commands using existing CLI structure
- Create new API endpoints following existing patterns
- Use existing error handling conventions
- Maintain existing bootstrap and initialization flow

CONFIGURATION EXTENSION:
Add OpenAI configuration block to existing config structure:
```go
type Configuration struct {
    // Existing fields...
    OpenAI OpenAIConfig `hcl:"openai,block"`
}

type OpenAIConfig struct {
    APIKey      string  `hcl:"api_key"`
    Model       string  `hcl:"model"`
    Temperature float64 `hcl:"temperature"`
    MaxTokens   int     `hcl:"max_tokens"`
    BaseURL     string  `hcl:"base_url"`
}
```

CLI COMMANDS TO ADD:
- `chat` - Interactive chat mode
- `ask` - Single question/response
- `history` - View conversation history

API ENDPOINTS TO ADD:
- `POST /v1/chat/completions` - Send message to ChatGPT
- `GET /v1/chat/conversations` - List conversations
- `GET /v1/chat/conversations/{id}` - Get conversation history
- `DELETE /v1/chat/conversations/{id}` - Delete conversation

IMPLEMENTATION APPROACH:
1. Analyze guide.md and README.md and existing package structure
2. Plan integration points without modifying core files
3. Create new packages/modules as needed (e.g., packages/openai/)
4. Extend CLI with new commands following existing patterns
5. Add API endpoints following existing v1 structure
6. Implement proper error handling and logging throughout
7. Add metrics for API calls, response times, and error rates
8. Create comprehensive tests for all new functionality

DELIVERABLES:
- Complete ChatGPT integration service
- Updated configuration schema with OpenAI settings
- New CLI commands for chat functionality
- HTTP API endpoints for programmatic access
- Conversation persistence and management
- Comprehensive error handling and logging
- Metrics collection and monitoring
- Unit and integration tests
- Updated documentation

Remember: Extend existing components, don't replace them. Follow established patterns and maintain the boilerplate's architectural integrity.
```

---

## 📝 Usage Instructions for Engineers

1. **Read the Guide**: Always start by thoroughly reading `guide.md`
2. **Customize the Prompt**: Adapt the template above for your specific application needs
3. **Provide Clear Requirements**: Be specific about features and technical requirements
4. **Emphasize Extension**: Remind the agent to extend, not replace existing code
5. **Request Complete Implementation**: Ask for full features with tests and documentation
6. **Review Integration**: Ensure the final application properly integrates with existing boilerplate components

## ⚠️ Important Notes

- **Never skip reading guide.md and README.md** - These files contain critical architectural information
- **Always emphasize extending existing packages** rather than creating duplicates
- **Require proper error handling, logging, and metrics** in all new code
- **Ask for complete implementations** including tests and documentation
- **Verify integration** with existing boilerplate components

This prompt template ensures consistent, high-quality applications that properly leverage the boilerplate's architecture and maintain code quality standards.
