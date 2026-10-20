package config

import (
	"time"
)

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JwtConfig struct {
	SecretKey             string
	ExpirationTime        time.Duration
	RefreshExpirationTime time.Duration
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
	S3Endpoint      string
}

type UploadConfig struct {
	MaxFileSize      int64
	Path             string
	// AllowedFileTypes []string
}


