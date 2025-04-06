package helpers

import "os"

var bannerDir = "web/static/banners" // Папка для баннеров

func GetBanners() []string {
	files, err := os.ReadDir(bannerDir)
	if err != nil {
		return []string{}
	}

	var banners []string
	for _, file := range files {
		if !file.IsDir() {
			banners = append(banners, "static/banners/"+file.Name())
		}
	}
	return banners
}
