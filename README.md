
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
