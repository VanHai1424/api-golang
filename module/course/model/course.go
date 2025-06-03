package coursemodel

import "time"

type Course struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Desc      string `json:"desc"`
	Version   string `json:"version"`
	Update    string `json:"update"`
	Developer string `json:"developer"`
	Category  string `json:"category"`
	PlayId    string `json:"playId" gorm:"column:playId"`
	Downloads string `json:"downloads"`
	Link      string `json:"link"`
	TitleCate string `json:"titleCate" gorm:"column:titleCate"`
	Type      string `json:"type" gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
