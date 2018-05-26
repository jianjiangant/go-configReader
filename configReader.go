package configReader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"strconv"
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

var configInst ConfigReaderI = &ConfigReader{confMap: make(DoubleMap)}

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
	r, _ := regexp.Compile(`^\[.+\]$`)
	str_match := r.FindString(line)
	if str_match != "" {
		r, _ := regexp.Compile(`[\w\.\-_]+`)
		areaName := r.FindString(str_match)
		if areaName == "" {
			return INVALIDLINE, nil
		}
		return AREA, areaName
	// BUG : if the head of line has several blank or Tab , the comment will not be matched
	// [changed]
	//  old : ^# 
	//  new : see below
	//  test : NO
	} else if match, _ := regexp.MatchString(`^\s*#`, line); match == true {
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
			// set last field map to the area map when the find a area type .
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

// init configuration file 
func Init(path string) error{
	if err := SetConfigPath(path); err != nil {
		return err
	}
	return scanner()
}

func GetInt(s , f string) (v int, e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e!=nil{
		return strconv.Atoi(vstr)
	} else {
		return
	}
}

func GetInt32(s , f string){
}

func GetInt64(s , f string){
}

func GetFloat(s , f string){
}

func GetString(s , f string){
}
