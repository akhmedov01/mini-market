package packages

import "time"

func FindAge(bDate string) int {

	date, error := time.Parse("2006-01-02", bDate)

	if error != nil {
		return 0
	}

	result := time.Now().Year() - date.Year()

	return result

}
