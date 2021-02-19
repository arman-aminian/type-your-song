package model

type Music struct {
	Name   string `json:"name"`
	Cover  string `json:"cover"`
	Artist Artist `json:"artist"`
}
