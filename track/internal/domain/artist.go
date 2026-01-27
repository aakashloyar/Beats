package domain

import "time"

type Artist struct {
	ID                 string
	Name               string 
	Bio                *string
	ProfileImageURL  *string
	CreatedAt          time.Time
}


