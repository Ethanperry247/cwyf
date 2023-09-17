package internal

import (
	"danger-dodgers/pkg/auth"
	"danger-dodgers/pkg/db"
	"danger-dodgers/pkg/passwords"
	"net/mail"
	"time"
)

type UserService struct {
	db     db.Database[User]
	hasher passwords.Hasher
	auth   auth.Authenticator
}

func NewUserService(
	db db.Database[User],
	hasher passwords.Hasher,
	auth auth.Authenticator,
) *UserService {
	return &UserService{
		db:     db,
		hasher: hasher,
		auth:   auth,
	}
}

const (
	REFRESH_TOKEN_LIFESPAN = time.Hour * 48
	MAX_EMAIL_LENGTH       = 100
	MAX_USERNAME_LENGTH    = 100
	MAX_NAME_LENGTH        = 100
	MAX_PASSWORD_LENGTH    = 50
	MIN_PASSWORD_LENGTH    = 7
)

func validateUsername(username string) error {
	// Check to see if username is blank.
	if username == "" {
		return &BlankUsernameError{}
	}

	// Check to see if username is too large.
	if len(username) > MAX_USERNAME_LENGTH {
		return &UsernameTooLargeError{
			max: MAX_USERNAME_LENGTH,
		}
	}

	return nil
}

func (service *UserService) check(user *User) error {
	err := validateUsername(user.Username)
	if err != nil {
		return err
	}

	// Check to see if user already exists.
	_, err = service.db.Get(user.Username)
	if err == nil {
		return &UserAlreadyExistsError{}
	}

	// Affirm the error returned is a "not found" error.
	_, ok := err.(*db.NotFoundError)
	if !ok {
		return err
	}

	return nil
}

func (service *UserService) password(user *User) error {
	err := validateUsername(user.Username)
	if err != nil {
		return err
	}

	// Check to see if user exists.
	stored, err := service.db.Get(user.Username)
	if err != nil {
		return err
	}

	// Compare user password.
	err = service.hasher.Compare(user.Password, stored.Password)
	if err != nil {
		return &InvalidPasswordError{}
	}

	return nil
}

func (service *UserService) Create(user *User) error {
	err := service.check(user)
	if err != nil {
		return err
	}

	// Check to see if name is blank.
	if user.Name == "" {
		return &BlankNameError{}
	}

	// Check to see if name is too large.
	if len(user.Name) > MAX_NAME_LENGTH {
		return &NameTooLargeError{
			max: MAX_NAME_LENGTH,
		}
	}

	// Check to see if password is too small.
	if len(user.Password) < MIN_PASSWORD_LENGTH {
		return &PasswordTooSmallError{
			min: MIN_PASSWORD_LENGTH,
		}
	}

	// Check to see if password is too large.
	if len(user.Password) > MAX_PASSWORD_LENGTH {
		return &PasswordTooLargeError{
			max: MAX_PASSWORD_LENGTH,
		}
	}

	// Validate email is not blank.
	if user.Email == "" {
		return &BlankEmailError{}
	}

	// Validate email is not too large.
	if len(user.Email) > MAX_EMAIL_LENGTH {
		return &EmailTooLargeError{
			max: MAX_EMAIL_LENGTH,
		}
	}

	// Validate email address structure.
	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		return &InvalidEmailError{}
	}

	// Hash password and replace password in struct.
	hash, err := service.hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash

	return service.db.Create(user.Username, user)
}

func (service *UserService) Token(user *User) (*Token, error) {

	// Check password validity.
	err := service.password(user)
	if err != nil {
		return nil, err
	}

	// Generate new token.
	token, err := service.auth.Generate(user.Username, REFRESH_TOKEN_LIFESPAN)
	if err != nil {
		return nil, err
	}

	return &Token{
		Token: token,
	}, err
}

func (service *UserService) Delete(user *User) error {
	// Check password validity.
	err := service.password(user)
	if err != nil {
		return err
	}

	return service.db.Delete(user.Username)
}

func (service *UserService) Get(id string) (*User, error) {
	err := validateUsername(id)
	if err != nil {
		return nil, err
	}

	// Check to see if user exists.
	user, err := service.db.Get(id)
	if err != nil {
		return nil, err
	}

	// Redact password.
	user.Password = ""

	return user, nil
}
