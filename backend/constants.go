package main

const (
	HeaderContentType     = "Content-Type"
	HeaderCORSOrigin      = "Access-Control-Allow-Origin"
	HeaderCORSMethods     = "Access-Control-Allow-Methods"
	HeaderCORSHeaders     = "Access-Control-Allow-Headers"
	HeaderOrigin          = "Origin"
	DefaultTimeFormat     = "2006-01-02"
	DefaultAllowedOrigins = "*"

	ContentTypeJSON = "application/json"

	CORSMethods = "GET, POST, OPTIONS"
	CORSHeaders = "Content-Type, Authorization"

	APIBasePath            = "/api/v1"
	HealthPath             = "/health"
	FestivalsPath          = "/api/festivals"
	LoginPath              = "/api/auth/login"
	VerifyPath             = "/api/auth/verify"
	FestivalsBreweriesPath = "/api/festivals/"

	AppVersion = "1.0.0"

	DefaultErrorMessage = "Internal server error"
)
