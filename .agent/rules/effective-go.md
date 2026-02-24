---
trigger: always_on
---

# Go Development Standards for Antigravity Agents

## 1. Idiomatic Go (Effective Go)
- Use 'gofmt' and 'goimports' automatically on every save.
- Variable naming: Short, camelCase (e.g., 'srv' instead of 'serverInstance').
- Error handling: Never ignore errors. Use 'if err != nil { return fmt.Errorf("context: %w", err) }'.
- No package-level global state unless strictly necessary for hardware/os hooks.

## 2. Performance & Low-Level Optimization
- Memory Management: Avoid unnecessary heap allocations. Prefer stack allocation.
- Pointers: Only use pointers when the struct is large or requires mutation.
- Slices: Pre-allocate slice capacity with 'make([]T, 0, capacity)' to avoid re-allocations.
- Concurrency: Use 'sync.Pool' for frequently allocated objects (GC pressure reduction).
- Goroutines: Every goroutine must have a clear lifecycle and a way to exit (Context).

## 3. Security (Secure Coding)
- OWASP: Sanitize all inputs from external APIs/DBs.
- SQL: Always use parameterized queries (no string concatenation for SQL).
- Sensitive Data: Use 'mlock' or zero-out buffers containing credentials before they are garbage collected.
- Dependencies: Use 'go mod tidy' and check for vulnerabilities via 'govulncheck'.

## 4. Architecture
- Follow 'Clean Architecture' or 'Standard Layout' (cmd/, internal/, pkg/).
- Interface Segregation: Define interfaces where they are USED, not where they are implemented.