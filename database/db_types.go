package database

type Document struct {
	Date  string
	Time  string
	Price string
}

type MongoStamp struct {
	Time  string `bson:"time" json:"time"`
	Price string `bson:"price" json:"price"`
}

type MongoDocument struct { //the one that decodes raw bson
	Date   string       `bson:"date" json:"date"`
	Stamps []MongoStamp `bson:"stamps" json:"stamps"`
}
