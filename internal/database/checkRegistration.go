package database

import (
	"database/sql"
	"errors"
	forum "forum/internal"
	"strings"
)

//CanRegister ..
func CanRegister(db *sql.DB, user forum.User) error {
	if !checkForSpace(user.Password) || !checkForSpace(user.Name) {
		return errors.New("Please, do not use space")
	}
	if user.Password != user.Confirm {
		return errors.New("Password confirmation must match Password")
	}
	if err := ValidatingName(user.Name); err != nil {
		return err
	}

	if err := ValidatingPassword(user.Password); err != nil {
		return err
	}

	if err := ValidatingLogin(user.Login); err != nil {
		return err
	}

	return InsertNewUserIntoDB(db, user)
}

func checkForSpace(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if v == 32 {
			return false
		}
	}
	return true
}

//ValidatingName ..
func ValidatingName(name string) error {
	if len(name) > 30 {
		return errors.New("Name must contain less than 30 characters")
	}
	return nil
}

//ValidatingPassword ..
func ValidatingPassword(password string) error {
	var checked bool
	for _, l := range password {
		if l < 32 || l > 126 {
			checked = true
		}
	}
	if checked {
		return errors.New("Please, use only latin alphabet")
	}
	if len(password) < 8 {
		return errors.New("Password must contain more than 8 characters")
	}
	var numbers, bigletters, smallletters, specialsymbols bool
	for _, v := range password {
		if v >= 'a' && v <= 'z' {
			smallletters = true
		} else if v >= 'A' && v <= 'Z' {
			bigletters = true
		} else if v >= '0' && v <= '9' {
			numbers = true
		} else {
			specialsymbols = true
		}
	}
	if !numbers {
		return errors.New("Password must contains numbers")
	}
	if !bigletters {
		return errors.New("Password must contains uppercase")
	}
	if !smallletters {
		return errors.New("Password must contains lowercase")
	}
	if !specialsymbols {
		return errors.New("Password must contains symbols")
	}
	return nil
}

//ValidatingLogin ..
func ValidatingLogin(login string) error {
	var checked bool
	if len(login) > 30 {
		return errors.New("Email must contain less than 30 characters")
	}
	splitedLogin := strings.Split(login, "@")
	runes := []rune(splitedLogin[0])
	for _, l := range login {
		if l < 32 || l > 126 {
			checked = true
		}
	}
	if checked {
		return errors.New("Please, use only latin alphabet")
	}
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] <= ' ' || runes[i] >= '~' {
			checked = true
		}
		if runes[i] >= '!' && runes[i] <= '/' {
			if runes[i] == '-' {
				continue
			} else if (runes[i] == '.') && (runes[i] != runes[0] && runes[i] != runes[len(runes)-1]) {
				if runes[i+1] == '.' {
					return errors.New("Please, do not reapeat '.'")
				}
				continue
			} else {
				checked = true
			}
		}
		if runes[i] >= 'A' && runes[i] <= 'Z' {
			return errors.New("Please, use lowercase in the email")
		}
		if runes[i] >= '0' && runes[i] <= '9' {
			continue
		}
	}
	if checked {
		return errors.New("Please, use only '-', '_', '.' symbols in the name of email address. ( '.' must not be first or last)")
	}
	runes = []rune(splitedLogin[1])
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] <= ' ' || runes[i] >= '~' {
			checked = true
		}
		if runes[i] >= 'A' && runes[i] <= 'Z' {
			return errors.New("Please, use lowercase in the email")
		}
		if runes[i] >= '!' && runes[i] <= '/' {
			if (runes[i] == '.') && (runes[i] != runes[0] && runes[i] != runes[len(runes)-1]) {
				if runes[i+1] == '.' {
					return errors.New("Please, do not reapeat '.'")
				}
				continue
			} else {
				checked = true
			}
		}
	}
	if checked {
		return errors.New("Please, use only '.' symbol in the email domen. ( '.' must not be first or last)")
	}
	return nil
}

/*

	CanRegister проверяем на уникальность


*/
