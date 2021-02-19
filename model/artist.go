package model

type Artist struct {
	Name   string  `json:"name"`
	Musics []Music `json:"musics"`
}
