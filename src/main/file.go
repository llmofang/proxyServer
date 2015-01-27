package main
import(
	"os"
	"log"
)

type File struct {
	path string
}

func (f *File) cd(path string){
	f.path = path
}

func (f *File)  touch(name string,perm os.FileMode) string{
	_p := f.path+name
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

func (f *File)  mkdir(path string,perm os.FileMode){
	err :=os.Mkdir(path,perm) //750
	if err != nil {
		if !os.IsExist(err){
			log.Fatal("Error Occured When Create New Dir: ", err)
		}
	}	
}
func (f *File)  pwd() {

}