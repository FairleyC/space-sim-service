#### Passing context between layers
Context is implemented in the services of this application to allow for passing context between layers. This is useful for observability, logging, and security.

```go
func (s *Service) GetCommodity() (ctx context.Context, id string) (Commodity, error) {    
    // Example of using context to pass data between layers
    ctx = context.WithValue(ctx, "request_id", "unique-string")
    fmt.Println("request_id: ", ctx.Value("request_id"))

    ...
}
```

#### Owner-Player-Organization Polymorphism 
[stack overflow article](https://stackoverflow.com/questions/28222533/polymorphism-for-foreign-key-constraints)

For when I need to implement this model.

