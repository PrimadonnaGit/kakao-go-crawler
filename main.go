package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/tebeka/selenium"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func sleepSecond(second time.Duration) {
	time.Sleep(second * time.Second)
}

func loopPlaceElements(placeItems []webdriver.WebElement) {
	file, err := os.OpenFile("./output.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkErr(err)

	wr := csv.NewWriter(bufio.NewWriter(file))

	for _, placeItem := range placeItems {
		placeTitleElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".head_item .tit_name .link_name")
		checkErr(err)
		placeSubCategoryElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".head_item .subcategory")
		checkErr(err)

		placeScoreElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating .score .num")
		checkErr(err)

		placeScoreCountElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating .score a")
		checkErr(err)

		placeReviewCountElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating a em")
		checkErr(err)

		placeAddressElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".info_item .addr p[data-id='address']")
		checkErr(err)

		placeTitle, _ := placeTitleElement.Text()
		placeSubCategory, _ := placeSubCategoryElement.Text()
		placeScore, _ := placeScoreElement.Text()
		placeScoreCount, _ := placeScoreCountElement.Text()
		placeReviewCount, _ := placeReviewCountElement.Text()
		placeAddress, _ := placeAddressElement.Text()

		wr.Write([]string{placeTitle, placeSubCategory, placeScore, placeScoreCount, placeReviewCount, placeAddress})
	}
	wr.Flush()

}

func main() {

	const (
		seleniumPath  = "./chromedriver.exe"
		searchKeyword = "서울 카페"
		searchURL     = "https://map.kakao.com/"
		EnterKey      = string('\ue007')
	)

	chromeDriver := webdriver.NewChromeDriver(seleniumPath)
	defer chromeDriver.Stop()
	err := chromeDriver.Start()
	checkErr(err)

	desired := webdriver.Capabilities{"Platform": "Windows"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	defer session.Delete()
	checkErr(err)

	err = session.Url(searchURL)
	checkErr(err)

	// 검색 키워드 입력
	keywordInput, _ := session.FindElement(selenium.ByCSSSelector, ".box_searchbar > input.query")
	err = keywordInput.SendKeys(searchKeyword)
	checkErr(err)

	err = keywordInput.SendKeys(selenium.EnterKey) // Enter key
	checkErr(err)

	sleepSecond(1)

	// 더보기
	moreBtn, _ := session.FindElement(selenium.ByCSSSelector, ".places > .more")
	err = moreBtn.SendKeys(selenium.EnterKey)
	checkErr(err)

	sleepSecond(1)

	pageBtns, _ := session.FindElements(selenium.ByCSSSelector, ".keywordSearch .pages .pageWrap a")

	// 페이지 순회
	n := 0
	for n < 1000 {
		for _, pageBtn := range pageBtns {
			pageBtn.SendKeys(selenium.EnterKey)
			sleepSecond(1)
			placeItems, _ := session.FindElements(selenium.ByCSSSelector, ".PlaceItem")

			loopPlaceElements(placeItems)
		}

		nextBtn, _ := session.FindElement(selenium.ByCSSSelector, ".keywordSearch .pages .pageWrap .next")
		err = nextBtn.SendKeys(selenium.EnterKey)

		checkErr(err)
		n++
	}

}
