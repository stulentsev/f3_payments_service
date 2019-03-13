package conversion

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Lazy and wrong way to convert between structs
//
// To save some time and effort a lot of things serialise to the same JSON,
// while having different internal structure.
//
// Having different structures for API params, Command Data, Event Data and Read Models
// gives us a lot of flexibility, but that comes of a cost of extremely tedious conversions, as
// all those structures are incidentally similar.
//
// It's easy to imagine a V2 of the API, or a GRPC api with completely different structure,
// that needs to be completely decoupled from all other things.
//
// Used for:
//   - Given a swagger API definition we need to map the
//     payment request attributes to a command to CREATE or UPDATE domain commands
//   - Once we have domain commands we need to store their data to structured events
//   - When we want to create reader persistence we need to map the event data to the projection model
func Map(dest interface{}, src interface{}) error {
	jsonString, err := json.Marshal(src)

	if err != nil {
		return errors.Wrap(err, "marshal of source object failed")
	}

	err = json.Unmarshal(jsonString, dest)

	return err
}
