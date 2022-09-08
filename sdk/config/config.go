package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/jeffotoni/gcolor"
	"gopkg.in/yaml.v2"
)

type C struct {
	Cluster struct {
		NameSpace []string `yaml:"namespace,flow"`
	}
}

var (
	c    *C
	once sync.Once
)

func Config() *C {
	once.Do(func() {
		if c == nil {
			data, err := ioutil.ReadFile("./config.yaml")
			if err != nil {
				fmt.Println(gcolor.RedCor("..........................................."))
				fmt.Println(gcolor.RedCor("The config.yaml file needs to be in the root."))
				fmt.Println(gcolor.RedCor("..........................................."))
				os.Exit(0)
			}

			err = yaml.Unmarshal(data, &c)
			if err != nil {
				fmt.Println(gcolor.RedCor("..........................................."))
				fmt.Println(gcolor.RedCor("Error parsing config.yaml"))
				fmt.Println(gcolor.RedCor("..........................................."))
				os.Exit(0)
			}
		}
	})

	return c
}
