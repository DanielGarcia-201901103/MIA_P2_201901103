package comandos

import (
	"encoding/binary"
	"fmt"
	"math"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

func Login(usr string, pass string, id string) {

	enviroment.Command("login")
	fmt.Println(enviroment.BBLUE + "\n   login -user=" + usr + " -pass=" + pass + " -id=" + id + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   login -user=" + usr + " -pass=" + pass + " -id=" + id + "\n"

	if usr != "" {

		if pass != "" {

			if id != "" {

				loginUser(usr, pass, id)

			} else {

				enviroment.Error("en -id, debe colocar un id valido")
			}

		} else {

			enviroment.Error("en -pass, debe colocar una contraseña valida")
		}

	} else {

		enviroment.Error("en -user, debe colocar un ususario valido")
	}
}

func loginUser(user string, pass string, id string) {

	if !userLoggedExist() {

		index := IdExist(id)

		if index != -1 {

			credential := validateCredentials(user, pass, index)

			if credential == 1 {

				enviroment.UserLogged_.IdDisk = id
				enviroment.UserLogged_.User = user

				enviroment.Message("- Bienvenido usuario: " + user)
				enviroment.ContentLogin = "ok"

			} else if credential == -1 {

				enviroment.Error("El usuario no existe")
				enviroment.ContentLogin = "usuario"

			} else if credential == -4 {

				enviroment.Error("El disco no posee el formato EXT2")
				enviroment.ContentLogin = "ext"
			} else {

				enviroment.Error("Auntenticacion fallida")
				enviroment.ContentLogin = "fallida"
			}

		} else {

			enviroment.Error("No existe una particion montada asociada al id")
			enviroment.ContentLogin = "particion"
		}

	} else {
		enviroment.Error("ya existe un usuario logueado, debe primero debe cerrar la sesion actual")
		enviroment.ContentLogin = "logueado"
	}

}

func LoginFront(user string, pass string, id string) {

	enviroment.Command("login client")
	fmt.Println(enviroment.BBLUE + "\n   login -user=" + user + " -pass=" + pass + " -id=" + id + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   login -user=" + user + " -pass=" + pass + " -id=" + id + "\n"

	index := IdExist(id)

	if index != -1 {

		credential := validateCredentials(user, pass, index)

		if credential == 1 {

			enviroment.Message("- Bienvenido usuario: " + user + " puede visualizar los reportes")
			enviroment.ContentLogin = "ok"

		} else if credential == -1 {

			enviroment.Error("El usuario no existe")
			enviroment.ContentLogin = "usuario"

		} else if credential == -4 {

			enviroment.Error("El disco no posee el formato EXT2")
			enviroment.ContentLogin = "ext"
		} else {

			enviroment.Error("Auntenticacion fallida")
			enviroment.ContentLogin = "fallida"
		}

	} else {

		enviroment.Error("No existe una particion montada asociada al id")
		enviroment.ContentLogin = "particion"
	}
}

func userLoggedExist() bool {

	if enviroment.UserLogged_.IdDisk == "" && enviroment.UserLogged_.User == "" {

		return false
	} else {

		return true
	}
}

func validateCredentials(user string, pass string, index int) int {

	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	if _superBlock.S_magic != 61267 {

		return -4
	}

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += string(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") && !strings.Contains(data[0], "0") {

			if strings.Contains(data[3], user) {

				if strings.Contains(data[4], pass) {

					return 1

				} else {

					return -2
				}
			}
		}
	}

	return -1
}

func Logout() {

	enviroment.Command("logout")

	if userLoggedExist() {

		enviroment.UserLogged_.IdDisk = ""
		enviroment.UserLogged_.User = ""

		enviroment.Message("- Sesion cerrada correctamente")
		enviroment.ContentLogin = "ok"

	} else {

		enviroment.Error("No existe ningun usuario logueado")
		enviroment.ContentLogin = "no"
	}
}

func Mkgrp(name string) {

	enviroment.Command("mkgrp")
	fmt.Println(enviroment.BBLUE + "\n   mkrgp -name=" + name + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkrgp -name=" + name + "\n"

	if name != "" {

		if len(name) <= 10 {

			createGroup(name)
		} else {

			enviroment.Error("en -name, el nombre del grupo no puede ser mayor a 10 caracteres")
		}
	} else {

		enviroment.Error("en -name, el nombre del grupo no puede ser vacio")
	}

}

func createGroup(usr string) {

	if userLoggedExist() {

		if strings.Contains(enviroment.UserLogged_.User, "root") {

			if existGroup(usr) {

				if isActiveGroup(usr) {

					enviroment.Error("Este grupo de usuarios ya existe")

				} else {

					activateGroup(usr)
					enviroment.Message("- Grupo de usuarios: " + usr + " creado correctamente")
				}

			} else {

				tempId := getIdGroup()
				if tempId != 0 {

					_createGroup(usr, tempId)
					enviroment.Message("- Grupo de usuarios: " + usr + " creado correctamente")

				}
			}
		} else {

			enviroment.Error("Solo el usuario root puede crear grupos")
		}
	} else {

		enviroment.Error("Debe iniciar sesion como root para crear un grupo")
	}
}

func existGroup(user string) bool {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	// fmt.Println(lines)

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "G") {

			if strings.Contains(data[2], user) {

				return true
			}
		}
	}

	return false
}

func isActiveGroup(user string) bool {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "G") && !strings.Contains(data[0], "0") {

			if strings.Contains(data[2], user) {

				return true
			}
		}
	}

	return false
}

func getIdGroup() int {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	cont := 0

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "G") {

			cont++
		}
	}

	return cont + 1
}

func activateGroup(user string) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	newContent := ""

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "G") {

			if strings.Contains(data[2], user) {

				data[0] = strconv.Itoa(getIdGroup())
			}
			newContent += data[0] + ",G," + data[2] + "\n"
		} else {

			newContent += lines[i] + "\n"
		}
	}

	cantBlock := int(math.Ceil(float64(len(newContent)) / float64(64)))
	inicio := 0
	fin := 64

	if len(data) <= 64 {
		fin = len(data)
	}

	for i := 0; i < cantBlock; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := enviroment.FileBlock{}
			copy(_fileBlock.B_content[:], newContent[inicio:fin])
			WriteFileBlock(path, &_fileBlock, _inode.I_block[i])
			inicio += 64
			if len(newContent)-inicio <= 64 {

				fin = len(newContent)
			} else {

				fin += 64
			}
		}
	}
}

func _createGroup(user string, id int) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	indexBlock := -1

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] == -1 {

			indexBlock = i
			break
		}
	}

	if indexBlock == -1 {

		enviroment.Error("No hay espacio para crear el grupo")
	}

	_fileBlock := ReadFileBlock(path, _inode.I_block[indexBlock-1])
	data := enviroment.ByteToStr(_fileBlock.B_content[:])

	newLine := strconv.Itoa(id) + ",G," + user + "\n"

	if len(data) <= 64 {

		if (len(data) + len(newLine)) <= 64 {

			data += newLine
			copy(_fileBlock.B_content[:], data)
			WriteFileBlock(path, &_fileBlock, _inode.I_block[indexBlock-1])

		} else {

			restante := 64 - len(data)
			newLine1 := newLine[:restante]
			newLine2 := newLine[restante:]

			data1 := data + newLine1

			_fileBlock1 := enviroment.FileBlock{}
			copy(_fileBlock1.B_content[:], data1)
			WriteFileBlock(path, &_fileBlock1, _inode.I_block[indexBlock-1])

			_fileBlock2 := enviroment.FileBlock{}
			copy(_fileBlock2.B_content[:], newLine2)
			WriteFileBlock(path, &_fileBlock2, _superBlock.S_first_blo)

			// Actualizar inodo

			_inode.I_block[indexBlock] = _superBlock.S_first_blo
			copy(_inode.I_mtime[:], enviroment.GetTime())
			WriteInode(path, &_inode, inodeUserTxt)

			// Actualizar superbloque

			_superBlock.S_first_blo += int32(binary.Size(enviroment.FileBlock{}))
			_superBlock.S_free_blocks_count--
			WriteSuperBlock(path, &_superBlock, mounted.Start)

			// Actualizar bitmap de bloques

			WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))
		}
	} else {

		_fileBlock1 := enviroment.FileBlock{}
		copy(_fileBlock1.B_content[:], newLine)
		WriteFileBlock(path, &_fileBlock1, _superBlock.S_first_blo)

		// Actualizar inodo

		_inode.I_block[indexBlock] = _superBlock.S_first_blo
		copy(_inode.I_mtime[:], enviroment.GetTime())
		WriteInode(path, &_inode, inodeUserTxt)

		// Actualizar superbloque
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FileBlock{}))
		_superBlock.S_free_blocks_count--
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizar bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))
	}
}

func Rmgrp(name string) {

	enviroment.Command("rmgrp")
	fmt.Println(enviroment.BBLUE + "\n   rmgrp -name=" + name + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   rmgrp -name=" + name + "\n"

	if name != "" {

		deleteGroup(name)

	} else {

		enviroment.Error("en -name, debe colocar un nombre para el grupo")
	}
}

func deleteGroup(usr string) {

	if userLoggedExist() {

		if strings.Contains(enviroment.UserLogged_.User, "root") {

			if existGroup(usr) {

				if isActiveGroup(usr) {

					_deleteGroup(usr)
					enviroment.Message("- Grupo de usuarios: " + usr + " eliminado correctamente")
				} else {

					enviroment.Error("Este grupo de usuarios no existe")
				}
			} else {

				enviroment.Error("Este grupo de usuarios no existe")
			}
		} else {

			enviroment.Error("Solo el usuario root tiene los permisos para crear grupos")
		}
	} else {

		enviroment.Error("Debe iniciar sesion con el usuario root")
	}
}

func _deleteGroup(usr string) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	newContent := ""

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "G") {

			if strings.Contains(data[2], usr) {

				data[0] = "0"
			}
			newContent += data[0] + "," + data[1] + "," + data[2] + "\n"
		} else {

			newContent += lines[i] + "\n"
		}
	}

	cantBlock := int(math.Ceil(float64(len(newContent)) / float64(64)))
	inicio := 0
	fin := 64

	// fmt.Println(newContent)

	if len(data) <= 64 {
		fin = len(data)
	}

	for i := 0; i < cantBlock; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := enviroment.FileBlock{}
			copy(_fileBlock.B_content[:], newContent[inicio:fin])
			WriteFileBlock(path, &_fileBlock, _inode.I_block[i])
			inicio += 64
			if len(newContent)-inicio <= 64 {

				fin = len(newContent)
			} else {

				fin += 64
			}
		}
	}
}

func Mkusr(usr string, pass string, grp string) {

	enviroment.Command("mkusr")
	fmt.Println(enviroment.BBLUE + "\n   mkusr -usr=" + usr + " -pass=" + pass + " -grp=" + grp + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkusr -usr=" + usr + " -pass=" + pass + " -grp=" + grp + "\n"

	if usr != "" && pass != "" && grp != "" {

		if len(usr) <= 10 && len(pass) <= 10 {

			createUser(usr, pass, grp)
		} else {

			enviroment.Error("en -usr y -pass, solo se permiten 10 caracteres")
		}

	} else {

		enviroment.Error("en -usr, -pass y -grp, son parametros obligatorios")
	}
}

func createUser(usr string, pass string, grp string) {

	if userLoggedExist() {

		if strings.Contains(enviroment.UserLogged_.User, "root") {

			if existGroup(grp) {

				if existUser(usr) {

					if isActiveUser(usr) {

						enviroment.Error("Este usuario ya existe")
					} else {

						activateUser(usr, pass, grp)
						enviroment.Message("- Usuario: " + usr + " creado correctamente")
					}
				} else {

					tempId := getIdUser()

					if tempId != 0 {

						_createUser(usr, tempId, pass, grp)
						enviroment.Message("- Usuario: " + usr + " creado correctamente")
					}
				}
			} else {

				enviroment.Error("El grupo al que quiere asociar al usuario no existe")
			}
		} else {

			enviroment.Error("Solo el usuario root tiene los permisos para crear usuarios")
		}
	} else {

		enviroment.Error("Debe iniciar sesion con el usuario root")
	}
}

func existUser(user string) bool {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	// fmt.Println(lines)

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") {

			if strings.Contains(data[3], user) {

				return true
			}
		}
	}

	return false
}

func isActiveUser(user string) bool {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") && !strings.Contains(data[0], "0") {

			if strings.Contains(data[3], user) {

				return true
			}
		}
	}

	return false
}

func getIdUser() int {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	lines := strings.Split(data, "\n")

	cont := 0

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") {

			cont++
		}
	}

	return cont + 1
}

func _createUser(user string, id int, pass string, grp string) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	indexBlock := -1

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] == -1 {

			indexBlock = i
			break
		}
	}

	if indexBlock == -1 {

		enviroment.Error("No hay espacio para crear el grupo")
	}

	_fileBlock := ReadFileBlock(path, _inode.I_block[indexBlock-1])
	data := enviroment.ByteToStr(_fileBlock.B_content[:])

	newLine := strconv.Itoa(id) + ",U," + grp + "," + user + "," + pass + "\n"

	if len(data) <= 64 {

		if (len(data) + len(newLine)) <= 64 {

			data += newLine
			copy(_fileBlock.B_content[:], data)
			WriteFileBlock(path, &_fileBlock, _inode.I_block[indexBlock-1])

		} else {

			restante := 64 - len(data)
			newLine1 := newLine[:restante]
			newLine2 := newLine[restante:]

			data1 := data + newLine1

			_fileBlock1 := enviroment.FileBlock{}
			copy(_fileBlock1.B_content[:], data1)
			WriteFileBlock(path, &_fileBlock1, _inode.I_block[indexBlock-1])

			_fileBlock2 := enviroment.FileBlock{}
			copy(_fileBlock2.B_content[:], newLine2)
			WriteFileBlock(path, &_fileBlock2, _superBlock.S_first_blo)

			// Actualizar inodo

			_inode.I_block[indexBlock] = _superBlock.S_first_blo
			copy(_inode.I_mtime[:], enviroment.GetTime())
			WriteInode(path, &_inode, inodeUserTxt)

			// Actualizar superbloque

			_superBlock.S_first_blo += int32(binary.Size(enviroment.FileBlock{}))
			_superBlock.S_free_blocks_count--
			WriteSuperBlock(path, &_superBlock, mounted.Start)

			// Actualizar bitmap de bloques

			WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))
		}
	} else {

		_fileBlock1 := enviroment.FileBlock{}
		copy(_fileBlock1.B_content[:], newLine)
		WriteFileBlock(path, &_fileBlock1, _superBlock.S_first_blo)

		// Actualizar inodo

		_inode.I_block[indexBlock] = _superBlock.S_first_blo
		copy(_inode.I_mtime[:], enviroment.GetTime())
		WriteInode(path, &_inode, inodeUserTxt)

		// Actualizar superbloque
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FileBlock{}))
		_superBlock.S_free_blocks_count--
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizar bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))
	}
}

func activateUser(user string, pass string, grp string) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	newContent := ""

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") {

			if strings.Contains(data[3], user) {

				data[0] = strconv.Itoa(getIdGroup())
				data[2] = grp
				data[4] = pass
			}
			newContent += data[0] + "," + data[1] + "," + data[2] + "," + data[3] + "," + data[4] + "\n"
		} else {

			newContent += lines[i] + "\n"
		}
	}

	cantBlock := int(math.Ceil(float64(len(newContent)) / float64(64)))
	inicio := 0
	fin := 64

	if len(data) <= 64 {
		fin = len(data)
	}

	for i := 0; i < cantBlock; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := enviroment.FileBlock{}
			copy(_fileBlock.B_content[:], newContent[inicio:fin])
			WriteFileBlock(path, &_fileBlock, _inode.I_block[i])
			inicio += 64
			if len(newContent)-inicio <= 64 {

				fin = len(newContent)
			} else {

				fin += 64
			}
		}
	}
}

func Rmusr(usr string) {

	enviroment.Command("rmusr")
	fmt.Println(enviroment.BBLUE + "\n   rmusr -user=" + usr + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   rmusr -user=" + usr + "\n"

	if usr != "" {

		deleteUser(usr)
	} else {

		enviroment.Error("en -user, debe colocar un usuario para eliminarlo")
	}
}

func deleteUser(usr string) {

	if userLoggedExist() {

		if strings.Contains(enviroment.UserLogged_.User, "root") {

			if existUser(usr) {

				if isActiveUser(usr) {

					_deleteUser(usr)
					enviroment.Message("- Usuario: " + usr + " eliminado correctamente")
				} else {

					enviroment.Error("Este usuario no existe")
				}
			} else {

				enviroment.Error("Este usuario no existe")
			}
		} else {

			enviroment.Error("Solo el usuario root tiene los permisos para crear grupos")
		}
	} else {

		enviroment.Error("Debe iniciar sesion con el usuario root")
	}
}

func _deleteUser(usr string) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	inodeUserTxt := _superBlock.S_inode_start + int32(binary.Size(enviroment.Inode{}))

	_inode := ReadInodo(path, inodeUserTxt)

	data := ""

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := ReadFileBlock(path, _inode.I_block[i])
			data += enviroment.ByteToStr(_fileBlock.B_content[:])
		}
	}

	newContent := ""

	lines := strings.Split(data, "\n")

	for i := 0; i < len(lines)-1; i++ {

		data := strings.Split(lines[i], ",")

		if strings.Contains(data[1], "U") {

			if strings.Contains(data[3], usr) {

				data[0] = "0"
			}
			newContent += data[0] + "," + data[1] + "," + data[2] + "," + data[3] + "," + data[4] + "\n"
		} else {

			newContent += lines[i] + "\n"
		}
	}

	cantBlock := int(math.Ceil(float64(len(newContent)) / float64(64)))
	inicio := 0
	fin := 64

	// fmt.Println(newContent)

	if len(data) <= 64 {
		fin = len(data)
	}

	for i := 0; i < cantBlock; i++ {

		if _inode.I_block[i] != -1 {

			_fileBlock := enviroment.FileBlock{}
			copy(_fileBlock.B_content[:], newContent[inicio:fin])
			WriteFileBlock(path, &_fileBlock, _inode.I_block[i])
			inicio += 64
			if len(newContent)-inicio <= 64 {

				fin = len(newContent)
			} else {

				fin += 64
			}
		}
	}
}
