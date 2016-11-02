package ubm

import (
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

type (
	metrics struct{}

	UserMetric struct {
		UserID  string         `json:"userID"`
		Metrics map[string]int `json:"metrics"`
	}

	CompareStatement struct {
		Metric   string `json:"metric"`
		Operator string `json:"operator"`
		Value    int    `json:"value"`
	}
)

var (
	Metrics metrics
)

func (_ metrics) Push(userID string, key string, value int) (err error) {
	_, err = Models.Upsert(
		bson.M{"id": userID},
		bson.M{"$inc": bson.M{key: value}},
	)
	refresh("ubm", err)
	return
}

func (_ metrics) Get(userID string, keys []string) (answer map[string]int, err error) {
	var result map[string]int
	answer = make(map[string]int)

	project := bson.M{}
	for i, key := range keys {
		project[strconv.Itoa(i)] = "$" + key // ugly hack
	}

	err = Models.Pipe([]bson.M{
		bson.M{"$match": bson.M{"id": userID}},
		bson.M{"$project": project},
	}).One(&result)

	for i, key := range keys {
		if value, ok := result[strconv.Itoa(i)]; ok {
			answer[key] = value
		}
	}

	refresh("umb", err)
	return
}
