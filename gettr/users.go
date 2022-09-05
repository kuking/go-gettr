package gettr

import (
	"errors"
	"fmt"
	"github.com/dghubble/sling"
)

// User is an Gettr user details with description, followers/followings count, etc.
type User struct {
	Description string `json:"dsc"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	Following   int    `json:"flw"`
	Followers   int    `json:"flg"`
	Language    string `json:"lang"`
	UpdateDate  uint64 `json:"udate"`
	CreateDate  uint64 `json:"cdate"`
	ID          string `json:"_id"`
}

// UserService is an API for interacting with Users details
type UserService struct {
	sling  *sling.Sling
	client *Client
}

type loginQuery struct {
	Content loginContent `json:"content"`
}

type loginContent struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Sms   string `json:"sms"`
}

type cursorQueryParameters struct {
	Max     uint   `url:"max,omitempty"`
	Include string `url:"incl,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
}

func newUserService(sling *sling.Sling, client *Client) *UserService {
	return &UserService{
		sling:  sling,
		client: client,
	}
}

// Info retrieves an user detail bu username (not ID)
func (s *UserService) Info(username string) (*User, error) {
	result := new(resultData)
	user := new(User)
	result.Data = resultDataAux{user, resultAuxiliary{Users: nil, Cursor: nil}}
	apiErrorWrap := new(aPIErrorWrap)
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
	result := resultData{resultDataAux{Data: nil, Aux: resultAuxiliary{}}}
	apiError := aPIErrorWrap{Payload: APIError{}}

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
				Max:     50, // current limit is 20 but leaving 50 maybe it becomes more lax
				Include: "userinfo",
				Cursor:  cursor,
			}).
		Receive(&result, &apiError)

	if err != nil {
		return nil, relevantError(resp, err, apiError.Payload)
	}
	usersCursor := UsersCursor{us: s, UserID: id, what: cursorType}
	switch result.Data.Aux.Cursor.(type) {
	case string:
		usersCursor.cursor = result.Data.Aux.Cursor.(string)
	case float64:
		usersCursor.cursor = ""
	default:
		panic(errors.New("unexpected cursor type"))
	}

	for _, v := range result.Data.Aux.Users {
		usersCursor.Users = append(usersCursor.Users, v)
	}
	return &usersCursor, relevantError(resp, err, apiError.Payload)
}

// Follows the provided username, credentials must be set with Client.SetAuthToken
func (s *UserService) Follows(username string) error {
	return s.buildAndDoUnFollowsRequest("follows", username)
}

// Unfollows the provided username, credentials must be set with Client.SetAuthToken
func (s *UserService) Unfollows(username string) error {
	return s.buildAndDoUnFollowsRequest("unfollows", username)
}

func (s *UserService) buildAndDoUnFollowsRequest(action, username string) error {
	req, err := s.sling.New().Post("u/user/").
		Path(s.client.username+"/").
		Path(action+"/").Path(username).
		Set("x-app-auth", s.client.authHeader).
		Request()
	if err != nil {
		return err
	}
	res, err := s.client.httpClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode == 429 {
		return errors.New("Rate-Limit [http_result: 429]")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected error [http_result: %v]", res.StatusCode)
	}
	return nil
}

func (s *UserService) Login(email, sms, pwd string) error {
	result := resultLogin{resultLoginPayload{User: User{}, Token: "", Rtoken: ""}}
	request := loginQuery{Content: loginContent{Email: email, Pwd: pwd, Sms: sms}}
	apiError := new(aPIErrorWrap)

	req, err := s.sling.New().
		BodyJSON(request).
		Post("/u/user/v2/login").
		Receive(&result, &apiError)
	defer req.Request.Body.Close()
	if err != nil {
		return err
	}
	if req.StatusCode >= 400 && req.StatusCode < 500 {
		return fmt.Errorf("Acccess error [http_result %v]", req.StatusCode)
	}
	s.client.SetAuthToken(result.Result.User.Username, result.Result.User.ID, result.Result.Token)
	return nil
}
