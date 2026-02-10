# Example: TODO API

Build a RESTful TODO API with CRUD operations: POST/GET/PUT/DELETE at `/v1/todos`

## Implementation Steps

### 1. Models (`packages/api/v1/models.go`)
```go
type Todo struct {
    ID, Title, Description string
    Completed bool
    CreatedAt, UpdatedAt time.Time
}
type CreateTodoRequest struct { Title, Description string }
type UpdateTodoRequest struct { Title, Description string; Completed *bool }
```

### 2. Metrics (`packages/stats/stats.go`)
Add counters: `TodoCreateCounter`, `TodoListCounter`, `TodoGetCounter`, `TodoUpdateCounter`, `TodoDeleteCounter`
Initialize in `InitMetrics()` with pattern: `meter.Int64Counter("todo_create_operations", ...)`

### 3. Storage (`packages/api/v1/todo_storage.go`)
```go
type TodoStorage struct {
    mu sync.RWMutex
    todos map[string]*Todo
}
// Methods: Create, List, Get, Update, Delete (all with proper locking)
// Use uuid.New().String() for IDs
```

### 4. Handlers (`packages/api/v1/todos.go`)
Each handler pattern:
1. Check HTTP method
2. Increment metric counter (`stats.TodoCreateCounter.Add(r.Context(), 1)`)
3. Parse request/extract ID from path
4. Call storage method
5. Log operation (`l.Info("Created TODO: %s", id)`)
6. Return JSON response

Handlers: `CreateTodoHandler`, `ListTodosHandler`, `GetTodoHandler`, `UpdateTodoHandler`, `DeleteTodoHandler`

### 5. Register Routes (`packages/api/server.go`)
```go
http.HandleFunc("/v1/todos", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost: v1.CreateTodoHandler(w, r)
    case http.MethodGet: v1.ListTodosHandler(w, r)
    }
})
http.HandleFunc("/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
    // Route GET/PUT/DELETE by method
})
```

Dependencies: `github.com/google/uuid v1.3.0`

## Testing
```bash
# Create
curl -X POST http://localhost:8080/v1/todos -d '{"title":"Learn Go"}'

# List
curl http://localhost:8080/v1/todos

# Get/Update/Delete
curl http://localhost:8080/v1/todos/{id}
curl -X PUT http://localhost:8080/v1/todos/{id} -d '{"completed":true}'
curl -X DELETE http://localhost:8080/v1/todos/{id}
```

## Pattern Applied
✅ Extended existing packages (stats, API)
✅ Followed handler patterns (method check, metrics, logging, JSON response)
✅ Added new files without modifying core boilerplate
✅ Used logger (`l.Info`) and metrics (`stats.*Counter.Add`) throughout
