package main

//TODO
//
// more sys"call"s
// compile mode
// make it gc.py compatible (it is, but use include)

include "gutil"
import (
	//"strings"
	"os"
	"math"
)


func main() {
	var argv []string = os.Args
	var argc int = len(argv)
	var compile bool
	var flname = "program.e"
	if (argc == 2){
		flname = argv[1]
		compile = false
	} else if (argc == 3){
		flname = argv[1]
		compile = (argv[2] == "true" || argv[2] == "1")
	}
	var file []string = split(ReadFile(flname), "\n")
	var BinFile [][]bool
	var Ret []bool
	var char string
	for linei := 0 ; linei < len(file) ; linei++{
		Ret = []bool{}
		for chari := 0 ; chari < len(file[linei]) ; chari++ {
			char = string(file[linei][chari])
			if char == "e" {
				Ret = append(Ret, false)
			} else if char == "E" {
				Ret = append(Ret, true)
			} else if char == "#" {
				break
			}
		}
		if ( len(Ret) > 0 ){
			BinFile = append(BinFile, Ret)
		}
	}
	os.Exit(Execute(BinFile))
	print(compile)
}


func MakeBin( line []bool ) int {
	var rt int = 0
	var ll float64 = float64(len(line)-1)
	var mul int = int(math.Pow(2, ll))
	for i := 0 ; i < len(line) ; i++ {
		if line[i] {
			rt+=mul
		}
		mul /= 2
	}
	return rt
}

func pop(stack []int) ([]int, int) {
	var ls = len(stack)-1
	var ret = stack[ls]
	var stackr []int = stack[:ls]
	return stackr, ret
}

func Execute( file [][]bool) int {
	var line []bool
	var stack = []int{}
	var op string
	var new int
	var mem int
	var scall int
	var zero = false
	for i := 0 ; i < len(file) ; i++{
		line = file[i]
		if len(line) != 0 {
			if line[0] == zero { // var stuff
				if len(line) == 1 { // set var 0
					stack = append(stack, 0)
					op = sprintf("append 0")
				} else if line[1] == zero { // +, -
					if line[2] == zero {
						if line[3] == zero {
							stack, mem = pop(stack)
							stack, new = pop(stack)
							op = sprintf("add %d to %d", mem, new)
							stack = append(stack, new+mem)
						} else {
							stack, mem = pop(stack)
							stack, new = pop(stack)
							op = sprintf("sub %d to %d", mem, new)
							stack = append(stack, new-mem)
						}
					} else { // *, /
						if line[3] == zero {
							stack, mem = pop(stack)
							stack, new = pop(stack)
							op = sprintf("mul %d to %d", mem, new)
							stack = append(stack, new*mem)
						} else {
							stack, mem = pop(stack)
							stack, new = pop(stack)
							op = sprintf("div %d to %d", mem, new)
							stack = append(stack, new/mem)
						}
					}
				} else { // set var bin
					stack = append(stack, MakeBin(line[1:]))
					op = sprintf("append %d", stack[len(stack)-1])
				}

			} else { // call
				if len(line) != 1 {
					scall = MakeBin(line[1:])
					switch scall{
						case 1:
							stack, mem = pop(stack)
							return mem
							stdout.Flush()
						case 3:
							stack, mem = pop(stack)
							stdout.Write([]byte{byte(mem)})
							op = sprintf("write %d, \"%s\"", mem, string(byte(mem)))
							if mem == 10 {
								op+=" and flush"
								stdout.Flush()
							}
					}
				} else { // dup
					stack, mem = pop(stack)
					stack = append(stack, mem)
					stack = append(stack, mem)
					op = sprintf("dup %d", mem)
				}
			}
		}
	}
	stdout.Flush()
	return 0
}
