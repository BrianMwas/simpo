package gapi

import (
	"context"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"
	"simplebank/worker"
	"time"

	"github.com/hibiken/asynq"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser is grpc implementation for CreateUser API
func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Handle request violations
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentErr(violations)
	}

	// Hash password for secrecy
	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	// Create transaction params for failure rollback
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			FullName: req.GetFullName(), Email: req.GetEmail(), Username: req.GetUsername(), HashedPassword: hashPassword,
		},
		AfterCreate: func(user db.User) error {
			// Create a background task to run on the background
			// After successful user creation
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			// Send verification mail
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	// Access the transaction success result
	txResult, err := server.store.CreateUserTX(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Return an error if the action finds an existing account
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)
			}
		}
		// Return error for create user failure
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)

	}

	// Parse our response from the results sent.
	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

// Validate CreateUserRequest data
func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("mail", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}
	return violations
}
