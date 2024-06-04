package tests

import (
	"testing"
	"time"

	ssov1 "github.com/4aykovski/grpc_auth_protos/gen/go/sso"
	"github.com/4aykovski/grpc_auth_sso/tests/suite"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "my_secret_app1_key_for_tests"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	registerResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registerResp.GetUserId())

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResp.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, int(claims["app_id"].(float64)), appID)
	assert.Equal(t, int64(claims["user_id"].(float64)), registerResp.GetUserId())
	assert.Equal(t, claims["email"].(string), email)

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.AccessTokenTtl).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	registerResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registerResp.GetUserId())

	registerResp, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, registerResp.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "invalid email",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "invalid password",
		},
		{
			name:        "invalid email",
			email:       "invalid",
			password:    randomFakePassword(),
			expectedErr: "invalid email",
		},
		{
			name:        "invalid password and email",
			email:       "invalid",
			password:    "",
			expectedErr: "invalid email;invalid password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "empty email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid email",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "invalid password",
		},
		{
			name:        "invalid email",
			email:       "invalid",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid email",
		},
		{
			name:        "invalid password and email",
			email:       "invalid",
			password:    "",
			appID:       appID,
			expectedErr: "invalid email;invalid password",
		},
		{
			name:        "login failed",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid credentials",
		},
		{
			name:        "invalid app id",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "invalid app id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)
			require.NotEmpty(t, loginResp.GetUserId())

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}

}
