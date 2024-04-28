package enviroment

// Varialbes de colores
const BRED string = "\033[1;31m"
const BGREEN string = "\033[1;32m"
const BCYN string = "\033[1;36m"
const BBLUE string = "\033[1;34m"
const BORANGE string = "\033[1;30m"
const BWHITE string = "\033[1;37m"
const DEFAULT string = "\033[0m"

var UserLogged_ = UserLogged{}

var MountedPartitionsList = make([]MountedPartitions, 0)
var ContDisk int = 1
var IdentDisk = make(map[string]int)

var ContentConsola string = ""
var ContentLogin string = ""
var ArrGraph = make([]ContentGraph, 0)

var Letras = [...]string {
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}