package users

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"

	"time"

	"github.com/georgi-georgiev/passport"
	"github.com/georgi-georgiev/passport/facade"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/georgi-georgiev/passport/responses"
	"github.com/golang-jwt/jwt"
	"github.com/rotisserie/eris"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserService struct {
	notificationFacade *facade.NotificationFacade
	repository         *UserRepository
	roleService        *permissions.RoleService
	rightService       *permissions.RightService
	conf               *passport.Config
	log                *zap.Logger
}

type UserClaims struct {
	ID     string
	Role   string
	Rights []string
}

func NewUserService(notificationFacade *facade.NotificationFacade, repository *UserRepository, roleService *permissions.RoleService, rightService *permissions.RightService, conf *passport.Config, log *zap.Logger) *UserService {
	return &UserService{notificationFacade: notificationFacade, repository: repository, roleService: roleService, rightService: rightService, conf: conf, log: log}
}

func (s *UserService) CreateUser(ctx context.Context, username string, email string, password string, r string, isAdmin bool, rr []string) (*User, error) {

	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, eris.Wrap(err, "Could not get user by username")
	}

	if user != nil {
		return nil, eris.New("Username already exists")
	}

	dbUser, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, eris.Wrap(err, "Could not get user by email")
	}

	if dbUser != nil {
		return nil, eris.New("Email already exists")
	}

	hashedPassword, err := passport.Hash(password)
	if err != nil {
		return nil, eris.Wrap(err, "Could not hash password")
	}

	role, err := s.roleService.GetByName(ctx, r)
	if err != nil {
		return nil, eris.Wrap(err, "could not get role by name")
	}

	if role == nil {
		return nil, eris.New("could not find role")
	}

	rights := make([]*permissions.Right, 0)
	for _, rightName := range rr {
		right, err := s.rightService.GetByName(ctx, rightName)
		if err != nil {
			return nil, eris.Wrap(err, "could not get right by name")
		}

		if right == nil {
			return nil, eris.New("could not find right")
		}

		rights = append(rights, right)
	}

	token, err := generateCode(32)
	if err != nil {
		return nil, eris.Wrap(err, "could not generate token")
	}

	newUser := NewUser(token, username, email, hashedPassword, role, rights)

	if isAdmin {
		newUser.IsVerified = true
	}

	ID, err := s.repository.Create(ctx, newUser)
	if err != nil {
		return nil, eris.Wrap(err, "Cannot create user")
	}

	newUser.ID = ID

	if !isAdmin {
		s.notificationFacade.Publish(ctx, "email", facade.Message{
			Topic:     "email_verification",
			Header:    "Email Verification",
			Body:      fmt.Sprintf("Verification link: http://%s//verify/%s\n", s.conf.Swagger.Host, token),
			Params:    map[string]string{"email": newUser.Email},
			Meta:      nil,
			Timestamp: time.Now().Unix(),
		}, newUser.ID.Hex())
	}

	return newUser, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	err := s.repository.Verify(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, id primitive.ObjectID, email string, username string, isActive bool, changePassword bool) (*User, error) {

	u := &User{}
	found, err := s.repository.GetById(ctx, id, u)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	if u.Email != email {
		u.Email = email
	}

	if u.Username != username {
		u.Username = username
	}

	u.IsActive = isActive

	now := time.Now().UTC()
	u.UpdatedOn = &now

	if changePassword {
		go s.SendRecoveryEmail(ctx, u.Username)
	}

	err = s.repository.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id primitive.ObjectID) (bool, error) {
	isDeleted, err := s.repository.DeleteById(ctx, id)
	if err != nil {
		return false, err
	}

	return isDeleted, nil
}

func (s *UserService) BasicAuthToken(ctx context.Context, username, password string) (string, string, int, error) {

	//TODO: implement with only 1 query for optimization
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", 0, err
	}

	if user == nil || !passport.Match(password, user.Password) {
		return "", "", 0, eris.New("Username or password is wrong")
	}

	userClaims := s.MapToUserClaims(user)

	tokenString, exp, err := s.IssueAccessToken(userClaims)
	if err != nil {
		return "", "", 0, err
	}

	refreshTokenString, err := s.IssueRefreshToken(userClaims)
	if err != nil {
		return "", "", 0, err
	}

	return tokenString, refreshTokenString, exp, nil
}

func (s *UserService) IssueAccessToken(userClaims *UserClaims) (string, int, error) {

	keyData, err := os.ReadFile(s.conf.App.PrivKeyPath)
	if err != nil {
		return "", 0, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", 0, err
	}

	exp := int(time.Now().Add(time.Hour).UTC().Unix())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId": userClaims.ID,
		"exp":    exp,
		"role":   userClaims.Role,
		"rights": userClaims.Rights,
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", 0, err
	}

	return tokenString, exp, nil
}

func (s *UserService) IssueRefreshToken(userClaims *UserClaims) (string, error) {
	keyData, err := os.ReadFile(s.conf.App.PrivKeyPath)
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", err
	}

	exp := time.Now().AddDate(1, 0, 0).UTC().Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId": userClaims.ID,
		"exp":    exp,
		"role":   userClaims.Role,
		"rights": userClaims.Rights,
	})

	refreshTokenString, err := refreshToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

// RefreshToken refreshes existing token
func (s *UserService) RefreshToken(ctx context.Context, t string) (string, int, error) {
	userClaims, err := s.ValidateToken(ctx, t)
	if err != nil {
		accessToken, exp, err := s.IssueAccessToken(userClaims)
		return accessToken, exp, err
	}

	return "", 0, err
}

func (s *UserService) GetPublicKey() (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(s.conf.App.PubKeyPath)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

func (s *UserService) ValidateToken(ctx context.Context, t string) (*UserClaims, error) {
	key, err := s.GetPublicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["userId"].(string)
		if !ok {
			return nil, eris.New("user id is not string")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, eris.New("role is not string")
		}

		rights, ok := claims["rights"].([]interface{})
		if !ok {
			return nil, eris.New("rights is not interface array")
		}

		r := make([]string, 0)
		for _, right := range rights {
			r = append(r, right.(string))
		}

		return &UserClaims{
			ID:     userId,
			Role:   role,
			Rights: r,
		}, nil
	}

	return nil, nil
}

func (s *UserService) GetUserByToken(ctx context.Context, t string) (*User, error) {
	key, err := s.GetPublicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdHex, ok := claims["userId"].(string)
		if !ok {
			return nil, eris.New("user id is not string")
		}

		userId, err := primitive.ObjectIDFromHex(userIdHex)
		if err != nil {
			return nil, err
		}

		u, err := s.GetById(ctx, userId)

		if err != nil {
			return nil, err
		}

		return u, nil
	}
	return nil, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u, err := s.repository.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUserById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	u := &User{}
	found, err := s.repository.GetById(ctx, id, u)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return u, nil
}

func (s *UserService) GetById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	u := &User{}
	found, err := s.repository.GetById(ctx, id, u)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return u, nil
}

func (s *UserService) LoadPermisions(ctx context.Context, user *User) (*permissions.Role, []permissions.Right, error) {
	role, err := s.roleService.GetById(ctx, user.Role)
	if err != nil {
		return nil, nil, err
	}

	rights := make([]permissions.Right, 0)

	if len(user.Rights) > 0 {
		rights, err = s.rightService.GetManyByIds(ctx, user.Rights)
		if err != nil {
			return nil, nil, err
		}
	}

	return role, rights, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u, err := s.repository.GetByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*User, error) {
	us, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "could no get all users")
	}

	return us, nil
}

func (s *UserService) SendRecoveryEmail(ctx context.Context, email string) {
	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("could not get user by email")
		return
	}

	if u == nil {
		s.log.Error("Email does not registered or email is absent")
		return
	}

	code, err := generateCode(6)
	if err != nil {
		s.log.Error("could not generate code")
		return
	}

	s.notificationFacade.Publish(ctx, "email", facade.Message{
		Topic:     "identity",
		Header:    "Password Recovery",
		Body:      fmt.Sprintf("Your recovery code is: %s", code),
		Meta:      nil,
		Timestamp: time.Now().Unix(),
	}, u.ID.Hex())

	hashedCode, err := passport.Hash(code)
	if err != nil {
		s.log.Error("could not hash code")
		return
	}

	err = s.repository.SetRecoveryCode(ctx, u.ID, hashedCode)
	if err != nil {
		s.log.Error("could not set recovery code")
		return
	}
}

// ExchangeRecoveryCode exchanges recovery code for a password resetting one
func (s *UserService) ExchangeRecoveryCode(ctx context.Context, email string, code string) (string, error) {
	generalErrorMsg := "Username does not registered or recovery process has not been initiated"

	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return "", eris.Wrap(err, "could not get user by email")
	}
	if u == nil {
		return "", eris.New(generalErrorMsg)
	}

	existingCodeHash, err := s.repository.GetRecoveryCode(ctx, u.ID)
	if err != nil {
		return "", eris.Wrap(err, "could not get recovery code")
	}

	if existingCodeHash == "" {
		return "", eris.New(generalErrorMsg)
	}

	if !passport.Match(code, existingCodeHash) {
		return "", eris.New("Provided recovery code does not match")
	}

	resettingCode, err := generateCode(10)
	if err != nil {
		return "", eris.Wrap(err, "could no generate code")
	}

	resettingCodeHash, err := passport.Hash(resettingCode)
	if err != nil {
		return "", eris.Wrap(err, "could no hash resetting code")
	}

	err = s.repository.SetResettingCode(ctx, u.ID, resettingCodeHash)
	if err != nil {
		return "", eris.Wrap(err, "could not set resetting code")
	}

	return resettingCode, nil
}

func (s *UserService) ResetPassword(ctx context.Context, email string, code string, newPassword string) error {
	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if u == nil {
		return eris.New("Username does not registered or recovery process has not been initiated")
	}

	existingCodeHash, err := s.repository.GetResettingCode(ctx, u.ID)
	if err != nil {
		return err
	}

	if !passport.Match(code, existingCodeHash) {
		return eris.New("Provided recovery code does not match")
	}

	hashedPassword, err := passport.Hash(newPassword)
	if err != nil {
		return err
	}

	err = s.repository.ResetPassword(ctx, u.ID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func generateCode(size int) (string, error) {
	token := make([]byte, size)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

func (s *UserService) MapToUserResponse(ctx context.Context, u *User) (*responses.UserResponse, error) {

	role, err := s.roleService.GetById(ctx, u.Role)
	if err != nil {
		return nil, err
	}

	righsNames := make([]string, 0)

	if len(u.Rights) > 0 {
		rights, err := s.rightService.GetManyByIds(ctx, u.Rights)
		if err != nil {
			return nil, err
		}

		for _, right := range rights {
			righsNames = append(righsNames, right.Name)
		}

	}

	return &responses.UserResponse{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Role:     role.Name,
		Rights:   righsNames,
	}, nil
}

func (s *UserService) MapToUserClaims(u *User) *UserClaims {
	rightsIds := make([]string, 0)
	for _, right := range u.Rights {
		rightsIds = append(rightsIds, right.Hex())
	}

	userClaims := &UserClaims{
		ID:     u.ID.Hex(),
		Role:   u.Role.Hex(),
		Rights: rightsIds,
	}

	return userClaims
}
