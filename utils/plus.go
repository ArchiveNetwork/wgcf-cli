package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// UGLY CODE , Pull Request is welcome
func Plus(filePath string, test bool) error {
	var times sync.WaitGroup
	var err error
	var id string
	var currentStep int = 1
	ctx, cancel := context.WithCancel(context.Background())
	if fileType, err := GetFileType(filePath); fileType == "json" && err == nil {
		if _, id, err = GetTokenID(filePath); err != nil {
			panic(err)
		}
	} else if fileType == "ini" && err == nil {
		if _, id, err = IniGetTokenID(filePath); err != nil {
			panic(err)
		}
	}
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		<-signalCh
		cancel()
		fmt.Println("\nWaiting for Response...")
		go func() {
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			<-signalCh
			fmt.Println()
			os.Exit(1)
		}()
		times.Wait()
		fmt.Printf("Total added %d GB\n", currentStep)
		fmt.Println("Updating config file...")
		if fileType, err := GetFileType(filePath); fileType == "json" && err == nil {
			UpdateConfigFile(filePath)
		} else if fileType == "ini" && err == nil {
			UpdateIniConfig(filePath)
		}
		fmt.Println("Updated config file successfully")
		os.Exit(0)
	}()

	if !test {
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				if i := currentStep % 10; i == 0 {
					if fileType, err := GetFileType(filePath); fileType == "json" && err == nil {
						UpdateConfigFile(filePath)
					} else if fileType == "ini" && err == nil {
						UpdateIniConfig(filePath)
					}
				}
			}
		}()
	}

	for {
		time.Sleep(500 * time.Millisecond)
		times.Add(1)
		if test {
			select {
			case <-ctx.Done():
				times.Done()
				return nil
			default:
			}
		}
		select {
		case <-ctx.Done():
			times.Done()
		default:
			go func(index int) {
				var publicKey string
				currentStep++
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
					times.Done()
					fmt.Println(err)
					fmt.Println("Waiting for 30 seconds...")
					select {
					case <-ctx.Done():
						return
					default:
					}
					cancel()
					time.Sleep(30 * time.Second)
					currentStep--
					ctx, cancel = context.WithCancel(context.Background())
				} else {
					times.Done()
					fmt.Println("						Added", index, "GB")
				}
			}(currentStep)
		}
	}
}
