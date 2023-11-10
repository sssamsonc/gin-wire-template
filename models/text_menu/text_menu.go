package text_menu

type TextMenu struct {
	Type string `json:"type" bson:"type"`
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}
