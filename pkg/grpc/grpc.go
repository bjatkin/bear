package beargrpc

// grpc return codes
// https://pkg.go.dev/google.golang.org/grpc/codes
const (
	OK                 = 0
	Canceled           = 1
	Unknown            = 2
	InvalidArggument   = 3
	DeadlineExceeded   = 4
	NotFound           = 5
	AlreadyExists      = 6
	PermissionDenied   = 7
	ResourceExhausted  = 8
	FailedPrecondition = 9
	Aborted            = 10
	OutOfRange         = 11
	Unimplemented      = 12
	Internal           = 13
	Unavailable        = 14
	DataLoss           = 15
	Unauthenticated    = 16
)
