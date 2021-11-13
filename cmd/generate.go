package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/neobaran/csac/config"
	"github.com/neobaran/csac/lets"
	"github.com/neobaran/csac/tencent"
	"gopkg.in/yaml.v2"
)

func Generate(configFile string, debug bool) {
	appConfig := &config.Config{}
	appConfig.TTL = 600

	if _, err := os.Stat(configFile); err == nil {
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(yamlFile, appConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		envNamespace := "CSAC_"
		appConfig.Email = os.Getenv(envNamespace + "EMAIL")
		appConfig.Tencent.SecretId = os.Getenv(envNamespace + "TENCENT_SECRET_ID")
		appConfig.Tencent.SecretKey = os.Getenv(envNamespace + "TENCENT_SECRET_KEY")
		appConfig.Domains = append(appConfig.Domains, os.Getenv(envNamespace+"DOMAIN"))
		if TTL, err := strconv.ParseUint(os.Getenv(envNamespace+"TTL"), 10, 64); err == nil {
			appConfig.TTL = TTL
		}
		appConfig.KeyType = os.Getenv(envNamespace + "KEY_TYPE")
	}

	if len(appConfig.Email) == 0 {
		log.Fatal("Email is missing")
	}
	if len(appConfig.Domains) == 0 {
		log.Fatal("Domains is missing")
	}
	if len(appConfig.Tencent.SecretId) == 0 || len(appConfig.Tencent.SecretKey) == 0 {
		log.Fatal("SecretId or SecretKey is missing")
	}
	if len(appConfig.KeyType) == 0 {
		appConfig.KeyType = "RSA2048"
	}
	if _, ok := lets.KeyTypes[strings.ToUpper(appConfig.KeyType)]; !ok {
		log.Fatal("Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384.(default: \"rsa2048\")")
	}

	cloudHelper := tencent.NewTencentCloudHelp(appConfig.Tencent)

	letsHelper, err := lets.NewCSACHelper(appConfig, cloudHelper, debug)
	if err != nil {
		log.Fatalln(err)
	}

	cert, err := letsHelper.CreateSSL(appConfig.Domains)
	if err != nil {
		log.Fatalln(err)
	}

	if debug {
		log.Println("success", cert.CertURL)
		return
	}

	if err := letsHelper.UploadToCloud(cert); err != nil {
		log.Fatalln(err)
	}
}
