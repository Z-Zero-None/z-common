package codes

const (
	Internal           = 10000 // internal error
	Unknown            = 10001 // unknown error
	BadRequest         = 10002 // bad request
	Unauthorized       = 10003 // unauthorized
	Forbidden          = 10004 // forbidden
	NotFound           = 10005 // not found
	Conflict           = 10006 // already exists
	ResourceExhausted  = 10007 // resource exhausted
	ServiceMaintenance = 10008 // service maintenance
	FrequencyLimit     = 10009
)
