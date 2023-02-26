package rule

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/chaos-io/chaos/valid/v2/inspection"
)

// MustMatch matches given regular expression with inspected value as string
func MustMatch(regex *regexp.Regexp) Rule {
	return func(v *inspection.Inspected) error {
		if k := v.Indirect.Kind(); k != reflect.String {
			return fmt.Errorf("%s: %w", k, ErrInvalidType)
		}
		if !regex.MatchString(v.Indirect.String()) {
			return ErrPatternMismatch
		}
		return nil
	}
}
