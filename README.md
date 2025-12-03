# Let's Go by Alex Edwards 

Quick reference guide for concepts, patterns, and best practices

---
## Table of Contents

- [Module Paths & Naming Convention](#module-paths--naming-convention)
  - [Module Path Basics](#module-path-basics)
  - [Naming Guidelines](#naming-guidelines)
- [Running Go](#running-go)
  - [Equivalent Commands](#equivalent-commands)
- [Networking Fundamentals](#networking-fundamentals)
  - [Network Addresses](#network-addresses)
- [HTTP Architecture Components](#http-architecture-components)
  - [Controllers](#controllers)
  - [Router (ServeMux)](#router-servemux)
  - [Web Server](#web-server)
  - [HTTP ResponseWriter](#http-responsewriter)
- [URL Pattern Matching](#url-pattern-matching)
  - [Two Pattern Types](#two-pattern-types)
- [DefaultServeMux](#defaultservemux)
- [HTTP Response Headers & Status Codes](#http-response-headers--status-codes)
  - [Writing Response Headers](#writing-response-headers)
  - [HTTP Error Helper Function](#http-error-helper-function)
- [File Serving & Security](#file-serving--security)
  - [Serving Static Files with http.ServeFile()](#serving-static-files-with-httpservefile)
- [Concurrency & Race Conditions](#concurrency--race-conditions)
  - [Request Handling](#request-handling)
- [Managing Configuration Settings](#managing-configuration-settings)
  - [Command-Line Flags](#command-line-flags)
  - [Port Restrictions](#port-restrictions)
- [Alternative Configuration Methods](#alternative-configuration-methods)
  - [Environment Variables](#environment-variables)
- [Configuration Structs with flag.StringVar()](#configuration-structs-with-flagstringvar)
  - [Pre-Existing Variable Pattern](#pre-existing-variable-pattern)
  - [Available Var Functions](#available-var-functions)
  - [Benefits](#benefits)
- [Logging](#logging)
  - [Leveled Logging](#leveled-logging)
  - [Logging Methods Reference](#logging-methods-reference)
  - [Decoupled Logging](#decoupled-logging)
  - [Concurrent Logging](#concurrent-logging)
  - [Logging to a File](#logging-to-a-file)
- [Dependency Injection](#dependency-injection)
  - [Overview](#overview)
  - [Single Package Pattern (Recommended)](#single-package-pattern-recommended)
  - [Multi-Package Pattern (Closures)](#multi-package-pattern-closures)
  - [Pattern Comparison](#pattern-comparison)
- [Quick Reference Checklist](#quick-reference-checklist)


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

## Logging

### Leveled Logging

#### Two Custom Loggers (INFO & ERROR)
Create separate loggers for different log levels using `log.New()`:

```go
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

infoLog.Printf("Starting server on %s", *addr)
err := http.ListenAndServe(*addr, mux)
errorLog.Fatal(err)
```

**Parameters:**
- **First argument:** Destination (os.Stdout or os.Stderr)
- **Second argument:** Prefix (e.g., "INFO\t", "ERROR\t")
- **Third argument:** Flags for formatting (see table below)

#### Log Formatting Flags
| Flag | Description |
|------|-------------|
| `log.Ldate` | Include date (YYYY-MM-DD) |
| `log.Ltime` | Include time (HH:MM:SS) |
| `log.Lshortfile` | Include short filename (e.g., main.go:42) |
| `log.Llongfile` | Include full file path (e.g., /home/user/project/main.go:42) |
| `log.LUTC` | Use UTC datetimes instead of local time |

#### Combining Flags
Use the bitwise OR operator (`|`) to combine multiple flags:

```go
log.Ldate | log.Ltime | log.Lshortfile  // Date, time, and short filename
```

---

### Logging Methods Reference

| Method | Output | Use Case | Notes |
|--------|--------|----------|-------|
| `Print()` / `Printf()` / `Println()` | Normal message | General logging | No special behavior |
| `Fatal()` / `Fatalf()` / `Fatalln()` | Error message, then exit | Critical errors | Calls `os.Exit(1)` after logging |
| `Panic()` / `Panicf()` / `Panicln()` | Error message, then panic | Exceptional conditions | Panics after logging; can be recovered |

**Best Practice:** Avoid using `Panic()` and `Fatal()` outside of `main()` function. Return errors instead and only exit/panic from main.

---

### Decoupled Logging

#### Benefits of Logging to Standard Streams
Logging to stdout and stderr **decouples** your application from log storage and routing.

**Advantages:**
- Application doesn't need to manage log files
- Easy to redirect logs to different destinations depending on environment
- Follows Unix philosophy (small, focused tools)

#### Redirecting Logs at Runtime
Redirect stdout and stderr to files when starting the application:

```bash
$ go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
```

**Breakdown:**
- `>>` - Append to file (create if doesn't exist)
- `1` - stdout (implicit, so `>>/tmp/info.log` is same as `1>>/tmp/info.log`)
- `2` - stderr
- `/tmp/info.log` - INFO logger output destination
- `/tmp/error.log` - ERROR logger output destination

---

### Concurrent Logging

#### Thread-Safety
Custom loggers created by `log.New()` are **concurrency-safe**. You can safely:
- Share a single logger across multiple goroutines
- Use the same logger in all handlers
- No need to worry about race conditions on the logger itself

```go
// ✅ Safe to use across multiple goroutines
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Use in multiple handlers concurrently
handler1(w, r) { infoLog.Println("Handler 1") }
handler2(w, r) { infoLog.Println("Handler 2") }
```

#### Multiple Loggers, Same Destination
**⚠️ Important:** If you have **multiple loggers writing to the same destination**, ensure the destination's `Write()` method is also safe for concurrent use.

**Example of Potential Issue:**
```go
// ⚠️ Risky - Multiple loggers to same file
f, _ := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE, 0666)
log1 := log.New(f, "LOG1\t", log.Ldate|log.Ltime)
log2 := log.New(f, "LOG2\t", log.Ldate|log.Ltime)
// log1 and log2 may interfere with each other
```

---

### Logging to a File

#### Manual File-Based Logging
Alternative to redirecting streams at runtime - open a file directly in Go.

```go
f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
if err != nil {
    log.Fatal(err)
}
defer f.Close()

infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
```

#### File Flags Explanation
| Flag | Meaning |
|------|---------|
| `os.O_RDWR` | Open file for reading and writing |
| `os.O_CREATE` | Create file if it doesn't exist |
| `0666` | File permissions (rw-rw-rw-) |

#### Recommendation
**General Best Practice:** Log to standard streams and redirect at runtime (more flexible). Only log directly to files if you have specific requirements.

---

## Dependency Injection

### Overview
Most web applications need multiple dependencies that handlers access, such as:
- Database connection pools
- Centralized error handlers
- Template caches
- Logging utilities

**Problem with Global Variables:** Less explicit, more error-prone, harder to unit test.

**Solution:** Inject dependencies into handlers through structured patterns.

---

### Single Package Pattern (Recommended)

#### Using an Application Struct
For applications where all handlers are in the same package, define an application struct containing all dependencies and define handler functions as **methods** on that struct.

```go
type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger
}
```

#### In main() Function
Create and initialize your application struct with all dependencies:

```go
func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()
    
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    
    // Initialize application struct with dependencies
    app := &application{
        errorLog: errorLog,
        infoLog:  infoLog,
    }
    
    // Pass to server...
}
```

#### Defining Handler Methods
Define handlers as methods on the application receiver:

```go
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Access dependencies via app receiver
    app.infoLog.Println("Creating a new snippet...")
}
```

#### Benefits
- ✅ Dependencies explicitly available via struct fields
- ✅ Less error-prone than global variables
- ✅ Easier to unit test (mock dependencies)
- ✅ Clean, organized code
- ✅ Type-safe

---

### Multi-Package Pattern (Closures)

#### Problem
The struct method pattern doesn't work when handlers are spread across multiple packages.

#### Solution: Closure-Based Injection
Create a config package exporting an Application struct and have handler functions **close over** this to form a closure.

```go
// config/application.go
package config

type Application struct {
    ErrorLog *log.Logger
    InfoLog  *log.Logger
}
```

#### In main() Function
Create the application and pass it to handler factory functions:

```go
func main() {
    app := &config.Application{
        ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
        InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
    }
    
    mux := http.NewServeMux()
    mux.Handle("/", examplePackage.ExampleHandler(app))
}
```

#### Handler Factory Function
Create a function that accepts dependencies and returns an `http.HandlerFunc`:

```go
// examplePackage/handlers.go
package examplePackage

func ExampleHandler(app *config.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // The handler closes over 'app', accessing dependencies
        ts, err := template.ParseFiles(files...)
        if err != nil {
            app.ErrorLog.Println(err.Error())
            http.Error(w, "Internal Server Error", 500)
            return
        }
        
        // Handler logic...
    }
}
```

#### How It Works
1. Handler factory function accepts dependencies (`*config.Application`)
2. Returns an `http.HandlerFunc` closure
3. The closure **captures** the `app` variable in its scope
4. When the handler runs, it has access to all dependencies

#### Benefits
- ✅ Works across multiple packages
- ✅ Dependencies still injected (not global)
- ✅ Factory function pattern is flexible
- ✅ Each handler can have different dependencies

#### Pattern Comparison

| Aspect | Single Package Struct | Multi-Package Closure |
|--------|----------------------|----------------------|
| Works across packages? | ❌ No | ✅ Yes |
| Code clarity | ✅ Very clear | ✅ Clear |
| Testing | ✅ Easy | ✅ Easy |
| Setup complexity | ✅ Simple | ⚠️ More setup |
| Use case | Small to medium apps | Larger projects |

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
- [ ] Using leveled logging (separate loggers for INFO and ERROR)
- [ ] Logging to stdout/stderr with runtime redirection
- [ ] Aware of logger concurrency safety and multiple logger issues
- [ ] Using dependency injection pattern (struct methods or closures)
- [ ] Dependencies passed explicitly, not via global variables
