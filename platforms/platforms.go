package platforms

import (
	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/platforms/atcoder"
	"github.com/mxwell/wac/platforms/codeforces"
)

var platformsList []model.Platform

func initPlatforms() []model.Platform {
	if len(platformsList) == 0 {
		platformsList = append(platformsList, atcoder.InitAtCoder())
		platformsList = append(platformsList, codeforces.InitCodeforces())
	}
	return platformsList
}

func FindPlatform(url string) model.Platform {
	for _, platform := range initPlatforms() {
		if platform.ValidUrl(url) {
			return platform
		}
	}
	return nil
}
