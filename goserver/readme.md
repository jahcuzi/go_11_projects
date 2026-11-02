# Notes

## My Notes

- `/form` is not the same as `/form.html`
- `/form` gets methods and we can allow or disallow in the function itself using `switch` case: 

```go
    switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, <file_name>)
	case http.MethodPost:
		// Handle Method
	default:
        // only allows get and post fucntions
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
```

- `http.FileServer(http.Dir("./static"))` only serves files when it is requested in the `url` itself i.e `/form.html` is handled by file Server but `/form` is not.



## AI suggestions

1. Instead of calling `http.ListenAndServe` directly, configure an `http.Server` with timeouts:

```go
server := &http.Server{
	Addr:              ":8080",
	Handler:           nil, // Uses http.DefaultServeMux.
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       10 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       60 * time.Second,
}

log.Fatal(server.ListenAndServe())
```

2. Tests using `httptest`

## Thoughts on AI suggestions

1. If server settings need to come from environment variables or be reused across multiple servers, group them in an application `Config` struct. For a single small server, configuring `http.Server` directly is usually simpler.

2. Done.

- `httptest.NewRequest(method, url, body)` gives a request
- `httptest.NewRecorder()` is used in place of a writer to record responses and that can tested
