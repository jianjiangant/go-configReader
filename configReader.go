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
	BLANK
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
	if len(line) == 0 {
		return BLANK , nil
	}

	if match, _ := regexp.MatchString(`^#[.]*`, line); match == true {
		return COMMENT, line
	}
//	r, _ := regexp.Compile(`^\[[\w\s]*[\d\w\-_]*\][\s]*$`)
	r, _ := regexp.Compile(`^\[[\s]*[a-zA-Z]+[\w-\.]*[\s]*\]$`)
	strMatch := r.FindString(line)
	if strMatch != "" {
		strMatch = strings.TrimSpace(strMatch)
		r, _ := regexp.Compile(`[\w-]+`)
		areaName := r.FindString(strMatch)
		if areaName == "" {
			return INVALIDLINE, nil
		}
		return AREA, areaName
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
		//fmt.Println(scanner.Text())
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
		case INVALIDLINE :
			panic("the word setting is not illegal , words=" + scanner.Text())
		}

	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if len(field_map) != 0 && currentArea != "" {
		this.confMap[currentArea] = field_map
		field_map = make(map[string]interface{})
	}
	//fmt.Println(this.confMap)
	return nil

}

// ------------ default Instances ---------------

// init configuration file 
func Init(path string) error{
	if err := SetConfigPath(path); err != nil {
		return err
	}
	return scanner()
}

func GetInt(s , f string) (v int, e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e ==nil{
		return strconv.Atoi(vstr)
	}
	return 
}

func GetInt32(s , f string) (v int32 , e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e == nil{
		vInt , e :=strconv.Atoi(vstr)
		return  int32(vInt) , e
	}
	return 
}

func GetInt64(s , f string) (v int64 , e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e == nil{
		vInt , e :=strconv.Atoi(vstr)
		return  int64(vInt) , e
	}
	return
}

func GetFloat32(s , f string) (v float32 , e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e == nil{
		vInt , e :=strconv.ParseFloat(vstr , 32)
		return  float32(vInt) , e
	}
	return

}

func GetFloat64(s , f string) (v float64 , e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e == nil{
		vInt , e :=strconv.ParseFloat(vstr , 64)
		return  float64(vInt) , e
	}
	return

}

func GetBool(s , f string) (v bool, e error){
	var vstr string
	if vstr ,  e = GetField(s ,f) ; e == nil{
		vBool , e :=strconv.ParseBool(vstr)
		return  vBool , e
	}
	return

}

func GetString(s , f string) (v string , e error){
	return GetField(s ,f)
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


