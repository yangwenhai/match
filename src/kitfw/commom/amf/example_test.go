package amf

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

type MixedComplexObject struct {
	MapStringToInt       map[string]int
	MapStringToString    map[string]string
	MapStringToInterface map[string]interface{}
	MapStringToPointer   map[string]*string
	SliceOfInt           []int
	SliceOfInterface     []interface{}
	SliceOfPointer       []*string
	ArrayOfInt           [3]int
	ArrayOfInterface     [3]interface{}
	ArrayOfPointer       [3]*string
	Integer              int
	Uinteger             uint
	Float32              float32
	Float64              float64
	Bool                 bool
	Pointer              **string
	String               string
	dummy                int
}

func (obj *MixedComplexObject) GetFieldName(key string) string {
	if key == "mapStringToInt" {
		return "MapStringToInt"
	}
	return key
}

func (obj *MixedComplexObject) GetAmfName(key string) string {
	if key == "MapStringToInt" {
		return "mapStringToInt"
	}
	return key
}

/*
func TestEncode(t *testing.T) {
	ExampleEncode_decode(t)
}
*/

func TestEncodeMapInt(t *testing.T) {
	inputData := map[int]string{
		-123: "b",
		456:  "d",
	}

	t.Log(&inputData)
	buffer := bytes.NewBuffer(nil)

	err := Encode(buffer, binary.BigEndian, &inputData)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%X\n", buffer.Bytes())
	}

	err = Decode(buffer, binary.BigEndian, &inputData)
	fmt.Printf("%v\n", inputData)

}

func TestEncodeStruct(t *testing.T) {
	type ABC struct {
		AB map[string]string
	}

	//var inputData ABC

	//inputData.Uid = 230
	//inputData.Uname = "abcdefghi"
	//inputData.Abc = 10

	//str := "0a0b010541420a0b01054344060545460101"
	str := "0a0b010541420a0b01056344060565460101"

	var info ABC
	buf, _ := hex.DecodeString(str)
	buffer := bytes.NewBuffer(buf)
	Decode(buffer, binary.BigEndian, &info)

	println("token:")
	fmt.Println(info)
	/*
		err := Encode(buffer, binary.BigEndian, &inputData)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("struct: %X\n", buffer.Bytes())
		}
	*/

}

func ExampleEncode_decode(t *testing.T) {
	str := "hello, world"
	pstr := &str
	inputData := MixedComplexObject{
		MapStringToInt: map[string]int{
			"hello": 1234,
			"long":  123435678,
		},
		MapStringToString: map[string]string{
			"hell": "hello",
			"long": "long",
		},
		MapStringToInterface: map[string]interface{}{
			"hello": 1234,
			"long":  123435678,
		},
		MapStringToPointer: map[string]*string{
			"hello": &str,
		},
		SliceOfInt: []int{
			1, 2, 3, 4, -0x1fffff,
		},
		SliceOfInterface: []interface{}{
			1, 2, 3, 4, -0x1fffff,
		},
		SliceOfPointer: []*string{
			&str, &str, &str,
		},
		ArrayOfInt: [3]int{
			1, 2, -0x1ffffff,
		},
		ArrayOfInterface: [3]interface{}{
			1, 2, -0x1ffffff,
		},
		ArrayOfPointer: [3]*string{
			&str, &str, &str,
		},
		Integer:  12345678,
		Uinteger: 0x1234567,
		Float32:  3.1415926,
		Float64:  3.1415926,
		Bool:     true,
		Pointer:  &pstr,
		String:   str,
		dummy:    1234,
	}

	t.Log(&inputData)
	buffer := bytes.NewBuffer(nil)

	err := Encode(buffer, binary.BigEndian, &inputData)
	if err != nil {
		t.Error(err)
	} else {
		//fmt.Println(buffer.Bytes())
	}
	//fmt.Println(hex.EncodeToString(buffer.Bytes()))

	//var data MixedComplexObject
	//var data map[string]interface{}
	//var data interface{}
	err = Decode(buffer, binary.BigEndian, &inputData)
	if err != nil {
		t.Error(err)
	} else {
		//fmt.Println(data)
	}

	t.Log(&inputData)
}

func TestDecodeToken(t *testing.T) {
	//0a0b01077569640481a80309747970650605706b09726f6f6d060761626413626567696e54696d6504641164757261746f696e040a1163616c6c6261636b060561621963616c6c6261636b4172677606076173640f657874496e666f0a0b01056870040a011175736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101
	//str := "0a0b01077569640481a80309747970650605706b09726f6f6d060761626413626567696e54696d6504641164757261746f696e040a1163616c6c6261636b060561621963616c6c6261636b4172677606076173640f657874496e666f0a0b01056870040a011175736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101"
	//str := "0a0b01075569640481a80315426174746c65547970650605706b09526f6f6d060761626413426567696e54696d6504641144757261746f696e040a1143616c6c6261636b060561621943616c6c6261636b4172677606076173640f457874496e666f0a0b01056870040a011155736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101"
	// str := "0a0b01075569640481a80315426174746c65547970650605706b09526f6f6d060761626413426567696e54696d6504641144757261746f696e040a1143616c6c6261636b060563621943616c6c6261636b417267760a0b01100607636272010f457874496e666f0a0b01056870040a011155736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101"
	//buf := []byte(str)

	//	  str := "0a0b01075569640481a80315426174746c65547970650605706b09526f6f6d060761626413426567696e54696d650464114475726174696f6e04321143616c6c6261636b060563621943616c6c6261636b417267760a0b01100607636272010f457874496e666f0a0b01056870040a011155736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101"
	//str := "0a0b01075569640481a80315426174746c65547970650605706b09526f6f6d060761626413426567696e54696d650464114475726174696f6e04321143616c6c6261636b060563621943616c6c6261636b417267760a0b010543620607636272010f457874496e666f0a0b01054870040a011155736572496e666f0683550a0b01096865726f0a0b010b6c6576656c04010d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a8030b756e616d65060b3132336466116d6f64656c5f696404000101"
	//str := "0a0b0107556964046415426174746c655479706506097479706509526f6f6d0609726f6f6d13426567696e54696d650464114475726174696f6e040a1143616c6c6261636b061163616c6c6261636b1943616c6c6261636b417267760a0b01076362610607616264010f457874496e666f0a0b0105636406056566011155736572496e666f060761626301"
	//str := "0a0b011943616c6c6261636b417267760a0b01076362610607616264010f457874496e666f0a0b0105636406056566011155736572496e666f060761626301"
	//str := "0a0b011943616c6c6261636b417267760a0b01076362610607616264010f457874496e666f0a0b010563640a0b010561620605646505646609010101011155736572496e666f060761626301"
	str := "0a0b0109486173680641313561643535663238333638633034633237323937333439323530303730396615437265617465496e666f06885f0a0b01075569640481a87a15426174746c655479706506097465616d094e616d65064165353731393063653365306161363064663761663463333932383237666135660b47726f7570061367616d65313030343313426567696e54696d650541d4881e81800000114475726174696f6e043c1143616c6c6261636b0627636f70792e67726f7570426174746c655265731943616c6c6261636b4172677606821b613a343a7b733a373a22636f70795f6964223b693a313030313b733a31303a22626567696e5f74696d65223b693a313337373836303130323b733a343a2275696473223b613a313a7b693a303b693a32313632363b7d733a373a22726f6f6d5f6964223b733a33323a226535373139306365336530616136306466376166346333393238323766613566223b7d0f457874496e666f0901011155736572496e666f0683550a0b01096865726f0a0b010b6c6576656c040b0d61726d696e670a0b01053130090101033509010103330901010331090101033609010103370901010332090101033809010103390901010334090101011174616c69736d616e09010113617374726f6c61626504000d666f737465720a0b010d7374726f6e6704000b6d616769630400136465787465726974790400010b64656d6f6e0901010d736b696c6c730a0b0109313030360401093130303704010931303038040109313030390401093130313004010101077569640481a87a0b756e616d65060b6564666774116d6f64656c5f69640400010101"

	type CreateInfo struct {
		Uid          int
		BattleType   string
		Name         string
		Group        string
		BeginTime    int
		Duration     int
		Callback     string
		CallbackArgv string
		ExtInfo      map[string]interface{}
		UserInfo     string
	}

	type LoginInfo struct {
		Hash       string
		CreateInfo string
	}

	var loginInfo LoginInfo
	buf, _ := hex.DecodeString(str)
	buffer := bytes.NewBuffer(buf)
	err := Decode(buffer, binary.BigEndian, &loginInfo)
	if err != nil {
		fmt.Println(err)
	}

	var info CreateInfo
	buffer = bytes.NewBuffer([]byte(loginInfo.CreateInfo))
	err = Decode(buffer, binary.BigEndian, &info)
	fmt.Println(err)

	println("token:")
	fmt.Println(info)
	fmt.Println(info.Uid)

}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */