package ziphelper

import (
	"log"
	"testing"
)

func TestZip(t *testing.T) {

	zipHelper := NewZipHelper()
	// zipHelper.Add("1.bin", []byte("1test2test3test4test"))
	// zipHelper.Add("2.json", []byte("I have nothing in the world"))
	// data, err := zipHelper.Compress()
	// if nil != err {

	// 	log.Println("zip err:", err)
	// 	return
	// }

	// ioutil.WriteFile("1.zip", data, 0644)

	// files, _ := zipHelper.Uncompress(data)

	// for _, v := range files {

	// 	log.Println(string(v))
	// }

	zipHelper.AddDir("D:\\体温")
	data, err := zipHelper.Compress()
	log.Println("data：", data)
	log.Println("err", err)
}
