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
2. Done.
