# go-gettr: Golang Gettr client

work in progress -- some things working, under heavy development.
Come back in a few days, or dive in the code.

API is easy to use as the paginated calls are abstracted out. i.e.

```go
var client = gettr.NewClient(http.DefaultClient)

user, err := client.Users.Info("support")
if err != nil {
    panic(err)
}

followers, err := client.Users.Followers(user.ID)
if err != nil {
    panic(err)
}
for user := range followers.Iter(50) {
    fmt.Println(user.ID)
}
```
