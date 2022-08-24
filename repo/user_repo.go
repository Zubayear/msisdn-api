package repo

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"huspass/ent"
	user2 "huspass/ent/user"
	"huspass/external"
	"huspass/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, user *model.User) (*model.UserDto, error)
}

func (u *UserRepoImpl) GetUser(ctx context.Context, user *model.User) (*model.UserDto, error) {
	username := user.Username
	userFromRepo, err := u.client.User.Query().Where(user2.Username(username)).First(ctx)
	if err != nil {
		return nil, err
	}
	return model.NewUserDto(userFromRepo.Username, userFromRepo.Password), nil
}

func (u *UserRepoImpl) CreateUser(ctx context.Context, user *model.User) error {
	id := uuid.New()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	_, err = u.client.User.Create().
		SetID(id).
		SetUsername(user.Username).
		SetPassword(string(hashedPass)).
		SetRoles(user.Roles).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

type UserRepoImpl struct {
	client *ent.Client
}

func NewUserRepoImpl(client *ent.Client) *UserRepoImpl {
	return &UserRepoImpl{client: client}
}

func UserDBImplProvider(h *external.Host) (*UserRepoImpl, error) {
	client, err := ent.Open("mysql", h.ConnString)
	if err != nil {
		return nil, err
	}
	return NewUserRepoImpl(client), nil
}
