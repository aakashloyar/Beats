package domain
import "time"
type Artist struct {
	ID                 string
	Name               string 
	Bio                *string
	Profile_image_url  *string
	Create_At          time.Time

}

//address bcz these are big 