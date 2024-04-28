package comandos

import (
	"fmt"
	"io"
	"os"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

func Fdisk(driveletterr string, path string, size int, unit string, fit string, typee string, name string) {
	// visual
	enviroment.Command("fdisk")
	fmt.Println(enviroment.BBLUE + "\n   fdisk -path=" + path + " -size=" + strconv.Itoa(size) + " -unit=" + unit + " -fit=" + fit + " -name=" + name + " -type=" + typee + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   fdisk -path=" + path + " -size=" + strconv.Itoa(size) + " -unit=" + unit + " -fit=" + fit + " -name=" + name + " -type=" + typee + "\n"

	if driveletterr != "" {
		if size > 0 {
			if name != "" {
				if len(name) < 16 {
					if enviroment.ExistDisk(path) {
						fdiskHandler(path, size, unit, fit, typee, name)
					} else {
						enviroment.Error("en -path, el disco no existe")
					}
				} else {
					enviroment.Error("en -name, el nombre del disco no puede tener mas de 16 caracteres")
				}
			} else {
				enviroment.Error("en -name, parametro no especificado")
			}
		} else {
			enviroment.Error("en -size, el parametro debe ser mayor a 0")
		}
	} else {
		enviroment.Error("en -driveletter, parametro no especificado")
	}

}

func fdiskHandler(path string, size int, unit string, fit string, typee string, name string) {

	enviroment.Message("- Creando particion: " + name)

	mbrInfo := ReadMBR(path)

	tamanioParticion := 0
	if strings.Contains(unit, "K") || unit == "" {
		tamanioParticion = size * 1024
	} else if strings.Contains(unit, "M") {
		tamanioParticion = size * 1024 * 1024
	} else if strings.Contains(unit, "B") {
		tamanioParticion = size
	} else {
		enviroment.Error("en >unit, parametro no valido")
		return
	}

	if enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) >= tamanioParticion {

		primary := 0
		extended := 0
		freePartitions := 0

		for i := 0; i < 4; i++ {
			if enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:]) == 0 {

				freePartitions++
			} else if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "P" {

				primary++
			} else if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "E" {

				extended++
			}

			if strings.Contains(string(mbrInfo.Mbr_partitions[i].Part_name[:]), name) {

				enviroment.Error("en >name, ya existe una particion con ese nombre")
				return
			}
		}

		// ^OJO AGREGAR VALIDACION SI EXISTE PARTICION CON MISMO NOMBRE EN LOGICAS LINEA DE CODIGO 85-92 C++
		if extended != 0 {

			if existeLogica(path, name) {

				enviroment.Error("en >name, ya existe una particion con ese nombre")
				return
			}
		}

		if freePartitions == 0 && !strings.Contains(typee, "L") {
			enviroment.Error("no se pueden crear mas particiones debido a que no hay particiones libres")
			return
		}

		if strings.Contains(typee, "P") || typee == "" {

			if typee == "" {
				typee = "P"
			}
			if primary == 3 && extended == 1 {

				enviroment.Error("no se pueden crear mas particiones primarias debido a que ya existen 3 y una extendida")
				return
			} else if primary == 4 {

				enviroment.Error("no se pueden crear mas particiones primarias debido a que ya existen 4")
				return
			} else if primary == 3 && extended == 0 {

				enviroment.Message("- Creando particion primaria")
				createPrimaryPartition(path, tamanioParticion, fit, name, typee)
			} else {

				enviroment.Message("- Creando particion primaria")
				createPrimaryPartition(path, tamanioParticion, fit, name, typee)
			}
		} else if strings.Contains(typee, "E") {

			if extended >= 1 {

				enviroment.Error("no se pueden crear mas particiones extendidas debido a que ya existe una")
				return
			} else {

				enviroment.Message("- Creando particion extendida")
				createExtendedPartition(path, tamanioParticion, fit, name, typee)
			}
		} else if strings.Contains(typee, "L") {

			if extended == 0 {

				enviroment.Error("no se puede crear una particion logica sin una particion extendida")
				return
			} else {

				enviroment.Message("- Creando particion logica")
				createLogicalPartition(path, tamanioParticion, fit, name, typee)
				// ^------------------------------------------------ DEBUGER EBR -------------------------------------------------
				// ShowLogicals(path)
			}
		} else {

			enviroment.Error("en >type, debe colocar un tipo valido")
			return
		}

	} else {
		enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
	}
}

func createPrimaryPartition(path string, tamanioParticion int, f string, name string, t string) {

	if strings.Contains(f, "WF") || f == "" {

		f = "W"
	} else if strings.Contains(f, "BF") {

		f = "B"
	} else if strings.Contains(f, "FF") {

		f = "F"
	} else {

		enviroment.Error("en >fit, parametro no valido")
		return
	}

	mbrInfo := ReadMBR(path)

	start := len(enviroment.StructToBytes(enviroment.MBR{}))
	end := -1
	espacioAnterior := false
	freeSpace := 0
	spaceUsed := len(enviroment.StructToBytes(enviroment.MBR{}))
	indexPartition := 0
	noSpace := 0

	// fmt.Println("start: ", start)
	// fmt.Println("end: ", end)
	// fmt.Println("espacioAnterior: ", espacioAnterior)
	// fmt.Println("freeSpace: ", freeSpace)
	// fmt.Println("spaceUsed: ", spaceUsed)
	// fmt.Println("indexPartition: ", indexPartition)
	// fmt.Println("noSpace", noSpace)

	if string(mbrInfo.Mbr_disk_fit[:]) == "F" {

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						_partition := enviroment.Particion{}
						copy(_partition.Part_status[:], "0")
						copy(_partition.Part_type[:], t)
						copy(_partition.Part_fit[:], f)
						copy(_partition.Part_start[:], strconv.Itoa(start))
						copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_partition.Part_name[:], name)
						mbrInfo.Mbr_partitions[indexPartition] = _partition
						WriteMBR(path, &mbrInfo)

						enviroment.Message("- Particion: " + name + " creada con exito")

						return
					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
							return
						}
						noSpace = 0
					}
				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
						return
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if end == -1 {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				_partition := enviroment.Particion{}
				copy(_partition.Part_status[:], "0")
				copy(_partition.Part_type[:], t)
				copy(_partition.Part_fit[:], f)
				copy(_partition.Part_start[:], strconv.Itoa(start))
				copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
				copy(_partition.Part_name[:], name)
				mbrInfo.Mbr_partitions[indexPartition] = _partition
				WriteMBR(path, &mbrInfo)

				enviroment.Message("- Particion: " + name + " creada con exito")
			} else {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
			}
		}
	} else if string(mbrInfo.Mbr_disk_fit[:]) == "B" {

		smallestSpace := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start
		smallestIndex := 0
		smallestStart := 0
		partEnded := false

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						if smallestSpace >= freeSpace {

							smallestSpace = freeSpace
							smallestIndex = indexPartition
							smallestStart = start
						}

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0

					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0
					}

				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						partEnded = true
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if !partEnded {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				if smallestSpace >= spaceLeft {

					smallestSpace = spaceLeft
					smallestIndex = indexPartition
					smallestStart = start
				}
			} else if smallestStart == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		} else {

			if smallestSpace == enviroment.ByteToInt(mbrInfo.Mbr_tamano[:])-len(enviroment.StructToBytes(enviroment.MBR{})) {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		}

		_partition := enviroment.Particion{}
		copy(_partition.Part_status[:], "0")
		copy(_partition.Part_type[:], t)
		copy(_partition.Part_fit[:], f)
		copy(_partition.Part_start[:], strconv.Itoa(smallestStart))
		copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
		copy(_partition.Part_name[:], name)

		mbrInfo.Mbr_partitions[smallestIndex] = _partition
		WriteMBR(path, &mbrInfo)

		enviroment.Message("- Particion: " + name + " creada con exito")

	} else if string(mbrInfo.Mbr_disk_fit[:]) == "W" {

		biggestSpace := 0
		biggestIndex := 0
		biggestStart := 0
		partEnded := false

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						if biggestSpace <= freeSpace {

							biggestSpace = freeSpace
							biggestIndex = indexPartition
							biggestStart = start
						}

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0

					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}
						noSpace = 0
					}
				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						partEnded = true
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if !partEnded {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				if biggestSpace <= spaceLeft {

					biggestSpace = spaceLeft
					biggestIndex = indexPartition
					biggestStart = start
				}
			} else if biggestStart == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		} else {

			if biggestSpace == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		}

		_partition := enviroment.Particion{}
		copy(_partition.Part_status[:], "0")
		copy(_partition.Part_type[:], t)
		copy(_partition.Part_fit[:], f)
		copy(_partition.Part_start[:], strconv.Itoa(biggestStart))
		copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
		copy(_partition.Part_name[:], name)
		mbrInfo.Mbr_partitions[biggestIndex] = _partition
		WriteMBR(path, &mbrInfo)

		enviroment.Message("- Particion: " + name + " creada con exito")

	}

}

func WriteEBR(path string, pos int, ebr *enviroment.EBR) {
	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	newpos, err := disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}
	ebrByte := enviroment.StructToBytes(ebr)

	_, err = disk.WriteAt(ebrByte, newpos)
	if err != nil {
		enviroment.Error("en escribir EBR en disco")
	}
	disk.Close()
}

func createExtendedPartition(path string, tamanioParticion int, f string, name string, t string) {

	if strings.Contains(f, "WF") || f == "" {

		f = "W"
	} else if strings.Contains(f, "BF") {

		f = "B"
	} else if strings.Contains(f, "FF") {

		f = "F"
	} else {

		enviroment.Error("en >fit, parametro no valido")
		return
	}

	mbrInfo := ReadMBR(path)

	start := len(enviroment.StructToBytes(enviroment.MBR{}))
	end := -1
	espacioAnterior := false
	freeSpace := 0
	spaceUsed := len(enviroment.StructToBytes(enviroment.MBR{}))
	indexPartition := 0
	noSpace := 0

	// fmt.Println("start: ", start)
	// fmt.Println("end: ", end)
	// fmt.Println("espacioAnterior: ", espacioAnterior)
	// fmt.Println("freeSpace: ", freeSpace)
	// fmt.Println("spaceUsed: ", spaceUsed)
	// fmt.Println("indexPartition: ", indexPartition)
	// fmt.Println("noSpace", noSpace)

	if string(mbrInfo.Mbr_disk_fit[:]) == "F" {

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						_partition := enviroment.Particion{}
						copy(_partition.Part_status[:], "0")
						copy(_partition.Part_type[:], t)
						copy(_partition.Part_fit[:], f)
						copy(_partition.Part_start[:], strconv.Itoa(start))
						copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_partition.Part_name[:], name)
						mbrInfo.Mbr_partitions[indexPartition] = _partition
						WriteMBR(path, &mbrInfo)

						_ebr := enviroment.EBR{}
						copy(_ebr.Part_status[:], "0")
						copy(_ebr.Part_fit[:], "-")
						copy(_ebr.Part_start[:], strconv.Itoa(start))
						copy(_ebr.Part_size[:], "0")
						copy(_ebr.Part_next[:], "-1")
						copy(_ebr.Part_name[:], "-")

						WriteEBR(path, start, &_ebr)

						enviroment.Message("- Particion: " + name + " creada con exito")

						// *--------------------------------------------------------------------------------
						fmt.Println()
						fmt.Println("   - status: ", string(_ebr.Part_status[:]))
						fmt.Println("   - fit: ", string(_ebr.Part_fit[:]))
						fmt.Println("   - start: ", enviroment.ByteToInt(_ebr.Part_start[:]))
						fmt.Println("   - size: ", enviroment.ByteToInt(_ebr.Part_size[:]))
						fmt.Println("   - next: ", enviroment.ByteToInt(_ebr.Part_next[:]))
						fmt.Println("   - name: ", string(_ebr.Part_name[:]))

						fmt.Println()
						fmt.Println("   - TAMANIO EBR: ", len(enviroment.StructToBytes(enviroment.EBR{})))
						fmt.Println("   - TAMANIO EBR2: ", len(enviroment.StructToBytes(_ebr)))
						// *--------------------------------------------------------------------------------

						return
					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
							return
						}
						noSpace = 0
					}
				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
						return
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if end == -1 {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				_partition := enviroment.Particion{}
				copy(_partition.Part_status[:], "0")
				copy(_partition.Part_type[:], t)
				copy(_partition.Part_fit[:], f)
				copy(_partition.Part_start[:], strconv.Itoa(start))
				copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
				copy(_partition.Part_name[:], name)
				mbrInfo.Mbr_partitions[indexPartition] = _partition
				WriteMBR(path, &mbrInfo)

				_ebr := enviroment.EBR{}
				copy(_ebr.Part_status[:], "0")
				copy(_ebr.Part_fit[:], "-")
				copy(_ebr.Part_start[:], strconv.Itoa(start))
				copy(_ebr.Part_size[:], "0")
				copy(_ebr.Part_next[:], "-1")
				copy(_ebr.Part_name[:], "-")

				WriteEBR(path, start, &_ebr)

				enviroment.Message("- Particion: " + name + " creada con exito")

				// *--------------------------------------------------------------------------------
				fmt.Println()
				fmt.Println("   - status: ", string(_ebr.Part_status[:]))
				fmt.Println("   - fit: ", string(_ebr.Part_fit[:]))
				fmt.Println("   - start: ", enviroment.ByteToInt(_ebr.Part_start[:]))
				fmt.Println("   - size: ", enviroment.ByteToInt(_ebr.Part_size[:]))
				fmt.Println("   - next: ", enviroment.ByteToInt(_ebr.Part_next[:]))
				fmt.Println("   - name: ", string(_ebr.Part_name[:]))

				fmt.Println()
				fmt.Println("   - TAMANIO EBR: ", len(enviroment.StructToBytes(enviroment.EBR{})))
				fmt.Println("   - TAMANIO EBR2: ", len(enviroment.StructToBytes(_ebr)))
				// *--------------------------------------------------------------------------------

			} else {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
			}
		}
	} else if string(mbrInfo.Mbr_disk_fit[:]) == "B" {

		smallestSpace := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start
		smallestIndex := 0
		smallestStart := 0
		partEnded := false

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						if smallestSpace >= freeSpace {

							smallestSpace = freeSpace
							smallestIndex = indexPartition
							smallestStart = start
						}

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0

					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0
					}

				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						partEnded = true
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if !partEnded {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				if smallestSpace >= spaceLeft {

					smallestSpace = spaceLeft
					smallestIndex = indexPartition
					smallestStart = start
				}
			} else if smallestStart == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		} else {

			if smallestSpace == enviroment.ByteToInt(mbrInfo.Mbr_tamano[:])-len(enviroment.StructToBytes(enviroment.MBR{})) {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		}

		_partition := enviroment.Particion{}
		copy(_partition.Part_status[:], "0")
		copy(_partition.Part_type[:], t)
		copy(_partition.Part_fit[:], f)
		copy(_partition.Part_start[:], strconv.Itoa(smallestStart))
		copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
		copy(_partition.Part_name[:], name)

		mbrInfo.Mbr_partitions[smallestIndex] = _partition
		WriteMBR(path, &mbrInfo)

		_ebr := enviroment.EBR{}
		copy(_ebr.Part_status[:], "0")
		copy(_ebr.Part_fit[:], "-")
		copy(_ebr.Part_start[:], strconv.Itoa(smallestStart))
		copy(_ebr.Part_size[:], "0")
		copy(_ebr.Part_next[:], "-1")
		copy(_ebr.Part_name[:], "-")

		WriteEBR(path, smallestStart, &_ebr)

		enviroment.Message("- Particion: " + name + " creada con exito")

		// *--------------------------------------------------------------------------------
		fmt.Println()
		fmt.Println("   - status: ", string(_ebr.Part_status[:]))
		fmt.Println("   - fit: ", string(_ebr.Part_fit[:]))
		fmt.Println("   - start: ", enviroment.ByteToInt(_ebr.Part_start[:]))
		fmt.Println("   - size: ", enviroment.ByteToInt(_ebr.Part_size[:]))
		fmt.Println("   - next: ", enviroment.ByteToInt(_ebr.Part_next[:]))
		fmt.Println("   - name: ", string(_ebr.Part_name[:]))

		fmt.Println()
		fmt.Println("   - TAMANIO EBR: ", len(enviroment.StructToBytes(enviroment.EBR{})))
		fmt.Println("   - TAMANIO EBR2: ", len(enviroment.StructToBytes(_ebr)))
		// *--------------------------------------------------------------------------------

	} else if string(mbrInfo.Mbr_disk_fit[:]) == "W" {

		biggestSpace := 0
		biggestIndex := 0
		biggestStart := 0
		partEnded := false

		for i := 0; i < 4; i++ {

			spaceUsed += enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

			if ocupado != 0 {

				if espacioAnterior {

					end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
					freeSpace = end - start

					if freeSpace >= tamanioParticion {

						if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

							indexPartition = 0
						} else {

							indexPartition = i - noSpace
						}

						if biggestSpace <= freeSpace {

							biggestSpace = freeSpace
							biggestIndex = indexPartition
							biggestStart = start
						}

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}

						noSpace = 0

					} else {

						espacioAnterior = false
						start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
						end = -1

						if i == 3 {

							partEnded = true
						}
						noSpace = 0
					}
				} else {

					start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					end = -1

					if indexPartition <= 2 {

						indexPartition = i + 1
					} else {

						indexPartition = 3
					}

					if i == 3 {

						partEnded = true
					}
				}
			} else {

				espacioAnterior = true
				end = -1

				if start == len(enviroment.StructToBytes(enviroment.MBR{})) {

					indexPartition = 0
				} else {

					indexPartition = i - noSpace
				}
				noSpace++
			}
		}

		if !partEnded {

			spaceLeft := enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start

			if spaceLeft >= tamanioParticion {

				if biggestSpace <= spaceLeft {

					biggestSpace = spaceLeft
					biggestIndex = indexPartition
					biggestStart = start
				}
			} else if biggestStart == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		} else {

			if biggestSpace == 0 {

				enviroment.Error("el disco no tiene espacio suficiente para crear la particion")
				return
			}
		}

		_partition := enviroment.Particion{}
		copy(_partition.Part_status[:], "0")
		copy(_partition.Part_type[:], t)
		copy(_partition.Part_fit[:], f)
		copy(_partition.Part_start[:], strconv.Itoa(biggestStart))
		copy(_partition.Part_size[:], strconv.Itoa(tamanioParticion))
		copy(_partition.Part_name[:], name)
		mbrInfo.Mbr_partitions[biggestIndex] = _partition
		WriteMBR(path, &mbrInfo)

		_ebr := enviroment.EBR{}
		copy(_ebr.Part_status[:], "0")
		copy(_ebr.Part_fit[:], "-")
		copy(_ebr.Part_start[:], strconv.Itoa(biggestStart))
		copy(_ebr.Part_size[:], "0")
		copy(_ebr.Part_next[:], "-1")
		copy(_ebr.Part_name[:], "-")

		WriteEBR(path, biggestStart, &_ebr)

		enviroment.Message("- Particion: " + name + " creada con exito")

		// *--------------------------------------------------------------------------------
		fmt.Println()
		fmt.Println("   - status: ", string(_ebr.Part_status[:]))
		fmt.Println("   - fit: ", string(_ebr.Part_fit[:]))
		fmt.Println("   - start: ", enviroment.ByteToInt(_ebr.Part_start[:]))
		fmt.Println("   - size: ", enviroment.ByteToInt(_ebr.Part_size[:]))
		fmt.Println("   - next: ", enviroment.ByteToInt(_ebr.Part_next[:]))
		fmt.Println("   - name: ", string(_ebr.Part_name[:]))

		fmt.Println()
		fmt.Println("   - TAMANIO EBR: ", len(enviroment.StructToBytes(enviroment.EBR{})))
		fmt.Println("   - TAMANIO EBR2: ", len(enviroment.StructToBytes(_ebr)))
		// *--------------------------------------------------------------------------------
	}

}

func createLogicalPartition(path string, tamanioParticion int, f string, name string, t string) {

	if strings.Contains(f, "WF") || f == "" {

		f = "W"
	} else if strings.Contains(f, "BF") {

		f = "B"
	} else if strings.Contains(f, "FF") {

		f = "F"
	} else {

		enviroment.Error("en >fit, parametro no valido")
		return
	}

	mbrInfo := ReadMBR(path)

	startExtended := 0
	sizeExtended := 0
	typeFit := ""

	for i := 0; i < 4; i++ {

		if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "E" {

			startExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
			sizeExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
			typeFit = string(mbrInfo.Mbr_partitions[i].Part_fit[:])
			break
		}
	}

	if startExtended == 0 {

		enviroment.Error("no existe una particion extendida en el disco")
		return
	}

	_ebr := ReadEBR(path, startExtended)

	// start := startExtended
	end := startExtended + sizeExtended
	ftell := 0
	// nextVacio := false

	if typeFit == "F" {

		// Ojo verificar si es necesario convertir a int
		if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

			if sizeExtended >= tamanioParticion {

				copy(_ebr.Part_status[:], "0")
				copy(_ebr.Part_fit[:], f)
				copy(_ebr.Part_start[:], strconv.Itoa(startExtended))
				copy(_ebr.Part_size[:], strconv.Itoa(tamanioParticion))
				copy(_ebr.Part_next[:], "-1")
				copy(_ebr.Part_name[:], name)

				WriteEBR(path, startExtended, &_ebr)
				enviroment.Message("- Particion: " + name + " creada con exito")
				return

			} else if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

				sizeBetween := enviroment.ByteToInt(_ebr.Part_next[:]) - enviroment.ByteToInt(_ebr.Part_start[:])

				if sizeBetween >= tamanioParticion {

					copy(_ebr.Part_status[:], "0")
					copy(_ebr.Part_fit[:], f)
					copy(_ebr.Part_start[:], strconv.Itoa(startExtended))
					copy(_ebr.Part_size[:], strconv.Itoa(tamanioParticion))
					copy(_ebr.Part_name[:], name)

					WriteEBR(path, startExtended, &_ebr)
					enviroment.Message("- Particion: " + name + " creada con exito")
					return
				}
			} else {

				enviroment.Error("no hay espacio para crear la particion logica")
				return
			}
		} else {

			for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

				if (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])) != enviroment.ByteToInt(_ebr.Part_next[:]) {

					if (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) >= tamanioParticion {

						partNextTemp := enviroment.ByteToInt(_ebr.Part_next[:])
						copy(_ebr.Part_next[:], strconv.Itoa(enviroment.ByteToInt(_ebr.Part_start[:])+enviroment.ByteToInt(_ebr.Part_size[:])))
						WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], _ebr.Part_next[:])
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], strconv.Itoa(partNextTemp))
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]), &_ebrTemp)

						enviroment.Message("- Particion: " + name + " creada con exito")
						return
					}
				}

				_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
				ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

				if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

					break
				}
			}

			nextTemp := enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])

			if end > nextTemp {

				if end-nextTemp >= tamanioParticion {

					copy(_ebr.Part_next[:], strconv.Itoa(nextTemp))
					WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

					_ebrTemp := enviroment.EBR{}
					copy(_ebrTemp.Part_status[:], "0")
					copy(_ebrTemp.Part_fit[:], f)
					copy(_ebrTemp.Part_start[:], strconv.Itoa(nextTemp))
					copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
					copy(_ebrTemp.Part_next[:], "-1")
					copy(_ebrTemp.Part_name[:], name)

					WriteEBR(path, nextTemp, &_ebrTemp)
					enviroment.Message("- Particion: " + name + " creada con exito")
					return
				} else {

					enviroment.Error("no hay espacio para crear la particion logica")
					return
				}
			} else {

				enviroment.Error("no hay espacio para crear la particion logica")
				return
			}
		}
	} else if typeFit == "B" {

		if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

			if sizeExtended >= tamanioParticion {

				copy(_ebr.Part_status[:], "0")
				copy(_ebr.Part_fit[:], f)
				copy(_ebr.Part_start[:], strconv.Itoa(startExtended))
				copy(_ebr.Part_size[:], strconv.Itoa(tamanioParticion))
				copy(_ebr.Part_next[:], "-1")
				copy(_ebr.Part_name[:], name)

				WriteEBR(path, startExtended, &_ebr)
				enviroment.Message("- Particion: " + name + " creada con exito")
				return
			} else {

				enviroment.Error("no hay espacio para crear la particion logica")
				return
			}
		} else {

			ebrAnt := enviroment.EBR{}
			smallestSize := sizeExtended

			for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

				if (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])) != enviroment.ByteToInt(_ebr.Part_next[:]) {

					if (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) >= tamanioParticion {

						if smallestSize >= (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) {

							if smallestSize != (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) {

								smallestSize = enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))
								ebrAnt = _ebr
							}
						}
					}
				}

				_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
				ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

				if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

					break
				}
			}

			nextTemp := enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])

			if end > nextTemp {

				if end-nextTemp >= tamanioParticion {

					if smallestSize >= (end - nextTemp) {

						copy(_ebr.Part_next[:], strconv.Itoa(nextTemp))
						WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], strconv.Itoa(nextTemp))
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], "-1")
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, nextTemp, &_ebrTemp)

						enviroment.Message("- Particion: " + name + " creada con exito")
						return

					} else {

						partNextTemp := ebrAnt.Part_next
						copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], partNextTemp[:])
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)
						enviroment.Message("- Particion: " + name + " creada con exito")
						return
					}
				} else {

					if smallestSize != sizeExtended {

						partNextTemp := ebrAnt.Part_next
						copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], partNextTemp[:])
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)
						enviroment.Message("- Particion: " + name + " creada con exito")
						return

					} else {

						enviroment.Error("- No se pudo crear la particion: " + name + " por falta de espacio")
						return
					}
				}
			} else {

				if smallestSize != sizeExtended {

					partNextTemp := ebrAnt.Part_next
					copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
					WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

					_ebrTemp := enviroment.EBR{}
					copy(_ebrTemp.Part_status[:], "0")
					copy(_ebrTemp.Part_fit[:], f)
					copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
					copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
					copy(_ebrTemp.Part_next[:], partNextTemp[:])
					copy(_ebrTemp.Part_name[:], name)

					WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)
					enviroment.Message("- Particion: " + name + " creada con exito")
					return

				} else {

					enviroment.Error("- No se pudo crear la particion: " + name + " por falta de espacio")
					return
				}
			}
		}
	} else if typeFit == "W" {

		if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

			if sizeExtended >= tamanioParticion {

				copy(_ebr.Part_status[:], "0")
				copy(_ebr.Part_fit[:], f)
				copy(_ebr.Part_start[:], strconv.Itoa(startExtended))
				copy(_ebr.Part_size[:], strconv.Itoa(tamanioParticion))
				copy(_ebr.Part_next[:], "-1")
				copy(_ebr.Part_name[:], name)

				WriteEBR(path, startExtended, &_ebr)

				enviroment.Message("- Particion: " + name + " creada con exito")
				return
			} else {

				enviroment.Error("- No se pudo crear la particion: " + name + " por falta de espacio")
				return
			}
		} else {

			ebrAnt := enviroment.EBR{}
			biggestSize := 0

			for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

				if (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])) != enviroment.ByteToInt(_ebr.Part_next[:]) {

					if (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) >= tamanioParticion {

						if biggestSize <= (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) {

							if biggestSize != (enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))) {

								biggestSize = enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))
								ebrAnt = _ebr
							}
						}
					}
				}

				_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
				ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

				if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

					break
				}
			}

			nextTemp := enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:])

			if end > nextTemp {

				if end-nextTemp >= tamanioParticion {

					if biggestSize <= (end - nextTemp) {

						copy(_ebr.Part_next[:], strconv.Itoa(nextTemp))
						WriteEBR(path, enviroment.ByteToInt(_ebr.Part_start[:]), &_ebr)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], strconv.Itoa(nextTemp))
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], "-1")
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, nextTemp, &_ebrTemp)
						enviroment.Message("- Particion: " + name + " creada con exito")
						return

					} else {

						partNextTemp := ebrAnt.Part_next
						copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], partNextTemp[:])
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)
						enviroment.Message("- Particion: " + name + " creada con exito")
						return
					}
				} else {

					if biggestSize != 0 {

						partNextTemp := ebrAnt.Part_next
						copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

						_ebrTemp := enviroment.EBR{}
						copy(_ebrTemp.Part_status[:], "0")
						copy(_ebrTemp.Part_fit[:], f)
						copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
						copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
						copy(_ebrTemp.Part_next[:], partNextTemp[:])
						copy(_ebrTemp.Part_name[:], name)

						WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)

						enviroment.Message("- Particion: " + name + " creada con exito")
						return
					} else {

						enviroment.Error("- No se pudo crear la particion logica: " + name + " por falta de espacio")
						return
					}
				}
			} else {

				if biggestSize != 0 {

					partNextTemp := ebrAnt.Part_next
					copy(ebrAnt.Part_next[:], strconv.Itoa(enviroment.ByteToInt(ebrAnt.Part_start[:])+enviroment.ByteToInt(ebrAnt.Part_size[:])))
					WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_start[:]), &ebrAnt)

					_ebrTemp := enviroment.EBR{}
					copy(_ebrTemp.Part_status[:], "0")
					copy(_ebrTemp.Part_fit[:], f)
					copy(_ebrTemp.Part_start[:], ebrAnt.Part_next[:])
					copy(_ebrTemp.Part_size[:], strconv.Itoa(tamanioParticion))
					copy(_ebrTemp.Part_next[:], partNextTemp[:])
					copy(_ebrTemp.Part_name[:], name)

					WriteEBR(path, enviroment.ByteToInt(ebrAnt.Part_next[:]), &_ebrTemp)

					enviroment.Message("- Particion: " + name + " creada con exito")
					return
				} else {

					enviroment.Error("- No se pudo crear la particion logica: " + name + " por falta de espacio")
					return
				}
			}
		}
	}
}

func ReadEBR(path string, pos int) enviroment.EBR {
	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	newpos, err := disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	noBytes := enviroment.StructToBytes(enviroment.EBR{})
	ebrByte := make([]byte, len(noBytes))
	_, err = disk.ReadAt(ebrByte, newpos)
	if err != nil {
		enviroment.Error("en leer EBR en disco")
	}

	disk.Close()

	_ebr := enviroment.BytesToEBR(ebrByte)

	return _ebr

}

func ShowLogicals(path string) {

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

	if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

		fmt.Println()
		fmt.Println("   - EBR:")
		fmt.Println("   - Part_status:", string(_ebr.Part_status[:]))
		fmt.Println("   - Part_fit:", string(_ebr.Part_fit[:]))
		fmt.Println("   - Part_start:", string(_ebr.Part_start[:]))
		fmt.Println("   - Part_size:", string(_ebr.Part_size[:]))
		fmt.Println("   - Part_next:", string(_ebr.Part_next[:]))
		fmt.Println("   - Part_name:", string(_ebr.Part_name[:]))
		return
	} else {

		for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

			fmt.Println()
			fmt.Println("   - EBR:")
			fmt.Println("   - Part_status:", string(_ebr.Part_status[:]))
			fmt.Println("   - Part_fit:", string(_ebr.Part_fit[:]))
			fmt.Println("   - Part_start:", string(_ebr.Part_start[:]))
			fmt.Println("   - Part_size:", string(_ebr.Part_size[:]))
			fmt.Println("   - Part_next:", string(_ebr.Part_next[:]))
			fmt.Println("   - Part_name:", string(_ebr.Part_name[:]))

			_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
			ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

			if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

				break
			}
		}
	}

	fmt.Println()
	fmt.Println("   - EBR:")
	fmt.Println("   - Part_status:", string(_ebr.Part_status[:]))
	fmt.Println("   - Part_fit:", string(_ebr.Part_fit[:]))
	fmt.Println("   - Part_start:", string(_ebr.Part_start[:]))
	fmt.Println("   - Part_size:", string(_ebr.Part_size[:]))
	fmt.Println("   - Part_next:", string(_ebr.Part_next[:]))
	fmt.Println("   - Part_name:", string(_ebr.Part_name[:]))
}

func existeLogica(path string, name string) bool {

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
		return false
	}

	_ebr := ReadEBR(path, startExtended)

	end := startExtended + sizeExtended
	ftell := 0

	if strings.Contains(string(_ebr.Part_name[:]), name) {

		return true

	} else {

		for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

			if strings.Contains(string(_ebr.Part_name[:]), name) {

				return true
			}

			_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
			ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

			if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

				break
			}
		}
	}

	if strings.Contains(string(_ebr.Part_name[:]), name) {

		return true
	} else {

		return false
	}
}
