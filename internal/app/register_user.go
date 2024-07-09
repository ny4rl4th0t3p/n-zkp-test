package app

import (
	"context"
	"log/slog"

	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"
	"practical-case-test/internal/repository"
)

// RegisterUserExecuter is an interface that defines the method for executing user registration.
// Exec takes a context and a RegisterRequest and returns an error if any occurred during the execution.
type RegisterUserExecuter interface {
	Exec(ctx context.Context, req *interactor.RegisterRequest) error
}

// RegisterUser is a type that is responsible for registering a new user.
// It uses an AuthRepository to store user registration information.
type RegisterUser struct {
	ar repository.AuthRepository
}

// NewRegisterUser is a function that returns a RegisterUserExecuter.
// It takes in an AuthRepository and initializes a RegisterUser struct with the provided repository.
// The returned RegisterUserExecuter can be used to execute the registration operation.
func NewRegisterUser(ar repository.AuthRepository) RegisterUserExecuter {
	return &RegisterUser{ar: ar}
}

// Exec executes the register user use case.
// It takes in a context and a RegisterRequest object and returns an error.
// The function extracts user, y1, and y2 values from the request.
// It then logs the registration request.
// The function creates a new User object using auth.NewUser and the extracted values.
// If the user object is invalid, it returns an error.
// The function stores the user registration using the AuthRepository.
// If there is an error storing the registration, it returns the error.
// Finally, it returns nil if no errors occurred.
func (ru RegisterUser) Exec(ctx context.Context, req *interactor.RegisterRequest) error {
	user := req.GetUser()
	y1 := req.GetY1()
	y2 := req.GetY2()

	slog.Info("received registration request\n", "user", user, "y1", y1, "y2", y2)

	newUser, err := auth.NewUser(user, y1, y2)
	if err != nil {
		return err
	}

	err = ru.ar.StoreUserRegistration(ctx, *newUser)
	if err != nil {
		return err
	}

	return nil
}
