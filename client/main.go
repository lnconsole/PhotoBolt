package main

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lnconsole/photobolt/api/background"
	"github.com/lnconsole/photobolt/api/icon"
	"github.com/lnconsole/photobolt/api/overlay"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/http"
	"github.com/lnconsole/photobolt/service/ln"
	"github.com/lnconsole/photobolt/shared"
)

const (
	// CONST
	serverUrl             = "http://127.0.0.1:8081"
	replacedBg_outputFile = "replaced-bg.png"
	icon_outputFile       = "icon.png"
	combined_outputFile   = "combined.png"
)

func main() {
	if err := env.Init("../env/.env"); err != nil {
		log.Printf("env err: %s", err)
		return
	}

	if err := ln.Init(
		env.PhotoBolt.LNDMacaroonHex,
		env.PhotoBolt.LNDCertPath,
		env.PhotoBolt.LNDGrpcAddr,
		env.PhotoBolt.LndClientNetwork(),
		env.PhotoBolt.LnNetwork(),
	); err != nil {
		log.Printf("lnd init: %s", err)
		return
	}

	var (
		replacedBg_inputFile   string // /Users/ChongjinChua/Downloads/whiskey.png
		replacedBg_inputPrompt string // A whiskey bottle at a beach, surrounded by waves
		icon_inputFile         string // /Users/ChongjinChua/Downloads/ibex.png
		icon_inputPrompt       string // ibex with a whiskey glass
		budgetSat              int
		costMsat               int
		scanner                = bufio.NewScanner(os.Stdin)
	)
	fmt.Println("\nWelcome! I am a Designer Agent. I can generate a poster for your product")
	fmt.Println("\nProvide the path to your product image (png):")
	if scanner.Scan() {
		replacedBg_inputFile = scanner.Text()
	}

	fmt.Println("\nDescribe how the final poster should look like:")
	if scanner.Scan() {
		replacedBg_inputPrompt = scanner.Text()
	}

	fmt.Println("\nProvide the path to your company logo (png):")
	if scanner.Scan() {
		icon_inputFile = scanner.Text()
	}

	fmt.Println("\nDescribe how your company logo should look like:")
	if scanner.Scan() {
		icon_inputPrompt = scanner.Text()
	}

	fmt.Println("\nProvide the budget for this task (sats):")
	if scanner.Scan() {
		tmp := scanner.Text()
		var err error
		budgetSat, err = strconv.Atoi(tmp)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	rbL402, err := replaceBackground(replacedBg_inputFile, replacedBg_inputPrompt, "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	giL402, err := generateIcon(icon_inputFile, icon_inputPrompt, "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	ciL402, err := combineImages("")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if rbL402 != "" {
		l402Data, err := ln.ParseL402(rbL402)
		if err != nil {
			log.Fatalf(err.Error())
		}
		costMsat += l402Data.AmountMsat
	}
	if giL402 != "" {
		l402Data, err := ln.ParseL402(giL402)
		if err != nil {
			log.Fatalf(err.Error())
		}
		costMsat += l402Data.AmountMsat
	}
	if ciL402 != "" {
		l402Data, err := ln.ParseL402(ciL402)
		if err != nil {
			log.Fatalf(err.Error())
		}
		costMsat += l402Data.AmountMsat
	}

	if costMsat > (budgetSat * 1000) {
		log.Fatalf("\ncost(%d sats) is more than budget(%d sats). Will not proceed.", costMsat/1000, budgetSat)
	}

	fmt.Printf("\ncost(%d sats) is within budget(%d sats). Proceeding.\n", costMsat/1000, budgetSat)
	// replace background
	l402Data, err := ln.HandleL402(rbL402)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("\nLSAT Auth token obtained for 'replace background' endpoint: (%s)\n", fmt.Sprintf("%s...%s", l402Data.AuthToken[:15], l402Data.AuthToken[len(l402Data.AuthToken)-10:]))
	stopch := make(chan struct{})
	go printJobStatus("replace background", stopch)

	_, err = replaceBackground(replacedBg_inputFile, replacedBg_inputPrompt, l402Data.AuthToken)
	if err != nil {
		log.Fatalf(err.Error())
	}
	stopch <- struct{}{}

	// generate icon
	l402Data, err = ln.HandleL402(giL402)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("\nLSAT Auth token obtained for 'generate logo' endpoint: (%s)\n", fmt.Sprintf("%s...%s", l402Data.AuthToken[:15], l402Data.AuthToken[len(l402Data.AuthToken)-10:]))
	stopch = make(chan struct{})
	go printJobStatus("generate logo", stopch)

	_, err = generateIcon(icon_inputFile, icon_inputPrompt, l402Data.AuthToken)
	if err != nil {
		log.Fatalf(err.Error())
	}
	stopch <- struct{}{}

	// combine images
	l402Data, err = ln.HandleL402(ciL402)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("\nLSAT Auth token obtained for 'combine images' endpoint: (%s)\n", fmt.Sprintf("%s...%s", l402Data.AuthToken[:15], l402Data.AuthToken[len(l402Data.AuthToken)-10:]))
	stopch = make(chan struct{})
	go printJobStatus("combine images", stopch)

	_, err = combineImages(l402Data.AuthToken)
	if err != nil {
		log.Fatalf(err.Error())
	}
	stopch <- struct{}{}

	fmt.Printf("\nreplaced background png generated: %s/client/%s\n", env.PhotoBolt.RepoDirectory, replacedBg_outputFile)
	fmt.Printf("generated logo png generated: %s/client/%s\n", env.PhotoBolt.RepoDirectory, icon_outputFile)
	fmt.Printf("combined image png generated: %s/client/%s\n", env.PhotoBolt.RepoDirectory, combined_outputFile)
}

func replaceBackground(inputFilePath string, inputPrompt string, authToken string) (string, error) {
	var (
		output = &background.ReplaceBackgroundResponse{}
	)
	// pass token
	l402, err := http.PostForm(
		fmt.Sprintf("%s/api/background", serverUrl),
		authToken,
		map[string]string{
			"file":   inputFilePath,
			"prompt": inputPrompt,
		},
		output,
	)
	if err != nil {
		return "", fmt.Errorf("postForm: %s", err.Error())
	}

	if l402 != "" {
		return l402, nil
	}

	if len(output.Image) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	outputImage := output.Image[0]
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return "", fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(replacedBg_outputFile)
	png.Encode(f, img)
	if err != nil {
		return "", fmt.Errorf("create: %s", err.Error())
	}

	return "", nil
}

func generateIcon(inputFilePath string, inputPrompt string, authToken string) (string, error) {
	var (
		output = &icon.GenerateIconResponse{}
	)

	l402, err := http.PostForm(
		fmt.Sprintf("%s/api/icon", serverUrl),
		authToken,
		map[string]string{
			"file":   inputFilePath,
			"prompt": inputPrompt,
		},
		output,
	)
	if err != nil {
		return "", fmt.Errorf("postForm: %s", err.Error())
	}

	if l402 != "" {
		return l402, nil
	}

	if len(output.Image) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	outputImage := output.Image
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return "", fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(icon_outputFile)
	png.Encode(f, img)
	if err != nil {
		return "", fmt.Errorf("create: %s", err.Error())
	}

	return "", nil
}

func combineImages(authToken string) (string, error) {
	var (
		frontFilePath = fmt.Sprintf("%s/client/%s", env.PhotoBolt.RepoDirectory, icon_outputFile)
		backFilePath  = fmt.Sprintf("%s/client/%s", env.PhotoBolt.RepoDirectory, replacedBg_outputFile)
		output        = &overlay.CombineImagesResponse{}
	)

	l402, err := http.PostForm(
		fmt.Sprintf("%s/api/overlay", serverUrl),
		authToken,
		map[string]string{
			"front": frontFilePath,
			"back":  backFilePath,
		},
		output,
	)
	if err != nil {
		return "", fmt.Errorf("postForm: %s", err.Error())
	}

	if l402 != "" {
		return l402, nil
	}

	outputImage := output.Image
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return "", fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(combined_outputFile)
	png.Encode(f, img)
	if err != nil {
		return "", fmt.Errorf("create: %s", err.Error())
	}

	return "", nil
}

func printJobStatus(jobName string, stopCh chan struct{}) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	fmt.Printf("waiting for '%s' job to complete....\n", jobName)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("waiting for '%s' job to complete....\n", jobName)
		case <-stopCh:
			return
		}
	}
}
