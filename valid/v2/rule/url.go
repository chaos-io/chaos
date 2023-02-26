package rule

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/chaos-io/chaos/valid/v2/inspection"
)

// IsURL check if string value is a valid URL.
func IsURL(v *inspection.Inspected) error {
	if k := v.Indirect.Kind(); k != reflect.String {
		return fmt.Errorf("%s: %w", k, ErrInvalidType)
	}

	s := v.Indirect.String()
	_, err := url.Parse(s)
	return err
}
