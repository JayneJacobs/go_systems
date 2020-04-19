package profilesystem

import (
	_ "encoding/json" //needed for FileInfo
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

// FileInfo is a struct for the File Metadata
type FileInfo struct {
	Name string 		`json:"name"`
	Size int64			`json:"size"`
	Mode os.FileMode	`json:"mode"`
	ModTime time.Time	`json:"mod_time"`
	IsDir bool			`json:"is_dir"`
}

func fileInfoFromInterface(v os.FileInfo) *FileInfo {
		return &FileInfo{
			v.Name(), 
			v.Size(), 
			v.Mode(), 
			v.ModTime(), 
			v.IsDir(),
		}
}


// Node is a struct that describes the node directory
type Node struct {
	FullPath string 	`json:"path"`
	Info *FileInfo 		`json:"info"`
	Children []*Node 	`json:"children"`
	Parent *Node 		`json:"-"`
}

// NewTree takes a root string and returns the reference to Node Struct and an error
func NewTree(root string) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	fmt.Println("In NewTree path absRoot 44: ", absRoot)
	if err != nil {
		fmt.Println("In NewTree path Error 46: ", err)
		return 
	}
	parents := make(map[string]*Node)
	fmt.Println("In NewTree 50")
	walkFunc := func(path string, info os.FileInfo, err error) error  {
		fmt.Println("In walkFunc path: ", path)
		if err != nil {
			fmt.Println("In walkFunc first if error: ", err)
			return err
		}
		parents[path] = &Node{
			FullPath: 	path,
			Info: 		fileInfoFromInterface(info),
			Children: 	make([]*Node, 0),
		}
		fmt.Println("Return from wakFunc: ", err)
		return nil 
	}
	err = filepath.Walk(absRoot, walkFunc)
	if err != nil {
		fmt.Println("In Filepath.Walk first if error: ", err)
		return 
	}

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if exists {
			fmt.Println(" if exists")
			node.Parent = parent
			fmt.Println("After for loop assigned parentt: ", parent)
			parent.Children = append(parent.Children, node)
			return
		}
		result = node
		fmt.Println("In for loop if !exists: ", node)
	}
	// fmt.Println(result, err)
   return 
}

// GetFileSystemTask struct
type GetFileSystemTask struct {
	path string
	ws *websocket.Conn
}

// NewGetFileSystemTask takes a path string and websocket connection 
func NewGetFileSystemTask(path string, ws *websocket.Conn) *GetFileSystemTask  {
	
	return &GetFileSystemTask{path, ws}
}

// Perform provides the values for GetFileSystemTask struct
func (t *GetFileSystemTask ) Perform()  {

	root, err := NewTree(t.path)
	if err != nil {
		fmt.Printf("Error from NewTree in gws: %v\n", err)
	}
	err = t.ws.WriteJSON(root)
	if err != nil {
		fmt.Printf("Error from writeJSON to ws in gws: %v\n", err)
	}
}
