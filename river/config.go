package river

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type SourceConfig struct {
	Schema string
	Tables []string
}

type Config struct {
	MyAddr     string
	MyUser     string
	MyPassword string

	ESAddr string

	StatAddr string

	ServerID uint32
	Flavor   string
	DataDir  string

	DumpExec string

	Sources []SourceConfig

	Rules []*Rule
}

func NewConfigWithFile(name string) (*Config, error) {

	fmt.Println("NewConfigWithFile")

	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return NewConfig(data)
}

func NewConfig(data []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("Value: %#v\n", config)

	return &config, nil
}
