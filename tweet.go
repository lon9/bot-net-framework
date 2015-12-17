package main
import (
"time"
"strconv"
	"github.com/jinzhu/gorm"
	"github.com/go-martini/martini"
"github.com/martini-contrib/render"
	"net/http"
)

type Tweet struct  {
	ID        int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	BotID     int `json:"botId" sql:"index"`
	Bot       Bot `json:"bot"`
	TalkID   int `json:"talkId"`
	Sequence  int `json:"sequence"`
	Bots      []Bot `json:"targets" gorm:"many2many:belonging_bots;"`
	Text      string `json:"text" sql:"type:text"`
}

type Tweets []Tweet

func IndexTweet(r render.Render, req *http.Request, db gorm.DB){

	talkId := req.FormValue("talkId")
	if talkId == "" {
		r.JSON(404, Error{404, "Talk not found"})
	}

	var tweets Tweets
	db.Find(&tweets, "talk_id = ?", talkId)
	for i, v := range tweets{
		db.Model(&v).Related(&tweets[i].Bot);
	}
	r.JSON(200, tweets)
}

func GetTweet(r render.Render, params martini.Params, db gorm.DB){
	id := params["id"]
	var Tweet Tweet
	db.First(&Tweet, id)
	if Tweet.ID == 0{
		r.JSON(404, Error{404, "Tweet was not found."})
	}
	r.JSON(200, Tweet)
}

func CreateTweet(r render.Render, db gorm.DB, tweet Tweet){
	db.Create(&tweet)
	db.Find(&tweet.Bot, tweet.BotID)
	r.JSON(201, tweet)
}

func UpdateTweet(r render.Render, db gorm.DB, tweet Tweet){
	tweet.UpdatedAt = time.Now()
	db.Save(&tweet)
	db.Find(&tweet.Bot, tweet.BotID)
	r.JSON(200, tweet)
}

func DeleteTweet(r render.Render, params martini.Params, db gorm.DB){
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	db.Delete(&Tweet{ID:id})
	r.JSON(204, nil)
}
