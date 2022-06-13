package auth

import (
	"context"

	"firebase.google.com/go/auth"
)

type AuthUtil struct {
	firebaseClient *auth.Client
	ctx            context.Context
}

func NewUtil(ctx context.Context) *AuthUtil {
	return &AuthUtil{
		firebaseClient: client,
		ctx:            ctx,
	}
}

func (util *AuthUtil) NewAdmin(email string, pass string) error {
	admin := &auth.UserToCreate{}
	admin.Email(email)
	admin.Password(pass)
	admin.EmailVerified(true)
	user, err := util.firebaseClient.CreateUser(util.ctx, admin)
	if err != nil {
		return err
	}

	adminClaims := map[string]interface{}{"admin": true}
	return util.firebaseClient.SetCustomUserClaims(util.ctx, user.UID, adminClaims)
}

func (util *AuthUtil) GetUser(email string) (*auth.UserRecord, error) {
	return util.firebaseClient.GetUserByEmail(util.ctx, email)
}
