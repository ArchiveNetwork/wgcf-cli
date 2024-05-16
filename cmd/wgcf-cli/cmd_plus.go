//go:build with_plus

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var plusCmd = &cobra.Command{
	Use:   "plus",
	Short: "Recharge your account indefinitely",
	PreRun: func(cmd *cobra.Command, args []string) {
		client.New()
	},
	Run:     plus,
	PostRun: update,
}

func init() {
	rootCmd.AddCommand(plusCmd)
}

func plus(cmd *cobra.Command, args []string) {
	warn_str := "\033[1;93mWARN:\033[0m"
	error_str := "\033[1;31mERROR:\033[0m"
	info_str := "\033[1;36mINFO:\033[0m"
	var wg sync.WaitGroup
	var currentStep uint = 1
	var added uint8
	ctx, cancel := context.WithCancel(context.Background())
	_, id := utils.GetTokenID(configPath)
	return_chan := make(chan bool)
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		<-signalCh
		fmt.Println()
		cancel()
		log.Println(info_str, "Waiting for Response...")
		go func() {
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			<-signalCh
			fmt.Println()
			os.Exit(1)
		}()
		return_chan <- true
	}()
outter:
	for {
		time.Sleep(500 * time.Millisecond)
		wg.Add(1)
		select {
		case <-ctx.Done():
			wg.Done()
		default:
			go func(index uint) {
				var publicKey string
				currentStep++
				_, publicKey = utils.GenerateKey()

				installID := utils.RandStringRunes(22, nil)
				fcmtoken := utils.RandStringRunes(134, nil)
				r := utils.Request{
					Action: "register",
					Payload: []byte(
						`{
						"key":"` + publicKey + `",
						"install_id":"` + installID + `",
						"fcm_token":"` + installID + `:APA91b` + fcmtoken + `",
						"tos":"` + time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + `",
						"model":"Android",
						"referrer": "` + id + `",
						"serial_number":"` + installID + `"
					}`,
					),
					ID: id,
				}
				request, err := r.New()
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(info_str, `[`+"\033[1;36m"+strconv.FormatUint(uint64(index), 10)+"\033[0m"+`]`, "Sending request")
				if _, err = client.Do(request); err != nil {
					wg.Done()
					log.Println(error_str, `[`+"\033[1;31m"+strconv.FormatUint(uint64(index), 10)+"\033[0m"+`]`, err)
					log.Println(warn_str, "Waiting for 30 seconds...")
					select {
					case <-ctx.Done():
						return
					default:
					}
					cancel()
					time.Sleep(30 * time.Second)
					ctx, cancel = context.WithCancel(context.Background())
				} else {
					wg.Done()
					log.Println(info_str, `[`+"\033[1;33m"+strconv.FormatUint(uint64(index), 10)+"\033[0m"+`]`, "Request successed")
					added++
				}
			}(currentStep)
		}
		select {
		case <-return_chan:
			break outter
		default:
		}
	}
	wg.Wait()
	log.Println(info_str, "Total added "+strconv.FormatUint(uint64(added), 10)+" GB")
}
