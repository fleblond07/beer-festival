package main

const (
	DefaultPort           = "8080"
	DefaultFestivalsFile  = "./festivals.json"
	DefaultAllowedOrigins = "*"

	HeaderContentType = "Content-Type"
	HeaderCORSOrigin  = "Access-Control-Allow-Origin"
	HeaderCORSMethods = "Access-Control-Allow-Methods"
	HeaderCORSHeaders = "Access-Control-Allow-Headers"
	HeaderOrigin      = "Origin"

	ContentTypeJSON = "application/json"

	CORSMethods = "GET, OPTIONS"
	CORSHeaders = "Content-Type"

	APIBasePath   = "/api/v1"
	HealthPath    = "/health"
	FestivalsPath = "/api/festivals"

	AppVersion = "1.0.0"
)
