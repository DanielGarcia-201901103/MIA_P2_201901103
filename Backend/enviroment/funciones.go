package enviroment

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func Command(cmd string) {

	comando := ""

	if UserLogged_.User != "" {

		comando = BRED + "\n" + UserLogged_.User + "@" + UserLogged_.IdDisk + " ~-" + BGREEN + " Ejecutando comando: " + cmd + DEFAULT
		ContentConsola += "\n\n" + UserLogged_.User + "@" + UserLogged_.IdDisk + " ~- Ejecutando comando: " + cmd + "\n"
	} else {

		comando = BRED + "\n~-" + BGREEN + " Ejecutando comando: " + cmd + DEFAULT
		ContentConsola += "\n\n~- Ejecutando comando: " + cmd + "\n"
	}

	fmt.Println(comando)
}

func Error(msg string) {
	fmt.Println("\n" + BRED + "   ERROR: " + msg + DEFAULT)
	ContentConsola += "\n   ERROR: " + msg + "\n"
}

func Message(msg string) {
	fmt.Print("\n" + BCYN + "   " + msg + DEFAULT)
	ContentConsola += "\n   " + msg
}

func MessageWhite(msg string) {
	fmt.Print("\n\n" + BWHITE + "   - " + msg + DEFAULT + "\n")
	ContentConsola += "\n\n   - " + msg + "\n"
}

func Advertencia(msg string) {
	fmt.Println("\n" + BWHITE + "   - ADVERTENCIA: " + msg + DEFAULT)
	ContentConsola += "\n   - ADVERTENCIA: " + msg + "\n"
}

func MsgMounted(msg string) {
	fmt.Print("\n     " + msg + DEFAULT)
	ContentConsola += "\n     " + msg
}

func StructToBytes(str interface{}) []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(str)
	if err != nil {
		Error("no se pudo convertir la estructura a bytes")
	}
	return buffer.Bytes()
}

func BytesToMBR(bts []byte) MBR {
	str := MBR{}
	decoder := gob.NewDecoder(bytes.NewReader(bts))
	err := decoder.Decode(&str)
	if err != nil && err != io.EOF {
		Error("no se pudo convertir los bytes a estructura MBR")
	}
	return str
}

func BytesToEBR(bts []byte) EBR {
	str := EBR{}
	decoder := gob.NewDecoder(bytes.NewReader(bts))
	err := decoder.Decode(&str)
	if err != nil && err != io.EOF {
		Error("no se pudo convertir los bytes a estructura EBR")
	}
	return str
}

func GetTime() string {
	currentTime := time.Now()
	dateTimeString := currentTime.Format("2006-01-02 15:04:05")

	return dateTimeString
}

func ByteToStr(array []byte) string {
	// fmt.Println("------------------------")
	// fmt.Println(array)
	cont := 0
	str := ""
	for {
		if cont == len(array) {
			break
		} else {
			if array[cont] == uint8(0) {
				array[cont] = uint8(0)
				break
			} else if array[cont] != 0 {
				str += string(array[cont])
			}
		}
		cont++
	}

	return str
}

func ByteToInt(part []byte) int {

	fus := -8888888888888

	ff1 := ByteToStr(part[:])
	partSize, err := strconv.Atoi(ff1)
	if err != nil {
		Error("no se pudo convertir de bytes a int")
		return fus
	}
	fus = partSize

	return fus
}

func ExistDisk(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
