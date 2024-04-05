package logic

import (
	"context"
	"mywon/students_reports/constants"
	"mywon/students_reports/graph/model"
	"mywon/students_reports/jwt"
	"mywon/students_reports/repository"
	"mywon/students_reports/utility"
	"mywon/students_reports/validation"
	"os"
	"strings"
)

func CreateStudentDetails(ctx context.Context, input model.CreateStudentsInput )(*model.CreateStudentsResponse, error){
	PgSchema := os.Getenv(constants.POSTGRES_SCHEMA)
	response := &model.CreateStudentsResponse{}
	if repository.Pool == nil{
		repository.Pool = repository.GetPool()
	}
	conn := repository.SQLConnDetails{
		PgSchema: PgSchema,
		Pool:     repository.Pool,
	}

	err := validation.UpsertStudentValidation(ctx, input)
	if err != nil{
		return nil, err
	}
	res, err := conn.CreateStudentDetails(ctx, input)
	if err != nil{
		return nil, err
	}
	response = res
	return response, nil
}

func GetStudentDetails(ctx context.Context, input model.GetStudentDetailsInput )(*model.CreateStudentsResponse, error){
	PgSchema := os.Getenv(constants.POSTGRES_SCHEMA)
	response := &model.CreateStudentsResponse{}
	if repository.Pool == nil{
		repository.Pool = repository.GetPool()
	}
	conn := repository.SQLConnDetails{
		PgSchema: PgSchema,
		Pool:     repository.Pool,
	}

	res, err := conn.GetStudentDetails(ctx, input)
	if err != nil{
		return nil, err
	}
	response = res
	return response, nil
}

func CreateUser(ctx context.Context,input model.CreateUserInput)(*model.CreateUserResponse, error){
	input.Password = utility.Base64Decoder(input.Password)

	response, err := validation.CreateUserValidation(input)
	if err != nil || (!response.Success && response.Message != ""){
		return response, err
	}

	if !response.Success && response.Message != "" {
		return response, nil
	}

	input.Email = strings.Join(strings.Fields(strings.TrimSpace(input.Email)), "")
	input.UserName = strings.Join(strings.Fields(strings.TrimSpace(input.UserName)), "")
	input.UserName = strings.TrimSpace(input.UserName)
	input.Password = strings.TrimSpace(input.Password)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)

	PgSchema := os.Getenv(constants.POSTGRES_SCHEMA)
	if repository.Pool == nil{
		repository.Pool = repository.GetPool()
	}
	conn := repository.SQLConnDetails{
		PgSchema: PgSchema,
		Pool:     repository.Pool,
	}

	if conn.IsAlreadyExistUniqueField(ctx, input.Email, input.UserName) {
		response.Success = false
		response.Message = "This email is already registered. If you are the owner of this email, please click Forget Password to set your password."
		return response, nil
	}

	hashedPassword, err := utility.PasswordHashAndSalt(input.Password)
	if err != nil {
		return nil, err
	}

	message, userId, err := conn.CreateUser(ctx, input, hashedPassword)
	if err != nil {
		return nil, err
	}
	signUpDetails, err := conn.GetSignuPDetails(ctx, userId)

	if err != nil {
		return nil, err
	}

	generateJwtDetails := model.GenerateJWTDetails{
		Id:       int(signUpDetails.Id.Int64),
		Username: signUpDetails.UserName.String,
		Role:     "",
		IsAdmin:  false,
		Rights:   nil,
	}

	token, err := jwt.GenerateToken(generateJwtDetails)
	if err != nil {
		return nil, err
	}
	response.Success = true
	response.Message = message
	response.ID = int(signUpDetails.Id.Int64)
	response.Email = signUpDetails.Email.String
	response.UserName = signUpDetails.UserName.String
	response.FirstName = signUpDetails.FirstName.String
	response.MiddleName = &signUpDetails.MiddleName.String
	response.LastName = signUpDetails.LastName.String
	response.Gender = signUpDetails.Gender.String
	response.DateOfBirth = signUpDetails.DateOfBirth.String
	response.UserType = signUpDetails.UserType.String
	response.JwtToken = token
	return response, nil

}
