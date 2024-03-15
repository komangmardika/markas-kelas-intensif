package controller

import (
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"pertemuan5/model"
	"sync"
)

func GetCar(ctx *fiber.Ctx) error {
	// two channel
	cars := make(chan model.Car)
	done := make(chan bool)
	// load file csv
	fileCsv, err := os.Open("csv/cars_500.csv")
	if err != nil {
		fmt.Println("cant get file", err)
	}

	defer func(fileCsv *os.File) {
		err := fileCsv.Close()
		if err != nil {
			fmt.Println("cant close file")
		}
	}(fileCsv)

	reader := csv.NewReader(fileCsv)

	const numWorkers = 2
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processRows(reader, cars, &wg)
	}

	// Collect results from workers
	go func() {
		wg.Wait()
		close(cars)
		done <- true
	}()

	return ctx.Send(marshal)
}
