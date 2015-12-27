package main
import (
"time"
"strconv"
	"github.com/jinzhu/gorm"
	"github.com/go-martini/martini"
"github.com/martini-contrib/render"
	"net/http"
)

// Tweet is object for Tweet.
type Tweet struct  {
	ID        int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	BotID     int `json:"botId" sql:"index"`
	Bot       Bot `json:"bot"`
	TalkID   int `json:"talkId" sql:"index"`
	Sequence  int `json:"sequence"`
	Text      string `json:"text" sql:"type:text"`
	TweetId string `json:"tweetId"`
}

// Tweets is array of Tweet.
type Tweets []Tweet

// IndexTweet returns array of Tweet.
func IndexTweet(r render.Render, req *http.Request, db gorm.DB){

	talkId := req.FormValue("talkId")
	if talkId == "" {
		r.JSON(404, Error{404, "Talk not found"})
		return
	}

	var tweets Tweets
	db.Find(&tweets, "talk_id = ?", talkId)
	for i, v := range tweets{
		db.Model(&v).Related(&tweets[i].Bot);
	}
	r.JSON(200, tweets)
}

// GetTweet returns a Tweet.
func GetTweet(r render.Render, params martini.Params, db gorm.DB){
	id := params["id"]
	var Tweet Tweet
	db.First(&Tweet, id)
	if Tweet.ID == 0{
		r.JSON(404, Error{404, "Tweet was not found."})
		return
	}
	r.JSON(200, Tweet)
}

// CreateTweet inserts a tweet.
func CreateTweet(r render.Render, db gorm.DB, tweet Tweet){
	db.Create(&tweet)
	db.Find(&tweet.Bot, tweet.BotID)
	r.JSON(201, tweet)
}

// UpdateTweet updates a Tweet.
func UpdateTweet(r render.Render, db gorm.DB, tweet Tweet){
	tweet.UpdatedAt = time.Now()
	db.Save(&tweet)
	db.Find(&tweet.Bot, tweet.BotID)
	r.JSON(200, tweet)
}

// DeleteTweet deletes a Tweet.
func DeleteTweet(r render.Render, params martini.Params, db gorm.DB){
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	db.Delete(&Tweet{ID:id})
	r.JSON(204, nil)
}
