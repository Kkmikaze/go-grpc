package common

import (
	"errors"
	"fmt"
	playground "github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func TransformValidatorResponse(err error) error {
	s := status.New(codes.InvalidArgument, "Invalid Argument")
	objectMapper := strings.Split(err.Error(), ";")
	request := &errdetails.BadRequest{}
	for _, v := range objectMapper {
		clearCause := strings.Split(v, "|")[0]
		mapper := strings.Split(clearCause, ":")

		key := strings.Split(mapper[0], ".")[1]
		value := mapper[1]
		mapData := &errdetails.BadRequest_FieldViolation{
			Field:       key,
			Description: strings.TrimSpace(value),
		}
		request.FieldViolations = append(request.FieldViolations, mapData)
	}
	details, err := s.WithDetails(request)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return details.Err()
}

func TransformValidator(err error) error {
	if _, ok := err.(*playground.InvalidValidationError); ok {
		return status.Error(codes.Internal, "Invalid validation error")
	}

	var errMsg string
	for _, err2 := range err.(playground.ValidationErrors) {
		field := strings.ToLower(err2.Field())
		namespace := err2.Namespace()
		tag := err2.Tag()
		switch tag {
		case "required":
			errMsg += fmt.Sprintf("%s: Field '%s' is required;", namespace, field)
		case "min":
			errMsg += fmt.Sprintf("%s: Field '%s' is too short;", namespace, field)
		case "max":
			errMsg += fmt.Sprintf("%s: Field '%s' is too long;", namespace, field)
		default:
			errMsg += fmt.Sprintf("%s: Field '%s' validation failed;", namespace, field)
		}
	}

	err = TransformValidatorResponse(errors.New(strings.TrimRight(errMsg, ";")))
	if err != nil {
		return err
	}

	return nil
}