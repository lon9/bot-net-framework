package main
import (
	"github.com/ChimeraCoder/anaconda"
)

type TalkController struct {
	Talk Talk
	Seq int
}

func NewTalkController(talk Talk)*TalkController{
	return &TalkController{
		Talk:talk,
		Seq:0,
	}
}


func (tc *TalkController) PostOne() (Tweet, error){
	api := anaconda.NewTwitterApi(tc.Talk.Tweets[tc.Seq].Bot.AccessToken, tc.Talk.Tweets[tc.Seq].Bot.AccessTokenSecret)

	// Post tweet
	_, err := api.PostTweet(tc.Talk.Tweets[tc.Seq].Text, nil)

	// Increment sequence
	tc.Seq++
	return tc.Talk.Tweets[tc.Seq-1], err
}