package analizador

import (
	"bufio"
	"fmt"
	"os"
	"paquetes/comandos"
	"paquetes/enviroment"
	"strconv"
	"strings"
)

var path, fit, unit, driveletterr, typee, name, id, user, pass, grp, cont, fs, ruta string
var size int
var r bool
var ruta1 string = "/home/pjd/archivos/"
var contadorLetras = 0

func InputHandler() {
	var activo bool = true
	// Create a new reader.
	reader := bufio.NewReader(os.Stdin)

	// Loop that reads the input.
	for activo {

		var linea string
		fmt.Print("\n" + enviroment.BRED + " ~- " + enviroment.DEFAULT)
		linea, _ = reader.ReadString('\n')

		// Remove the new line character.
		linea = limpiarSaltosLinea(linea)
		fmt.Println()

		if strings.Compare(linea, "exit") == 0 {
			activo = false
		} else if strings.Compare(linea, "") == 0 {
			// Do nothing.
		} else {
			linea = quitarEspaciosLinea(linea)
			dataLine := strings.Split(linea, " ")
			if dataLine[0] == "execute" {
				fmt.Println("Analizar archivo: " + dataLine[1])
				analizarArchivo(dataLine[1])
			} else {
				// fmt.Println("Mandar a analizar: " + linea)
				analizarLinea(linea)
			}
		}
	}
}

func quitarEspaciosLinea(linea string) string {
	linea = strings.TrimLeft(linea, " ")
	linea = strings.TrimRight(linea, " ")
	return linea
}

func analizarLinea(linea string) {

	if strings.Contains(linea, "#") {
		fmt.Println("\n   " + linea)
		fmt.Println()
		enviroment.ContentConsola += "\n\n   " + linea + "\n\n"
		return
	}

	tokens := strings.Split(linea, " -")
	command := strings.ToLower(tokens[0])

	switch command {

	case "mkdisk":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			letraAsignada := string(rune(65 + contadorLetras))
			path = ruta1 + letraAsignada + ".dsk"
			//path := "/home/pjd/archivos" + ".dsk"
			comandos.Mkdisk(path, size, unit, fit)
			contadorLetras++
			limpiarVariables()
		}
	case "rmdisk":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			path = ruta1 + driveletterr + ".dsk"
			comandos.Rmdisk(path)
			limpiarVariables()
		}
	case "fdisk":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			path = ruta1 + driveletterr + ".dsk"
			comandos.Fdisk(driveletterr, path, size, unit, fit, typee, name)
			// ^------------------------------------------------ DEBUGER MBR -------------------------------------------------
			// comandos.ShowMBRInfo(path)
			limpiarVariables()
		}
	case "mount":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			path = ruta1 + driveletterr + ".dsk"
			comandos.Mount(path, name)
			limpiarVariables()
		}
	case "unmount":
		fmt.Println("Ejecutando comando: " + command)
		getParams(tokens)

	case "mkfs":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Mkfs(id, typee, fs)
			limpiarVariables()
		}
	case "login":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Login(user, pass, id)
			limpiarVariables()
		}
	case "logout":
		comandos.Logout()
		limpiarVariables()

	case "mkgrp":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Mkgrp(name)
			limpiarVariables()
		}
	case "rmgrp":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Rmgrp(name)
			limpiarVariables()
		}
	case "mkusr":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Mkusr(user, pass, grp)
			limpiarVariables()
		}
	case "rmusr":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Rmusr(user)
			limpiarVariables()
		}
	case "mkfile":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Mkfile(path, size, cont, r)
			limpiarVariables()
		}
	case "mkdir":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Mkdir(path, r)
			limpiarVariables()
		}
	case "rep":
		vvalida := getParams(tokens)
		if vvalida != "Error" {
			comandos.Rep(name, path, id, ruta)
			limpiarVariables()
		}
	case "pause":
		fmt.Println("Ejecutando comando: " + command)
		enviroment.Command("pause")
		// fmt.Println(enviroment.ContentConsola)

	default:
		enviroment.Error("Comando no encontrado: " + command)
		// fmt.Println("Error: Comando no encontrado " + command)
	}

	fmt.Println()
}

func getParams(tokens []string) string {
	for i := 1; i < len(tokens); i++ {
		param := quitarEspaciosLinea(tokens[i])
		atributes := strings.Split(param, "=")

		param = strings.ToLower(atributes[0])

		if param == "r" {
			r = true
			// fmt.Print("   " + param)
			// fmt.Println(": ", r)
			continue
		}

		switch param {

		case "path":
			// fmt.Print("   " + param)
			path = strings.Trim(atributes[1], "\"")
			path = quitarEspaciosLinea(path)
			path = limpiarSaltosLinea(path)
			// fmt.Println(": " + path)

		case "size":
			// fmt.Print("   " + param)
			tmp := strings.Trim(atributes[1], "\"")
			tmp = quitarEspaciosLinea(tmp)
			tmp = limpiarSaltosLinea(tmp)
			size, _ = strconv.Atoi(tmp)
			// fmt.Println(":", size)

		case "fit":
			// fmt.Print("   " + param)
			fit = strings.ToUpper(atributes[1])
			fit = quitarEspaciosLinea(fit)
			fit = limpiarSaltosLinea(fit)
			// fmt.Println(": " + fit)

		case "unit":
			// fmt.Print("   " + param)
			unit = strings.ToUpper(atributes[1])
			unit = quitarEspaciosLinea(unit)
			unit = limpiarSaltosLinea(unit)
			// fmt.Println(": " + unit)
		case "driveletter":
			// fmt.Print("   " + param)
			driveletterr = strings.ToUpper(atributes[1])
			driveletterr = quitarEspaciosLinea(driveletterr)
			driveletterr = limpiarSaltosLinea(driveletterr)
			// fmt.Println(": " + driveletterr)

		case "type":
			// fmt.Print("   " + param)
			typee = strings.ToUpper(atributes[1])
			typee = quitarEspaciosLinea(typee)
			typee = limpiarSaltosLinea(typee)
			// fmt.Println(": " + typee)

		case "name":
			// fmt.Print("   " + param)
			name = strings.Trim(atributes[1], "\"")
			name = quitarEspaciosLinea(name)
			name = limpiarSaltosLinea(name)
			// fmt.Println(": " + name)

		case "id":
			// fmt.Print("   " + param)
			id = strings.Trim(atributes[1], "\"")
			id = quitarEspaciosLinea(id)
			id = limpiarSaltosLinea(id)
			// fmt.Println(": " + id)

		case "user":
			// fmt.Print("   " + param)
			user = strings.Trim(atributes[1], "\"")
			user = quitarEspaciosLinea(user)
			user = limpiarSaltosLinea(user)
			// fmt.Println(": " + user)

		case "pass":
			// fmt.Print("   " + param)
			pass = strings.Trim(atributes[1], "\"")
			pass = limpiarSaltosLinea(pass)
			pass = quitarEspaciosLinea(pass)
			// fmt.Println(": " + pass)

		case "grp":
			// fmt.Print("   " + param)
			grp = strings.Trim(atributes[1], "\"")
			grp = quitarEspaciosLinea(grp)
			grp = limpiarSaltosLinea(grp)
			// fmt.Println(": " + grp)

		case "fs":
			// fmt.Print("   " + param)
			fs = strings.Trim(atributes[1], "\"")
			fs = quitarEspaciosLinea(fs)
			fs = limpiarSaltosLinea(fs)
			// fmt.Println(": " + grp)

		case "cont":
			// fmt.Print("   " + param)
			cont = strings.Trim(atributes[1], "\"")
			cont = quitarEspaciosLinea(cont)
			cont = limpiarSaltosLinea(cont)
			// fmt.Println(": " + cont)

		case "ruta":
			// fmt.Print("   " + param)
			ruta = strings.Trim(atributes[1], "\"")
			ruta = quitarEspaciosLinea(ruta)
			ruta = limpiarSaltosLinea(ruta)
			// fmt.Println(": " + ruta)

		default:
			enviroment.Error("Parametro incorrecto " + param)
			return "Error"
			// fmt.Println("Error: Parametro incorrecto " + param)
		}
	}
	return ""
}

func analizarArchivo(ruta string) {
	ruta = quitarEspaciosLinea(ruta)
	path := strings.Split(ruta, "=")
	ruta = quitarEspaciosLinea(path[1])
	ruta = strings.Trim(ruta, "\"")

	ruta = limpiarSaltosLinea(ruta)

	file, err := os.Open(ruta)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		return
	}
	var contenido string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contenido += scanner.Text()
		contenido += "\n"
	}

	file.Close()

	// re := regexp.MustCompile("\n+")
	// contenido = re.ReplaceAllString(contenido, "\n")
	// contenido = strings.TrimLeft(contenido, "\n")

	lineas := strings.Split(contenido, "\n")
	for i := 0; i < len(lineas); i++ {
		linea := quitarEspaciosLinea(lineas[i])
		if linea != "" {
			analizarLinea(linea)
			// fmt.Println(linea)
		}
	}
}

func AnalizarPeticion(contenido string) {

	lineas := strings.Split(contenido, "\n")
	for i := 0; i < len(lineas); i++ {
		linea := quitarEspaciosLinea(lineas[i])
		linea = strings.ReplaceAll(linea, "\r", "")
		if linea != "" {
			analizarLinea(linea)
			// fmt.Println(linea)
		}
	}
}

func limpiarVariables() {
	path = ""
	size = 0
	fit = ""
	unit = ""
	typee = ""
	name = ""
	id = ""
	user = ""
	pass = ""
	grp = ""
	cont = ""
	r = false
	ruta = ""
	driveletterr = ""
	fs = ""
}

func limpiarSaltosLinea(cadena string) string {
	cadena = strings.ReplaceAll(cadena, "\n", "")
	cadena = strings.ReplaceAll(cadena, "\r", "")
	return cadena
}
