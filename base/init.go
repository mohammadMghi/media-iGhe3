package base

import "os"

func Initialize(config *Config) {
	CurrentConfig = config
	CurrentConfig.Initialize()
	for _, dir := range []string{CurrentConfig.MediaDirectoryPath, CurrentConfig.ImageDirectoryPath} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}
