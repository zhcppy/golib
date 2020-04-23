package util

import (
	"os"
	"os/user"
)

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

func RunEnv() string {
	env := os.Getenv("RUN_ENV")
	if env != EnvDev && env != EnvTest && env != EnvProd {
		return EnvDev
	}
	return env
}

func HomeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		if curU, err := user.Current(); err == nil {
			home = curU.HomeDir
		}
	}
	if home == "" {
		panic("can't get home dir")
	}
	return home
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
