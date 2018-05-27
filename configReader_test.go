package configReader

import (
	"testing"
)

func assert(t *testing.T , refValue , testValue interface{} , info string){
	if refValue != testValue {
		t.Fatal(info)
	}
}

func TestGetString(t *testing.T) {
	Init("./testConfig.cfg")
	ip , e := GetString("common" , "ip")
	t.Logf("ip = %s \n", ip)
	if e != nil {
		//t.Error("configReader GetSection failed")
		t.Fatal("configReader GetSection failed")
	}
	assert(t , "127.0.0.1" , ip , "GetString test failed")
}

func TestGetInt(t *testing.T) {
	Init("./testConfig.cfg")
	port , e := GetInt("common" , "port")
	t.Logf("port = %d \n", port)
	if e != nil {
		t.Fatal("configReader GetSection failed")
	}
	assert(t , 8080 , port, "GetInt test failed")
}

func TestGetBool(t *testing.T) {
	Init("./testConfig.cfg")
	v, e := GetBool("common" , "enable")
	if e != nil {
		t.Fatal("configReader GetSection failed")
	}
	assert(t , true, v , "GetBool test failed")
}

func TestGetFloat32(t *testing.T) {
	Init("./testConfig.cfg")
	v, e := GetFloat32("common" , "time")
	t.Logf("time = %f \n", v)
	if e != nil {
		t.Fatal("configReader GetSection failed , error" , e)
	}
	assert(t , float32(11.19) , v , "GetFloat32 test failed")
}

func TestGetFloat64(t *testing.T) {
	Init("./testConfig.cfg")
	v, e := GetFloat64("common" , "time")
	t.Logf("time = %f \n", v)
	if e != nil {
		t.Fatal("configReader GetSection failed , error" , e)
	}
	assert(t , float64(11.19) , v , "GetFloat64 test failed")
}

func TestAnalysisConfigLine(t *testing.T) {
	type TestCase struct{
		LineType AreaType
		TestData string
	}

	testcase := []TestCase{
		TestCase{
			LineType : COMMENT,
			TestData: "#comment",
		},
		TestCase{
			LineType : COMMENT,
			TestData: "# comment",
		},
		TestCase{
			LineType : COMMENT,
			TestData: " #comment",
		},
		TestCase{
			LineType : COMMENT,
			TestData: "	#comment",
		},
		TestCase{
			LineType : COMMENT,
			TestData: "# comm ent",
		},
		TestCase{
			LineType : AREA,
			TestData: "[area]",
		},
		TestCase{
			LineType : AREA,
			TestData: " [area]",
		},
		TestCase{
			LineType : AREA,
			TestData: "	[area]",
		},
		TestCase{
			LineType : AREA,
			TestData: "[ area ]",
		},
		TestCase{
			LineType : AREA,
			TestData: "[area.test]",
		},
		TestCase{
			LineType : FIELD,
			TestData: "a=b",
		},
		TestCase{
			LineType : FIELD,
			TestData: " a = b ",
		},
		TestCase{
			LineType : FIELD,
			TestData: "	a = b ",
		},
		TestCase{
			LineType : INVALIDLINE,
			TestData: "invalid",
		},
		TestCase{
			LineType : INVALIDLINE,
			TestData: "	invalid",
		},
		TestCase{
			LineType : INVALIDLINE,
			TestData: "[invalid invalid]",
		},
	}
	cTestInstance := new(ConfigReader)
	for _ , v := range testcase {
		areaType , _ :=cTestInstance.analysisConfigLine(v.TestData)
		if areaType != v.LineType {
			t.Error("analysisConfigLine failed , testcase = " , v.TestData)
		}
	}
}
