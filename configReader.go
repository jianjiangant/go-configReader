package configReader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// config line type
const (
	AREA = iota
	FIELD
	COMMENT
	INVALIDLINE
)

type AreaType int

type DoubleMap map[string]map[string]interface{}

type ConfigReaderI interface {
	SetConfigPath(path string) error
	GetSection(sectionName string) (map[string]interface{}, error)
	GetField(sectionName, fieldName string) (string, error)
	Scanner() error
}

type ConfigReader struct {
	path    string
	confMap DoubleMap
}

func (this *ConfigReader) GetSection(sectionName string) (map[string]interface{}, error) {
	sectionMap, ok := this.confMap[sectionName]
	if !ok {
		return nil, fmt.Errorf("error : secion %s is not exsist\n", sectionName)
	}

	return sectionMap, nil
}

func (this *ConfigReader) GetField(sectionName, fieldName string) (string, error) {
	fieldMap, err := this.GetSection(sectionName)
	if err != nil {
		return "", err
	}
	field_v, ok := fieldMap[fieldName]
	if !ok {
		return "", fmt.Errorf("error : config value [%s](%s) is not exsist\n", sectionName, fieldName)
	}
	return field_v.(string), nil

}

func (this *ConfigReader) SetConfigPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("error: the path [%s] is not exsist\n", path)
	}
	this.path = path
	return nil
}

func (ConfigReader) analysisConfigLine(line string) (AreaType, interface{}) {
	line = strings.TrimSpace(line)
	r, _ := regexp.Compile(`^\[\w+\]$`)
	str_match := r.FindString(line)
	if str_match != "" {
		r, _ := regexp.Compile(`[a-zA-Z0-9]+`)
		areaName := r.FindString(str_match)
		if areaName == "" {
			return INVALIDLINE, nil
		}
		return AREA, areaName
	} else if match, _ := regexp.MatchString(`^#`, line); match == true {
		return COMMENT, line
	}

	config_slice := strings.SplitN(line, "=", 2)
	if len(config_slice) != 2 {
		return INVALIDLINE, line
	}
	//return FIELD, map[string]interface{}{strings.TrimSpace(config_slice[0]): strings.TrimSpace(config_slice[1])}
	return FIELD, []string{strings.TrimSpace(config_slice[0]), strings.TrimSpace(config_slice[1])}
}

func (this *ConfigReader) Scanner() error {
	var file *os.File
	var err error
	if file, err = os.Open(this.path); err != nil {
		return err
	}
	defer file.Close()
	var currentArea string = ""
	scanner := bufio.NewScanner(file)
	var field_map = make(map[string]interface{})
	for scanner.Scan() {
		// set DoubleMap
		areaType, data := this.analysisConfigLine(scanner.Text())
		switch areaType {
		case AREA:
			if len(field_map) != 0 && currentArea != "" {
				this.confMap[currentArea] = field_map
			}
			currentArea = data.(string)
			field_map = make(map[string]interface{})
		case FIELD:
			var field_v_slice []string = data.([]string)
			field_map[field_v_slice[0]] = field_v_slice[1]
		case INVALIDLINE, COMMENT:

		}

	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if len(field_map) != 0 && currentArea != "" {
		this.confMap[currentArea] = field_map
		field_map = make(map[string]interface{})
	}
	fmt.Println(this.confMap)
	return nil

}

var configInst ConfigReaderI = &ConfigReader{confMap: make(DoubleMap)}

func SetConfigPath(path string) error {
	return configInst.SetConfigPath(path)
}

func GetSection(sectionName string) (map[string]interface{}, error) {
	section, err := configInst.GetSection(sectionName)
	return section, err
}

func GetField(sectionName, fieldName string) (string, error) {
	return configInst.GetField(sectionName, fieldName)
}

func scanner() error {
	return configInst.Scanner()
}

func InitConfigReader(path string) error {
	if err := SetConfigPath(path); err != nil {
		return err
	}
	return scanner()
}
