package errors

import (
	"github.com/hasnain-zafar/go-microservices/common/logger"
	"github.com/hasnain-zafar/go-microservices/common/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorHandler struct {
	logger  *logger.Logger
	service string
}

func NewErrorHandler(log *logger.Logger) *ErrorHandler {
	serviceName := log.GetServiceName()
	return &ErrorHandler{
		logger:  log,
		service: serviceName,
	}
}

func (e *ErrorHandler) HandleNotFound(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "not_found")
	return status.Errorf(codes.NotFound, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleInvalidArgument(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "invalid_argument")
	logger.IncrementValidationErrorCount()
	return status.Errorf(codes.InvalidArgument, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleDatabaseError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "database")
	logger.IncrementDBErrorCount()
	return status.Errorf(codes.Internal, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleInternalError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "internal")
	return status.Errorf(codes.Internal, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleUnauthenticated(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "unauthenticated")
	return status.Errorf(codes.Unauthenticated, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandlePermissionDenied(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "permission_denied")
	return status.Errorf(codes.PermissionDenied, "%s: %v", msg, err)
}

func (e *ErrorHandler) HandleNetworkError(msg string, err error) error {
	e.logger.Error(msg, "error", err)
	metrics.IncrementErrorCounter(e.service, "network")
	logger.IncrementNetworkErrorCount()
	return status.Errorf(codes.Unavailable, "%s: %v", msg, err)
}

func (e *ErrorHandler) LogError(msg string, err error) {
	e.logger.Error(msg, "error", err)
}
