package beater

import "fmt"

type TwitterConfig struct {
	Period  *int64
	Twitter struct {
		Names          *[]string
		ConsumerKey    *string
		ConsumerSecret *string
		AccessKey      *string
		AccessSecret   *string
	}
}

type TwitterConfigMissing byte

func (e TwitterConfigMissing) Error() string {
	return fmt.Sprintf("Twitter config is missing")
}
