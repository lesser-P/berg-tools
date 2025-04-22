package utils

import (
	"log"
	"testing"
)

func TestCopyStruct(t *testing.T) {
	po := UserPO{
		ID:       1,
		Name:     "张三",
		Password: "123456",
		Address: Address{
			City:    "北京",
			Country: "中国",
			Hobby: Hobby{
				Name: "打篮球",
				Note: "喜欢打篮球",
			},
		},
	}

	var vo UserVO

	err := CopyStruct(&vo, po)
	if err != nil {
		log.Print(err)
	}
}

type UserPO struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Address  Address `json:"address"`
}

type UserVO struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Address Address2 `json:"address"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Hobby   Hobby  `json:"hobby"`
}
type Address2 struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Hobby   Hobby  `json:"hobby"`
}

type Hobby struct {
	Name string `json:"name"`
	Note string `json:"note"`
}
