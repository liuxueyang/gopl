package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Comic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Transcript string `json:"transcript"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Link       string `json:"link"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
}

func main() {
	xkcdURL := "https://xkcd.com/"

	for i := range 3104 {
		item := i + 1

		if item == 404 {
			continue // skip the 404 comic, it's a joke comic
		}

		url := xkcdURL + fmt.Sprintf("%d/info.0.json", item)
		fileName := fmt.Sprintf("comics/%d.json", item)

		// check if the directory exists, if not create it
		if err := os.MkdirAll("comics", os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating comics directory: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(fileName); err == nil {
			fmt.Printf("File %s already exists, skipping download.\n", fileName)

			// read the existing comic file
			file, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading existing comic file: %v\n", err)
				os.Exit(1)
			}
			var comic Comic
			if err := json.Unmarshal(file, &comic); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling existing comic file: %v\n", err)
				os.Exit(1)
			}

			// download the image if it doesn't exist
			imageName := getImageName(&comic)

			if len(imageName) == 0 {
				fmt.Fprintf(os.Stderr, "Error: Image name is empty for comic %d\n", comic.Num)
				os.Exit(1)
			}

			if err := os.MkdirAll("images", os.ModePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating images directory: %v\n", err)
				os.Exit(1)
			}

			if _, err := os.Stat(imageName); os.IsNotExist(err) {
				if err := downloadImage(&comic); err != nil {
					fmt.Fprintf(os.Stderr, "Error downloading image: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("Image %s already exists, skipping download.\n", imageName)
			}

			continue
		}

		comic, err := fetchComic(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comic: %v\n", err)
			os.Exit(1)
		}

		if err := saveComicToFile(comic, fileName); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving comic to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Comic %d saved to %s\n", comic.Num, fileName)
	}
}

func fetchComic(url string) (*Comic, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var comic Comic
	if err := json.Unmarshal(b, &comic); err != nil {
		// save the raw JSON to a file for debugging
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &comic, nil
}

func saveComicToFile(comic *Comic, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	defer file.Close()

	b, err := json.MarshalIndent(comic, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling comic to JSON: %w", err)
	}

	if _, err := file.Write(b); err != nil {
		return fmt.Errorf("error writing to file %s: %w", fileName, err)
	}

	return nil
}

func getImageName(comic *Comic) string {
	parts := strings.Split(comic.Img, "/")
	if len(parts) < 2 {
		return ""
	}
	return fmt.Sprintf("images/%s", parts[len(parts)-1])
}

func downloadImage(comic *Comic) error {
	resp, err := http.Get(comic.Img)
	if err != nil {
		return fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()

	fileName := getImageName(comic)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating image file %s: %w", fileName, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("error saving image to file %s: %w", fileName, err)
	}

	fmt.Printf("Image for comic %d saved to %s\n", comic.Num, fileName)
	return nil
}
