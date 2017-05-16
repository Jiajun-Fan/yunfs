package main

import (
	"encoding/json"
	"errors"
	"github.com/segmentio/go-prompt"
	"io/ioutil"
	"os"
)

const kConfigDirName = ".yunfs"
const kConfigFileName = "yunfs.json"

var ErrorNoConfig error = errors.New("no config file found")

type OssConfig struct {
	Type     string `json:"type"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Bucket   string `json:"bucket"`
	EndPoint string `json:"end_point"`
}

type EncryptConfig struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type FileSystemConfig struct {
	BlockSize  int    `json:"block_size"`
	CacheSize  int    `json:"cache_size"`
	Prefix     string `json:"meta_prefix"`
	MountPoint string `json:"mount_point"`
}

type Config struct {
	Oss   OssConfig        `json:"oss"`
	Cache OssConfig        `json:"cache"`
	Enc   EncryptConfig    `json:"encrypt"`
	Fs    FileSystemConfig `json:"file_system"`
}

func NewConfig() *Config {

	fileName, err := checkConfigFile()

	config := &Config{}
	if err != nil {
		fp, _ := os.Create(fileName)
		defer fp.Close()
		configWizard(config)
		if buff, err1 := json.MarshalIndent(config, "", "    "); err1 != nil {
			Fatal(err1.Error())
		} else {
			fp.Write(buff)
		}
	} else {
		fp, _ := os.Open(fileName)
		defer fp.Close()
		if buff, err1 := ioutil.ReadAll(fp); err1 != nil {
			Fatal(err1.Error())
		} else {
			if err2 := json.Unmarshal(buff, config); err2 != nil {
				Fatal(err2.Error())
			}
		}
	}
	return config
}

func configWizardOss(config *Config) {
	oss_types := []string{"local", "aliyun", "s3"}
	index := prompt.Choose("please choose storage type", oss_types)
	config.Oss.Type = oss_types[index]

	if config.Oss.Type != "local" {
		config.Oss.Key = prompt.StringRequired("please specify APP key")
		config.Oss.Secret = prompt.StringRequired("please specify APP Secret")
		config.Oss.Bucket = prompt.StringRequired("please specify storage bucket")
		config.Oss.EndPoint = prompt.StringRequired("please specify storage endpoint")
	} else {
		config.Oss.Bucket = prompt.StringRequired("please specify storage dir")
	}
}

func configWizard(config *Config) {
	configWizardOss(config)
}

func checkConfigDir() (dirName string) {
	home := os.Getenv("HOME")
	if home == "" {
		home = "~"
	}

	dirName = home + "/" + kConfigDirName
	if _, err := os.Stat(dirName); err != nil {
		os.Mkdir(dirName, 0755)
	}

	return
}

func checkConfigFile() (fileName string, err error) {
	fileName = checkConfigDir() + "/" + kConfigFileName
	_, err = os.Stat(fileName)
	return
}
