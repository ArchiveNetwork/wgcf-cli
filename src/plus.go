package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// UGLY CODE , Pull Request is welcome
func plus(filePath string) error {
	var times sync.WaitGroup
	var currentStep int32
	var err error
	var id string
	ctx, cancel := context.WithCancel(context.Background())

	if _, id, err = readConfigFile(filePath); err != nil {
		panic(err)
	}

	for {
		times.Add(1)
		atomic.AddInt32(&currentStep, 1)
		time.Sleep(500 * time.Millisecond)
		go func(index int) {
			var publicKey string
			defer times.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			if _, publicKey, err = GenerateKey(); err != nil {
				panic(err)
			}

			installID := RandStringRunes(22, nil)
			fcmtoken := RandStringRunes(134, nil)

			payload := []byte(
				`{
					"key":"` + publicKey + `",
					"install_id":"` + installID + `",
					"fcm_token":"` + installID + `:APA91b` + fcmtoken + `",
					"tos":"` + time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + `",
					"model":"Android",
					"referrer": "` + id + `",
					"serial_number":"` + installID + `"
				}`,
			)
			fmt.Println("Registering...", index, "times")

			if _, err = request(payload, "", "", "register"); err != nil {
				fmt.Println(err)
				fmt.Println("Waiting for 30 seconds...")
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				time.Sleep(30 * time.Second)
				plus(filePath)
			}
			fmt.Println("						Added", index, "GB")
		}(int(currentStep))
	}
}
