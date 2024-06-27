package twitter

type Retweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

type MockedAPI struct{}

func (a MockedAPI) GetParsedRetweets() ([]Retweet, error) {
	ret := []Retweet{
		{
			User: struct {
				ScreenName string `json:"screen_name"`
			}{
				ScreenName: "user1",
			},
		},
		{
			User: struct {
				ScreenName string `json:"screen_name"`
			}{
				ScreenName: "user2",
			},
		},
		{
			User: struct {
				ScreenName string `json:"screen_name"`
			}{
				ScreenName: "user3",
			},
		},
		{
			User: struct {
				ScreenName string `json:"screen_name"`
			}{
				ScreenName: "user4",
			},
		},
	}
	return ret, nil
}
