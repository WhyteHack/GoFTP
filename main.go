package main

import (
	"log"
	"os"
	"path"
	"sync"

	. "github.com/whytehack/goftp/pkg/constants"
	. "github.com/whytehack/goftp/pkg/file"
	"github.com/whytehack/goftp/pkg/goftp"
)

func areFilesCorrect(c *goftp.SSFTP) {
	remoteFileList := c.GetRemoteFileList("/")
	localFileList := GetLocalFiles("C:\\Users\\musta\\Desktop\\Files")
	for element, _ := range remoteFileList {
		if localFileList[path.Base(element)] != remoteFileList[element] {
			log.Printf(FAIL+"%s could not be downloaded properly.", element)
		}
	}

	if len(localFileList) == 0 {
		log.Printf("Local path is empty! Check your local directory.")
	}
}

func main() {
	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf(ERROR+"error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	var wg sync.WaitGroup

	client, err := goftp.New(
		"anonymous",
		"anonymous",
		"speedtest.tele2.net",
	)
	if err != nil {
		log.Println("Client could not be initialize")
	}
	defer client.Close()

	remoteFileList := client.GetRemoteFileList("/")
	if len(remoteFileList) == 0 {
		log.Printf(FAIL + "Remote path is empty! Exiting the program")
		return
	}

	localFileList := GetLocalFiles("C:\\Users\\musta\\Desktop\\Files")

	for element, _ := range remoteFileList {
		if localFileList[path.Base(element)] != remoteFileList[element] {
			wg.Add(1)
			go client.Copy(element, "C:\\Users\\musta\\Desktop\\Files", &wg)
		}
	}
	wg.Wait()

	areFilesCorrect(client)

}
