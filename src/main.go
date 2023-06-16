package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
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

func downloadImage(url string, savePath string, bar *pb.ProgressBar) error {
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

	bar.Increment()

	return nil
}

func sanitizeFilename(filename string) string {
	invalidChars := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	sanitized := filename
	for _, char := range invalidChars {
		sanitized = strings.ReplaceAll(sanitized, char, "!")
	}
	return sanitized
}

func extractImagesFromURL(urlStr string, savePath string) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("%sError while fetching webpage: %v%s", ColorRed, err, ColorReset)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("%sError while parsing HTML: %v%s", ColorRed, err, ColorReset)
	}

	title := doc.Find("title").Text()
	title = strings.TrimSpace(title)
	folderName := sanitizeFilename(title)
	folderPath := filepath.Join(savePath, folderName)

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%sError while creating folder: %v%s", ColorRed, err, ColorReset)
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if exists {
			imgURL, err := url.Parse(imgSrc)
			if err == nil && !imgURL.IsAbs() {
				baseURL, err := url.Parse(urlStr)
				if err == nil {
					imgURL = baseURL.ResolveReference(imgURL)
				}
			}
			images = append(images, imgURL.String())
		}
	})

	totalImages := len(images)

	bar := pb.StartNew(totalImages)
	bar.SetTemplateString(fmt.Sprintf("Downloading images from %s\n{{bar . }} {{counters . }}", urlStr))

	for _, img := range images {
		err := downloadImage(img, folderPath, bar)
		if err != nil {
			fmt.Printf("%sError while downloading image: %v%s\n", ColorRed, err, ColorReset)
		}
	}

	bar.Finish()

	fmt.Printf("%sDownloading completed from %s. Total images downloaded: %d%s\n", ColorGreen, urlStr, totalImages, ColorReset)
	return nil
}

func main() {
	banner()

	helpFlag := flag.Bool("h", false, "Display help")
	urlFlag := flag.String("u", "", "URL of the webpage")
	locationFlag := flag.String("l", "", "Location to save the extracted images")
	fileFlag := flag.String("f", "", "File containing URLs")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		fmt.Println("Example:")
		fmt.Println("go run src/main.go -u https://example.com -l /path/to/save/location")
		fmt.Println("or")
		fmt.Println("go run src/main.go -f urls.txt -l /path/to/save/location")
		return
	}

	if *fileFlag == "" {
		if *urlFlag == "" || *locationFlag == "" {
			fmt.Print("Enter the URL of the webpage: ")
			fmt.Scanln(urlFlag)

			fmt.Print("Enter the location to save the extracted images: ")
			fmt.Scanln(locationFlag)
		}

		err := extractImagesFromURL(*urlFlag, *locationFlag)
		if err != nil {
			fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		}
	} else {
		filePath := *fileFlag
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("%sError while opening file: %v%s\n", ColorRed, err, ColorReset)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urlStr := scanner.Text()
			err := extractImagesFromURL(urlStr, *locationFlag)
			if err != nil {
				fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
			}
		}

		if scanner.Err() != nil {
			fmt.Printf("%sError while reading file: %v%s\n", ColorRed, scanner.Err(), ColorReset)
		}
	}
}
