package upackage

import (
	"log"
	"testing"
)

func TestPackage(t *testing.T) {

	//build a package
	idpkg := NewBPackage()
	idpkg.AddPackageID(1)

	idpkg.AddInt8(8)
	idpkg.AddInt16(16)
	idpkg.AddInt32(32)
	idpkg.AddInt64(64)

	idpkg.AddUint8(18)
	idpkg.AddUint16(116)
	idpkg.AddUint32(132)
	idpkg.AddUint64(164)

	idpkg.AddBool(false)

	idpkg.AddFloat32(1.2345678, 7)
	idpkg.AddFloat64(1.234567890, 9)

	idpkg.AddString("not with length")
	idpkg.AddStringL("with length", 20)

	idpkg.AddInt64(6464)
	idpkg.Done()

	log.Println(string(idpkg.GetData()))

	//create a read package
	rpkg := NewBPackage()
	rpkg.AddBytes(idpkg.GetData())

	log.Println("mark: ", rpkg.ReadHeaderMark())
	log.Println("id: ", rpkg.ReadPackageID())
	log.Println("size: ", rpkg.ReadDataLength())

	log.Println("", rpkg.ReadInt8())
	log.Println("", rpkg.ReadInt16())
	log.Println("", rpkg.ReadInt32())
	log.Println("", rpkg.ReadInt64())

	log.Println("", rpkg.ReadUint8())
	log.Println("", rpkg.ReadUint16())
	log.Println("", rpkg.ReadUint32())
	log.Println("", rpkg.ReadUint64())

	log.Println("", rpkg.ReadBool())

	log.Println(rpkg.ReadFloat32())
	log.Println(rpkg.ReadFloat64())

	log.Println("", rpkg.ReadString())
	log.Println("", rpkg.ReadStringL(20))

	log.Println("", rpkg.ReadInt64())
}
