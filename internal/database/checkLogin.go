package database

import (
	"database/sql"
	"errors"
	"fmt"
	forum "forum/internal"

	"golang.org/x/crypto/bcrypt"
)

//CanLogin ..
func CanLogin(db *sql.DB, user forum.User) error {

	rows, err := db.Query("SELECT Login , Password  FROM SignInUser")
	if err != nil {
		err = fmt.Errorf("While selecting form SignInUser : %v", err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var login, password string
		if err := rows.Scan(&login, &password); err != nil {
			continue
		}
		errh := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
		if login == user.Login && errh == nil {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return errors.New("User with such data was not found")
}

//CanLoginBySocialNetworks ..
func CanLoginBySocialNetworks(db *sql.DB, user forum.User) error {
	rows, err := db.Query("SELECT Login, Confirm FROM SignInuser")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var login, confirm string
		if err := rows.Scan(&login, &confirm); err != nil {
			continue
		}
		if login == user.Login && confirm == user.Confirm {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return errors.New("Can not sign in")
}

/*

	CanLogin смотрим может ли пользователь зайти на сайт , сравниваем логин и хэш пароль
	с БД если все проверки пройдены добавляем его в сессию.

*/
