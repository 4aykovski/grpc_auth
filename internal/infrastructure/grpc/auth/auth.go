package auth

import (
	"context"
	"fmt"
	"strings"

	ssov1 "github.com/4aykovski/grpc_auth_protos/gen/go/sso"
	authservice "github.com/4aykovski/grpc_auth_sso/internal/service/auth"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(ctx context.Context, dto authservice.LoginDTO) (string, error)
	Register(ctx context.Context, dto authservice.RegisterDTO) (int64, error)
	IsAdmin(ctx context.Context, dto authservice.IsAdminDTO) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer

	validate    *validator.Validate
	authService AuthService
}

func Register(gRPC *grpc.Server, authService AuthService) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		validate:    validator.New(),
		authService: authService,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := validateLoginRequest(req, s.validate); err != nil {
		var errMsgs []string
		for _, err := range err {
			errMsgs = append(errMsgs, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, strings.Join(errMsgs[:], ";"))
	}

	token, err := s.authService.Login(ctx, authservice.LoginDTO{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppId:    int(req.GetAppId()),
	})
	if err != nil {
		// TODO: add error handling
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validateRegisterRequest(req, s.validate); err != nil {
		var errMsgs []string
		for _, err := range err {
			errMsgs = append(errMsgs, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, strings.Join(errMsgs[:], ";"))
	}

	userId, err := s.authService.Register(ctx, authservice.RegisterDTO{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		// TODO: add error handling
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}

func validateLoginRequest(req *ssov1.LoginRequest, validate *validator.Validate) []error {
	var errs []error

	email := req.GetEmail()
	if err := validate.Var(email, "required,email"); err != nil {
		errs = append(errs, fmt.Errorf("invalid email"))
	}

	password := req.GetPassword()
	if err := validate.Var(password, "required"); err != nil {
		errs = append(errs, fmt.Errorf("invalid password"))
	}

	appId := req.GetAppId()
	if err := validate.Var(appId, "required"); err != nil {
		errs = append(errs, fmt.Errorf("invalid appId"))
	}

	return errs
}

func validateRegisterRequest(req *ssov1.RegisterRequest, validate *validator.Validate) []error {
	var errs []error

	email := req.GetEmail()
	if err := validate.Var(email, "required,email"); err != nil {
		errs = append(errs, fmt.Errorf("invalid email"))
	}

	password := req.GetPassword()
	if err := validate.Var(password, "required"); err != nil {
		errs = append(errs, fmt.Errorf("invalid password"))
	}

	return errs
}
