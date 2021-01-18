package ziphelper

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

//ZipHelper Encapsulate compress and uncompress for zip
type ZipHelper struct {
	filelist map[string][]byte
}

//NewZipHelper create zip reader obj
func NewZipHelper() *ZipHelper {

	p := ZipHelper{}
	p.filelist = make(map[string][]byte)

	return &p
}

//Add add a file to zip package
func (z *ZipHelper) Add(fName string, data []byte) error {

	if 0 == len(fName) || nil == data || 0 == len(data) {

		return fmt.Errorf("Params error, file name or file body is nil")
	}

	z.filelist[fName] = data

	return nil
}

//Compress compress files to zip file
func (z *ZipHelper) Compress() ([]byte, error) {

	if 0 == len(z.filelist) {

		return nil, fmt.Errorf("no files for compress")
	}

	buf := bytes.Buffer{}
	zip := zip.NewWriter(&buf)

	for name, data := range z.filelist {

		context, err := zip.Create(name)
		if err != nil {

			return nil, err
		}

		context.Write(data)
	}

	zip.Close()

	return buf.Bytes(), nil
}

//Uncompress uncompress a zip file and return included files
func (z *ZipHelper) Uncompress(data []byte) (map[string][]byte, error) {

	if nil == data || 0 == len(data) {

		return nil, fmt.Errorf("uncompress failed, because the data is nil")
	}

	zipData := bytes.NewReader(data)
	zipReader, err := zip.NewReader(zipData, int64(len(data)))
	if nil != err {

		return nil, errors.New("zipReader fail.err:" + err.Error())
	}

	if 0 == len(zipReader.File) {

		return nil, errors.New("no files included in the zip")
	}

	files := make(map[string][]byte)
	for _, v := range zipReader.File {

		f, err := v.Open()
		if nil != err {

			log.Println("uncompress zip err: ", err)
			continue
		}
		defer f.Close()

		fdata, err := ioutil.ReadAll(f)
		if nil != err {

			log.Println("read zip file err:", err)
			continue
		}

		files[v.FileInfo().Name()] = fdata
	}

	return files, nil
}
