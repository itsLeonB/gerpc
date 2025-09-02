package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errorInterceptor handles errors and panics in gRPC handlers
type errorInterceptor struct {
	logger ezutil.Logger
}

func NewErrorInterceptor(logger ezutil.Logger) Interceptor {
	return &errorInterceptor{logger}
}

// Handle is the main interceptor function that processes errors and panics
func (ei *errorInterceptor) Handle(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			ei.handlePanic(r, ctx, info)
			appError := ungerr.InternalServerError()
			err = status.Error(codes.Code(appError.GrpcStatus()), appError.Error())
		}
	}()

	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	// Already a gRPC status → just return
	if _, ok := status.FromError(err); ok {
		return resp, err
	}

	// Check if it's already an AppError
	if appErr, ok := err.(ungerr.AppError); ok {
		return resp, status.Error(codes.Code(appErr.GrpcStatus()), appErr.Error())
	}

	// Handle other errors
	appError := ei.constructAppError(err, info)
	return resp, status.Error(codes.Code(appError.GrpcStatus()), appError.Error())
}

// constructAppError converts various error types into AppError
func (ei *errorInterceptor) constructAppError(err error, info *grpc.UnaryServerInfo) ungerr.AppError {
	// First, try to unwrap with eris to get the original error
	originalErr := eris.Unwrap(err)
	if originalErr == nil {
		// No eris wrapping found - this means the error wasn't properly wrapped
		// Log the location where the error occurred
		return ei.logUnwrappedError(err, info)
	}

	// Handle known error types from eris-wrapped errors
	switch originalErr := originalErr.(type) {
	case validator.ValidationErrors:
		var errors []string
		for _, e := range originalErr {
			errors = append(errors, e.Error())
		}
		return ungerr.ValidationError(errors)

	case *json.SyntaxError:
		return ungerr.BadRequestError("invalid json")

	case *json.UnmarshalTypeError:
		return ungerr.BadRequestError(fmt.Sprintf("invalid value for field %s", originalErr.Field))

	default:
		// Handle common error patterns
		errStr := originalErr.Error()

		// EOF error from json package is unexported
		if originalErr == io.EOF || errStr == "EOF" {
			return ungerr.BadRequestError("missing request body")
		}

		// Check for network-related errors that might be client errors
		if strings.Contains(errStr, "connection reset by peer") ||
			strings.Contains(errStr, "broken pipe") ||
			strings.Contains(errStr, "context canceled") ||
			strings.Contains(errStr, "context deadline exceeded") {
			return ungerr.BadRequestError("connection error")
		}

		// This is an eris-wrapped error but not a known type
		// Log with full stack trace and mask from user
		return ei.logAndMaskError(err)
	}
}

// logUnwrappedError handles errors that weren't properly wrapped with eris
func (ei *errorInterceptor) logUnwrappedError(err error, info *grpc.UnaryServerInfo) ungerr.AppError {
	// This function helps you identify where errors are being added without proper wrapping
	ei.logger.Error("UNWRAPPED ERROR DETECTED - Please add eris.Wrap() or return ungerr.AppError")
	ei.logger.Errorf("Error type: %T", err)
	ei.logger.Errorf("Error message: %s", err.Error())
	ei.logger.Errorf("gRPC method: %s", info.FullMethod)
	ei.logger.Errorf("Server: %s", info.Server)

	ei.logger.Error("Stack trace from error location:")
	ei.logger.Errorf("%+v", err)

	// Return a masked error to the user
	return ungerr.InternalServerError()
}

// logAndMaskError handles eris-wrapped errors that need to be masked from users
func (ei *errorInterceptor) logAndMaskError(err error) ungerr.AppError {
	ei.logger.Errorf("Unhandled eris-wrapped error of type: %T", err)
	ei.logger.Error("Full stack trace:")
	ei.logger.Error(eris.ToString(err, true))

	return ungerr.InternalServerError()
}

// handlePanic recovers from panics and converts them to structured errors
func (ei *errorInterceptor) handlePanic(r interface{}, ctx context.Context, info *grpc.UnaryServerInfo) {
	// Log the panic with full stack trace
	ei.logger.Error("PANIC RECOVERED in gRPC handler")
	ei.logger.Errorf("gRPC method: %s", info.FullMethod)
	ei.logger.Errorf("Server: %s", info.Server)
	ei.logger.Errorf("Panic value: %v", r)
	ei.logger.Errorf("Panic type: %T", r)

	// Log context information if available
	if deadline, ok := ctx.Deadline(); ok {
		ei.logger.Errorf("Context deadline: %v", deadline)
	}
	if ctx.Err() != nil {
		ei.logger.Errorf("Context error: %v", ctx.Err())
	}

	// Print stack trace
	ei.logger.Error("Stack trace:")
	ei.logger.Error(string(debug.Stack()))

	// Try to convert panic to a meaningful error
	switch panicValue := r.(type) {
	case string:
		// Handle string panics (often from panic("message"))
		if strings.Contains(panicValue, "index out of range") ||
			strings.Contains(panicValue, "slice bounds out of range") {
			ei.logger.Error("Array/slice bounds panic detected")
		} else if strings.Contains(panicValue, "nil pointer dereference") {
			ei.logger.Error("Nil pointer dereference panic detected")
		} else {
			ei.logger.Errorf("String panic: %s", panicValue)
		}

	case runtime.Error:
		// Handle runtime errors (nil pointer, index out of bounds, etc.)
		ei.logger.Errorf("Runtime error panic: %v", panicValue)
		switch panicValue.Error() {
		case "runtime error: invalid memory address or nil pointer dereference":
			ei.logger.Error("Nil pointer dereference detected")
		case "runtime error: index out of range":
			ei.logger.Error("Index out of range detected")
		case "runtime error: slice bounds out of range":
			ei.logger.Error("Slice bounds out of range detected")
		}

	default:
		// Unknown panic type
		ei.logger.Errorf("Unknown panic type: %T, value: %v", r, r)
	}
}
