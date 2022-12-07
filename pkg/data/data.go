package data

import (
	"sync"
	"fmt"
	"github.com/gookit/color"
)


type Data struct {
	Name            string
	Path            string
	ContentBytes	[]byte 
	ModifyLock		sync.Mutex
	IsDeleted		bool
}

func PrintResources(dataSlice []*Data) (outString string) {
	outString = "\n\n          Resource List\n" 
	green := color.FgGreen.Render
	for _, data := range dataSlice {
		curString := "          --------------------------------------\n"
		curString += fmt.Sprintf("          Name:%s\n", green(data.Name))
		curString += fmt.Sprintf("          Path:%s\n", green(data.Path))
		curString += fmt.Sprintf("          Content-length:%s\n", green(len(data.ContentBytes)))
		curString += fmt.Sprintf("          Servable:%s\n\n", green(data.IsDeleted))
		outString += curString
	}
	return outString
}