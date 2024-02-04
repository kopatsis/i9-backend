package stretches

import "fulli9/shared"

func FilterStretches(level float32, stretches map[string][]shared.Stretch, bodyparts map[int]bool) map[string][]shared.Stretch {
	newstatics := []shared.Stretch{}
	newdynamics := []shared.Stretch{}

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

			if allowed {
				newdynamics = append(newdynamics, str)
			}
		}
	}

	stretches["Dynamic"] = newdynamics
	stretches["Static"] = newstatics

	return stretches
}
