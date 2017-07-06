package helper

import (
	"strings"
	"os"
	"bufio"
	"io"
)

var appConfig *Config

func GetAppProperties() *Config {

	if appConfig == nil {
		path := "conf/app.properties"
		appConfig = &Config{}
		appConfig.InitConfig(path)
	}
	return  appConfig
}

type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet  + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) Read( key string) string {
	key = c.GetModal()  + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}

func (c Config) GetModal() string {
	key := "main"  + "modal.type"
	v, found := c.Mymap[key]
	if !found {
		return "default"
	}
	return v
}
