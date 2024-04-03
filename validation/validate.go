package validation

import (
	"context"
	"errors"
	"github.com/agrison/go-commons-lang/stringUtils"
	"mywon/students_reports/constants"
	"mywon/students_reports/graph/model"
	"mywon/students_reports/repository"
	"os"
	"regexp"
	"strings"
	"time"
)

func UpsertStudentValidation(ctx context.Context, input model.CreateStudentsInput) error {
	if input.ID == nil {
		if input.StudentName != nil {
			if stringUtils.IsBlank(*input.StudentName) {
				return errors.New("student name can't be empty")
			}
		} else {
			return errors.New("student name can't be nil")
		}
		if input.StudentMobileNumber != nil {
			if stringUtils.IsBlank(*input.StudentMobileNumber) {
				return errors.New("student mobile number can't be empty")
			} else {
				isExist := IsMobileNumberExist(*input.StudentMobileNumber)
				if isExist {
					return errors.New("student mobile number should be unique")
				}
			}
		} else {
			return errors.New("student mobile number can't be nil")
		}
		if input.StudentBloodGroup != nil {
			if stringUtils.IsBlank(*input.StudentBloodGroup) {
				return errors.New("student blood group can't be empty")
			} else {
				err := BloodGroupValidations(*input.StudentBloodGroup)
				if err != nil {
					return errors.New("please provide a valid blood group")
				}
			}
		} else {
			return errors.New("student blood group can't be nil")
		}
		if input.StudentAddress != nil {
			if stringUtils.IsBlank(*input.StudentAddress) {
				return errors.New("student address is required")
			}
		} else {
			return errors.New("student address can't be nil")
		}
		if input.DateOfBirth != nil {
			if stringUtils.IsBlank(*input.DateOfBirth) {
				return errors.New("date of birth is required")
			} else {
				re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
				if !re.MatchString(*input.DateOfBirth) {
					return errors.New("please enter valid DateOfBirth like DD/MM/YYYY")
				}
			}
		} else {
			return errors.New("student DateOfBirth can't be nil")
		}
		if input.StudentClass != nil {
			if *input.StudentClass == 0 {
				return errors.New("student class is required")
			}
		} else {
			return errors.New("student class can't be nil")
		}

		if input.StudentRoll != nil {
			if *input.StudentRoll == 0 {
				return errors.New("student roll is required")
			} else {
				isExist := IsRollNumberExist(*input.StudentRoll)
				if isExist {
					return errors.New("student roll number should be unique")
				}
			}
		} else {
			return errors.New("student roll can't be nil")
		}
	} else {
		
		if input.StudentMobileNumber != nil {
			if stringUtils.IsBlank(*input.StudentMobileNumber) {
				return errors.New("student mobile number can't be empty")
			} else {
				isExist := IsMobileNumberExist(*input.StudentMobileNumber)
				if isExist {
					return errors.New("student mobile number should be unique")
				}
			}
		}

		if input.StudentRoll != nil {
			if *input.StudentRoll == 0 {
				return errors.New("student roll is required")
			} else {
				isExist := IsRollNumberExist(*input.StudentRoll)
				if isExist {
					return errors.New("student Roll number should be unique")
				}
			}
		}

		if input.StudentBloodGroup != nil {
			if stringUtils.IsBlank(*input.StudentBloodGroup) {
				return errors.New("student blood group can't be empty")
			} else {
				err := BloodGroupValidations(*input.StudentBloodGroup)
				if err != nil {
					return errors.New("please provide a valid blood group")
				}
			}
		}

		if input.DateOfBirth != nil {
			if stringUtils.IsBlank(*input.DateOfBirth) {
				return errors.New("date of birth is required")
			} else {
				re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
				if !re.MatchString(*input.DateOfBirth) {
					return errors.New("please enter valid DateOfBirth like DD/MM/YYYY")
				}
			}
		}
	}
	return nil
}

func BloodGroupValidations(bloodGroup string) error {
	blood := make(map[string]bool)
	blood["A+"] = true
	blood["A-"] = true
	blood["B+"] = true
	blood["B-"] = true
	blood["O+"] = true
	blood["O-"] = true
	blood["AB+"] = true
	blood["AB-"] = true
	if _, ok := blood[bloodGroup]; !ok {
		return errors.New("please provide valid blood group")
	}
	return nil
}

func IsMobileNumberExist(mobielNumber string) bool {
	var hasValue int
	var response bool
	var sqlQuery string
	PgSchema := os.Getenv(constants.POSTGRES_SCHEMA)
	if repository.Pool == nil {
		repository.Pool = repository.GetPool()
	}
	conn := repository.SQLConnDetails{
		PgSchema: PgSchema,
		Pool:     repository.Pool,
	}
	sqlQuery = `SELECT 1 FROM ` + conn.PgSchema + `.student s WHERE s.student_mobile_no = $1`
	err := conn.Pool.QueryRowContext(context.Background(), sqlQuery, mobielNumber).Scan(&hasValue)
	if err != nil {
		return false
	}
	if hasValue == 1 {
		response = true
	}
	return response
}

func IsRollNumberExist(RollNumber int) bool {
	var hasValue int
	var response bool
	var sqlQuery string
	PgSchema := os.Getenv(constants.POSTGRES_SCHEMA)
	if repository.Pool == nil {
		repository.Pool = repository.GetPool()
	}
	conn := repository.SQLConnDetails{
		PgSchema: PgSchema,
		Pool:     repository.Pool,
	}
	sqlQuery = `SELECT 1 FROM ` + conn.PgSchema + `.student s WHERE s.student_roll_no = $1`
	err := conn.Pool.QueryRowContext(context.Background(), sqlQuery, RollNumber).Scan(&hasValue)
	if err != nil {
		return false
	}
	if hasValue == 1 {
		response = true
	}
	return response
}


func CreateUserValidation(input model.CreateUserInput) (*model.CreateUserResponse, error) {
	response := model.CreateUserResponse{}

	if strings.TrimSpace(input.Email) == "" {
		response.Success = false
		response.Message = "Email is empty!"
		return &response, nil
	}
	if strings.TrimSpace(input.UserName) == "" {
		response.Success = false
		response.Message = "User name is empty!"
		return &response, nil
	}
	if strings.TrimSpace(input.FirstName) == "" {
		response.Success = false
		response.Message = "First name is empty!"
		return &response, nil
	}
	if strings.TrimSpace(input.LastName) == "" {
		response.Success = false
		response.Message = "Last name is empty!"
		return &response, nil
	}
	if strings.TrimSpace(input.Password) == "" {
		response.Success = false
		response.Message = "Password is empty!"
		return &response, nil
	}
	if input.Gender != nil && strings.TrimSpace(*input.Gender) == "" {
		response.Success = false
		response.Message = "Gender is empty!"
		return &response, nil
	}

	if input.DateOfBirth != nil && strings.TrimSpace(*input.DateOfBirth) == "" {
		response.Success = false
		response.Message = "Date of birth is empty!"
		return &response, nil
	}
	if input.MobileNumber != nil && strings.TrimSpace(*input.MobileNumber) == "" {
		response.Success = false
		response.Message = "Mobile number is empty!"
		return &response, nil
	}
	if input.MobileCountryCode != nil && strings.TrimSpace(*input.MobileCountryCode) == "" {
		response.Success = false
		response.Message = "Mobile country code is empty!"
		return &response, nil
	}

	if input.UserType == nil {
		response.Success = false
		response.Message = "User type cannot be null!"
		return &response, nil
	}
	if input.UserType != nil && strings.TrimSpace(*input.UserType) == "" {
		response.Success = false
		response.Message = "User type is empty!"
		return &response, nil
	}

	if input.Gender != nil && (strings.ToLower(*input.Gender) != "male" && strings.ToLower(*input.Gender) != "female") {
		response.Success = false
		response.Message = "Gender is not valid!"
		return &response, nil
	}
		if input.DateOfBirth != nil {
			layout := "2006-01-02T15:04:05Z"
			t, err := time.Parse(layout, *input.DateOfBirth)
			if err != nil {
				response.Success = false
				response.Message = "Date of birth format is not correct!"
				return &response, nil
			}
			if t.After(time.Now()) {
				response.Success = false
				response.Message = "Date of birth can't be future date!"
				return &response, nil
			}
		}

	return &model.CreateUserResponse{}, nil
}
