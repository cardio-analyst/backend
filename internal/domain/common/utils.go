package common

import "time"

func GetCurrentAge(birthDate time.Time) int {
	today := time.Now().In(birthDate.Location())

	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)

	by, bm, bd := birthDate.Date()
	birthDate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)

	if today.Before(birthDate) {
		return 0
	}

	age := ty - by

	anniversary := birthDate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}

	return age
}
