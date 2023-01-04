package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"time"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertUserSession(sessionId string, accessToken string, accessTokenExpiresAt time.Time, refreshToken string, refreshTokenExpiresAt time.Time, user db.User) *pb.LoginUserResponse {
	return &pb.LoginUserResponse{
		SessionId:             sessionId,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenExpiresAt),
		User:                  convertUser(user),
	}
}
