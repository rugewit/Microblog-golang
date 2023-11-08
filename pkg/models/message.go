package models

import (
	"encoding/xml"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	XmlId            int                `bson:"xml_id" json:"xml_id" xml:"Id,attr"`
	PostTypeId       int                `bson:"post_type_id" json:"post_type_id" xml:"PostTypeId,attr"`
	AcceptedAnswerId int                `bson:"accepted_answer_id" json:"accepted_answer_id" xml:"AcceptedAnswerId,attr"`
	CreationDate     string             `bson:"creation_date" json:"creation_date" xml:"CreationDate,attr"`
	Score            int                `bson:"score" json:"score" xml:"Score,attr"`
	ViewCount        int                `bson:"view_count" json:"view_count" xml:"ViewCount,attr"`
	Body             string             `bson:"body" json:"body" xml:"Body,attr"`
	OwnerUserId      int                `bson:"owner_user_id" json:"owner_user_id" xml:"OwnerUserId,attr"`
	LastEditorUserId int                `bson:"last_editor_user_id" json:"last_editor_user_id" xml:"LastEditorUserId,attr"`
	LastEditDate     string             `bson:"last_edit_date" json:"last_edit_date" xml:"LastEditDate,attr"`
	LastActivityDate string             `bson:"last_activity_date" json:"last_activity_date" xml:"LastActivityDate,attr"`
	Title            string             `bson:"title" json:"title" xml:"Title,attr"`
	Tags             string             `bson:"tags" json:"tags" xml:"Tags,attr"`
	AnswerCount      int                `bson:"answer_count" json:"answer_count" xml:"AnswerCount,attr"`
	CommentCount     int                `bson:"comment_count" json:"comment_count" xml:"CommentCount,attr"`
	ContentLicense   string             `bson:"content_license" json:"content_license" xml:"ContentLicense,attr"`
}

type MessageCollection struct {
	XmlName xml.Name  `xml:"MessageCollection"`
	Posts   PostArray `xml:"posts"`
}

type PostArray struct {
	Messages []*Message `xml:"row"`
}
