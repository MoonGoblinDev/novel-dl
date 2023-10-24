package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func main() {
	fmt.Println(`
                      _           _ _ 
 _ __   _____   _____| |       __| | |
| '_ \ / _ \ \ / / _ \ |_____ / _\` + " `" + ` | 
| | | | (_) \ V /  __/ |_____| (_| | |
|_| |_|\___/ \_/ \___|_|      \__,_|_|
                                      
novel-dl
https://github.com/yeorinhieut/novel-dl
`)

	url := input("다운로드할 소설의 회차 목록 URL을 입력하세요: ")
	lastChapter := inputInt("소설의 마지막 회차 번호를 입력하세요: ")
	startDownload := input("다운로드를 시작하시겠습니까? (y/n): ")
	if startDownload != "y" && startDownload != "Y" {
		fmt.Println("프로그램을 종료합니다.")
		return
	}
	numThreads := inputInt("다운로드할 스레드 수를 입력하세요 (일반적으로 1을 권장합니다): ")

	// 사용자 에이전트 설정
	userAgent := randomUserAgent()

	// URL에서 HTML 가져오기
	fmt.Println("페이지에서 HTML을 가져오는 중...")
	html, err := fetchHTML(url, userAgent)
	if err != nil {
		log.Fatalf("HTML 가져오기 실패: %v\n", err)
	}

	// HTML을 파싱하여 링크 추출
	links := extractLinks(html, lastChapter)

	// 추출된 링크를 파일에 저장
	err = saveLinksToFile(links, "links.txt")
	if err != nil {
		log.Fatalf("링크를 파일에 저장하는데 실패했습니다: %v\n", err)
	}

	fmt.Printf("%d개의 링크를 찾았습니다.\n", len(links))


	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numThreads)

	// output 폴더 생성
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", os.ModePerm)
	}

	for i, link := range links {
		wg.Add(1)
		go func(i int, link string) {
			defer wg.Done()
			semaphore <- struct{}{}
			fmt.Printf("다운로드 중... %d/%d\n", i+1, len(links))
			downloadNovel(link, userAgent, i)
			<-semaphore
		}(i, link)
	}

	wg.Wait()
	fmt.Println(color.GreenString("다운로드가 완료되었습니다."))
}

func input(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func inputInt(prompt string) int {
	fmt.Print(prompt)
	var num int
	_, err := fmt.Scanf("%d\n", &num)
	if err != nil {
		log.Fatalf("입력 오류: %v\n", err)
	}
	return num
}

func randomUserAgent() string {
	UAlist := []string{
		"Mozilla/5.0 (Linux; Android 4.4.1; SM-J200G Build/KTU84P) AppleWebKit/601.9 (KHTML, like Gecko) Chrome/54.0.2322.256 Mobile Safari/533.9",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_3_1) AppleWebKit/603.26 (KHTML, like Gecko) Chrome/48.0.1152.123 Safari/600",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 9_5_0) Gecko/20100101 Firefox/71.4",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_8_7; like Mac OS X) AppleWebKit/533.1 (KHTML, like Gecko) Chrome/53.0.1144.134 Mobile Safari/603.9",
	}

	rand.Seed(time.Now().UnixNano())
	return UAlist[rand.Intn(len(UAlist))]
}

func fetchHTML(url, userAgent string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("HTTP 요청 생성 중 에러가 발생했습니다: %v\n", err)
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP 요청 중 에러가 발생했습니다: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP 요청 실패. 응답 코드: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTP 처리 중 에러가 발생했습니다: %v\n", err)
		return "", err
	}

	return string(body), nil
}

func extractLinks(html string, lastChapter int) []string {
	links := []string{}
	for i := 1; i <= lastChapter; i++ {
		selector := fmt.Sprintf("#serial-move > div > ul > li:nth-child(%d) > div.wr-subject > a", i)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		link, _ := doc.Find(selector).Attr("href")
		if link != "" {
			links = append(links, link)
		}
	}
	return links
}

func saveLinksToFile(links []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, link := range links {
		_, err := file.WriteString(link + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadNovel(link, userAgent string, index int) {
	// 랜덤한 딜레이 생성 (200ms에서 2000ms 사이)
	delay := 200 + rand.Intn(1800)
	time.Sleep(time.Millisecond * time.Duration(delay))

	html, err := fetchHTML(link, userAgent)
	if err != nil {
		log.Printf("%d번째 링크에서 HTML을 가져오는 중 에러가 발생했습니다: %v\n", index, err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("%d번째 링크의 HTML을 파싱하는 중 에러가 발생했습니다: %v\n", index, err)
		return
	}

	title := doc.Find("#content_wrapper > div.page-title > span").Text()
	content, _ := doc.Find("#novel_content").Html()

	cleanedContent := cleanText(content)

	outputDir := "./output"
	outputFile := fmt.Sprintf("%s/%s.txt", outputDir, strings.ReplaceAll(title, " ", ""))

	err = saveNovelToFile(outputFile, cleanedContent)
	if err != nil {
		log.Printf("%d번째 링크에서 파일을 저장하는 중 에러가 발생했습니다: %v\n", index, err)
	} else {
		log.Printf("%d번째 링크의 다운로드가 완료되었습니다.\n", index)
	}
}

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "<div>", "")
	text = strings.ReplaceAll(text, "</div>", "")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")

	text = stripHTML(text)

	text = strings.TrimSpace(text)
	return text
}

func stripHTML(html string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(html, "")
}

func saveNovelToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
