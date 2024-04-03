package model

import "database/sql"

type GenerateJWTDetails struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	IsAdmin  bool      `json:"isAdmin"`
	Rights   []*string `json:"rights"`
}

type SignupDetails struct {
	Id          sql.NullInt64
	FirstName   sql.NullString
	MiddleName  sql.NullString
	LastName    sql.NullString
	UserName    sql.NullString
	Email       sql.NullString
	DateOfBirth sql.NullString
	Gender      sql.NullString
	UserType    sql.NullString
}
