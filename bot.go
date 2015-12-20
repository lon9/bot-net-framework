package main
import (
	"time"
	"github.com/jinzhu/gorm"
	"strconv"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"net/http"
)

// Bot is struct of Bots
type Bot struct  {
	ID int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name string `json:"name"`
	ScreenName string `json:"screenName" sql:"unique"`
	AccessToken string `json:"-"`
	AccessTokenSecret string `json:"-"`
	TwitterId int64 `json:"twitterId" sql:"unique"`
	IconURL string `json:"iconUrl"`
}

// Bots is array of Bot
type Bots []Bot

// IndexBot returns all bot registered.
func IndexBot(r render.Render, req *http.Request, db gorm.DB){
	var Bots Bots
	db.Find(&Bots)
	r.JSON(200, Bots)
}

// GetBot return a bot.
func GetBot(r render.Render, params martini.Params, db gorm.DB){
	id := params["id"]
	var bot Bot
	db.First(&bot, id)
	if bot.ID == 0{
		r.JSON(404, Error{404, "Bot was not found."})
	}
	r.JSON(200, bot)
}

// CreateBot inserts bot to database.
func CreateBot(r render.Render, db gorm.DB, bot Bot){
	db.Create(&bot)
	r.JSON(201, bot)
}

// UpdateBot updates a bot.
func UpdateBot(r render.Render, db gorm.DB, bot Bot){
	bot.UpdatedAt = time.Now()
	db.Save(&bot)
	r.JSON(200, bot)
}

// DeleteBot deletes a bot.
func DeleteBot(r render.Render, params martini.Params, db gorm.DB){
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	db.Delete(&Bot{ID:id})
	r.JSON(204, nil)
}