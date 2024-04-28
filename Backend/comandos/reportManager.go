package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

// acá falta corregir y crear la imagen del dot y reporte mbr
func Rep(name string, path string, id string, ruta string) {

	if name == "" {
		enviroment.Error("en -name, no se ha ingresado el nombre del reporte a generar")
		return
	}

	if path == "" {
		enviroment.Error("en -path, no se ha ingresado la ruta del reporte a generar")
		return
	}

	if id == "" {
		enviroment.Error("en -id, no se ha ingresado el id del reporte a generar")
		return
	}

	name = strings.ToLower(name)
	id = strings.ToUpper(id)
	enviroment.Command("rep")
	fmt.Println(enviroment.BBLUE + "\n   rep -path=" + path + " -name=" + name + " -id=" + id + " -ruta=" + ruta + enviroment.DEFAULT)
	enviroment.ContentConsola += "\n   rep -path=" + path + " -name=" + name + " -id=" + id + " -ruta=" + ruta + "\n"

	//if !userLoggedExist() {
	//	enviroment.Advertencia("Debe iniciar sesión con algun usuario para poder visualizar los reportes")
	//}

	reportHandler(name, path, id, ruta)
}

func reportHandler(name string, path string, id string, ruta string) {

	if strings.Contains(name, "disk") {
		fmt.Println()
		reportDisk(path, id)
	} else if strings.Contains(name, "tree") {

		fmt.Println()
		reportTree(path, id)
	} else if strings.Contains(name, "file") {

		fmt.Println()
		reportFile(ruta, id, path)
	} else if strings.Contains(name, "sb") {

		fmt.Println()
		reportSuperBlock(path, id)
	} else {

		enviroment.Error("El reporte pedido no existe")
	}
}

func reportDisk(path2 string, id string) {

	index := IdExist(id)

	if index == -1 {
		enviroment.Error("No se pudo generar el reporte, el id no esta asociada a una particion montada existe")
		return
	}

	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	mbrInfo := ReadMBR(path)

	tem := strings.LastIndex(path, "/")
	nameDisk := path[tem+1:]

	dot := "digraph structs { label=\"" + nameDisk + "\" fontsize=20; \nnode [shape=plaintext]; \nstruct1 [label= <<TABLE CELLSPACING=\"3\"> \n<TR> \n<TD COLOR=\"#c23616\" BGCOLOR=\"#ff6b6b\">MBR\n</TD>"

	start := len(enviroment.StructToBytes(enviroment.MBR{}))
	end := -1
	espacioAnterior := false
	freeSpace := 0

	for i := 0; i < 4; i++ {

		ocupado := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])

		if ocupado != 0 {

			if espacioAnterior {

				end = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:])
				freeSpace = end - start

				porcentaje := (float64(freeSpace) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

				if freeSpace != 0 {

					dot += "<TD>Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
				}

				espacioAnterior = false
				start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
				end = -1
			} else {

				start = enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
				end = -1
			}

			porcentaje := (float64(enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

			if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "P" {

				dot += "<TD COLOR=\"#01a3a4\" BGCOLOR=\"#00d2d3\">Primaria<BR/><FONT POINT-SIZE=\"10\">" + enviroment.ByteToStr(mbrInfo.Mbr_partitions[i].Part_name[:]) + "</FONT><BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
			} else if string(mbrInfo.Mbr_partitions[i].Part_type[:]) == "E" {

				dot += "<TD BORDER=\"0\"  CELLPADDING=\"0\">\n <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n <TR>\n <TD COLSPAN=\"25\" COLOR=\"#27800e\" BGCOLOR=\"#65d446\">Extendida<BR/><FONT POINT-SIZE=\"10\">" + enviroment.ByteToStr(mbrInfo.Mbr_partitions[i].Part_name[:]) + "</FONT></TD> \n</TR>\n"
				dot += generarDotParticionesLogicas(path, i)
				dot += "\n</TABLE>\n </TD>\n"
			}

			if i < 3 {

				if enviroment.ByteToInt(mbrInfo.Mbr_partitions[i+1].Part_size[:]) != 0 {

					posFinalActual := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_start[:]) + enviroment.ByteToInt(mbrInfo.Mbr_partitions[i].Part_size[:])
					posInicialSiguiente := enviroment.ByteToInt(mbrInfo.Mbr_partitions[i+1].Part_start[:])

					if posFinalActual != posInicialSiguiente {

						spaceLength := posInicialSiguiente - posFinalActual
						porcentaje := (float64(spaceLength) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

						if spaceLength != 0 {

							dot += "<TD>Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
						}
					}
				}
			}
		} else {

			espacioAnterior = true
			end = -1
		}
	}

	if end == -1 {

		freeSpace = enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]) - start
		porcentaje := (float64(freeSpace) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

		if freeSpace != 0 {

			dot += "<TD>Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
		}
	}

	dot += "</TR>\n</TABLE>>]; \n}"

	// fmt.Println(dot)

	az := strings.LastIndex(path2, "/")
	nameAz := path2[az+1:]

	t1 := enviroment.ContentGraph{}
	t1.Type = "disk"
	t1.Name = nameAz
	t1.Path = path2
	t1.Graph = dot

	enviroment.ArrGraph = append(enviroment.ArrGraph, t1)
	txtTemp(dot, path2)
	enviroment.Message("- Reporte generado con exito")
}

func IdExist(id string) int {

	id = strings.ToUpper(id)

	for index, mounted := range enviroment.MountedPartitionsList {

		if strings.Contains(mounted.IdDisk, id) {
			return index
		}
	}
	return -1
}

func generarDotParticionesLogicas(path string, indexPartition int) string {

	dotAux := "<TR>\n"

	mbrInfo := ReadMBR(path)

	startExtended := 0
	sizeExtended := 0

	startExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[indexPartition].Part_start[:])
	sizeExtended = enviroment.ByteToInt(mbrInfo.Mbr_partitions[indexPartition].Part_size[:])

	_ebr := ReadEBR(path, startExtended)

	end := startExtended + sizeExtended
	ftell := 0

	if enviroment.ByteToInt(_ebr.Part_size[:]) == 0 && enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

		dotAux += "<TD ROWSPAN=\"1\" COLOR=\"#a36707\" BGCOLOR=\"#fcbf5b\">EBR</TD>\n"

		porcentaje := (float64(enviroment.ByteToInt(mbrInfo.Mbr_partitions[indexPartition].Part_size[:])) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

		dotAux += "<TD ROWSPAN=\"1\">Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
		dotAux += "</TR>\n"

		return dotAux
	} else {

		for end > ftell && enviroment.ByteToInt(_ebr.Part_next[:]) != -1 {

			dotAux += "<TD ROWSPAN=\"1\" COLOR=\"#a36707\" BGCOLOR=\"#fcbf5b\">EBR</TD>\n"

			porcentaje := (float64(enviroment.ByteToInt(_ebr.Part_size[:])) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

			if enviroment.ByteToInt(_ebr.Part_size[:]) != 0 {

				dotAux += "<TD ROWSPAN=\"1\" BGCOLOR=\"#c8d6e5\">Logica<BR/><FONT POINT-SIZE=\"10\">" + enviroment.ByteToStr(_ebr.Part_name[:]) + "</FONT><BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>"
			}

			if enviroment.ByteToInt(_ebr.Part_start[:])+enviroment.ByteToInt(_ebr.Part_size[:]) != enviroment.ByteToInt(_ebr.Part_next[:]) {

				freeSpace := enviroment.ByteToInt(_ebr.Part_next[:]) - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))
				porcentaje := (float64(freeSpace) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

				dotAux += "<TD ROWSPAN=\"1\">Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
			}

			_ebr = ReadEBR(path, enviroment.ByteToInt(_ebr.Part_next[:]))
			ftell = enviroment.ByteToInt(_ebr.Part_start[:]) + len(enviroment.StructToBytes(enviroment.EBR{}))

			if enviroment.ByteToInt(_ebr.Part_next[:]) == -1 {

				break
			}
		}
	}

	dotAux += "<TD ROWSPAN=\"1\" COLOR=\"#a36707\" BGCOLOR=\"#fcbf5b\">EBR</TD>\n"

	porcentaje := (float64(enviroment.ByteToInt(_ebr.Part_size[:])) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

	dotAux += "<TD ROWSPAN=\"1\" BGCOLOR=\"#c8d6e5\">Logica<BR/><FONT POINT-SIZE=\"10\">" + enviroment.ByteToStr(_ebr.Part_name[:]) + "</FONT><BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>"

	if enviroment.ByteToInt(_ebr.Part_start[:])+enviroment.ByteToInt(_ebr.Part_size[:]) < end {

		freeSpace := end - (enviroment.ByteToInt(_ebr.Part_start[:]) + enviroment.ByteToInt(_ebr.Part_size[:]))
		porcentaje := (float64(freeSpace) / float64(enviroment.ByteToInt(mbrInfo.Mbr_tamano[:]))) * 100

		dotAux += "<TD ROWSPAN=\"1\">Libre<BR/><FONT POINT-SIZE=\"8\">" + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "% del disco</FONT></TD>\n"
	}

	dotAux += "</TR>\n"
	return dotAux

}

func reportSuperBlock(path2 string, id string) {

	index := IdExist(id)

	if index == -1 {
		enviroment.Error("No se pudo generar el reporte, el id no esta asociada a una particion montada existe")
		return
	}

	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	tem := strings.LastIndex(path, "/")
	nameDisk := path[tem+1:]

	dot := "\ndigraph structs { label=\"" + mounted.IdDisk + " -> " + mounted.Name + "\" fontsize=20; node [shape=plaintext]; struct1 [label= <<TABLE CELLSPACING=\"1\"> <TR> <TD BORDER=\"0\"  CELLPADDING=\"1\"> <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\"> <TR> <TD COLSPAN=\"25\" BGCOLOR=\"#d35400\"><FONT COLOR=\"white\"><B>REPORTE DE SUPERBLOQUE</B></FONT></TD> </TR>"

	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">sb_nombre_hd</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + nameDisk + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_filesystem_type</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_filesystem_type) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_inodes_count</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_inodes_count) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_blocks_count</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_blocks_count) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_free_blocks_count</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_free_blocks_count) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_free_inodes_count</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_free_inodes_count) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_mtime</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + enviroment.ByteToStr(_superBlock.S_mtime[:]) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_umtime</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + enviroment.ByteToStr(_superBlock.S_umtime[:]) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_mnt_count</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_mnt_count) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_magic</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_magic) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_inode_s</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_inode_s) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_block_s</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_block_s) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_first_ino</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_first_ino) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_first_blo</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_first_blo) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_bm_inode_start</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_bm_inode_start) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_bm_block_start</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_bm_block_start) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_inode_start</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_inode_start) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">s_block_start</TD> <TD ALIGN=\"LEFT\" ROWSPAN=\"1\" BGCOLOR=\"#f5ba5f\">" + fmt.Sprint(_superBlock.S_block_start) + "</TD> </TR>"

	dot += "\n</TABLE></TD></TR></TABLE>>];}"

	// fmt.Println(dot)
	az := strings.LastIndex(path2, "/")
	nameAz := path2[az+1:]

	t1 := enviroment.ContentGraph{}
	t1.Type = "sb"
	t1.Name = nameAz
	t1.Path = path2
	t1.Graph = dot

	enviroment.ArrGraph = append(enviroment.ArrGraph, t1)
	txtTemp(dot, path2)
	enviroment.Message("- Reporte generado con exito")
}

func reportTree(path2 string, id string) {

	index := IdExist(id)

	if index == -1 {
		enviroment.Error("No se pudo generar el reporte, el id no esta asociada a una particion montada existe")
		return
	}

	mounted := enviroment.MountedPartitionsList[index]
	path := mounted.Path

	_superBlock := ReadSuperBlock(path, mounted.Start)

	cont := 0
	cont2 := 0

	dot := "digraph structs { label=\"Reporte TREE\" fontsize=20; rankdir=LR; splines=false; node [shape=plaintext];"
	dotEdge := ""

	_inode := ReadInodo(path, _superBlock.S_inode_start)

	_reportTree(_inode, path, _superBlock.S_inode_start, &cont, &cont2, &dot, &dotEdge)

	dot += dotEdge
	dot += "\n}"

	// fmt.Println(dot)

	az := strings.LastIndex(path2, "/")
	nameAz := path2[az+1:]

	t1 := enviroment.ContentGraph{}
	t1.Type = "tree"
	t1.Name = nameAz
	t1.Path = path2
	t1.Graph = dot

	enviroment.ArrGraph = append(enviroment.ArrGraph, t1)

	txtTemp(dot, path2)
	enviroment.Message("- Reporte generado con exito")
}

func _reportTree(_inode enviroment.Inode, path string, index int32, cont *int, cont2 *int, dot *string, dotEdge *string) {

	if strings.Contains(string(_inode.I_type[:]), "0") {

		*dot += printInode(_inode, index, cont, dotEdge)
		for i := 0; i < 16; i++ {

			if _inode.I_block[i] != -1 {

				_folderBlock := ReadFolderBlock(path, _inode.I_block[i])
				*dot += printFolderBlock(_folderBlock, _inode.I_block[i], cont2, dotEdge)

				for j := 0; j < 4; j++ {

					if j > 1 {

						if _folderBlock.B_content[j].B_inodo != -1 {

							_inode2 := ReadInodo(path, _folderBlock.B_content[j].B_inodo)
							_reportTree(_inode2, path, _folderBlock.B_content[j].B_inodo, cont, cont2, dot, dotEdge)
						}
					}
				}
			}
		}
	} else if strings.Contains(string(_inode.I_type[:]), "1") {

		*dot += printInode(_inode, index, cont, dotEdge)

		for i := 0; i < 16; i++ {

			if _inode.I_block[i] != -1 {

				_fileBlock := ReadFileBlock(path, _inode.I_block[i])
				*dot += printFileBlock(_fileBlock, _inode.I_block[i], cont2, dotEdge)
			}
		}
	}
}

func printInode(_inode enviroment.Inode, index int32, cont *int, dotEdge *string) string {

	dot := ""
	dot += "s" + fmt.Sprint(index) + " [label= <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"> <TR> <TD BORDER=\"0\"  CELLPADDING=\"1\"> <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n <TR> <TD COLSPAN=\"25\" BGCOLOR=\"#044875\"><FONT COLOR=\"white\"><B>INODO " + fmt.Sprint(*cont) + "</B></FONT></TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_uid</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + fmt.Sprint(_inode.I_uid) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_gid</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + fmt.Sprint(_inode.I_gid) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_s</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + fmt.Sprint(_inode.I_size) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_atime</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + enviroment.ByteToStr(_inode.I_atime[:]) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_ctime</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + enviroment.ByteToStr(_inode.I_ctime[:]) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_ctime</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + enviroment.ByteToStr(_inode.I_mtime[:]) + "</TD> </TR>"

	for aux := 0; aux < 16; aux++ {

		dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_block_" + fmt.Sprint(aux+1) + "</TD>\n <TD PORT=\"" + fmt.Sprint(_inode.I_block[aux]) + "\" ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + fmt.Sprint(_inode.I_block[aux]) + "</TD> </TR>"
		if _inode.I_block[aux] != -1 {

			*dotEdge += "\ns" + fmt.Sprint(index) + ":" + fmt.Sprint(_inode.I_block[aux]) + " -> s" + fmt.Sprint(_inode.I_block[aux]) + ";"
		}
	}

	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_type</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + enviroment.ByteToStr(_inode.I_type[:]) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">i_perm</TD>\n <TD ROWSPAN=\"1\" BGCOLOR=\"#9fd3f5\">" + fmt.Sprint(_inode.I_perm) + "</TD> </TR>"
	dot += "\n</TABLE></TD></TR></TABLE>>];"

	*cont = *cont + 1

	return dot
}

func printFolderBlock(_folderBlock enviroment.FolderBlock, index int32, cont *int, dotEdge *string) string {

	dot := ""
	dot += "s" + fmt.Sprint(index) + " [label= <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"> <TR> <TD BORDER=\"0\"  CELLPADDING=\"1\"> <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\"> "
	dot += "\n<TR> <TD COLSPAN=\"25\" BGCOLOR=\"#08963e\"><FONT COLOR=\"white\"><B>Folder Block " + fmt.Sprint(*cont) + "</B></FONT></TD> </TR> <TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#1fbf5c\">"
	dot += "\n<FONT COLOR=\"white\">Name</FONT></TD> <TD ROWSPAN=\"1\" BGCOLOR=\"#1fbf5c\"><FONT COLOR=\"white\">Inodo</FONT></TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#b3f5cc\">.</TD> <TD ROWSPAN=\"1\">" + fmt.Sprint(_folderBlock.B_content[0].B_inodo) + "</TD> </TR>"
	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#b3f5cc\">..</TD> <TD ROWSPAN=\"1\">" + fmt.Sprint(_folderBlock.B_content[1].B_inodo) + "</TD> </TR>"

	dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#b3f5cc\">" + enviroment.ByteToStr(_folderBlock.B_content[2].B_name[:]) + "</TD> <TD PORT=\"" + fmt.Sprint(_folderBlock.B_content[2].B_inodo) + "\" ROWSPAN=\"1\">" + fmt.Sprint(_folderBlock.B_content[2].B_inodo) + "</TD> </TR>"
	*dotEdge += "\ns" + fmt.Sprint(index) + ":" + fmt.Sprint(_folderBlock.B_content[2].B_inodo) + " -> s" + fmt.Sprint(_folderBlock.B_content[2].B_inodo) + ";"

	if strings.Contains(string(_folderBlock.B_content[3].B_name[:]), "-") {

		dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#b3f5cc\"> - </TD> <TD ROWSPAN=\"1\">" + fmt.Sprint(_folderBlock.B_content[3].B_inodo) + "</TD> </TR>"
	} else {

		dot += "\n<TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#b3f5cc\">" + enviroment.ByteToStr(_folderBlock.B_content[3].B_name[:]) + "</TD> <TD PORT=\"" + fmt.Sprint(_folderBlock.B_content[3].B_inodo) + "\" ROWSPAN=\"1\">" + fmt.Sprint(_folderBlock.B_content[3].B_inodo) + "</TD> </TR>"
		*dotEdge += "\ns" + fmt.Sprint(index) + ":" + fmt.Sprint(_folderBlock.B_content[3].B_inodo) + " -> s" + fmt.Sprint(_folderBlock.B_content[3].B_inodo) + ";"
	}

	dot += "\n</TABLE></TD></TR></TABLE>>];"

	*cont = *cont + 1

	return dot

}

func printFileBlock(_fileBlock enviroment.FileBlock, index int32, cont *int, dotEdge *string) string {

	dot := ""

	dot += "\ns" + fmt.Sprint(index) + " [label= <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"> <TR> <TD BORDER=\"0\"  CELLPADDING=\"1\"> <TABLE BORDER=\"0\" CELLBORDER=\"1\" "
	dot += "\nCELLSPACING=\"0\" CELLPADDING=\"4\"> <TR> <TD COLSPAN=\"25\" BGCOLOR=\"#141414\"><FONT COLOR=\"white\"><B>File Block " + fmt.Sprint(*cont) + "</B>"
	dot += "\n</FONT></TD> </TR> <TR> <TD ROWSPAN=\"1\" BGCOLOR=\"#4d4d4d\"><FONT COLOR=\"white\">Contenido</FONT></TD>"

	temp := enviroment.ByteToStr(_fileBlock.B_content[:])

	temp = strings.ReplaceAll(temp, "\n", "<BR/>")

	dot += "\n<TD ROWSPAN=\"1\">" + temp + "</TD> </TR>"

	dot += "\n</TABLE></TD></TR></TABLE>>];"

	*cont = *cont + 1

	return dot
}

func txtTemp(contenido string, path string) {

	r2 := strings.LastIndex(path, ".")
	ruta := path[:r2] + ".txt"

	fmt.Print(ruta)
	CreatePaths(ruta)

	archivo, err := os.Create(ruta)
	if err != nil {
		enviroment.Error("No se pudo crear el reporte")
		return
	}

	_, err = archivo.WriteString(contenido)
	if err != nil {
		enviroment.Error("No se pudo escribir el reporte")
		return
	}
	archivo.Close()

}

func reportFile(path2 string, id string, path3 string) {

	if path2 != "" {

		index := IdExist(id)

		if index == -1 {
			enviroment.Error("No se pudo generar el reporte, el id no esta asociada a una particion montada existente")
			return
		}

		carpetas := strings.Split(path2, "/")

		if strings.Contains(carpetas[1], "users.txt") {

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

			dotFile(data, path3)

		} else {

			if !isDirExist(path2) {
				a := strings.LastIndex(path2, "/")
				enviroment.Error("El archivo no existe: " + path2[a+1:])
				return
			}

			var indexInode int32 = -1

			carpetas := strings.Split(path2, "/")

			for i := 1; i < len(carpetas); i++ {

				if i == 1 {

					indexInode = searchRoot(carpetas[i])

					if indexInode == -1 {

						enviroment.Error("No existe el archivo " + carpetas[i])
						break
					}
				} else {

					indexInode = searchFolder(carpetas[i], indexInode)

					if indexInode == -1 {

						enviroment.Error("No existe el archivo: " + carpetas[i] + " en la carpeta: " + carpetas[i-1])
						break
					}
				}
			}

			if indexInode != -1 {

				mounted := enviroment.MountedPartitionsList[index]
				path := mounted.Path

				_inode := ReadInodo(path, indexInode)

				data := ""

				for i := 0; i < 16; i++ {

					if _inode.I_block[i] != -1 {

						_fileBlock := ReadFileBlock(path, _inode.I_block[i])
						data += enviroment.ByteToStr(_fileBlock.B_content[:])
					}
				}

				dotFile(data, path3)

			}

		}

	} else {

		enviroment.Error("en >path debe especificar una ruta valida")
	}
}

func dotFile(data string, path string) {

	aux := strings.LastIndex(path, "/")
	nameFile := path[aux+1:]

	temp := data
	data = strings.ReplaceAll(data, "\n", "\\n")

	dot := "digraph { label=\"" + nameFile + "\""
	dot += "\n	node [shape=note, style=filled, fillcolor = \"#00d2d3\", color=\"#01a3a4\" penwidth=2.5];"
	dot += "\n   \"a\" [label=\"" + data + "\", align=justify]; \n}"

	fmt.Println(dot)

	t1 := enviroment.ContentGraph{}
	t1.Type = "file"
	t1.Name = nameFile
	t1.Path = path
	t1.Graph = dot

	enviroment.ArrGraph = append(enviroment.ArrGraph, t1)

	txtTemp(temp, path)
	enviroment.Message("- Reporte generado con exito")
}
