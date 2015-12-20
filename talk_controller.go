package main
import (
	"github.com/ChimeraCoder/anaconda"
)

// TalkController is controller of talk.
type TalkController struct {
	Talk Talk
	Seq int
}

// NewTalkController is constructor of TalkController.
func NewTalkController(talk Talk)*TalkController{
	return &TalkController{
		Talk:talk,
		Seq:0,
	}
}

// PostOne posts one tweet and inclement sequence.
func (tc *TalkController) PostOne() (Tweet, error){
	api := anaconda.NewTwitterApi(tc.Talk.Tweets[tc.Seq].Bot.AccessToken, tc.Talk.Tweets[tc.Seq].Bot.AccessTokenSecret)

	// Post tweet
	_, err := api.PostTweet(tc.Talk.Tweets[tc.Seq].Text, nil)

	// Increment sequence
	tc.Seq++
	return tc.Talk.Tweets[tc.Seq-1], err
}