package database

import (
	"errors"
	"icas/pkg/api"
	"icas/pkg/password"
)

func UserLogin(userEmail string, userPassword string) (api.User, error) {
	db, err := Open()
	if err != nil {
		return api.User{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email, password, role FROM users WHERE email = ?", userEmail)
	if err != nil {
		return api.User{}, err
	}
	defer rows.Close()

	var u api.User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role); err != nil {
			return api.User{}, err
		}
		if password.VerifyPassword(userPassword, u.Password) {
			return u, nil
		}
	}
	if err = rows.Err(); err != nil {
		return api.User{}, err
	}
	return api.User{}, errors.New("user not found")
}
