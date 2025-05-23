package coursemodel

type Course struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Desc      string `json:"desc"`
	Version   string `json:"version"`
	Update    string `json:"update"`
	Developer string `json:"developer"`
	Category  string `json:"category"`
	PlayId    string `json:"playId"`
	Downloads string `json:"downloads"`
	Link      string `json:"link"`
	TitleCate string `json:"titleCate"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Type      string `json:"type"`
}
