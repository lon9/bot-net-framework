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
	"encoding/json"
)


func Index(r render.Render){
	r.HTML(200, "index", "")
}


func TwitterLogin(r render.Render, s sessions.Session, req *http.Request){
	url, tmpCred, err := anaconda.AuthorizationURL("http://"+ req.Host + "/twitter/callback")
	if err != nil{
		panic(err)
	}
	s.Set("token", tmpCred.Token)
	s.Set("secret", tmpCred.Secret)
	r.Redirect(url, 302)
}

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
	}

	// Insert new bot
	db.Create(&bot)

	r.HTML(200, "callback", bot)

}

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

	// Streaming response.
	hj, _:= res.(http.Hijacker)
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	bufrw.WriteString("HTTP/1.1 200 OK\r\n")
	bufrw.WriteString("Content-Type: application/json\r\n")
	bufrw.WriteString("Cache-Control: no-cache\r\n")
	bufrw.WriteString("\r\n")
	bufrw.Flush()

	enc := json.NewEncoder(bufrw)


	// Start talk
	for  range talk.Tweets{
		tweet, err := talkController.PostOne()
		if err != nil {
			break
		}

		enc.Encode(tweet)
		bufrw.Flush()
		time.Sleep(2*time.Second)
	}

	// End of response.
	bufrw.WriteString("EOR")
	bufrw.Flush()

	r.JSON(204, nil)
}