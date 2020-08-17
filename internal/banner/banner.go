package banner

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
)

func StartUpBanner() {
	if configs.Conf.IsDebug {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		bannerFile, err := ioutil.ReadFile(filepath.Join(dir, "internal", "banner") + "/banner.txt")
		if err != nil {
			bannerFile, err = ioutil.ReadFile("/etc/banner/cybermedium.txt")
			if err != nil {
				log.Println(err)
				return
			}
		}
		logger.CmdInfo.Println(string(bannerFile))
	}
}
