package bmodel

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	logs struct{}
)

var (
	limit = 128
	Logs  logs
)

func (_ logs) Push(userID string, key string, value interface{}) (err error) {
	_, err = Models.Upsert(
		bson.M{"id": userID},
		bson.M{
			"$push": bson.M{
				"logs." + key: bson.M{
					"$each":  []interface{}{value},
					"$slice": -limit,
				},
			},
		})
	refresh("bmodel", err)
	return
}
