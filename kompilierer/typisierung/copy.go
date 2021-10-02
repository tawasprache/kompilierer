package typisierung

import "github.com/barkimedes/go-deepcopy"

func copy(v interface{}) interface{} {
	return deepcopy.MustAnything(v)
}
