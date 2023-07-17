package models

type SiteDetail struct {
	BaseModel
	Name 		string 		`json:"name" default:"Kay's Auction House"`
	Email 		string 		`json:"email" default:"kayprogrammer1@gmail.com"`
	Phone 		string		`json:"phone" default:"+2348133831036"`
	Address 	string		`json:"address" default:"234, Lagos, Nigeria"`
	Fb 			string		`json:"fb" default:"https://facebook.com"`
	Tw 			string		`json:"tw" default:"https://twitter.com"`
	Wh 			string		`json:"wh" default:"https://wa.me/2348133831036"`
	Ig 			string		`json:"ig" default:"https://instagram.com"`
}