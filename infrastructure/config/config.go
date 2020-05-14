package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
)

// MainConfig is the main config file
type MainConfig struct {
	Server struct {
		Name string
		Port string
	}
	API     API
	Datadog DatadogConfig
}

type API struct {
	NormalPrefix   string
	DefaultTimeout int
}

type DatadogConfig struct {
	Endpoint string
}

func ReadConfig(cfg interface{}, service string, module string, rootServicePath ...string) interface{} {
	configPath := ""
	dir, _ := os.Getwd()

	// added for common realize
	if rootServicePath != nil {
		dir = rootServicePath[0]
	}
	dirNames := strings.Split(dir, "/")

	for _, dirName := range dirNames {
		configPath += dirName + "/"
		if dirName == service {
			break
		}
	}

	log.Println("configPath: ", configPath)
	environ := os.Getenv("ENTENV")
	if environ == "" {
		environ = "development"
	}

	configPath += ""
	ok := ReadModuleConfig(cfg, configPath, module)
	if !ok {
		configPathNew := fmt.Sprintf("/etc/%s", service)
		log.Printf("failed to read config for module:%s at Path:%s. So trying to read from %s for environ=%s\n", module, configPath, configPathNew, environ)
		configPath = configPathNew
		ok = ReadModuleConfig(cfg, configPath, module)
		if !ok {
			log.Fatalf("failed to read config for module:%s at Path:%s for environ=%s\n", module, configPath, environ)
		}
	}

	return cfg
}

func ReadModuleConfig(cfg interface{}, path string, module string) bool {
	environ := os.Getenv("ENTENV")
	if environ == "" {
		environ = "development"
	}

	fname := path + "/" + module + "." + environ + ".ini"

	fmt.Println("fname: ", fname)

	/* #nosec G304 */
	config, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Println(err)
		return false
	}
	err = gcfg.FatalOnly(gcfg.ReadStringInto(cfg, string(config)))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
