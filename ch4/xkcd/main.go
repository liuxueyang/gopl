package main

import (
	"fmt"
	"net/http"
	"log"
)

var URL = "https://xkcd.com/"

var _ = `
// 20200908214535
// https://xkcd.com/571/info.0.json

{
  "month": "4",
  "num": 571,
  "link": "",
  "year": "2009",
  "news": "",
  "safe_title": "Can't Sleep",
  "transcript": "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}",
  "alt": "If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.",
  "img": "https://imgs.xkcd.com/comics/cant_sleep.png",
  "title": "Can't Sleep",
  "day": "20"
}
`

type Xkcd struct {
	Num        int32
	Month      string
	Year       string
	Link       string
	News       string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

func fetch(num int32) {
	url := fmt.
		Sprintf("%s%d/info.0.json", URL, num)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
}

func main() {

}
