package configReader_test

import (
	"fmt"
	"go-configReader"
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
