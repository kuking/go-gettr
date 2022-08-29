# go-gettr: Golang Gettr client

work in progress -- some things working, under heavy development.
Come back in a few days, or dive in the code.

API is easy to use as the paginated calls are abstracted out. i.e.

```go
package main

import (
	"fmt"
	"github.com/kuking/go-gettr/gettr"
	"net/http"
)

func main() {
	var client = gettr.NewClient(http.DefaultClient)

	user, err := client.User.Info("your-user")
	if err != nil {
		panic(err)
	}

	followers, err := client.User.Followers(user.ID)
	if err != nil {
		panic(err)
	}
	for user := range followers.Iter(50) {
		fmt.Println(user.ID)
	}

	// follow-unfollows
	user, err = client.User.Info("your-user")
	client.SetAuthToken(user.Username, user.ID, "browsers' local_storage LS_SESSION_INFO.userinfo.token")
	err = client.User.Follows("support")
	err = client.User.Unfollows("support")
	
}
```
