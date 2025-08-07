# Software Engineering Guidelines

**Error Handling:**  
- All new code must include robust error handling and clear error messages, following Go best practices.
- Return errors, don't panic. Use `fmt.Errorf()` for context wrapping.

**Logging:**  
- Use the existing logger module for logging. Follow the same logging patterns and conventions used in the boilerplate.
- Log at appropriate levels: Debug for development, Info for operations, Error for failures.

**Modularity:**  
- Structure code into small, focused modules with clear interfaces.
- Single responsibility principle: one function, one purpose.

**Reusability:**  
- Write functions and components so they can be reused across the project.
- Avoid hardcoded values; use configuration or parameters.

**Readability:**  
- Code should be easy to read, with descriptive names and clear formatting.
- Use `gofmt` and follow Go naming conventions (camelCase, PascalCase).

**Performance:**  
- Avoid premature optimization, but be mindful of memory allocations and goroutine usage.
- Use context for cancellation and timeouts.

**Security:**  
- Validate all inputs and sanitize outputs.
- Never log sensitive data (passwords, tokens, PII).

**Composition:**  
- Prefer composition over inheritance.
- Use structs and interfaces to define objects and their behaviors.
