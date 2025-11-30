
This readme contains my notes and highlight from the book, somewhere I can really fast check meaning of terms or refresh a certain aspect that I may need without having to open the book
## Module paths for downloadable packages | Naming Convention
If you’re creating a project which can be downloaded and used by other people and
programs, then it’s good practice for your module path to equal the location that the code
can be downloaded from.

For instance, if your package is hosted at https://github.com/foo/bar then the module path
for the project should be github.com/foo/bar.

you want to pick a module path that is globally unique and unlikely to be
used by anything else. In the Go community, a common convention is to base your module
paths on a URL that you own.
In my case, a clear, succinct and unlikely-to-be-used-by-anything-else module path for this
project would be `snippetbox.alexedwards.net`


## Running Go

- The three following commands are all equivalent:
```go
$ go run .
$ go run main.go
$ go run snippetbox.alexedwards.ne
```


## Networking

- In other Go projects or documentation you might sometimes see network addresses written using named ports like ":http" or ":http-alt" instead of a number. If you use a named port then Go will attempt to look up the relevant port number from your /etc/services file when starting the server, or will return an error if a match can’t be found.
#### Controllers:
They’re responsible for executing your
application logic and for writing HTTP response headers and bodies.

#### Router
(or servemux in Go terminology). 
This stores a mapping between the URL patterns for your application and the corresponding handlers. Usually
you have one servemux for your application containing all your routes.

#### Web server
One of the great things about Go is that you can
establish a web server and listen for incoming requests as part of your application itself.
You don’t need an external third-party server like Nginx or Apache.

`http.ResponseWriter` parameter provides methods for assembling a HTTP response
and sending it to the user

! If you omit the host (like we did with ":4000") then the server will listen on all
your computer’s available network interfaces. 

### Fixed path and subtree patterns
Now that the two new routes are up and running let’s talk a bit of theory.
Go’s servemux supports two different types of URL patterns: fixed paths and subtree paths.
Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.
Our two new patterns — "/snippet/view" and "/snippet/create" — are both examples of
fixed paths. In Go’s servemux, fixed path patterns like these are only matched (and the
corresponding handler called) when the request URL path exactly matches the fixed path.
In contrast, our pattern "/" is an example of a subtree path (because it ends in a trailing
slash). Another example would be something like "/static/". Subtree path patterns are
matched (and the corresponding handler called) whenever the start of a request URL path
matches the subtree path. If it helps your understanding, you can think of subtree paths as
acting a bit like they have a wildcard at the end, like "/**" or "/static/**".
This helps explain why the "/" pattern is acting like a catch-all. The pattern essentially means
match a single slash, followed by anything (or nothing at all)

#### DefaultServeMux
Because DefaultServeMux is a global variable, any package can access it and register a route
— including any third-party packages that your application imports. If one of those third-
party packages is compromised, they could use DefaultServeMux to expose a malicious
handler to the web

- If you don’t call w.WriteHeader() explicitly, then the first call to w.Write() will automatically send a 200 OK status code to the user. So, if you want to send a non-200 status code, you must call w.WriteHeader() before any call to w.Write().

- http.Error() shortcut. This is a
lightweight helper function which takes a given message and status code, then calls the
w.WriteHeader() and w.Write() methods behind-the-scenes for us.