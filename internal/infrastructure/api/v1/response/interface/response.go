package responseinterface

import "io"

// Responder - response service interface
type Responder interface {
	Respond(w io.Writer, dataOrErr any)
}
