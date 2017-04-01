package configReader_test

import (
	"configReader"
	"fmt"
	"testing"
)

func TestConfigReader(t *testing.T) {
	configReader.InitConfigReader("./testConfig.cfg")
	mapSection, err := configReader.GetSection("redis")
	if err != nil {
		panic(err)
	}
	fmt.Println(mapSection)
}
