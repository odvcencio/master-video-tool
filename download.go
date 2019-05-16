package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	var assetID, downloadLink, fileName string

	flag.StringVar(&assetID, "AssetID", "", "the asset ID we want the master of")
	flag.StringVar(&downloadLink, "DownloadLink", "", "the link to download the master")
	flag.StringVar(&fileName, "Name", "", "the name of the video file we are saving")
	flag.Parse()

	if assetID != "" {
		if err := enableMasterAccess(assetID); err != nil {
			log.Println("Error on asset ID: ", err)
		}
	} else {
		log.Println("No asset ID passed")
	}

	if downloadLink != "" && fileName != "" {
		if err := downloadMasterCopy(downloadLink, fileName); err != nil {
			log.Println("Error on asset ID: ", err)
		}
	} else {
		log.Println("No link or file name was passed.")
	}
}

func downloadMasterCopy(link, fileName string) error {
	// Get the data
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	path := fmt.Sprintf("%s.mp4", fileName)

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func enableMasterAccess(assetID string) error {
	url := fmt.Sprintf("https://api.mux.com/video/v1/assets/%s/master-access", assetID)

	m := echo.Map{
		"master_access": "temporary",
	}
	jsonBody, _ := json.Marshal(m)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("538b240d-f948-4fcf-bcd9-a8f01fcc5052", "Ts90G4LOAQCiYf7+bCgIOovSppNxUB5CF2iT0iTE1dZRSSRQ42ELQEBCMRaYHlUggXQSColRrtE")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseMap := echo.Map{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return err
	}

	log.Println(responseMap)

	return nil
}
