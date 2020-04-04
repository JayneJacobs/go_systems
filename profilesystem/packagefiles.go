package profilesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

// FileInfo is a struct for the File Metadata
type FileInfo struct {
	Name string 		`json:"name"`
	Size int64 		`json:"size"`
	Mode os.FileMode	`json:"mode"`
	ModTime time.Time 	`json:"mod_time"`
	IsDir bool 			`json:"is_dir"`
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
	if err != nil {
		return nil, err
	}
	parents := make(map[string]*Node)
	
	walkFunc := func(path string, info os.FileInfo, err error) error  {
		if err != nil {
			return err
		}
		parents[path] = &Node{
			FullPath: 	path,
			Info: 		fileInfoFromInterface(info),
			Children: 	make([]*Node, 0),
		}
		return nil 
	}
	err = filepath.Walk(absRoot, walkFunc)
	if err != nil {
		return 
	}

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists {
			result = node
		}
		node.Parent = parent
		parent.Children = append(parent.Children, node)
	}
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
