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
)

var (
	UserVerifier       = FieldVerifier{}.WithModel(USER)
	VerifyUserID       = REQUIRED_STANDARD_VERIFIER(UserVerifier).WithField("ID").Build()
	VerifyUserEmail    = REQUIRED_STANDARD_VERIFIER(UserVerifier).WithField("email").Build()
	VerifyUserName     = REQUIRED_STANDARD_VERIFIER(UserVerifier).WithField("name").Build()
	VerifyUserPassword = PASSWORD_VERIFIER(UserVerifier).Build()
)

func (service *UserService) password(user *User) error {
	err := VerifyUserID(user.Email)
	if err != nil {
		return err
	}

	// Check to see if user exists.
	stored, err := service.db.Get(user.Email)
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

	err := VerifyUserEmail(user.Email)
	if err != nil {
		return err
	}

	// Validate email address structure.
	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		return &InvalidEmailError{}
	}
	// Make sure that ID does not already exist.
	err = Peek(user.Email, USER, service.db)
	if err != nil {
		return err
	}

	err = VerifyUserName(user.Name)
	if err != nil {
		return err
	}

	err = VerifyUserPassword(user.Password)
	if err != nil {
		return err
	}

	// Hash password and replace password in struct.
	hash, err := service.hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash

	return service.db.Create(user.Email, user)
}

func (service *UserService) Update(users *ModifyUser) error {
	err := service.password(users.CurrentUser)
	if err != nil {
		return err
	}

	if users.CurrentUser.Email != users.NewUser.Email {
		return &CannotUpdateEmailError{}
	}

	return service.db.Update(users.CurrentUser.Email, users.NewUser)
}

func (service *UserService) Token(user *User) (*Token, error) {

	// Check password validity.
	err := service.password(user)
	if err != nil {
		return nil, err
	}

	// Generate new token.
	token, err := service.auth.Generate(user.Email, REFRESH_TOKEN_LIFESPAN)
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

	return service.db.Delete(user.Email)
}

func (service *UserService) Get(id string) (*User, error) {
	err := VerifyUserID(id)
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
