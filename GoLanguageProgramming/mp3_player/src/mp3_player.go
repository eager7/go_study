package main
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"mplayer"
	"manager"
)
/***************************************local var*****************************************/
var lib *manager.MusicManager
var id int = 0
//var ctrl, signal chan int
/*****************************************main func***************************************/
func main() {
	fmt.Println(`
		Enter following commands to control the player:
		lib list -- View the existing music lib
		lib add <name><artist><source><type> -- Add a music to the music lib
		lib remove <name> -- Remove the specified music from the lib
		play <name> -- Play the specified music
		`)

	lib = manager.NewMusicManager()

	r := bufio.NewReader(os.Stdin)

	for{
		fmt.Println("Enter command->")
		rawLine,_,_ := r.ReadLine()
		line := string(rawLine)
		if line == "q" || line == "e"{
			println("Enter Exit Command")
			break
		}

		tokens := strings.Split(line, " ")

		if tokens[0] == "lib"{
			println("manager operate")
			handleLibCommands(tokens)
		}else if tokens[0] == "play"{
			println("play operate")
			handlePlayCommand(tokens)
		}else{
			println("unkonw operate")
		}

	}
}

/*****************************************local func***************************************/
func handleLibCommands(tokens []string){
	switch tokens[1]{
		case "list":
			for i:=0;i < lib.Len();i++{
				e,_ := lib.Get(i)
				fmt.Println(e.Id, e.Name, e.Artist, e.Source, e.Type)
			}
		case "add":
			if len(tokens) == 6{
				id++
				lib.Add(&manager.MusicEntry{strconv.Itoa(id),tokens[2],tokens[3],tokens[4],tokens[5]})
			}else{
				fmt.Println("USAGE: lib add <name><artist><source><type>")
			}
		case "remove":
			if len(tokens) == 3{
				lib.RemoveByName(tokens[2])
			}else{
				fmt.Println("USAGE: lib remove <id>")
			}
		default:
			fmt.Println("unkonw command")
	}
}

func handlePlayCommand(tokens []string){
	if len(tokens) != 2{
		fmt.Println("USAGE: play <name>")
		return
	}
	_,e := lib.Find(tokens[1])
	if e == nil{
		fmt.Println("The music does not exist",tokens[1])
		return
	}
	mplayer.Play(e.Source, e.Type)
}
