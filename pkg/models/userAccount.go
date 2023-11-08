package models

import (
	"encoding/xml"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAccount struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	XmlId          int                `bson:"xml_id" json:"xml_id" xml:"Id,attr"`
	Reputation     int                `bson:"reputation" json:"reputation" xml:"Reputation,attr"`
	CreationDate   string             `bson:"creation_date" json:"creation_date" xml:"CreationDate,attr"`
	DisplayName    string             `bson:"display_name" json:"display_name" xml:"DisplayName,attr"`
	LastAccessDate string             `bson:"last_access_date" json:"last_access_date" xml:"LastAccessDate,attr"`
	WebsiteUrl     string             `bson:"website_url" json:"website_url" xml:"WebsiteUrl,attr"`
	Location       string             `bson:"location" json:"location" xml:"Location,attr"`
	AboutMe        string             `bson:"about_me" json:"about_me" xml:"AboutMe,attr"`
	Views          int                `bson:"views" json:"views" xml:"Views,attr"`
	UpVotes        int                `bson:"up_votes" json:"up_votes" xml:"UpVotes,attr"`
	DownVotes      int                `bson:"down_votes" json:"down_votes" xml:"DownVotes,attr"`
	AccountId      int                `bson:"account_id" json:"account_id" xml:"AccountId,attr"`
}

type UserCollection struct {
	XmlName xml.Name `xml:"UserCollection"`
	Users   Users    `xml:"users"`
}

type Users struct {
	UserAccounts []*UserAccount `xml:"row"`
}
