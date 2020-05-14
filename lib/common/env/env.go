package env

import (
	"net"
	"os"
)

// Environment List
const (
	EnvDevelopment = "development"
	EnvAlpha       = "alpha"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// Get return string of current environment flag
func Get() string {
	env := os.Getenv("ENTENV")

	if env != EnvAlpha && env != EnvStaging && env != EnvProduction {
		env = "development"
	}

	return env
}

// GetServerIpAddress is function to get the ip address from instance
func GetServerIpAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip += ipnet.IP.String()
			}
		}
	}
	return ip, err
}
