```
go run ./serviceA

# other terminal
go run ./serviceB

# other terminal
go run ./serviceC

# other terminal
curl http://localhost:8001/foo
```

```
User -> A(:8001/foo) -> B(:8002/bar) -> C(:8003/baz)
```
