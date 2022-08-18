package bear

// grpc return codes
// https://pkg.go.dev/google.golang.org/grpc/codes
const (
	GRPCOK                 = 0
	GRPCCanceled           = 1
	GRPCUnknown            = 2
	GRPCInvalidArggument   = 3
	GRPCDeadlineExceeded   = 4
	GRPCNotFound           = 5
	GRPCAlreadyExists      = 6
	GRPCPermissionDenied   = 7
	GRPCResourceExhausted  = 8
	GRPCFailedPrecondition = 9
	GRPCAborted            = 10
	GRPCOutOfRange         = 11
	GRPCUnimplemented      = 12
	GRPCInternal           = 13
	GRPCUnavailable        = 14
	GRPCDataLoss           = 15
	GRPCUnauthenticated    = 16
)
