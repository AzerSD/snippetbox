# Let's Go by Alex Edwards 

Quick reference guide for concepts, patterns, and best practices

---

## Module Paths & Naming Convention

### Module Path Basics
A module path is the identifier for your Go project when it's downloadable and used by others.

**Best Practice:** Base your module path on a URL that you own.

```
Example: github.com/foo/bar → use as module path if hosted at https://github.com/foo/bar
```

### Naming Guidelines
- Ensure the module path is **globally unique**
- Avoid conflicts with existing packages
- Choose something clear and succinct
- Common pattern: Use a domain you own (e.g., `snippetbox.alexedwards.net`)

---

## Running Go

### Equivalent Commands
These three commands are functionally identical:

```go
$ go run .
$ go run main.go
$ go run snippetbox.alexedwards.net
```

All will compile and execute your Go program.

---

## Networking Fundamentals

### Network Addresses

#### Named Ports vs Numeric Ports
Go allows you to specify network addresses using **named ports** or **numeric ports**.

| Type | Example | Behavior |
|------|---------|----------|
| Named Port | `:http` or `:http-alt` | Go looks up port number in `/etc/services` file; returns error if not found |
| Numeric Port | `:4000` | Directly uses the specified port number |

#### Host Specification
When omitting the host (e.g., `:4000`), the server will listen on **all available network interfaces** on your computer.

```go
// Listens on all network interfaces
server := http.ListenAndServe(":4000", mux)
```

---

## HTTP Architecture Components

### Controllers
**Role:** Execute application logic and write HTTP responses

**Responsibilities:**
- Process requests
- Execute business logic
- Write HTTP response headers
- Write HTTP response bodies

### Router (ServeMux)
**Definition:** Stores a mapping between URL patterns and their corresponding handlers

**Key Points:**
- Usually one servemux per application
- Contains all routes for your application
- Go terminology: called a "servemux" instead of "router"

### Web Server
**Advantage of Go:** You can create and manage a web server directly within your application

- No external third-party server required (unlike Nginx or Apache)
- Built-in HTTP server capabilities

### HTTP ResponseWriter
**Type:** `http.ResponseWriter` parameter in handler functions

**Purpose:** Provides methods for assembling and sending HTTP responses to users

---

## URL Pattern Matching

### Two Pattern Types

#### 1. Fixed Path Patterns
- **Format:** Do NOT end with a trailing slash
- **Example:** `/snippet/view`, `/snippet/create`
- **Matching:** Only matches when request URL path **exactly matches** the fixed path
- **Behavior:** Exact matching only

#### 2. Subtree Path Patterns
- **Format:** END with a trailing slash
- **Example:** `/`, `/static/`
- **Matching:** Matches whenever the start of the request URL path **matches the subtree path**
- **Behavior:** Acts like a wildcard pattern (e.g., `/**` or `/static/**`)
- **Catch-all:** The `/` pattern matches any request path

| Pattern Type | Example | Matches | Does Not Match |
|--------------|---------|---------|----------------|
| Fixed Path | `/snippet/view` | `/snippet/view` only | `/snippet/view/` |
| Subtree Path | `/static/` | `/static/`, `/static/css`, `/static/css/style.css` | `/static` (no trailing slash) |
| Root Catch-all | `/` | Any path | Nothing (catches all) |

---

## DefaultServeMux

### ⚠️ Security Warning
`DefaultServeMux` is a **global variable** accessible to any package, including third-party packages.

**Security Risk:**
- Any package can register routes on the global DefaultServeMux
- Compromised third-party packages could expose malicious handlers
- **Recommendation:** Avoid using DefaultServeMux in production applications

**Best Practice:** Create your own custom servemux instance instead.

---

## HTTP Response Headers & Status Codes

### Writing Response Headers

#### Automatic Status Code
- If you don't explicitly call `w.WriteHeader()`, the first call to `w.Write()` automatically sends **200 OK**
- This is the default behavior

#### Explicit Status Code
- To send a **non-200 status code**, you MUST call `w.WriteHeader()` before any `w.Write()` call
- Order matters: headers must be written before the body

```go
w.WriteHeader(http.StatusNotFound) // Must come before w.Write()
w.Write([]byte("Not found"))
```

### HTTP Error Helper Function

#### http.Error()
**Purpose:** Lightweight shortcut for sending error responses

**What It Does Behind the Scenes:**
- Calls `w.WriteHeader()` with your specified status code
- Calls `w.Write()` with your error message

```go
http.Error(w, "Page not found", http.StatusNotFound)
// Equivalent to:
// w.WriteHeader(http.StatusNotFound)
// w.Write([]byte("Page not found"))
```

---

## File Serving & Security

### Serving Static Files with http.ServeFile()

#### ⚠️ Critical Security Warning
`http.ServeFile()` does **NOT automatically sanitize file paths**.

**Vulnerability:** Directory Traversal Attacks
- Untrusted user input can be manipulated to access files outside intended directory
- Attacker could use `../` to navigate to sensitive files

**Mitigation:** Always sanitize input with `filepath.Clean()`

```go
// ❌ UNSAFE - Vulnerable to directory traversal
filePath := userInput // e.g., "../../../etc/passwd"
http.ServeFile(w, r, filePath)

// ✅ SAFE - Sanitized with filepath.Clean()
filePath := filepath.Clean(userInput)
http.ServeFile(w, r, filePath)
```

**Remember:** If constructing file paths from any untrusted user input, always use `filepath.Clean()` first.

---

## Concurrency & Race Conditions

### Request Handling

#### Goroutine-Based Concurrency
**Important:** All incoming HTTP requests are served in their own goroutine.

**Implications:**
- Go is **blazingly fast** due to concurrent request handling
- Multiple handlers can execute simultaneously
- Code in or called by handlers will likely run concurrently on busy servers

#### ⚠️ Race Condition Risk
When accessing **shared resources** from handlers, you must be aware of and protect against race conditions.

**Key Consideration:** If multiple goroutines access the same data simultaneously without synchronization, data corruption or unexpected behavior can occur.

**Best Practice:** Use synchronization primitives (mutexes, channels) when sharing data between handlers.

---

## Managing Configuration Settings

### Command-Line Flags

#### Basic Usage
Define flags using the `flag` package to accept command-line arguments.

```go
addr := flag.String("addr", ":4000", "HTTP network address")
flag.Parse() // Must call to parse the flags

// Usage: $ go run ./cmd/web -addr=":80"
```

**Benefits:**
- Flexible runtime configuration
- Default values provided
- Automatic `-help` flag support
- Type-safe (flag.String, flag.Int, flag.Bool, etc.)

#### Automated Help
Use the `-help` flag to list all available command-line flags and their descriptions.

```bash
$ go run ./cmd/web -help
Usage of /tmp/go-build3672328037/b001/exe/web:
	-addr string
	HTTP network address (default ":4000")
```

### Port Restrictions

#### Privileged Ports (0-1023)
- **Restriction:** Ports 0-1023 are reserved for system services
- **Access:** Typically requires root/administrator privileges
- **Error:** If you try to use a restricted port without privileges, you'll get: `bind: permission denied`

```go
// ❌ Will fail on most systems without root
addr := flag.String("addr", ":80", "HTTP network address")

// ✅ Safe - Use ports above 1023
addr := flag.String("addr", ":4000", "HTTP network address")
```

---

## Alternative Configuration Methods

### Environment Variables

#### Usage
```go
addr := os.Getenv("SNIPPETBOX_ADDR")

// Usage: $ SNIPPETBOX_ADDR=":80" go run ./cmd/web
```

#### Drawbacks vs Command-Line Flags
- ❌ No default value support (returns empty string if not set)
- ❌ No automatic `-help` functionality
- ❌ Less convenient for users

**Recommendation:** Use command-line flags for better user experience.

---

## Configuration Structs with flag.StringVar()

### Pre-Existing Variable Pattern
Parse flag values directly into the memory addresses of pre-existing variables.

**Useful for:** Storing all configuration settings in a single struct

```go
type config struct {
    addr      string
    staticDir string
}

var cfg config
flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
flag.Parse()

// Access configuration: cfg.addr, cfg.staticDir
```

### Available Var Functions
| Function | Type | Example |
|----------|------|---------|
| `flag.StringVar()` | string | `flag.StringVar(&s, "name", "default", "help")` |
| `flag.IntVar()` | int | `flag.IntVar(&i, "port", 4000, "help")` |
| `flag.BoolVar()` | bool | `flag.BoolVar(&b, "debug", false, "help")` |

### Benefits
- Single struct for all configuration
- Clean, organized code
- Easy to pass configuration around your application
- Type-safe configuration management

---

## Quick Reference Checklist

- [ ] Module path matches project repository URL
- [ ] Using custom servemux, not DefaultServeMux
- [ ] Fixed paths used for exact matches (no trailing slash)
- [ ] Subtree paths used for wildcard-like matching (trailing slash)
- [ ] Response headers written before response body
- [ ] File paths sanitized with `filepath.Clean()` if from user input
- [ ] Named ports checked against `/etc/services` when used
- [ ] Aware of goroutine concurrency and potential race conditions
- [ ] Using command-line flags for configuration (not environment variables)
- [ ] Avoiding privileged ports (0-1023) unless running with root privileges
- [ ] Configuration stored in struct using flag.*Var() functions