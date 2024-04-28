package comandos

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

func Mkdisk(path string, size int, unit string, fit string) {
	// visual
	enviroment.Command("mkdisk")
	fmt.Println(enviroment.BBLUE + "\n   mkdisk -path=" + path + " -size=" + strconv.Itoa(size) + " -unit=" + unit + " -fit=" + fit + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkdisk -path=" + path + " -size=" + strconv.Itoa(size) + " -unit=" + unit + " -fit=" + fit + "\n"

	if path != "" {
		if size > 0 {
			enviroment.Message("- Creando disco")
			createDisk(path, size, unit, fit)
		} else {
			enviroment.Error("en -size, el parametro debe ser mayor a 0")
		}
	} else {
		enviroment.Error("en -path, parametro no especificado")
	}

}

func createDisk(path string, size int, unit string, fit string) {

	_size := size

	if strings.Contains(unit, "K") {
		size = size * 1024
	} else if strings.Contains(unit, "M") || unit == "" {
		size = size * 1024 * 1024
		_size = _size * 1024
	} else {
		enviroment.Error("en >unit, parametro no valido")
		return
	}

	if strings.Contains(fit, "BF") {
		fit = "B"
	} else if strings.Contains(fit, "FF") || fit == "" {
		fit = "F"
	} else if strings.Contains(fit, "WF") {
		fit = "W"
	} else {
		enviroment.Error("en >fit, parametro no valido")
		return
	}

	// crear carpetas
	CreatePaths(path)

	//crear bloque
	bloque := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		bloque[i] = 0
	}

	// crear disco
	disk, err := os.Create(path)
	if err != nil {
		enviroment.Error("en crear disco")
		fmt.Println(err)
		return
	}

	// escribir en disco
	for i := 0; i < _size; i++ {
		disk.Write(bloque)
	}

	// cerrar disco
	disk.Close()

	_MBR := enviroment.MBR{}

	//
	copy(_MBR.Mbr_tamano[:], strconv.Itoa(size))
	copy(_MBR.Mbr_fecha_creacion[:], enviroment.GetTime())
	signature := rand.Intn(1000) * 1000
	copy(_MBR.Mbr_disk_signature[:], strconv.Itoa(signature))
	copy(_MBR.Mbr_disk_fit[:], fit)

	// fmt.Println(string(_MBR.Mbr_tamano[:]))
	// fmt.Println(string(_MBR.Mbr_fecha_creacion[:]))
	// fmt.Println(string(_MBR.Mbr_disk_signature[:]))
	// fmt.Println(string(_MBR.Mbr_disk_fit[:]))

	_Partition := enviroment.Particion{}
	copy(_Partition.Part_status[:], "0")
	copy(_Partition.Part_type[:], "-")
	copy(_Partition.Part_fit[:], "-")
	copy(_Partition.Part_start[:], "0")
	copy(_Partition.Part_size[:], "0")
	// _Partition.Part_start = 0
	// _Partition.Part_size = 0
	copy(_Partition.Part_name[:], "-")

	for i := 0; i < 4; i++ {
		_MBR.Mbr_partitions[i] = _Partition
	}

	// escribir MBR
	WriteMBR(path, &_MBR)
	enviroment.Message("- Disco creado corectamente")
	// fmt.Println("LECTURA MBR DEL DISCO")
	// readMBR(path)
}

func CreatePaths(path string) {
	carpetas := strings.Split(path, "/")

	var ruta string

	for i := 1; i < len(carpetas)-1; i++ {
		ruta += "/" + carpetas[i]
	}

	command := "mkdir -p \"" + ruta + "\""

	// fmt.Println(command)

	_, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		enviroment.Error("en crear rutas del disco")
		fmt.Println(err)
	}
}

func WriteMBR(path string, mbr *enviroment.MBR) {
	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	newpos, err := disk.Seek(0, io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}
	mbrByte := enviroment.StructToBytes(mbr)

	// -----------------------------------------------------------------------------------------
	// aux := enviroment.MBR{}
	// aux2 := enviroment.StructToBytes(aux)
	// fmt.Println("Vacio",len(aux2))
	// fmt.Println("NoBytes",len(mbrByte))
	// fmt.Println("NoBytes2", int(binary.Size(enviroment.MBR{})))
	// -----------------------------------------------------------------------------------------

	_, err = disk.WriteAt(mbrByte, newpos)
	if err != nil {
		enviroment.Error("en escribir MBR en disco")
	}
	disk.Close()
}

func ReadMBR(path string) enviroment.MBR {
	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	newpos, err := disk.Seek(0, io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	noBytes := enviroment.StructToBytes(enviroment.MBR{})
	mbrByte := make([]byte, len(noBytes))
	_, err = disk.ReadAt(mbrByte, newpos)
	if err != nil {
		enviroment.Error("en leer MBR en disco")
	}

	disk.Close()

	_MBR := enviroment.BytesToMBR(mbrByte)

	return _MBR

	// var sise int = enviroment.ByteToInt(_MBR.Mbr_tamano[:])
	// fmt.Println(sise, sise+100)
	// fmt.Println(string(_MBR.Mbr_fecha_creacion[:]))
	// fmt.Println(string(_MBR.Mbr_disk_signature[:]))
	// fmt.Println(string(_MBR.Mbr_disk_fit[:]))

}

func Rmdisk(path string) {
	enviroment.Command("rmdisk")
	fmt.Println(enviroment.BBLUE + "\n   rmdisk -path=" + path + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   rmdisk -path=" + path + "\n"

	if path == "" {
		enviroment.Error("en -driveletter, parametro no especificado")
		return
	}

	if !enviroment.ExistDisk(path) {
		tem := strings.LastIndex(path, "/")
		name := path[tem+1:]
		enviroment.Error("el disco: " + name + " no existe")
		return
	}

	command := "rm \"" + path + "\""

	_, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		enviroment.Error("en elimanar disco disco")
		// fmt.Println(err)
		return
	}

	enviroment.Message("- Disco eliminado corectamente")
}

func ShowMBRInfo(path string) {
	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	newpos, err := disk.Seek(0, io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	noBytes := enviroment.StructToBytes(enviroment.MBR{})
	mbrByte := make([]byte, len(noBytes))
	_, err = disk.ReadAt(mbrByte, newpos)
	if err != nil {
		enviroment.Error("en leer MBR en disco")
	}

	disk.Close()

	_MBR := enviroment.BytesToMBR(mbrByte)

	fmt.Println()
	fmt.Println("   - Tamanio", string(_MBR.Mbr_tamano[:]))
	fmt.Println("   - Fecha creacion", string(_MBR.Mbr_fecha_creacion[:]))
	fmt.Println("   - Signature", string(_MBR.Mbr_disk_signature[:]))
	fmt.Println("   - Fit", string(_MBR.Mbr_disk_fit[:]))

	for i := 0; i < 4; i++ {
		println("\n  - Particion: ", i)
		fmt.Println("   - Status", string(_MBR.Mbr_partitions[i].Part_status[:]))
		fmt.Println("   - Type", string(_MBR.Mbr_partitions[i].Part_type[:]))
		fmt.Println("   - Fit", string(_MBR.Mbr_partitions[i].Part_fit[:]))
		fmt.Println("   - Start", string(_MBR.Mbr_partitions[i].Part_start[:]))
		fmt.Println("   - Size", string(_MBR.Mbr_partitions[i].Part_size[:]))
		fmt.Println("   - Name", string(_MBR.Mbr_partitions[i].Part_name[:]))
	}

}
