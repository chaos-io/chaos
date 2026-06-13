package errorx

import (
	"errors"
	"fmt"
	"sync"
)

var ErrRegisterConflict = errors.New("errorx: register conflict")

var (
	registryMu sync.RWMutex
	registry   = make(map[int32]Definition)
)

func Register(defs ...Definition) error {
	registryMu.Lock()
	defer registryMu.Unlock()

	for _, def := range defs {
		def = def.normalized()
		if current, ok := registry[def.Code]; ok {
			if sameDefinition(current, def) {
				continue
			}
			return fmt.Errorf(
				"%w: code=%d current_message=%q new_message=%q current_affects_stability=%t new_affects_stability=%t",
				ErrRegisterConflict,
				def.Code,
				current.Message,
				def.Message,
				current.AffectsStability,
				def.AffectsStability,
			)
		}
		registry[def.Code] = def
	}
	return nil
}

func MustRegister(defs ...Definition) {
	if err := Register(defs...); err != nil {
		panic(err)
	}
}

func sameDefinition(left, right Definition) bool {
	return left.Code == right.Code &&
		left.Message == right.Message &&
		left.AffectsStability == right.AffectsStability
}
