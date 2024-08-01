package stretches

import (
	"errors"
	"fulli9/shared"
	"math"
	"slices"
)

func FilterStretches(level float32, stretches map[string][]shared.Stretch, bodyparts map[int]bool, user shared.User) (map[string][]shared.Stretch, error) {
	newstatics := []shared.Stretch{}
	newdynamics := []shared.Stretch{}

	if _, ok := stretches["Static"]; !ok {
		return nil, errors.New("no static stretches in FilterStretches")
	}
	if _, ok := stretches["Dynamic"]; !ok {
		return nil, errors.New("no dynamic stretches in FilterStretches")
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

			if slices.Contains(user.BannedStretches, str.ID.Hex()) {
				allowed = false
			}

			if allowed {
				if val, ok := user.StrFavoriteRates[str.ID.Hex()]; ok {
					str.Weight *= float32(math.Min(1.625, math.Max(0.6667, float64((1+val)/2))))
				}
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

			if slices.Contains(user.BannedStretches, str.ID.Hex()) {
				allowed = false
			}

			if allowed {
				if val, ok := user.StrFavoriteRates[str.ID.Hex()]; ok {
					str.Weight *= float32(math.Min(1.625, math.Max(0.6667, float64((1+val)/2))))
				}
				newdynamics = append(newdynamics, str)
			}
		}
	}

	stretches["Dynamic"] = newdynamics
	stretches["Static"] = newstatics

	return stretches, nil
}
