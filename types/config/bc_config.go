package config

import (
	klaatoo "github.com/hessegg/nikto/types/config"
)

func SealConfig() {
	config := klaatoo.GetConfig()
	config.Seal()
}
