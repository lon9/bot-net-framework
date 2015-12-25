package main
import (
	"github.com/ChimeraCoder/anaconda"
	"runtime"
	"fmt"
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
func (tc *TalkController) PostOne() (Tweets, error){

	finishCh := make(chan bool, runtime.NumCPU())
	errCh := make(chan error, runtime.NumCPU())
	numTweet := len(tc.Talk.Tweets[tc.Seq])


	for _,v  := range tc.Talk.Tweets[tc.Seq]{
		go postTweet(v, finishCh, errCh)
	}

	count := 0
	var err error
	L1:
	for{
		select {
		case <-finishCh:
			count++
			if count == numTweet {
				break L1
			}
		case err = <-errCh:
			fmt.Println(err)
			break L1
		default:
		}
	}

	// Increment sequence
	tc.Seq++
	return tc.Talk.Tweets[tc.Seq-1], err
}

func postTweet(tweet Tweet, resultCh chan bool, errCh chan error) {
	api := anaconda.NewTwitterApi(tweet.Bot.AccessToken, tweet.Bot.AccessTokenSecret)

	_, err := api.PostTweet(tweet.Text, nil)
	if err != nil {
		errCh <- err
		return
	}
	resultCh <- true
}