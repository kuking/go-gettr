package gettr

import (
	"github.com/dghubble/sling"
)

// User is an Gettr user details with description, followers/followings count, etc.
type User struct {
	Description string `json:"dsc"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	Following   uint32 `json:"flw"`
	Followers   uint32 `json:"flg"`
	Language    string `json:"lang"`
	UpdateDate  uint64 `json:"udate"`
	CreateDate  uint64 `json:"cdate"`
	ID          string `json:"_id"`
}

// UserService is an API for interacting with Users details
type UserService struct {
	sling *sling.Sling
}

type cursorQueryParameters struct {
	Max     uint16 `url:"max,omitempty"`
	Include string `url:"incl,omitempty"`
	Cursor  string `url:"Cursor,omitempty"`
}

func newUserService(sling *sling.Sling) *UserService {
	return &UserService{
		sling: sling,
	}
}

// Info retrieves an user detail bu username (not ID)
func (s *UserService) Info(username string) (*User, error) {
	result := new(result)
	user := new(User)
	result.Data = resultData{user, resultAuxiliary{Users: nil, Cursor: nil}}
	apiErrorWrap := aPIErrorWrap{Payload: APIError{}}
	resp, err := s.sling.New().
		Path("s/uinf/").Path(username).
		Receive(result, apiErrorWrap)
	return user, relevantError(resp, err, apiErrorWrap.Payload)
}

// Followers returns a Cursor to the list of followers for the user id provided
func (s *UserService) Followers(id string) (*UsersCursor, error) {
	return s.userCursor(id, "", followers)
}

// Following returns a Cursor to the list of Users being followed for the user id provided
func (s *UserService) Following(id string) (*UsersCursor, error) {
	return s.userCursor(id, "", following)
}

func (s *UserService) userCursor(id string, cursor string, cursorType userCursorType) (*UsersCursor, error) {
	result := new(result)
	result.Data = resultData{Data: new(map[string]interface{}), Aux: resultAuxiliary{}}
	apiError := new(aPIErrorWrap)

	var queryTypePath string
	switch cursorType {
	case following:
		queryTypePath = "followings"
	case followers:
		queryTypePath = "followers"
	default:
		panic("I don't know how to lookup for this.")
	}

	resp, err := s.sling.New().
		Path("u/user/").Path(id+"/"). // if slash is not added it gets removed by next one
		Path(queryTypePath).
		QueryStruct(
			cursorQueryParameters{
				Max:     20,
				Include: "userinfo",
				Cursor:  cursor,
			}).
		Receive(result, apiError)

	if err != nil {
		return nil, relevantError(resp, err, apiError.Payload)
	}
	usersCursor := UsersCursor{us: s, UserID: id, what: cursorType}
	switch result.Data.Aux.Cursor.(type) {
	case string:
		usersCursor.cursor = result.Data.Aux.Cursor.(string)
	}

	for _, v := range result.Data.Aux.Users {
		usersCursor.Users = append(usersCursor.Users, v)
	}
	return &usersCursor, relevantError(resp, err, apiError.Payload)
}
