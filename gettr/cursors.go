package gettr

import (
	"fmt"
)

type userCursorType int

const (
	followers userCursorType = iota
	following userCursorType = iota
)

// UsersCursor is an iterable Cursor of Users
type UsersCursor struct {
	us     *UserService
	what   userCursorType
	cursor string
	UserID string
	Users  []User
}

// Iter allows a Cursor to be iterated in a for loop line any collection, abstracting out the API calls
func (uc *UsersCursor) Iter(limit int) <-chan User {
	var ch = make(chan User)
	go func() {
		var count = 0
		var err error
		var cursor = uc
		for {
			for _, user := range cursor.Users {
				ch <- user
				if count++; count >= limit && limit != -1 {
					break
				}
			}
			if count >= limit && limit != -1 {
				break
			}
			if !cursor.HasNext() {
				break
			}
			cursor, err = cursor.Next()
			if err != nil {
				fmt.Printf("Failed Cursor: %v", err)
				break
			}
		}
		close(ch)
	}()
	return ch
}

// Next returns the next Cursor in a paginated resul-set
func (uc *UsersCursor) Next() (*UsersCursor, error) {
	return uc.us.userCursor(uc.UserID, uc.cursor, uc.what)
}

// HasNext returns if there is any other page left in this result-set
func (uc *UsersCursor) HasNext() bool {
	return uc.cursor != ""
}
