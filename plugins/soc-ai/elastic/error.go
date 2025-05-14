package elastic

import (
	"errors"

	"github.com/threatwinds/go-sdk/catcher"
)

func RegisterError(message string, id string) {
	err := IndexStatus(id, "Error", "update")
	if err != nil {
		_ = catcher.Error("error while indexing error in elastic: %v", err, nil)
	}
	_ = catcher.Error("%s", errors.New(message), nil)
}
