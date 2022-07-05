package crawler

import (
	"fmt"
	"github.com/imthaghost/goclone/pkg/parser"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// file extension map for directing files to their proper directory in O(1) time
var (
	extensionDir = map[string]string{
		".css":  "css",
		".js":   "js",
		".jpg":  "imgs",
		".jpeg": "imgs",
		".gif":  "imgs",
		".png":  "imgs",
		".svg":  "imgs",
	}
)

// Extractor visits a link determines if its a page or sublink
// downloads the contents to a correct directory in project folder
// TODO add functionality for determining if page or sublink
func Extractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)

	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	// Closure
	defer resp.Body.Close()
	// file base
	base := parser.URLFilename(link)
	// store the old ext, in special cases the ext is weird ".css?a134fv"
	oldExt := filepath.Ext(base)
	// new file extension
	ext := parser.URLExtension(link)

	// checks if there was a valid extension
	if ext != "" {
		// checks if that extension has a directory path name associated with it
		// from the extensionDir map
		dirPath := extensionDir[ext]
		if dirPath != "" {
			dirPath2 := link[strings.Index(link, parser.GetDomain(link))+len(parser.GetDomain(link)):]
			dirPath3 := strings.Replace(dirPath2, base, "", -1)
			err := CreateMutiDir(projectPath + dirPath3)
			if err != nil {
				return
			}

			if len(dirPath3) > 2 {
				dirPath = dirPath3[1 : len(dirPath3)-1]
			}
			// If extension and path are valid pass to writeFileToPath
			writeFileToPath(projectPath, base, oldExt, ext, dirPath, resp)
		}
	}
}

func writeFileToPath(projectPath, base, oldFileExt, newFileExt, fileDir string, resp *http.Response) {
	var name = base[0 : len(base)-len(oldFileExt)]
	document := name + newFileExt

	// get the project name and path we use the path to
	f, err := os.OpenFile(projectPath+"/"+fileDir+"/"+document, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	f.Write(htmlData)
}

// CreateMutiDir 调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
