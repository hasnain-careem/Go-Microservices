package errors

import (
	"github.com/hasnain-zafar/go-microservices/common/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorHandler struct {
	logger *logger.Logger
}

func NewErrorHandler(log *logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: log,
	}
}

func (e *ErrorHandler) HandleNotFound(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	return status.Errorf(codes.NotFound, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleInvalidArgument(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	logger.IncrementValidationErrorCount()
	return status.Errorf(codes.InvalidArgument, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleDatabaseError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	logger.IncrementDBErrorCount()
	return status.Errorf(codes.Internal, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleInternalError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	return status.Errorf(codes.Internal, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleUnauthenticated(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	return status.Errorf(codes.Unauthenticated, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandlePermissionDenied(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	return status.Errorf(codes.PermissionDenied, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleNetworkError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	logger.IncrementNetworkErrorCount()
	return status.Errorf(codes.Unavailable, "%s: %v", msg, err)
}

func (e *ErrorHandler) LogError(msg string, err error) {
	e.logger.Error(msg, "error", err)
}