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
	talkName := req.FormValue("talkName")
	if talkName == ""{
		r.JSON(400, Error{400, "talkName must be set."})
		return
	}

	var talk Talk
	db.Where("title = ?", talkName).First(&talk)
	db.Model(&talk).Order("sequence", true).Related(&talk.Tweets)
	for i, _ := range talk.Tweets{
		db.Model(&talk.Tweets[i]).Related(&talk.Tweets[i].Bot)
	}

	if talk.ID == 0{
		r.JSON(400, Error{404, "Talk not found."})
		return
	}

	talkController := NewTalkController(talk)

	// Start talk
	for  range talk.Tweets{
		_, err := talkController.PostOne()
		if err != nil {
			break
		}
	}

	r.JSON(200, talk)
}

// StartTalkSocket is websocket for talking in Web UI.
func StartTalkSocket(r render.Render, w http.ResponseWriter, req *http.Request, db gorm.DB){

	ws, err := websocket.Upgrade(w, req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}


	talkName := req.FormValue("talkName")
	if talkName == ""{
		r.JSON(400, Error{400, "talkName must be set."})
		return
	}

	var talk Talk
	db.Where("title = ?", talkName).First(&talk)
	db.Model(&talk).Order("sequence", true).Related(&talk.Tweets)
	for i, _ := range talk.Tweets{
		db.Model(&talk.Tweets[i]).Related(&talk.Tweets[i].Bot)
	}

	if talk.ID == 0{
		r.JSON(400, Error{404, "Talk not found."})
		return
	}

	talkController := NewTalkController(talk)

	for range talk.Tweets{
		tweet, err := talkController.PostOne()
		if err != nil {
			return
		}
		if err := ws.WriteJSON(tweet); err !=nil{
			r.JSON(400, Error{400,"Cant send message."})
		}
		time.Sleep(1*time.Second)
	}
}