package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const kConfigDirName = ".yunfs"
const kConfigFileName = "yunfs.json"

var ErrorNoConfig error = errors.New("no config file found")

type OssConfig struct {
	Base string `json:"base"`
	Type string `json:"type"`
}

type EncryptConfig struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type FileSystemConfig struct {
	BlockSize int    `json:"block_size"`
	CacheSize int    `json:"cache_size"`
	Prefix    string `json:"meta_prefix"`
}

type Config struct {
	Oss OssConfig        `json:"oss"`
	Enc EncryptConfig    `json:"encrypt"`
	Fs  FileSystemConfig `json:"file_system"`
}

func DefaultConfig() *Config {
	config := &Config{}
	config.Oss.Type = "local"
	config.Oss.Base = ""

	config.Enc.Type = "aes"
	config.Enc.Key = "abcdefg"

	config.Fs.BlockSize = 1024 * 16
	config.Fs.Prefix = "y0k99t"
	return config
}

func NewConfig() (*Config, error) {
	home := os.Getenv("HOME")
	defaultConfig := DefaultConfig()
	if home == "" {
		return defaultConfig, ErrorNoConfig
	}

	dirName := home + "/" + kConfigDirName
	if _, err := os.Stat(dirName); err != nil {
		os.Mkdir(dirName, 0755)
	}

	fileName := dirName + "/" + kConfigFileName
	if _, err := os.Stat(fileName); err != nil {
		fp, _ := os.Create(fileName)
		defer fp.Close()

		if buff, err1 := json.MarshalIndent(defaultConfig, "", "    "); err1 != nil {
			panic(err1.Error())
		} else {
			fp.Write(buff)
		}
	} else {
		fp, _ := os.Open(fileName)
		defer fp.Close()
		if buff, err1 := ioutil.ReadAll(fp); err1 != nil {
			panic(err1.Error())
		} else {
			config := &Config{}
			if err2 := json.Unmarshal(buff, config); err2 != nil {
				panic(err2.Error())
			}
			return config, nil
		}
	}
	return defaultConfig, ErrorNoConfig
}
