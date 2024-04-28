package comandos

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

func Mkdir(path string, r bool) {

	enviroment.Command("mkdir")
	fmt.Println(enviroment.BBLUE + "\n   mkdir -path=" + path + " -r=" + fmt.Sprint(r) + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkdir -path=" + path + " -r=" + fmt.Sprint(r) + "\n"

	if path != "" {
		mkdirHandler(path, r)
	} else {

		enviroment.Error("en -path no se ha especificado un valor valido")
	}
}

func mkdirHandler(path string, r bool) {

	if r {

		mkdirRecursive(path)
		enviroment.Message("- Se ha creado el directorio " + path + " y sus subdirectorios")
	} else {
		mkdirNoRecursive(path)
	}
}

func mkdirRecursive(path2 string) {

	if userLoggedExist() {

		var indexInode int32 = -1

		carpetas := strings.Split(path2, "/")

		for i := 1; i < len(carpetas); i++ {

			if i == 1 {

				uid := getIdUserLogged(enviroment.UserLogged_.User)
				indexInode = buscarRaiz(carpetas[i], uid)
			} else {

				uid := getIdUserLogged(enviroment.UserLogged_.User)
				indexInode = crearCarpetas(carpetas[i], uid, indexInode)
			}
		}
	} else {

		enviroment.Error("No hay ningun usuario logueado")
	}
}

func getIdUserLogged(user string) int32 {

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

		if strings.Contains(data[1], "U") {

			if strings.Contains(data[3], user) {

				if !strings.Contains(data[0], "0") {

					num, _ := strconv.Atoi(data[0])
					return int32(num)
				} else {

					return -1
				}
			}
		}
	}

	return -1
}

func buscarRaiz(carpeta string, uid int32) int32 {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	_inode := ReadInodo(path, _superBlock.S_inode_start)

	encontrado := false
	indexInodo := int32(-1)
	indexVacio := int32(-1)

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_folderBlock := ReadFolderBlock(path, _inode.I_block[i])

			if enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[2].B_inodo
				encontrado = true
				break

			} else if enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[3].B_inodo
				encontrado = true
				break

			} else if _folderBlock.B_content[3].B_inodo == -1 {

				_inode2 := enviroment.Inode{}
				inicialiarInodo(&_inode2)

				_inode2.I_uid = uid
				_inode2.I_gid = 1
				_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
				copy(_inode2.I_atime[:], enviroment.GetTime())
				copy(_inode2.I_ctime[:], enviroment.GetTime())
				copy(_inode2.I_mtime[:], enviroment.GetTime())
				copy(_inode2.I_type[:], "0")
				_inode2.I_perm = 664

				// Escribo el nuevo inodo
				WriteInode(path, &_inode2, _superBlock.S_first_ino)

				// Actualizo el folder block

				_folderBlock.B_content[3].B_inodo = _superBlock.S_first_ino
				copy(_folderBlock.B_content[3].B_name[:], carpeta)
				WriteFolderBlock(path, &_folderBlock, _inode.I_block[i])

				encontrado = true
				indexInodo = _superBlock.S_first_ino

				// Actualizo el superblock

				_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
				_superBlock.S_free_inodes_count -= 1
				WriteSuperBlock(path, &_superBlock, mounted.Start)

				// Actualizo el bitmap de inodos
				WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

				break

			}

		} else if _inode.I_block[i] == -1 {

			indexVacio = int32(i)
			break
		}
	}

	if encontrado {

		return indexInodo

	} else {

		// Creo el folderblock
		_folderBlock := enviroment.FolderBlock{}
		_folderBlock.B_content[0].B_inodo = _superBlock.S_inode_start
		copy(_folderBlock.B_content[0].B_name[:], "/")
		_folderBlock.B_content[1].B_inodo = _superBlock.S_inode_start
		copy(_folderBlock.B_content[1].B_name[:], "/")
		_folderBlock.B_content[2].B_inodo = _superBlock.S_first_ino
		copy(_folderBlock.B_content[2].B_name[:], carpeta)
		_folderBlock.B_content[3].B_inodo = -1
		copy(_folderBlock.B_content[3].B_name[:], "-")

		// Escribo el folderblock
		WriteFolderBlock(path, &_folderBlock, _superBlock.S_first_blo)

		// Actualizo el inodo
		_inode.I_block[indexVacio] = _superBlock.S_first_blo
		WriteInode(path, &_inode, _superBlock.S_inode_start)

		// Actualizo el superblock
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FolderBlock{}))
		_superBlock.S_free_blocks_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))

		// Creo el nuevo inodo
		_inode2 := enviroment.Inode{}
		inicialiarInodo(&_inode2)
		_inode2.I_uid = uid
		_inode2.I_gid = 1
		_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
		copy(_inode2.I_atime[:], enviroment.GetTime())
		copy(_inode2.I_ctime[:], enviroment.GetTime())
		copy(_inode2.I_mtime[:], enviroment.GetTime())
		copy(_inode2.I_type[:], "0")
		_inode2.I_perm = 664

		// Escribo el nuevo inodo
		WriteInode(path, &_inode2, _superBlock.S_first_ino)

		indexInodo = _superBlock.S_first_ino

		// Actualizo el superblock
		_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
		_superBlock.S_free_inodes_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de inodos
		WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

		return indexInodo
	}
}

func inicialiarInodo(inodo *enviroment.Inode) {

	inodo.I_uid = 1
	inodo.I_gid = 1
	inodo.I_size = int32(binary.Size(enviroment.Inode{}))

	for i := 0; i < 16; i++ {

		inodo.I_block[i] = -1
	}

	inodo.I_perm = 664
}

func crearCarpetas(carpeta string, uid int32, indexInodee int32) int32 {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	_inode := ReadInodo(path, indexInodee)

	encontrado := false
	indexInodo := int32(-1)
	indexVacio := int32(-1)

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_folderBlock := ReadFolderBlock(path, _inode.I_block[i])

			if enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[2].B_inodo
				encontrado = true
				break

			} else if enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[3].B_inodo
				encontrado = true
				break

			} else if _folderBlock.B_content[3].B_inodo == -1 {

				_inode2 := enviroment.Inode{}
				inicialiarInodo(&_inode2)

				_inode2.I_uid = uid
				_inode2.I_gid = 1
				_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
				copy(_inode2.I_atime[:], enviroment.GetTime())
				copy(_inode2.I_ctime[:], enviroment.GetTime())
				copy(_inode2.I_mtime[:], enviroment.GetTime())
				copy(_inode2.I_type[:], "0")
				_inode2.I_perm = 664

				// Escribo el nuevo inodo
				WriteInode(path, &_inode2, _superBlock.S_first_ino)

				// Actualizo el folder block

				_folderBlock.B_content[3].B_inodo = _superBlock.S_first_ino
				copy(_folderBlock.B_content[3].B_name[:], carpeta)
				WriteFolderBlock(path, &_folderBlock, _inode.I_block[i])

				encontrado = true
				indexInodo = _superBlock.S_first_ino

				// Actualizo el superblock

				_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
				_superBlock.S_free_inodes_count -= 1
				WriteSuperBlock(path, &_superBlock, mounted.Start)

				// Actualizo el bitmap de inodos
				WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

				break

			}

		} else if _inode.I_block[i] == -1 {

			indexVacio = int32(i)
			break
		}
	}

	if encontrado {

		return indexInodo

	} else {

		// Creo el folderblock
		_folderBlock := enviroment.FolderBlock{}
		_folderBlock.B_content[0].B_inodo = indexInodee
		copy(_folderBlock.B_content[0].B_name[:], "/")
		_folderBlock.B_content[1].B_inodo = indexInodee - int32(binary.Size(enviroment.Inode{}))
		copy(_folderBlock.B_content[1].B_name[:], "/")
		_folderBlock.B_content[2].B_inodo = _superBlock.S_first_ino
		copy(_folderBlock.B_content[2].B_name[:], carpeta)
		_folderBlock.B_content[3].B_inodo = -1
		copy(_folderBlock.B_content[3].B_name[:], "-")

		// Escribo el folderblock
		WriteFolderBlock(path, &_folderBlock, _superBlock.S_first_blo)

		// Actualizo el inodo
		_inode.I_block[indexVacio] = _superBlock.S_first_blo
		WriteInode(path, &_inode, indexInodee)

		// Actualizo el superblock
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FolderBlock{}))
		_superBlock.S_free_blocks_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))

		// Creo el nuevo inodo
		_inode2 := enviroment.Inode{}
		inicialiarInodo(&_inode2)
		_inode2.I_uid = uid
		_inode2.I_gid = 1
		_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
		copy(_inode2.I_atime[:], enviroment.GetTime())
		copy(_inode2.I_ctime[:], enviroment.GetTime())
		copy(_inode2.I_mtime[:], enviroment.GetTime())
		copy(_inode2.I_type[:], "0")
		_inode2.I_perm = 664

		// Escribo el nuevo inodo
		WriteInode(path, &_inode2, _superBlock.S_first_ino)

		indexInodo = _superBlock.S_first_ino

		// Actualizo el superblock
		_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
		_superBlock.S_free_inodes_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de inodos
		WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

		return indexInodo
	}
}

func searchRoot(carpeta string) int32 {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	_inode := ReadInodo(path, _superBlock.S_inode_start)

	encontrado := false
	indexInodo := int32(-1)

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_folderBlock := ReadFolderBlock(path, _inode.I_block[i])

			if enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[2].B_inodo
				encontrado = true
				break

			} else if enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[3].B_inodo
				encontrado = true
				break

			} else if _folderBlock.B_content[3].B_inodo == -1 {

				encontrado = false
				indexInodo = int32(-1)
				break

			}

		} else if _inode.I_block[i] == -1 {

			encontrado = false
			indexInodo = int32(-1)
			break
		}
	}

	if encontrado {

		return indexInodo

	} else {

		return int32(-1)
	}
}

func searchFolder(carpeta string, indexInodee int32) int32 {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_inode := ReadInodo(path, indexInodee)

	encontrado := false
	indexInodo := int32(-1)

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_folderBlock := ReadFolderBlock(path, _inode.I_block[i])

			if enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[2].B_inodo
				encontrado = true
				break

			} else if enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) == carpeta {

				indexInodo = _folderBlock.B_content[3].B_inodo
				encontrado = true
				break

			} else if _folderBlock.B_content[3].B_inodo == -1 {

				encontrado = false
				indexInodo = int32(-1)
				break

			}

		} else if _inode.I_block[i] == -1 {

			encontrado = false
			indexInodo = int32(-1)
			break
		}
	}

	if encontrado {

		return indexInodo

	} else {

		return int32(-1)
	}
}

func mkdirNoRecursive(path2 string) {

	if userLoggedExist() {

		if isDirExist(path2) {

			enviroment.Error("Ya existe el directorio: " + path2)
			return
		}

		var indexInode int32 = -1

		carpetas := strings.Split(path2, "/")

		for i := 1; i < len(carpetas)-1; i++ {

			if i == 1 {

				indexInode = searchRoot(carpetas[i])

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta " + carpetas[i])
					break
				}
			} else {

				indexInode = searchFolder(carpetas[i], indexInode)

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta: " + carpetas[i] + " en la carpeta: " + carpetas[i-1])
					break
				}
			}
		}

		if indexInode != -1 {

			mkdirRecursive(path2)
			enviroment.Message("- Se ha creado el directorio " + path2)

		}
	} else {

		enviroment.Error("No hay ningun usuario logueado")
	}
}

func isDirExist(path2 string) bool {

	var indexInode int32 = -1
	res := false

	carpetas := strings.Split(path2, "/")

	for i := 1; i < len(carpetas); i++ {

		if i == 1 {

			indexInode = searchRoot(carpetas[i])

			if indexInode == -1 {

				res = false
				break
			}
		} else {

			indexInode = searchFolder(carpetas[i], indexInode)

			if indexInode == -1 {

				res = false
				break
			}
		}
	}

	if indexInode != -1 {

		res = true
	}

	return res
}

func Mkfile(path string, size int, cont string, r bool) {

	enviroment.Command("mkfile")
	fmt.Println(enviroment.BBLUE + "\n   mkfile -path=" + path + " -size=" + fmt.Sprint(size) + " -cont=" + cont + " -r=" + fmt.Sprint(r) + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkfile -path=" + path + " -size=" + fmt.Sprint(size) + " -cont=" + cont + " -r=" + fmt.Sprint(r) + "\n"

	if path != "" {

		if cont != "" {

			mkfileHandler(path, size, cont, r)
		} else if size != 0 {

			if size > 0 {

				mkfileHandler(path, size, cont, r)
			} else {

				enviroment.Error("el -size deber ser un valor mayor a 0")
			}
		} else if size == 0 {

			mkfileHandler(path, size, cont, r)
		} else {
			enviroment.Error("no se ha especificado un valor para -size o -cont")
		}
	} else {

		enviroment.Error("en -path no se ha especificado un valor valido")
	}
}

func mkfileHandler(path string, size int, cont string, r bool) {

	if r {

		fileRecursive(path, size, cont)
		asd := strings.LastIndex(path, "/")
		enviroment.Message("- Se ha creado el archivo " + path[asd+1:])
	} else {

		fileNoRecursive(path, size, cont)
		asd := strings.LastIndex(path, "/")
		enviroment.Message("- Se ha creado el archivo " + path[asd+1:])
	}
}

func fileNoRecursive(path2 string, size int, cont string) {

	if userLoggedExist() {

		if isDirExist(path2) {

			enviroment.Error("Ya existe el el archivo: " + path2)
			return
		}

		var indexInode int32 = -1

		carpetas := strings.Split(path2, "/")

		for i := 1; i < len(carpetas)-1; i++ {

			if i == 1 {

				indexInode = searchRoot(carpetas[i])

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta " + carpetas[i])
					break
				}
			} else {

				indexInode = searchFolder(carpetas[i], indexInode)

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta: " + carpetas[i] + " en la carpeta: " + carpetas[i-1])
					break
				}
			}
		}

		if indexInode != -1 {
			uid := getIdUserLogged(enviroment.UserLogged_.User)
			createInodeFile(carpetas[len(carpetas)-1], indexInode, size, cont, uid)

		}
	} else {

		enviroment.Error("No hay ningun usuario logueado")
	}
}

func fileRecursive(path2 string, size int, cont string) {

	if userLoggedExist() {

		if isDirExist(path2) {

			enviroment.Error("Ya existe el el archivo: " + path2)
			return
		} else {

			index := strings.LastIndex(path2, "/")
			pathTemp := path2[:index]
			mkdirRecursive(pathTemp)
		}

		var indexInode int32 = -1

		carpetas := strings.Split(path2, "/")

		for i := 1; i < len(carpetas)-1; i++ {

			if i == 1 {

				indexInode = searchRoot(carpetas[i])

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta " + carpetas[i])
					break
				}
			} else {

				indexInode = searchFolder(carpetas[i], indexInode)

				if indexInode == -1 {

					enviroment.Error("No existe la carpeta: " + carpetas[i] + " en la carpeta: " + carpetas[i-1])
					break
				}
			}
		}

		if indexInode != -1 {
			uid := getIdUserLogged(enviroment.UserLogged_.User)
			createInodeFile(carpetas[len(carpetas)-1], indexInode, size, cont, uid)

		}
	} else {

		enviroment.Error("No hay ningun usuario logueado")
	}
}

func createInodeFile(carpeta string, indexInodee int32, size int, cont string, uid int32) {

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	_inode := ReadInodo(path, indexInodee)

	encontrado := false
	indexVacio := int32(-1)
	indexInodo := int32(-1)

	for i := 0; i < 16; i++ {

		if _inode.I_block[i] != -1 {

			_folderBlock := ReadFolderBlock(path, _inode.I_block[i])

			if enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) == carpeta {

				encontrado = true
				break

			} else if enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) == carpeta {

				encontrado = true
				break

			} else if _folderBlock.B_content[3].B_inodo == -1 {

				_inode2 := enviroment.Inode{}
				inicialiarInodo(&_inode2)

				_inode2.I_uid = uid
				_inode2.I_gid = 1
				_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
				copy(_inode2.I_atime[:], enviroment.GetTime())
				copy(_inode2.I_ctime[:], enviroment.GetTime())
				copy(_inode2.I_mtime[:], enviroment.GetTime())
				copy(_inode2.I_type[:], "1")
				_inode2.I_perm = 664

				// Escribo el nuevo inodo
				WriteInode(path, &_inode2, _superBlock.S_first_ino)

				// Actualizo el folder block

				_folderBlock.B_content[3].B_inodo = _superBlock.S_first_ino
				copy(_folderBlock.B_content[3].B_name[:], carpeta)
				WriteFolderBlock(path, &_folderBlock, _inode.I_block[i])

				encontrado = true
				indexInodo = _superBlock.S_first_ino

				// Actualizo el superblock

				_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
				_superBlock.S_free_inodes_count -= 1
				WriteSuperBlock(path, &_superBlock, mounted.Start)

				// Actualizo el bitmap de inodos
				WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

				createFile(indexInodo, size, cont)

				break

			}

		} else if _inode.I_block[i] == -1 {

			indexVacio = int32(i)
			break
		}
	}

	if !encontrado {

		// Creo el folderblock
		_folderBlock := enviroment.FolderBlock{}
		_folderBlock.B_content[0].B_inodo = indexInodee
		copy(_folderBlock.B_content[0].B_name[:], "/")
		_folderBlock.B_content[1].B_inodo = indexInodee - int32(binary.Size(enviroment.Inode{}))
		copy(_folderBlock.B_content[1].B_name[:], "/")
		_folderBlock.B_content[2].B_inodo = _superBlock.S_first_ino
		copy(_folderBlock.B_content[2].B_name[:], carpeta)
		_folderBlock.B_content[3].B_inodo = -1
		copy(_folderBlock.B_content[3].B_name[:], "-")

		// Escribo el folderblock
		WriteFolderBlock(path, &_folderBlock, _superBlock.S_first_blo)

		// Actualizo el inodo
		_inode.I_block[indexVacio] = _superBlock.S_first_blo
		WriteInode(path, &_inode, indexInodee)

		// Actualizo el superblock
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FolderBlock{}))
		_superBlock.S_free_blocks_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))

		// Creo el nuevo inodo
		_inode2 := enviroment.Inode{}
		inicialiarInodo(&_inode2)
		_inode2.I_uid = uid
		_inode2.I_gid = 1
		_inode2.I_size = int32(binary.Size(enviroment.Inode{}))
		copy(_inode2.I_atime[:], enviroment.GetTime())
		copy(_inode2.I_ctime[:], enviroment.GetTime())
		copy(_inode2.I_mtime[:], enviroment.GetTime())
		copy(_inode2.I_type[:], "1")
		_inode2.I_perm = 664

		// Escribo el nuevo inodo
		WriteInode(path, &_inode2, _superBlock.S_first_ino)

		indexInodo = _superBlock.S_first_ino

		// Actualizo el superblock
		_superBlock.S_first_ino += int32(binary.Size(enviroment.Inode{}))
		_superBlock.S_free_inodes_count -= 1
		WriteSuperBlock(path, &_superBlock, mounted.Start)

		// Actualizo el bitmap de inodos
		WriteOneBM(path, _superBlock.S_bm_inode_start+(_superBlock.S_inodes_count-_superBlock.S_free_inodes_count-1))

		createFile(indexInodo, size, cont)

	}
}

func createFile(indexInode int32, size int, cont string) {

	if cont != "" {

		file, err := os.Open(cont)
		if err != nil {

			enviroment.Error("Error al abrir el archivo: " + cont)
			return
		}

		// Lee el tamaño del archivo
		stat, err := file.Stat()
		if err != nil {

			enviroment.Error("Error al tener el tamaño del archivo: " + cont)
			return
		}

		// Crea un slice de bytes con el tamaño del archivo
		contenido := make([]byte, stat.Size())

		// Lee el contenido del archivo en el slice de bytes
		_, err = file.Read(contenido)
		if err != nil {

			enviroment.Error("Error al leer el archivo: " + cont)
			return
		}

		file.Close()

		writeOnFile(indexInode, string(contenido))
	} else if size == 0 {

		writeOnFile(indexInode, "")
	} else if size > 0 {

		content := ""
		cont := 0

		for i := 0; i < size; i++ {

			content += strconv.Itoa(cont)
			cont++
			if cont > 9 {

				cont = 0
			}
		}

		writeOnFile(indexInode, content)
	}
}

func writeOnFile(indexInodee int32, data string) {

	// fmt.Println(len(data))

	if len(data) > 1024 {

		enviroment.Error("El texto es demasiado largo")
		return
	}

	index := IdExist(enviroment.UserLogged_.IdDisk)
	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	// Obtengo el inodo
	_inode := ReadInodo(path, indexInodee)

	cantBlock := int(math.Ceil(float64(len(data)) / float64(64)))
	inicio := 0
	fin := 64

	if len(data) <= 64 {
		fin = len(data)
	}

	for i := 0; i < cantBlock; i++ {

		_fileBlock := enviroment.FileBlock{}
		copy(_fileBlock.B_content[:], data[inicio:fin])

		// Le asigno una direccion al inodo
		_inode.I_block[i] = _superBlock.S_first_blo

		// Actualizando info del superbloque
		_superBlock.S_first_blo += int32(binary.Size(enviroment.FileBlock{}))
		_superBlock.S_free_blocks_count -= 1

		// Escribiendo el bitmap de bloques
		WriteOneBM(path, _superBlock.S_bm_block_start+(_superBlock.S_blocks_count-_superBlock.S_free_blocks_count-1))

		WriteFileBlock(path, &_fileBlock, _inode.I_block[i])
		inicio += 64
		if len(data)-inicio <= 64 {

			fin = len(data)
		} else {

			fin += 64
		}
	}

	// Actualizo el inodo
	_inode.I_size = int32(len(data))
	WriteInode(path, &_inode, indexInodee)

	// Actualizo el superblock
	WriteSuperBlock(path, &_superBlock, mounted.Start)
}
