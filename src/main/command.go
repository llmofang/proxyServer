package main
import(
	"os"
	"log"
)

type Command struct {
	path string
}

func (c *Command) cd(path string){
	c.path = path
}

func (c *Command)  touch(name string,perm os.FileMode) string{
	_p := c.path+name
	file, err := os.Open(_p) 
	if err != nil {
		newfile, err := os.Create(_p)
		if err != nil {
			log.Fatal("Error Occured When  Create New File: ", err)
		}		
		newfile.Chmod(640)
		defer newfile.Close()
	}
	defer file.Close()
	return _p
}

func (c *Command)  mkdir(path string,perm os.FileMode){
	err :=os.Mkdir(path,perm) //750
	if err != nil {
		if !os.IsExist(err){
			log.Fatal("Error Occured When Create New Dir: ", err)
		}
	}	
}
func (c *Command)  pwd()  string{
	return c.path
}