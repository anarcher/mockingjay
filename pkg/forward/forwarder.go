package forward

import (
	"net/http"
)

type Forwarder interface {
	Forward(http.ResponseWriter, *http.Request) error
}
