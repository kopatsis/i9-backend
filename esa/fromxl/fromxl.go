package fromxl

import (
	"strconv"
	"strings"

	"fmt"
	"fulli9/workoutgen2/datatypes"

	"github.com/xuri/excelize/v2"
)

func EnterEx() []datatypes.Exercise {
	f, err := excelize.OpenFile("esa/i9ea.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	end := []datatypes.Exercise{}

	row := 2
	for row < 252 {

		name, err := f.GetCellValue("Main", "A"+strconv.Itoa(row))
		if err != nil {
			fmt.Println(err)
			row++
			continue
		}

		blocked, err := f.GetCellValue("Main", "M"+strconv.Itoa(row))
		if err != nil {
			fmt.Println(err)
			row++
			continue
		}

		if name == "" {
			row = 253
		} else if blocked != "" {
			row++
			continue
		} else {
			parent, err := f.GetCellValue("Main", "B"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			minlevelSt, err := f.GetCellValue("Main", "C"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			minlevel, err := strconv.ParseFloat(minlevelSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			maxlevelSt, err := f.GetCellValue("Main", "D"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			maxlevel, err := strconv.ParseFloat(maxlevelSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			plyoSt, err := f.GetCellValue("Main", "E"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			plyo, err := strconv.Atoi(plyoSt)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			startSt, err := f.GetCellValue("Main", "F"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			start, err := strconv.ParseFloat(startSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			bodyparts := []int{}
			bodyString, err := f.GetCellValue("Main", "G"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}
			splitSt := strings.Split(strings.ReplaceAll(bodyString, " ", ""), ",")
			for _, str := range splitSt {
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println(err)
					row++
					continue
				}
				bodyparts = append(bodyparts, num)
			}

			minRepSt, err := f.GetCellValue("Main", "H"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			minRep, err := strconv.Atoi(minRepSt)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcASt, err := f.GetCellValue("Main", "I"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcA, err := strconv.ParseFloat(calcASt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcBSt, err := f.GetCellValue("Main", "J"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcB, err := strconv.ParseFloat(calcBSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcCSt, err := f.GetCellValue("Main", "K"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			calcC, err := strconv.ParseFloat(calcCSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			blockSplitSt, err := f.GetCellValue("Main", "L"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			inSplit := true
			if blockSplitSt != "" {
				inSplit = false
			}

			current := datatypes.Exercise{
				Name:         name,
				Parent:       parent,
				MinLevel:     float32(minlevel),
				MaxLevel:     float32(maxlevel),
				MinReps:      minRep,
				PlyoRating:   plyo,
				StartQuality: float32(start),
				BodyParts:    bodyparts,
				RepVars:      [3]float32{float32(calcA), float32(calcB), float32(calcC)},
				InSplits:     inSplit,
			}
			end = append(end, current)
		}
		row++
	}

	return end
}

func EnterSt() []datatypes.Stretch {
	f, err := excelize.OpenFile("esa/i9sa.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	end := []datatypes.Stretch{}

	row := 2
	for row < 252 {

		name, err := f.GetCellValue("Main", "A"+strconv.Itoa(row))
		if err != nil {
			fmt.Println(err)
			row++
			continue
		}

		if name == "" {
			row = 253
		} else {
			minlevelSt, err := f.GetCellValue("Main", "B"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			minlevel, err := strconv.ParseFloat(minlevelSt, 32)
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			status, err := f.GetCellValue("Main", "C"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}

			bodyparts := []int{}
			bodyString, err := f.GetCellValue("Main", "D"+strconv.Itoa(row))
			if err != nil {
				fmt.Println(err)
				row++
				continue
			}
			splitSt := strings.Split(strings.ReplaceAll(bodyString, " ", ""), ",")
			for _, str := range splitSt {
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println(err)
					row++
					continue
				}
				bodyparts = append(bodyparts, num)
			}

			current := datatypes.Stretch{
				Name:      name,
				Status:    status,
				MinLevel:  float32(minlevel),
				BodyParts: bodyparts,
			}
			end = append(end, current)
		}
		row++
	}

	return end
}
