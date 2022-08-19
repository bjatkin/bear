package bearhttp

// http 400 and 500 status codes
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#redirection_messages
const (
	// 400 error codes
	BadRequest                  = 400
	Unauthorized                = 401
	PaymentRequired             = 402
	Forbidden                   = 403
	NotFound                    = 404
	MethodNotAllowed            = 405
	NotAcceptable               = 406
	ProxyAuthenticatoinRequired = 407
	RequestTimeout              = 408
	Conflict                    = 409
	Gone                        = 410
	LengthRequired              = 411
	ReconditionFailed           = 412
	PayloadTooLarge             = 412
	URITooLong                  = 414
	UnsupportedMediaType        = 415
	RangeNotSatisfiable         = 416
	ExpectationFailed           = 417
	ImATeapot                   = 418
	MisdirectedRequest          = 421
	UnprocessableEntity         = 422
	Locked                      = 423
	FailedDependency            = 424
	TooEarly                    = 425
	UpgradeRequired             = 426
	PreconditionRequired        = 428
	TooManyRequests             = 429
	RequestHeaderFieldsTooLarge = 431
	UnavailableForLegalReasons  = 451

	// 500 error codes
	InternalServerError           = 500
	NotImplemented                = 501
	BadGateway                    = 502
	ServiceUnavailable            = 503
	GatewatyTimeout               = 504
	VersionNotSupported           = 505
	VariantAlsoNegotiates         = 506
	InsufficientStorage           = 507
	LoopDetected                  = 508
	NotExtended                   = 510
	NetworkAuthenticationRequired = 511
)
