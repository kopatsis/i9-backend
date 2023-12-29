package fromxl

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func AddStrToXL(createdSt map[string]string) {
	f, err := excelize.OpenFile("i9sa.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := f.Save(); err != nil {
			fmt.Println(err)
		}
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	insertrow, placeFound := 2, false
	for !placeFound {
		compat, err := f.GetCellValue("Ref", "A"+strconv.Itoa(insertrow))
		if err != nil {
			fmt.Println(err)
			insertrow++
			continue
		}
		if compat == "" {
			placeFound = true
		} else {
			insertrow++
		}
	}
	for id, stretch := range createdSt {
		f.SetCellValue("Ref", "A"+strconv.Itoa(insertrow), stretch)
		f.SetCellValue("Ref", "B"+strconv.Itoa(insertrow), id)
		insertrow++
	}

}

func AddExerToXL(createdEx map[string]string) {
	f, err := excelize.OpenFile("i9ea.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := f.Save(); err != nil {
			fmt.Println(err)
		}
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	insertrow, placeFound := 2, false
	for !placeFound {
		compat, err := f.GetCellValue("Ref", "C"+strconv.Itoa(insertrow))
		if err != nil {
			fmt.Println(err)
			insertrow++
			continue
		}
		if compat == "" {
			placeFound = true
		} else {
			insertrow++
		}
	}
	for id, exer := range createdEx {
		f.SetCellValue("Ref", "C"+strconv.Itoa(insertrow), exer)
		f.SetCellValue("Ref", "D"+strconv.Itoa(insertrow), id)
		insertrow++
	}
}
