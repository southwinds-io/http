/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Control represents a user access control based on a URI that is part of a realm
type Control struct {
	// the control realm (e.g. application)
	Realm string
	// the control resource URI (e.g. typically but not exclusively, a restful endpoint URI)
	URI string
	// the method(s) used by the resource (e.g. POST for create, PUT for update, DELETE for delete)
	Method []string
}

func (c *Control) hasMethod(method string) bool {
	for _, m := range c.Method {
		if strings.ToUpper(strings.Trim(m, " ")) == strings.ToUpper(strings.Trim(method, " ")) {
			return true
		}
	}
	return false
}

func (c *Control) equal(ctl Control) bool {
	return strings.EqualFold(ctl.Realm, c.Realm) && ctl.URI == c.URI && equalMethods(c.Method, ctl.Method)
}

func equalMethods(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type Controls []Control

func (controls Controls) Add(ctl ...Control) Controls {
	var (
		found bool
		r     = controls
	)
	for _, c := range ctl {
		found = false
		for _, control := range controls {
			if control.equal(c) {
				found = true
				break
			}
		}
		if !found {
			r = append(r, c)
		}
	}
	return r
}

// Allowed returns true if the specified control matches one of the controls granted to the user
func (controls Controls) allowed(realm, uri, method string) bool {
	for _, c := range controls {
		matched, err := regexp.MatchString(c.URI, uri)
		if err != nil {
			log.Printf("WARNING: cannot match URI path, %s\n", err)
		}
		if (c.Realm == realm || c.Realm == "*") &&
			(matched || c.URI == "*") &&
			(c.hasMethod(method) || c.hasMethod("*")) {
			return true
		}
	}
	return false
}

// RequestAllowed returns true if the http request matches one of the controls granted to the user for the given realm
func (controls Controls) RequestAllowed(realm string, r *http.Request) bool {
	// extract user principal from the request context
	if principal := r.Context().Value("User"); principal != nil {
		if value, ok := principal.(*UserPrincipal); ok {
			return value.Rights.allowed(realm, r.RequestURI, r.Method)
		}
	}
	return false
}

func newControls(acl string) (Controls, error) {
	var result Controls
	// if acl is empty then return an empty list of controls
	if len(strings.Trim(acl, " ")) == 0 {
		return Controls{}, nil
	}
	parts := strings.Split(acl, ",")
	for _, part := range parts {
		control, err := newControl(part)
		if err != nil {
			return nil, err
		}
		result = append(result, control)
	}
	return result, nil
}

func newControl(ac string) (Control, error) {
	parts := strings.Split(ac, ":")
	if len(parts) != 3 {
		return Control{}, fmt.Errorf("Invalid control format '%s', it should be realm:uri:method\n", ac)
	}
	return Control{
		Realm:  parts[0],
		URI:    parts[1],
		Method: strings.Split(parts[2], "|"),
	}, nil
}
