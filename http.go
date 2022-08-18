package bear

// http 400 and 500 status codes
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#redirection_messages
const (
	// 400 error codes
	HTTPBadRequest                  = 400
	HTTPUnauthorized                = 401
	HTTPPaymentRequired             = 402
	HTTPForbidden                   = 403
	HTTPNotFound                    = 404
	HTTPMethodNotAllowed            = 405
	HTTPNotAcceptable               = 406
	HTTPProxyAuthenticatoinRequired = 407
	HTTPRequestTimeout              = 408
	HTTPConflict                    = 409
	HTTPGone                        = 410
	HTTPLengthRequired              = 411
	HTTPReconditionFailed           = 412
	HTTPPayloadTooLarge             = 412
	HTTPURITooLong                  = 414
	HTTPUnsupportedMediaType        = 415
	HTTPRangeNotSatisfiable         = 416
	HTTPExpectationFailed           = 417
	HTTPImATeapot                   = 418
	HTTPMisdirectedRequest          = 421
	HTTPUnprocessableEntity         = 422
	HTTPLocked                      = 423
	HTTPFailedDependency            = 424
	HTTPTooEarly                    = 425
	HTTPUpgradeRequired             = 426
	HTTPPreconditionRequired        = 428
	HTTPTooManyRequests             = 429
	HTTPRequestHeaderFieldsTooLarge = 431
	HTTPUnavailableForLegalReasons  = 451

	// 500 error codes
	HTTPInternalServerError           = 500
	HTTPNotImplemented                = 501
	HTTPBadGateway                    = 502
	HTTPServiceUnavailable            = 503
	HTTPGatewatyTimeout               = 504
	HTTPVersionNotSupported           = 505
	HTTPVariantAlsoNegotiates         = 506
	HTTPInsufficientStorage           = 507
	HTTPLoopDetected                  = 508
	HTTPNotExtended                   = 510
	HTTPNetworkAuthenticationRequired = 511
)
