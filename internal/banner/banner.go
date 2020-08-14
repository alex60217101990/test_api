package banner

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alex60217101990/test_api/internal/logger"
)

func StartUpBanner() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	bannerFile, err := ioutil.ReadFile(filepath.Join(dir, "internal", "banner") + "/banner.txt")
	if err != nil {
		log.Println(err)
	}
	logger.CmdInfo.Println(string(bannerFile))
}
