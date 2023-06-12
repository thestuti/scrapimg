package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

func banner() {
	fmt.Printf(`
	
███████  ██████ ██████   █████  ██████  ██ ███    ███  ██████  
██      ██      ██   ██ ██   ██ ██   ██ ██ ████  ████ ██       
███████ ██      ██████  ███████ ██████  ██ ██ ████ ██ ██   ███ 
     ██ ██      ██   ██ ██   ██ ██      ██ ██  ██  ██ ██    ██ 
███████  ██████ ██   ██ ██   ██ ██      ██ ██      ██  ██████  
                                                               
                                                               
`)
}

func downloadImage(url string, savePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%sError while downloading image: %v%s", ColorRed, err, ColorReset)
	}
	defer response.Body.Close()

	fileName := filepath.Base(url)
	saveLocation := filepath.Join(savePath, fileName)

	file, err := os.Create(saveLocation)
	if err != nil {
		return fmt.Errorf("%sError while creating file: %v%s", ColorRed, err, ColorReset)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("%sError while saving image: %v%s", ColorRed, err, ColorReset)
	}

	fmt.Printf("%sDownloading completed%s\n", ColorGreen, ColorReset)
	return nil
}

func main() {
	banner()
	var webURL, saveLocation string

	helpFlag := flag.Bool("h", false, "Display help")
	urlFlag := flag.String("u", "url", "URL of the webpage")
	locationFlag := flag.String("l", "location", "Location to save the extracted images")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		fmt.Println("Example:")
		fmt.Println("go run src/main.go -u https://example.example.com -l /path/to/save/location")
		return
	}

	if *urlFlag == "" || *locationFlag == "" {
		fmt.Print("Enter the URL of the webpage:")
		fmt.Scanln(&webURL)

		fmt.Print("Enter the location to save the extracted images: ")
		fmt.Scanln(&saveLocation)
	} else {
		webURL = *urlFlag
		saveLocation = *locationFlag
	}

	resp, err := http.Get(webURL)
	if err != nil {
		fmt.Printf("%sError while fetching webpage: %v%s\n", ColorRed, err, ColorReset)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("%sError while parsing HTML: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if exists {

			imgURL, err := url.Parse(imgSrc)
			if err != nil {
				fmt.Printf("%sError while parsing image URL: %v%s\n", ColorRed, err, ColorReset)
				return
			}

			if !imgURL.IsAbs() {
				baseURL, err := url.Parse(webURL)
				if err != nil {
					fmt.Printf("%sError while parsing base URL: %v%s\n", ColorRed, err, ColorReset)
					return
				}
				imgURL = baseURL.ResolveReference(imgURL)
			}

			images = append(images, imgURL.String())
		}
	})

	for _, img := range images {
		err := downloadImage(img, saveLocation)
		if err != nil {
			fmt.Printf("%sError while downloading image: %v%s\n", ColorRed, err, ColorReset)
		}
	}
}
