package comandos

import (
	"fmt"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

func Mount(path string, name string) {

	enviroment.Command("mount")
	fmt.Println(enviroment.BBLUE + "\n   mount -path=" + path + " -name=" + name + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mount -path=" + path + " -name=" + name + "\n"

	if path != "" {

		if name != "" {

			if enviroment.ExistDisk(path) {

				mountHandler(path, name)
			} else {

				enviroment.Error("Disco no encontrado")
			}
		} else {

			enviroment.Error("en -name, no se ha especificado el nombre de la particion a montar")
		}
	} else {

		enviroment.Error("en -path, no se ha especificado una ruta valida para el disco (error en parametro de -driveletter)")
	}
}

func mountHandler(path string, name string) {

	mbrInfo := ReadMBR(path)

	index := -1

	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbrInfo.Mbr_partitions[i].Part_name[:]), name) {

			index = i
			break
		}
	}

	if index != -1 {

		mountPrimaryPartition(path, name, index)
	} else {

		if existeLogica(path, name) {

			mountLogicalPartition(path, name)
		} else {

			enviroment.Error("no existe la particion a montar")
		}
	}
}

func mountPrimaryPartition(path string, name string, index int) {

	mbrInfo := ReadMBR(path)

	if string(mbrInfo.Mbr_partitions[index].Part_type[:]) == "P" {

		if !isMounted(path, name) {

			copy(mbrInfo.Mbr_partitions[index].Part_status[:], "1")
			WriteMBR(path, &mbrInfo)

			tempId := getIdDisk(path)

			partStart := enviroment.ByteToInt(mbrInfo.Mbr_partitions[index].Part_start[:])

			_mountedPartition := enviroment.MountedPartitions{}

			_mountedPartition.Path = path
			_mountedPartition.IdDisk = tempId
			_mountedPartition.Name = name
			_mountedPartition.Index = index
			_mountedPartition.Typee = 0
			_mountedPartition.Start = partStart

			enviroment.MountedPartitionsList = append(enviroment.MountedPartitionsList, _mountedPartition)
			enviroment.Message("Particion montada con exito")
			showMounts()

		} else {

			enviroment.Error("la particion ya esta montada")
			showMounts()
		}
	} else {
		//cambiar de advertencia a Error y no continuar con la ejecución
		enviroment.Advertencia("No es ideal montar particiones extendidas")

		if !isMounted(path, name) {

			copy(mbrInfo.Mbr_partitions[index].Part_status[:], "1")
			WriteMBR(path, &mbrInfo)

			tempId := getIdDisk(path)

			partStart := enviroment.ByteToInt(mbrInfo.Mbr_partitions[index].Part_start[:])

			_mountedPartition := enviroment.MountedPartitions{}

			_mountedPartition.Path = path
			_mountedPartition.IdDisk = tempId
			_mountedPartition.Name = name
			_mountedPartition.Index = index
			_mountedPartition.Typee = 2
			_mountedPartition.Start = partStart

			enviroment.MountedPartitionsList = append(enviroment.MountedPartitionsList, _mountedPartition)
			enviroment.Message("- Particion montada con exito")
			showMounts()
		} else {

			enviroment.Error("la particion ya esta montada")
			showMounts()
		}
	}
}

func isMounted(path string, name string) bool {

	for _, mountedPartition := range enviroment.MountedPartitionsList {

		if strings.Contains(mountedPartition.Path, path) && strings.Contains(mountedPartition.Name, name) {

			return true
		}
	}

	return false
}

func getIdDisk(path string) string {

	id := "03"
	cont := 0

	for _, mountedPartition := range enviroment.MountedPartitionsList {

		if strings.Contains(mountedPartition.Path, path) {

			cont++
		}
	}

	tem := strings.LastIndex(path, "/")
	nameDisk := path[tem+1:]
	temp2 := strings.LastIndex(nameDisk, ".")
	nameDisk = nameDisk[:temp2]

	if dsk, encontrado := enviroment.IdentDisk[nameDisk]; encontrado {

		id += strconv.Itoa(dsk)
	} else {

		id += strconv.Itoa(enviroment.ContDisk)
		enviroment.IdentDisk[nameDisk] = enviroment.ContDisk
		enviroment.ContDisk++
	}

	id += enviroment.Letras[cont]

	return id

}

func showMounts() {

	enviroment.MessageWhite("Particiones montadas:")
	for _, mountedPartition := range enviroment.MountedPartitionsList {

		enviroment.MsgMounted(mountedPartition.IdDisk + " -> " + mountedPartition.Name)
	}
}

func mountLogicalPartition(path string, name string) {

	mbrInfo := ReadMBR(path)

	startExtended := 0
	sizeExtended := 0

	for i := 0; i < 4; i++ {

		if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "E" {

			startExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
			sizeExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
			break
		}
	}

	if startExtended == 0 {

		enviroment.Error("- No se encontro una particion extendida en el disco")
		return
	}

	_ebr := ReadEBR(path, startExtended)

	end := startExtended + sizeExtended
	ftell := 0

	if strings.Contains(string(_ebr.Part_name[:]), name) {

		if !isMounted(path, name) {

			copy(_ebr.Part_status[:], "1")
			WriteEBR(path, startExtended, &_ebr)

			tempId := getIdDisk(path)

			_mountedPartition := enviroment.MountedPartitions{}
			_mountedPartition.Path = path
			_mountedPartition.IdDisk = tempId
			_mountedPartition.Name = name
			_mountedPartition.Index = startExtended
			_mountedPartition.Typee = 1
			_mountedPartition.Start = startExtended + len(enviroment.StructToBytes(enviroment.EBR{}))

			enviroment.MountedPartitionsList = append(enviroment.MountedPartitionsList, _mountedPartition)

			enviroment.Message("- Particion montada con exito")
			showMounts()
		} else {

			enviroment.Error("la particion ya esta montada")
			showMounts()
		}
		return

	} else {

		for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

			if strings.Contains(string(_ebr.Part_name[:]), name) {

				if !isMounted(path, name) {

					copy(_ebr.Part_status[:], "1")
					WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

					tempId := getIdDisk(path)

					_mountedPartition := enviroment.MountedPartitions{}
					_mountedPartition.Path = path
					_mountedPartition.IdDisk = tempId
					_mountedPartition.Name = name
					_mountedPartition.Index = enviroment.ByteToInt(_ebr.Part_start[:])
					_mountedPartition.Typee = 1
					_mountedPartition.Start = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

					enviroment.MountedPartitionsList = append(enviroment.MountedPartitionsList, _mountedPartition)

					enviroment.Message("- Particion montada con exito")
					showMounts()
				} else {

					enviroment.Error("la particion ya esta montada")
					showMounts()
				}
				return
			}

			_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
			ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

			if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

				break
			}
		}
	}

	if strings.Contains(string(_ebr.Part_name[:]), name) {

		if !isMounted(path, name) {

			copy(_ebr.Part_status[:], "1")
			WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

			tempId := getIdDisk(path)

			_mountedPartition := enviroment.MountedPartitions{}
			_mountedPartition.Path = path
			_mountedPartition.IdDisk = tempId
			_mountedPartition.Name = name
			_mountedPartition.Index = enviroment.ByteToInt(_ebr.Part_start[:])
			_mountedPartition.Typee = 1
			_mountedPartition.Start = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

			enviroment.MountedPartitionsList = append(enviroment.MountedPartitionsList, _mountedPartition)

			enviroment.Message("- Particion montada con exito")
			showMounts()
		} else {

			enviroment.Error("la particion ya esta montada")
			showMounts()
		}
		return
	} else {

		enviroment.Error("No se encontro la particion")
		return
	}
}
