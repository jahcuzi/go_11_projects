# Notes

## My Notes

- Using Gorilla Mux, we can restrict methods that can be used and route it to a functions i.e `Router`
- ~~random Ids can have collisions but it is okay for now as practice~~ UUID V7 used.
- Get url params using `mux.Vars(r)`
- Might need to better error handling for 
    - everytime json is encoded or decoded there can be an error
- Update can use delete but since it is a http function handler, I have to send a request to actually delete the movie
    - write the delete logic again. OR
    - The compromise is to separate the logic from the handlers
- `w http.ResponseWriter` is the only way to send stuff to the client
- Encoding errors must be logged not responsed with `http.Error` as Encoding can fail due to failing to write to the response.
- encoding can also fail before writing anything if the value can’t be represented as JSON. But with movie structs, that’s unlikely. 
- In a larger application, we might encode to a buffer first; if that fails, we can still send a 500 before writing the response. But if writing the buffer and then seding to the client fails, logging is still the right thing to do.

## AI suggestions

1. protect `movies` from concurrent users using `sync.RWMutex`
2. Handle function tests

## My thoughts on AI suggestions

1. Done
    - `sync.Mutex` closes both Read and write but `sync.RWMutex` we can specify which lock should be used
    - The code execution between `Lock()` and `Unlock()` should be minimal as it will block execution of other executors
        - Avoid Async ops, if possible
        - extract variables by copying, rather than doing operations in the critical section
        - *NO* `return` in the critical section, Deadlock occurs
        - Make sure `Unlock()` is reached at some point
2. Done.
    - Create a reader of string literal `s` with `string.NewReader(s)`
    - Got tried of checking every subfields in the response, using `reflect.DeepEqual(want, got)` for concise coding
    - Tests use `setMovies(t, testMovies)` to replace the in-memory movie database with test data, then restore its previous contents during test cleanup.
    - `t.Helper` if an assertion fails inside it, it points to the test that called it not the fucntion it resides in.
    - `t.Cleanup(func(){...})` runs the function insdie when test finishes.
