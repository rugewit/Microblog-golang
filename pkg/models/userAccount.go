package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAccount struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Reputation     int                `bson:"reputation" json:"reputation"`
	CreationDate   string             `bson:"creation_date" json:"creation_date"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	LastAccessDate string             `bson:"last_access_date" json:"last_access_date"`
	WebsiteUrl     string             `bson:"website_url" json:"website_url"`
	Location       string             `bson:"location" json:"location"`
	AboutMe        string             `bson:"about_me" json:"about_me"`
	Views          int                `bson:"views" json:"views"`
	UpVotes        int                `bson:"up_votes" json:"up_votes"`
	DownVotes      int                `bson:"down_votes" json:"down_votes"`
	AccountId      int                `bson:"account_id" json:"account_id"`
}
