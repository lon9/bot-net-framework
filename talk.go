package main
import (
	"time"
	"github.com/martini-contrib/render"
	"github.com/jinzhu/gorm"
	"strconv"
	"net/http"
	"github.com/go-martini/martini"
)

type Talk struct  {
	ID int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title string `json:"title" sql:"unique"`
	Tweets []Tweet `json:"tweets"`
}

type Talks []Talk

func IndexTalk(r render.Render, req *http.Request, db gorm.DB){
	limit := 40
	offset := 0
	rawPage, rawMaxResults := req.FormValue("page"), req.FormValue("maxResults")
	page, err := strconv.Atoi(rawPage)
	maxResults, err := strconv.Atoi(rawMaxResults)
	if err != nil {
		r.JSON(400, Error{400, "page and maxResults must be integer."})
		return
	}

	limit = maxResults
	offset = (page - 1) * maxResults

	var talks Talks
	db.Order("id desc").Limit(limit).Offset(offset).Find(&talks)
	r.JSON(200, talks)
}

func GetTalk(r render.Render, params martini.Params, db gorm.DB){
	id := params["id"]
	var talk Talk
	db.First(&talk, id)
	if talk.ID == 0{
		r.JSON(404, Error{404, "Talk was not found."})
		return
	}
	r.JSON(200, talk)
}

func CreateTalk(r render.Render, db gorm.DB, talk Talk){
	db.Create(&talk)
	r.JSON(201, talk)
}

func UpdateTalk(r render.Render, db gorm.DB, talk Talk){
	talk.UpdatedAt = time.Now()
	db.Save(&talk)
	r.JSON(200, talk)
}

func DeleteTalk(r render.Render, params martini.Params, db gorm.DB){
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	db.Delete(&Talk{ID:id})
	r.JSON(204, nil)
}