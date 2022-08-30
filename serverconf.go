/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	VarMetricsEnabled      = "HTTP_METRICS_ENABLED"
	VarOpenApiEnabled      = "HTTP_OPEN_API_ENABLED"
	VarHTTPPort            = "HTTP_PORT"
	VarHTTPUser            = "HTTP_USER"
	VarHTTPPassword        = "HTTP_PASSWORD"
	VarHTTPRealm           = "HTTP_REALM"
	VarHTTPUploadLimit     = "HTTP_UPLOAD_LIMIT"
	VarHTTPUploadInMemSize = "HTTP_UPLOAD_IN_MEM_SIZE"
)

type ServerConfig struct {
	includeOpenAPI bool
}

// IncludeOpenAPI add swagger capability
func (c *ServerConfig) IncludeOpenAPI(value bool) {
	c.includeOpenAPI = value
}

func (c *ServerConfig) MetricsEnabled() bool {
	return c.getBoolean(VarMetricsEnabled, true)
}

// SwaggerEnabled if swagger capability exists then enable or disable it based on configuration
func (c *ServerConfig) SwaggerEnabled() bool {
	return c.getBoolean(VarOpenApiEnabled, true)
}

func (c *ServerConfig) HttpPort() string {
	return c.getString(VarHTTPPort, "8080")
}

func (c *ServerConfig) HttpUser() string {
	return c.getString(VarHTTPUser, "admin")
}

func (c *ServerConfig) HttpPwd() string {
	value := os.Getenv(VarHTTPPassword)
	if len(value) == 0 {
		// set as default value
		value = RandomPwd(20, true)
		log.Printf("undefined http password, setting it to: '%s'\n", value)
	}
	return value
}

func (c *ServerConfig) HttpRealm() string {
	return c.getString(VarHTTPRealm, "*")
}

func (c *ServerConfig) getBoolean(varName string, defaultValue bool) bool {
	value := os.Getenv(varName)
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		// set as default value
		enabled = defaultValue
	}
	return enabled
}

func (c *ServerConfig) getString(varName string, defaultValue string) string {
	value := os.Getenv(varName)
	if len(value) == 0 {
		// set as default value
		value = defaultValue
	}
	return value
}

func (c *ServerConfig) BasicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.HttpUser(), c.HttpPwd()))))
}

func (c *ServerConfig) HttpUploadLimit() int64 {
	limit, err := strconv.ParseInt(c.getString(VarHTTPUploadLimit, "250"), 0, 0)
	if err != nil {
		log.Fatalf("invalid upload limit specified")
	}
	return limit
}

func (c *ServerConfig) HttpUploadInMemorySize() int64 {
	limit, err := strconv.ParseInt(c.getString(VarHTTPUploadInMemSize, "150"), 0, 0)
	if err != nil {
		log.Fatalf("invalid upload limit specified")
	}
	return limit
}
