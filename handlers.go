package main
import (
	"github.com/martini-contrib/render"
	"github.com/ChimeraCoder/anaconda"
	"net/http"
	"github.com/jinzhu/gorm"
	"net/url"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/martini-contrib/sessions"
	"time"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
	"fmt"
	"strconv"
	"errors"
)

// Index returns top page.
func Index(r render.Render){
	r.HTML(200, "index", "")
}

// TwitterLogin login twitter with OAuth.
func TwitterLogin(r render.Render, s sessions.Session, req *http.Request){
	url, tmpCred, err := anaconda.AuthorizationURL("http://"+ req.Host + "/twitter/callback")
	if err != nil{
		panic(err)
	}
	s.Set("token", tmpCred.Token)
	s.Set("secret", tmpCred.Secret)
	r.Redirect(url, 302)
}

// TwitterCallback is callback of Twitter login.
func TwitterCallback(r render.Render, s sessions.Session, req *http.Request, db gorm.DB){

	// Make instance of Twitter's credentials
	tempCred := &oauth.Credentials{
		Token:s.Get("token").(string),
		Secret:s.Get("secret").(string),
	}

	// Getting accessToken for Twitter API
	tokenCred,_,  err := anaconda.GetCredentials(tempCred, req.FormValue("oauth_verifier"))
	if err != nil {
		panic(err)
	}

	// Delete sessions
	s.Delete("token")
	s.Delete("secret")

	// Make API instance
	api := anaconda.NewTwitterApi(tokenCred.Token, tokenCred.Secret)

	// Getting user info of me
	me, err := api.GetSelf(url.Values{})
	if err != nil {
		panic(err)
	}

	// Create bot instance
	bot := Bot{
		AccessToken:tokenCred.Token,
		AccessTokenSecret:tokenCred.Secret,
		Name:me.Name,
		TwitterId:me.Id,
		ScreenName:me.ScreenName,
		IconURL:me.ProfileImageUrlHttps,
	}

	// Insert new bot
	db.Create(&bot)

	r.HTML(200, "callback", bot)

}

// StartTalk starts talk with Json based API.
func StartTalk(r render.Render, req *http.Request, res http.ResponseWriter, db gorm.DB){

	// Getting talk name from query string.
	talkName := req.FormValue("talkName")
	if talkName == ""{
		r.JSON(400, Error{400, "talkName must be set."})
		return
	}

	// Getting talk from database.
	talk, err := getTalkFromDB(talkName, &db)
	if err != nil {
		r.JSON(400, err)
		return
	}


	// Start talk.
	talkController := NewTalkController(talk, &db)
	for  range talk.Tweets{
		_, err := talkController.PostOne()
		if err != nil {
			r.JSON(400, err)
			break
		}
		time.Sleep(1*time.Second)
	}

	r.JSON(200, talk)
}

// StartTalkSocket is websocket for talking in Web UI.
func StartTalkSocket(r render.Render, w http.ResponseWriter, req *http.Request, db gorm.DB){

	// Getting talk name from query string
	talkName := req.FormValue("talkName")
	if talkName == ""{
		r.JSON(400, Error{400, "talkName must be set."})
		return
	}

	// Upgrading connection to WebSocket.
	ws, err := websocket.Upgrade(w, req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	// Getting talk from database.
	talk, err := getTalkFromDB(talkName, &db)
	if err != nil {
		r.JSON(400, err)
		return
	}

	// Posting tweets
	talkController := NewTalkController(talk, &db)
	for range talk.Tweets{
		tweets, err := talkController.PostOne()
		if err != nil {
			fmt.Println(err)
			r.JSON(400, err)
			return
		}
		for _, v := range tweets{
			if err := ws.WriteJSON(v); err != nil{
				r.JSON(400, Error{400,"Cant send message."})
			}
		}
		time.Sleep(1*time.Second)
	}
}

func getTalkFromDB(talkName string, db *gorm.DB) (talk Talk, err error){

	// Utility function to get talk by talk name.

	// Getting talk
	db.Where("title = ?", talkName).First(&talk)

	// If talk was not found return 400 error.
	if talk.ID == 0{
		err = errors.New("Talk was not found")
		return
	}
	var tweets Tweets

	// Getting tweet sorted by sequence.
	db.Model(&talk).Order("sequence", true).Related(&tweets)
	talk.Tweets = make([]Tweets, len(tweets))
	index  := -1
	prevSeq := 0
	var prevBot Bot
	for _, v := range tweets{
		if v.Sequence == prevSeq {
			v.Bot = prevBot
			talk.Tweets[index] = append(talk.Tweets[index], v)
		}else{
			index++
			db.Model(&v).Related(&prevBot)
			v.Bot = prevBot
			talk.Tweets[index] = make(Tweets, 0)
			talk.Tweets[index] = append(talk.Tweets[index], v)
			prevSeq = v.Sequence
		}
	}

	return

}

// DelTalkTweets delete tweets of a talk from Twitter
func DelTalkTweets(r render.Render, db gorm.DB, req *http.Request){
	talkId := req.FormValue("talkId")
	var tweets Tweets

	db.Where("talk_id = ?", talkId).Find(&tweets)

	resultCh := make(chan Tweet, len(tweets))
	errCh := make(chan error, runtime.NumCPU())
	routineNum := 0
	for i, _ := range tweets{
		if tweets[i].TweetIdStr != ""{
			db.Model(&tweets[i]).Related(&tweets[i].Bot)
			go deleteTweets(tweets[i], resultCh, errCh)
			routineNum++
		}
	}

	finished := 0

	L:
	for{
		select{
		case result := <- resultCh:
			finished++
			db.Save(&result)
			if finished == routineNum{
				break L
			}
		case err := <- errCh:
			fmt.Println(err)
			break L
		default:
			if finished == routineNum{
				break L
			}
		}
	}

	r.JSON(204, nil)
}

func deleteTweets(tweet Tweet, resultCh chan Tweet, errCh chan error){

	// This function is for removing tweets asynchronously.

	api := anaconda.NewTwitterApi(tweet.Bot.AccessToken, tweet.Bot.AccessTokenSecret)


	// Converting TweetIdStr to int64.
	tweetIdInt, err := strconv.ParseInt(tweet.TweetIdStr, 10, 64)
	if err != nil {
		errCh <- err
		return
	}

	// Request Api to remove tweet from Twitter.
	_, err = api.DeleteTweet(tweetIdInt, true)
	if err != nil {
		errCh <- err
		return
	}

	// Resetting tweet id.
	tweet.TweetIdStr = ""
	resultCh <- tweet
}