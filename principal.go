/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import "time"

// UserPrincipal represents a logged on user and the access controls granted to them
type UserPrincipal struct {
	// the user Username used as a unique identifier (typically the user email address)
	Username string `json:"username"`
	// a list of rights or access controls granted to the user
	Rights Controls `json:"acl,omitempty"`
	// the time the principal was Created
	Created time.Time `json:"created"`
	// any context associated to the principal
	Context interface{}
}
