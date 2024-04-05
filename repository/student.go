package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"mywon/students_reports/graph/model"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/valyala/fastjson"
)

func (conn SQLConnDetails) CreateStudentDetails(ctx context.Context, input model.CreateStudentsInput) (*model.CreateStudentsResponse, error) {
	//cretaed by Dibyendu()
	var (
		id             int
		studentName    string
		studentClass   int
		studentRoll    int
		studentAddress string
		bloodGroup     string
		mobileNumber   string
		score          int
		dateOfBirth    string
		subject1       string
		createdAt      string
		updatedAt      string
	)
	response := model.CreateStudentsResponse{}
	tx, err := conn.Pool.Begin()
	if err != nil {
		return nil, err
	}
	var sqlErr error
	var sqlQuery string
	var inputArgs []interface{}
	subject, err := json.Marshal(input.Subject)
	if err != nil {
		return nil, err
	}

	parsedSubject, err := fastjson.ParseBytes(subject)
	if err != nil {
		return nil, err
	}
	var totalScore float64 = 0
	totalScore += parsedSubject.GetFloat64("Bengali")
	totalScore += parsedSubject.GetFloat64("English")
	totalScore += parsedSubject.GetFloat64("Mathematics")
	totalScore += parsedSubject.GetFloat64("Physics")
	totalScore += parsedSubject.GetFloat64("Biology")
	totalScore += parsedSubject.GetFloat64("Chemistry")
	average := totalScore / 6
	roundedAverage := int(math.Round(average))
	if input.ID == nil {
		sqlQuery = `INSERT INTO ` + conn.PgSchema + `.student (student_name, student_class, student_roll_no, student_address, student_blood_group, student_mobile_no, score, date_of_birth, subjects, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) returning id, student_name, student_class, student_roll_no, student_address, student_blood_group, student_mobile_no, score, date_of_birth, subjects, created_at, updated_at;`
		inputArgs = append(inputArgs, strings.TrimSpace(*input.StudentName), *input.StudentClass, *input.StudentRoll, strings.TrimSpace(*input.StudentAddress), strings.TrimSpace(*input.StudentBloodGroup), strings.TrimSpace(*input.StudentMobileNumber), roundedAverage, strings.TrimSpace(*input.DateOfBirth), string(subject))
		sqlErr = tx.QueryRowContext(ctx, sqlQuery, inputArgs...).Scan(&id, &studentName, &studentClass, &studentRoll, &studentAddress, &bloodGroup, &mobileNumber, &score, &dateOfBirth, &subject1, &createdAt, &updatedAt)
		response.ID = id
		response.StudentName = studentName
		response.StudentClass = studentClass
		response.StudentRoll = studentRoll
		response.StudentAddress = studentAddress
		response.StudentBloodGroup = bloodGroup
		response.StudentMobileNumber = mobileNumber
		response.Score = score
		response.DateOfBirth = dateOfBirth
		err := json.Unmarshal([]byte(subject1), &response.Subject)
		if err != nil {
			return nil, err
		}
		response.CreatedAt = createdAt
		response.UpdatedAt = updatedAt
		if sqlErr != nil {
			tx.Rollback()
			return nil, sqlErr
		}
		txErr := tx.Commit()
		if txErr != nil {
			tx.Rollback()
			return nil, txErr
		} else {
			return &response, nil
		}
	} else {
		sqlQuery = `UPDATE ` + conn.PgSchema + `.student SET updated_at = NOW()`
		if input.StudentName != nil {
			sqlQuery += `, student_name = ?`
			inputArgs = append(inputArgs, strings.TrimSpace(*input.StudentName))
		}
		if input.StudentClass != nil && *input.StudentClass != 0 {
			sqlQuery += `, student_class = ?`
			inputArgs = append(inputArgs, *input.StudentClass)
		}
		if input.StudentRoll != nil && *input.StudentRoll != 0 {
			sqlQuery += `, student_roll_no = ?`
			inputArgs = append(inputArgs, *input.StudentRoll)
		}
		if input.StudentAddress != nil {
			sqlQuery += `, student_address = ?`
			inputArgs = append(inputArgs, strings.TrimSpace(*input.StudentAddress))
		}
		if input.StudentBloodGroup != nil {
			sqlQuery += `, student_blood_group = ?`
			inputArgs = append(inputArgs, strings.TrimSpace(*input.StudentBloodGroup))
		}
		if input.StudentMobileNumber != nil {
			sqlQuery += `, student_mobile_no = ?`
			inputArgs = append(inputArgs, strings.TrimSpace(*input.StudentMobileNumber))
		}

		if input.DateOfBirth != nil {
			sqlQuery += `, date_of_birth = ?`
			inputArgs = append(inputArgs, strings.TrimSpace(*input.DateOfBirth))
		}

		if input.Subject != nil {
			subject, err := json.Marshal(input.Subject)
			if err != nil {
				return nil, err
			}
			sqlQuery += `, subjects = ?`
			inputArgs = append(inputArgs, subject)
			parsedSubject, err := fastjson.ParseBytes(subject)
			if err != nil {
				return nil, err
			}
			var totalScore float64 = 0
			totalScore += parsedSubject.GetFloat64("Bengali")
			totalScore += parsedSubject.GetFloat64("English")
			totalScore += parsedSubject.GetFloat64("Mathematics")
			totalScore += parsedSubject.GetFloat64("Physics")
			totalScore += parsedSubject.GetFloat64("Biology")
			totalScore += parsedSubject.GetFloat64("Chemistry")
			average := totalScore / 6
			avg := int(math.Round(average))

			// Assign avg to score field
			sqlQuery += `, score = ?`
			inputArgs = append(inputArgs, avg)
		}
		sqlQuery += ` WHERE id = ? RETURNING id,student_name, student_class, student_roll_no, student_address,student_blood_group, student_mobile_no, score, date_of_birth, subjects, created_at, updated_at`
		inputArgs = append(inputArgs, *input.ID)
		sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
		log.Println("SQL Query@@", sqlQuery, "inputArgs$$$", inputArgs)
		sqlErr = tx.QueryRowContext(ctx, sqlQuery, inputArgs...).Scan(&id, &studentName, &studentClass, &studentRoll, &studentAddress, &bloodGroup, &mobileNumber, &score, &dateOfBirth, &subject1, &createdAt, &updatedAt)
		response.ID = id
		response.StudentName = studentName
		response.StudentClass = studentClass
		response.StudentRoll = studentRoll
		response.StudentAddress = studentAddress
		response.StudentBloodGroup = bloodGroup
		response.StudentMobileNumber = mobileNumber
		response.Score = score
		response.DateOfBirth = dateOfBirth
		err := json.Unmarshal([]byte(subject1), &response.Subject)
		if err != nil {
			return nil, err
		}
		response.CreatedAt = createdAt
		response.UpdatedAt = updatedAt
		if sqlErr != nil {
			tx.Rollback()
			return nil, sqlErr
		}

		txErr := tx.Commit()
		if txErr != nil {
			tx.Rollback()

			return nil, txErr
		} else {
			return &response, nil
		}
	}
}

func (conn SQLConnDetails) GetStudentDetails(ctx context.Context, input model.GetStudentDetailsInput) (*model.CreateStudentsResponse, error) {
	response := model.CreateStudentsResponse{}
	var sqlQuery string
	var inputArgs []interface{}

	sqlQuery = `select id, student_name, student_class, student_roll_no, student_address, student_blood_group, student_mobile_no, score, date_of_birth, subjects, created_at, updated_at from ` + conn.PgSchema + `.student where id = ?`
	inputArgs = append(inputArgs, input.ID)
	sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	rows, sqlErr := conn.Pool.QueryContext(ctx, sqlQuery, inputArgs...)
	if sqlErr != nil {
		return nil, sqlErr
	}
	defer rows.Close()
	for rows.Next() {
		var jsonSubject sql.NullString
		err := rows.Scan(
			&response.ID,
			&response.StudentName,
			&response.StudentClass,
			&response.StudentRoll,
			&response.StudentAddress,
			&response.StudentBloodGroup,
			&response.StudentMobileNumber,
			&response.Score,
			&response.DateOfBirth,
			&jsonSubject,
			&response.CreatedAt,
			&response.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		//if err := json.Unmarshal([]byte(jsonSubject), &response.Subject); err != nil {
		//	return nil, err
		//}
		if jsonSubject.Valid && jsonSubject.String != "" {
			parsedSubject := fastjson.MustParse(jsonSubject.String)
			response.Subject = &model.Subject{
				Bengali:     parsedSubject.GetFloat64("Bengali"),
				English:     parsedSubject.GetFloat64("English"),
				Mathematics: parsedSubject.GetFloat64("Mathematics"),
				Physics:     parsedSubject.GetFloat64("Physics"),
				Biology:     parsedSubject.GetFloat64("Biology"),
				Chemistry:   parsedSubject.GetFloat64("Chemistry"),
			}
		} else {
			response.Subject = &model.Subject{
				Bengali:     0,
				English:     0,
				Mathematics: 0,
				Physics:     0,
				Biology:     0,
				Chemistry:   0,
			}
		}
	}
	return &response, nil
}

func (conn SQLConnDetails) CreateUser(eventCtx context.Context, input model.CreateUserInput, Password string) (string, int, error) {
	ctx := context.Background()
	tx, err := conn.Pool.Begin()
	if err != nil {
		return "", 0, err
	}
	var query string
	var inputArgs []interface{}
	var message string
	var userID int
	query = `insert into ` + conn.PgSchema + `.users (user_name, email, pass_word, created_at, first_name, middle_name, last_name, 
	mobile_no, user_type, date_of_birth, mobile_country_code, gender, updated_at) values(?, ?, ?, now(), ?, ?, ?, ?, ?, ?, ?, ?, now())returning(id)`
	inputArgs = append(inputArgs, input.UserName, input.Email, Password, input.FirstName, input.MiddleName, input.LastName, input.MobileNumber, input.UserType, input.DateOfBirth, input.MobileCountryCode, input.Gender)

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	sqlErr := tx.QueryRowContext(ctx, query, inputArgs...).Scan(&userID)
	if sqlErr != nil {
		tx.Rollback()
		return "", 0, sqlErr
	} else {
		message = "User created successfully!"
	}

	txErr := tx.Commit()
	if txErr != nil {
		tx.Rollback()
		return "", 0, txErr
	}

	return message, userID, nil
}

func (conn SQLConnDetails) GetSignuPDetails(eventCtx context.Context, userID int) (model.SignupDetails, error) {
	signupDetails := model.SignupDetails{}
	ctx := context.Background()
	query := `select upu.id, user_name, email, first_name, middle_name, last_name, 
	user_type, date_of_birth, gender
	from ` + conn.PgSchema + `.users upu where upu.id = $1`

	rows, sqlErr := conn.Pool.QueryContext(ctx, query, userID)
	if sqlErr != nil {
		return signupDetails, sqlErr
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&signupDetails.Id,
			&signupDetails.UserName,
			&signupDetails.Email,
			&signupDetails.FirstName,
			&signupDetails.MiddleName,
			&signupDetails.LastName,
			&signupDetails.UserType,
			&signupDetails.DateOfBirth,
			&signupDetails.Gender,
		)
		if err != nil {
			return signupDetails, err
		}
	}
	return signupDetails, nil

}

func (conn SQLConnDetails) IsAlreadyExistUniqueField(eventCtx context.Context, email string, username string) bool {
	querystring := `SELECT 1 FROM ` + conn.PgSchema + `.users WHERE email = $1 or user_name = $2`
	var hasValue int
	err := conn.Pool.QueryRowContext(context.Background(), querystring, email, username).Scan(&hasValue)
	if err != nil {
		return false
	}
	return true
}
