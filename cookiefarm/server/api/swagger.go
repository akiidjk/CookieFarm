package api

//go:generate swag init -g swagger.go -o docs --parseDependency --parseInternal --outputTypes json,yaml

// @title CookieFarm Server API
// @version 1.0
// @description CookieFarm REST API for authentication, configuration, flags, and stats.
// @BasePath /api/v1
// @schemes http https
// @accept json
// @produce json
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name token
// @description Session JWT cookie returned by `/api/v1/auth/login`.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use `Bearer <jwt-token>`.
