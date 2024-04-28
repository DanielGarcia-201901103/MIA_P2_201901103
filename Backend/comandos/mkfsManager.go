package comandos

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"paquetes/enviroment"
)

func Mkfs(id string, typee string, fs string) {

	enviroment.Command("mkfs")
	fmt.Println(enviroment.BBLUE + "\n   mkfs -id=" + id + " -type=" + typee + " -fs=" + fs + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   mkfs -id=" + id + " -type=" + typee + " -fs=" + fs + "\n"

	if id != "" {

		mkfsHandler(id)

	} else {

		enviroment.Error("en -id, debe colocar un id valido")
	}
}

func mkfsHandler(id string) {

	index := IdExist(id)

	if index != -1 {

		mountedPartition := enviroment.MountedPartitionsList[index]

		if mountedPartition.Typee == 0 {

			mbrInfo := ReadMBR(mountedPartition.Path)

			ext2Format(enviroment.ByteToInt(mbrInfo.Mbr_partitions[mountedPartition.Index].Part_start[:]), enviroment.ByteToInt(mbrInfo.Mbr_partitions[mountedPartition.Index].Part_size[:]), mountedPartition.Path)
		} else if mountedPartition.Typee == 1 {

			_ebr := ReadEBR(mountedPartition.Path, mountedPartition.Index)

			ext2Format((enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))), (enviroment.ByteToInt(_ebr.Part_size[:]) - len(enviroment.StructToBytes(enviroment.EBR{}))), mountedPartition.Path)
		} else if mountedPartition.Typee == 2 {

			enviroment.Error("no se puede aplicar un formato a una particion extendida")
		}

	} else {

		enviroment.Error("el id no esta asociado a una particion montada")
	}
}

func ext2Format(startPartition int, tamanioparticion int, paht string) {

	aux := (tamanioparticion - int(binary.Size(enviroment.SuperBlock{}))) / (4 + int(binary.Size(enviroment.Inode{})) + 3*int(binary.Size(enviroment.FileBlock{})))
	var n int32 = int32(aux)

	_superBlock := enviroment.SuperBlock{}

	_superBlock.S_filesystem_type = 2
	_superBlock.S_inodes_count = n
	_superBlock.S_blocks_count = 3 * n
	_superBlock.S_free_blocks_count = 3 * n
	_superBlock.S_free_inodes_count = n
	copy(_superBlock.S_mtime[:], enviroment.GetTime())
	copy(_superBlock.S_umtime[:], enviroment.GetTime())
	_superBlock.S_mnt_count = 1
	_superBlock.S_magic = 61267
	_superBlock.S_inode_s = int32(binary.Size(enviroment.Inode{}))
	_superBlock.S_block_s = int32(binary.Size(enviroment.FileBlock{}))

	_superBlock.S_bm_inode_start = int32(startPartition + binary.Size(_superBlock))
	_superBlock.S_bm_block_start = _superBlock.S_bm_inode_start + n
	_superBlock.S_inode_start = _superBlock.S_bm_block_start + (3 * n)
	_superBlock.S_block_start = _superBlock.S_inode_start + (n * int32(binary.Size(enviroment.Inode{})))

	_superBlock.S_first_ino = _superBlock.S_inode_start
	_superBlock.S_first_blo = _superBlock.S_block_start

	// ^Eliminar todo el contenido de la particion

	// Escribir el superbloque
	WriteSuperBlock(paht, &_superBlock, startPartition)

	// Creando el archivo user.txt
	// Creando la carpeta root

	_inode := enviroment.Inode{}

	_inode.I_uid = 1
	_inode.I_gid = 1
	_inode.I_size = 240
	copy(_inode.I_atime[:], enviroment.GetTime())
	copy(_inode.I_ctime[:], enviroment.GetTime())
	copy(_inode.I_mtime[:], enviroment.GetTime())
	var auxArray = [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	copy(_inode.I_block[:], auxArray[:])
	_inode.I_block[0] = _superBlock.S_first_blo
	copy(_inode.I_type[:], "0")
	_inode.I_perm = 664

	WriteInode(paht, &_inode, _superBlock.S_inode_start)

	// Actualizando el superbloque
	_superBlock.S_free_inodes_count--
	_superBlock.S_first_ino = _superBlock.S_first_ino + int32(binary.Size(enviroment.Inode{}))

	// Apuntador de carpetas

	_folderBlock := enviroment.FolderBlock{}
	copy(_folderBlock.B_content[0].B_name[:], "/")
	_folderBlock.B_content[0].B_inodo = _superBlock.S_first_ino - int32(binary.Size(enviroment.Inode{}))
	copy(_folderBlock.B_content[1].B_name[:], "/")
	_folderBlock.B_content[1].B_inodo = _superBlock.S_first_ino - int32(binary.Size(enviroment.Inode{}))
	copy(_folderBlock.B_content[2].B_name[:], "user.txt")
	_folderBlock.B_content[2].B_inodo = _superBlock.S_first_ino

	_folderBlock.B_content[3].B_inodo = -1
	copy(_folderBlock.B_content[3].B_name[:], "-")

	WriteFolderBlock(paht, &_folderBlock, _superBlock.S_first_blo)

	// Actualizando el superbloque
	_superBlock.S_free_blocks_count--
	_superBlock.S_first_blo = _superBlock.S_first_blo + int32(binary.Size(enviroment.FileBlock{}))

	content := "1,G,root\n1,U,root,root,123\n"

	// Inode del archivo user.txt

	_inodeUser := enviroment.Inode{}
	_inodeUser.I_uid = 1
	_inodeUser.I_gid = 1
	_inodeUser.I_size = int32(len(content)) + int32(binary.Size(enviroment.FolderBlock{}))
	copy(_inodeUser.I_atime[:], enviroment.GetTime())
	copy(_inodeUser.I_ctime[:], enviroment.GetTime())
	copy(_inodeUser.I_mtime[:], enviroment.GetTime())
	copy(_inodeUser.I_block[:], auxArray[:])
	_inodeUser.I_block[0] = _superBlock.S_first_blo
	copy(_inodeUser.I_type[:], "1")
	_inodeUser.I_perm = 664

	WriteInode(paht, &_inodeUser, _superBlock.S_first_ino)

	// Actualizando el superbloque
	_superBlock.S_free_inodes_count--
	_superBlock.S_first_ino = _superBlock.S_first_ino + int32(binary.Size(enviroment.Inode{}))

	// Contenido de archivo
	_fileBlock := enviroment.FileBlock{}
	copy(_fileBlock.B_content[:], content)

	WriteFileBlock(paht, &_fileBlock, _superBlock.S_first_blo)

	// Actualizando el superbloque
	_superBlock.S_free_blocks_count--
	_superBlock.S_first_blo = _superBlock.S_first_blo + int32(binary.Size(enviroment.FileBlock{}))

	// Escribiendo el superbloque
	WriteSuperBlock(paht, &_superBlock, startPartition)

	// Actualizando el inodo de la carpeta root
	_inode.I_size = _inodeUser.I_size + int32(binary.Size(enviroment.FolderBlock{})) + int32(binary.Size(enviroment.Inode{}))
	WriteInode(paht, &_inode, _superBlock.S_inode_start)

	// Escribiendo el bitmap de inodos
	WriteOneBM(paht, _superBlock.S_bm_inode_start)
	WriteOneBM(paht, _superBlock.S_bm_inode_start+1)

	// Escribiendo el bitmap de bloques
	WriteOneBM(paht, _superBlock.S_bm_block_start)
	WriteOneBM(paht, _superBlock.S_bm_block_start+1)

	enviroment.Message("- Formateo EXT2 realizado con exito")

}

func WriteSuperBlock(path string, superBlock *enviroment.SuperBlock, pos int) {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", n)

	err = binary.Write(disk, binary.LittleEndian, superBlock)

	if err != nil {
		enviroment.Error("en escribir el superbloque")
	}

	disk.Close()
}

func ReadSuperBlock(path string, pos int) enviroment.SuperBlock {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", newpos)

	_superBlock := enviroment.SuperBlock{}
	err = binary.Read(disk, binary.LittleEndian, &_superBlock)
	if err != nil {
		enviroment.Error("en leer el Supebloque en la particion")
	}

	disk.Close()

	return _superBlock
}

func WriteInode(path string, inode *enviroment.Inode, pos int32) {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", n)

	err = binary.Write(disk, binary.LittleEndian, inode)

	if err != nil {
		enviroment.Error("en escribir el inodo")
	}

	disk.Close()
}

func ReadInodo(path string, pos int32) enviroment.Inode {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", newpos)

	_inode := enviroment.Inode{}
	err = binary.Read(disk, binary.LittleEndian, &_inode)
	if err != nil {
		enviroment.Error("en leer el Inodo")
	}

	disk.Close()

	return _inode
}

func WriteFolderBlock(path string, folderBlock *enviroment.FolderBlock, pos int32) {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", n)

	err = binary.Write(disk, binary.LittleEndian, folderBlock)

	if err != nil {
		enviroment.Error("en escribir el FolderBlock")
	}

	disk.Close()
}

func ReadFolderBlock(path string, pos int32) enviroment.FolderBlock {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", newpos)

	_folderBlock := enviroment.FolderBlock{}
	err = binary.Read(disk, binary.LittleEndian, &_folderBlock)
	if err != nil {
		enviroment.Error("en leer el FolderBlock")
	}

	disk.Close()

	return _folderBlock
}

func WriteFileBlock(path string, fileBlock *enviroment.FileBlock, pos int32) {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", n)

	err = binary.Write(disk, binary.LittleEndian, fileBlock)

	if err != nil {
		enviroment.Error("en escribir el FileBlock")
	}

	disk.Close()
}

func ReadFileBlock(path string, pos int32) enviroment.FileBlock {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	// fmt.Println("posicion en disco: ", newpos)

	_fileBlock := enviroment.FileBlock{}
	err = binary.Read(disk, binary.LittleEndian, &_fileBlock)
	if err != nil {
		enviroment.Error("en leer el FileBlock")
	}

	disk.Close()

	return _fileBlock
}

func WriteOneBM(path string, pos int32) {

	disk, err := os.OpenFile(string(path), os.O_RDWR, 0777)
	if err != nil {
		enviroment.Error("en abrir disco")
	}
	_, err = disk.Seek(int64(pos), io.SeekCurrent)
	if err != nil {
		enviroment.Error("en encontrar posicion dentro del disco")
	}

	var caracter [1]byte
	copy(caracter[:], "1")
	// fmt.Println("Peso: ", binary.Size(caracter))
	// fmt.Println("posicion en disco escritura uno: ", n2)

	err = binary.Write(disk, binary.LittleEndian, caracter)
	if err != nil {
		enviroment.Error("en escribir el bitmap")
	}

	disk.Close()
}
