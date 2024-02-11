package stretches

import (
	"errors"
	"fulli9/shared"
	"slices"
)

func FilterStretches(level float32, stretches map[string][]shared.Stretch, bodyparts map[int]bool, bannedStretches []string) (map[string][]shared.Stretch, error) {
	newstatics := []shared.Stretch{}
	newdynamics := []shared.Stretch{}

	if _, ok := stretches["Static"]; !ok {
		return nil, errors.New("No static stretches in FilterStretches")
	}
	if _, ok := stretches["Dynamic"]; !ok {
		return nil, errors.New("No dynamic stretches in FilterStretches")
	}

	for _, str := range stretches["Static"] {
		if str.MinLevel <= level {

			allowed := false
			if bodyparts == nil {
				allowed = true
			}

			for _, part := range str.BodyParts {
				if _, ok := bodyparts[part]; ok {
					allowed = true
				}
			}

			if slices.Contains(bannedStretches, str.ID.Hex()) {
				allowed = false
			}

			if allowed {
				newstatics = append(newstatics, str)
			}

		}
	}

	for _, str := range stretches["Dynamic"] {
		if str.MinLevel <= level {

			allowed := false
			if bodyparts == nil {
				allowed = true
			}

			for _, part := range str.BodyParts {
				if _, ok := bodyparts[part]; ok {
					allowed = true
				}
			}

			if slices.Contains(bannedStretches, str.ID.Hex()) {
				allowed = false
			}

			if allowed {
				newdynamics = append(newdynamics, str)
			}
		}
	}

	stretches["Dynamic"] = newdynamics
	stretches["Static"] = newstatics

	return stretches, nil
}
