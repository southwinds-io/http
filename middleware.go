/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

// LoggingMiddleware log http requests to stdout
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request
		path, _ := url.PathUnescape(r.URL.Path)
		fmt.Printf("request from: %s %s %s\n", r.RemoteAddr, r.Method, path)
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware determines if the request is authenticated
func (s *Server) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// holds user principal
		var (
			user    *UserPrincipal
			matched bool
			err     error
		)
		// loop through specific authentication by URL path
		for urlPattern, authenticate := range s.Auth {
			// if the request URL match the authentication function pattern
			matched, err = regexp.Match(urlPattern, []byte(r.URL.Path))
			// regex error?
			if err != nil {
				// Write an error and stop the handler chain
				log.Printf("authentication function error: %s\n", err)
				http.Error(w, "Authentication Error", http.StatusInternalServerError)
				return
			}
			// if the regex matched the URL path
			if matched {
				// if we have an authentication function defined
				if authenticate != nil {
					// then try and authenticate using the specified function
					user, err = authenticate(*r)
					// if authentication fails the there is no user principal returned
					if user == nil {
						var errMsg string
						if err != nil {
							errMsg = err.Error()
						} else {
							errMsg = "Unauthorised"
						}
						// Write an error and stop the handler chain
						http.Error(w, errMsg, http.StatusUnauthorized)
						return
					} else {
						// exit loop
						break
					}
				} else {
					break
				}
			}
		}
		// if not authenticated by a custom handler then use default handler
		if user == nil && !matched {
			// no specific authentication function matched the request URL, so tries
			// the default authentication function if it has been defined
			// if no function has been defined then do not authenticate the request
			if s.DefaultAuth != nil {
				// Don't need to authorize options
				if r.Method == http.MethodOptions {
					next.ServeHTTP(w, r)
				}
				// NOTE: any anonymous access user will be prompted is accessing via browser
				// if anonymous access using browser is required, a revision of design should be considered
				// if no Authorization header is found and the user agent is a browser
				if r.Header.Get("Authorization") == "" && isBrowser(r.Header.Get("User-Agent")) {
					// prompts a client to authenticate by setting WWW-Authenticate response header
					w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, s.realm))
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Printf("! unauthorised http request from: '%v'\n", r.RemoteAddr)
					return
				} else {
					// authenticate the request using the default handler
					user, err = s.DefaultAuth(*r)
					if user == nil {
						var errMsg string
						if err != nil {
							errMsg = err.Error()
						} else {
							errMsg = "Unauthorised"
						}
						// if the authentication failed, write an error and stop the handler chain
						http.Error(w, errMsg, http.StatusUnauthorized)
						return
					}
				}
			}
		}
		// create a user context containing the user principal
		userContext := context.WithValue(r.Context(), "User", user)
		// create a shallow copy of the request with the user context added to it
		req := r.WithContext(userContext)
		// pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, req)
	})
}

// AuthorisationMiddleware authorises the http request based on the rights in user principal in the request context
func (s *Server) AuthorisationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserPrincipal(r)
		// if no principal is found reject the request
		if user == nil || !user.Rights.RequestAllowed(s.realm, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CorsMiddleware(origin string, headers string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CorsMiddleware(): origin = %v", origin)
			log.Printf("CorsMiddleware(): headers = %v", headers)

			if len(origin) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if r.Method == http.MethodOptions {
				log.Printf("CorsMiddleware(): process OPTIONS")
				if len(headers) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", headers)
					w.WriteHeader(200)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Authorise handler functions using user principal access control lists
// wraps the authorization middleware to be used when wrapping specific handler functions
func (s *Server) Authorise(handler http.HandlerFunc) http.Handler {
	return handler
}

func isBrowser(userAgent string) bool {
	safari := strings.Contains(userAgent, "Safari")
	opera := strings.Contains(userAgent, "OP")
	edge := strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Edge")
	firefox := strings.Contains(userAgent, "Firefox")
	chrome := strings.Contains(userAgent, "Chrome")
	return safari || opera || edge || firefox || chrome
}
