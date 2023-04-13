package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"github.com/jung-kurt/gofpdf"
)

type HQs struct {
	Links []string
	Names []string
}

type HQCaps struct {
	Links []string
	Caps  []string
}

func Search2HQ(query string) HQs {
	res, err := http.Get(fmt.Sprintf("https://hqdragon.com/pesquisa?nome_hq=%s", query))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln("Error loading HTML:", err)
	}

	titlesNode := doc.Find("div.lista-hqs")

	linksNode := doc.Find("div.lista-hqs")

	names := []string{}

	links := []string{}

	// fmt.Println(doc.Text())

	// Use XPath to select elements
	titlesNode.Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Find("a").Text())
	})

	linksNode.Each(func(i int, s *goquery.Selection) {

		link, status := s.Find("a").Attr("href")

		if !status {
			log.Fatalln("Erro")
		}

		links = append(links, link)
		//fmt.Println(s.Find("a").Attr("href"))
	})

	return HQs{
		Links: links,
		Names: names,
	}
}

func GetCaps(link string) HQCaps {
	hqIndexPage, err := http.Get(link)

	if err != nil {
		log.Fatal(err)
	}

	defer hqIndexPage.Body.Close()
	if hqIndexPage.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", hqIndexPage.StatusCode, hqIndexPage.Status)
	}

	caps := []string{}

	capsLinks := []string{}

	doc2, err := goquery.NewDocumentFromReader(hqIndexPage.Body)

	if err != nil {
		log.Fatal(err)
	}

	capsLinks = append(capsLinks, "All")

	caps = append(caps, "All")

	doc2.Find("table.table.table-bordered").Find("tbody").Find("tr").Find("td").Find("a").Each(func(i int, s *goquery.Selection) {
		caps = append(caps, s.Text())

		link, status := s.Attr("href")

		if !status {
			log.Fatalln("Erro")
		}

		capsLinks = append(capsLinks, link)

	})

	return HQCaps{
		Caps:  caps,
		Links: capsLinks,
	}

}

func DownloadHQ(capsLink []string, name string, label *widget.Label, path2Output string) {
	for idx, capLink := range capsLink {
		if capLink != "All" {
			hqCapPage, err := http.Get(capLink)

			capLinkArray := strings.Split(capLink, "/")

			capStr := capLinkArray[len(capLinkArray)-1]

			if err != nil {
				log.Fatal(err)
			}

			defer hqCapPage.Body.Close()

			if hqCapPage.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", hqCapPage.StatusCode, hqCapPage.Status)
			}

			doc3, err := goquery.NewDocumentFromReader(hqCapPage.Body)

			if err != nil {
				log.Fatal(err)
			}

			pags := doc3.Find("select#paginas").Find("option")

			pdf := gofpdf.New("P", "mm", "A4", "")

			pags.Each(func(i int, s *goquery.Selection) {
				pag := strings.Trim(strings.Replace(s.Text(), "Pag.", "", -1), " \t\n\r")
				pagInt, err := strconv.Atoi(pag)
				pagStr := pag
				if err != nil {
					log.Fatal(err)
				}

				if pagInt < 10 {
					pagStr = fmt.Sprintf("%s%d", "0", pagInt)
				}

				imgSrc, status := doc3.Find(fmt.Sprintf("img.pag_%s", pag)).First().Attr("src")

				if !status {
					log.Fatalln("Erro")
				}

				imgContent, err := http.Get(imgSrc)

				if err != nil {
					log.Fatal(err)
				}

				defer imgContent.Body.Close()

				if len(capsLink) > 1 {
					label.Text = fmt.Sprintf("Baixando %v de %v", idx, len(capsLink)-1)
					label.Refresh()
				} else {
					label.Text = fmt.Sprintf("Baixando %v de %v", idx+1, len(capsLink))
					label.Refresh()
				}

				imageFileName := fmt.Sprintf("image_%s.jpg", pagStr)

				file, err := os.Create(imageFileName)

				if err != nil {
					log.Fatal(err)
				}

				defer file.Close()

				_, err = io.Copy(file, imgContent.Body)

				if err != nil {
					log.Fatal(err)
				}

				pdf.AddPage()

				// Add the image to the PDF document
				pdf.Image(imageFileName, 0, 0, 210, 297, false, "", 0, "")

			})

			err = pdf.OutputFileAndClose(fmt.Sprintf("%s/%s - %s.pdf", path2Output, capStr, name))
			if err != nil {
				log.Fatal(err)
			}

			imagesFile, err := os.ReadDir(".")

			if err != nil {
				log.Fatal(err)
			}

			for _, image := range imagesFile {
				if strings.Contains(image.Name(), "image_") {
					err = os.Remove(image.Name())

					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}
}
