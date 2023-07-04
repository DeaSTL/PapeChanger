package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/npmaile/papeChanger/linux"
	macos "github.com/npmaile/papeChanger/macOS"
)

const errPrefix = "😭😭oOpSy DoOpSiE, you made a frickey-wickey 😭😭: "

func main() {
	// parse command line arguments
	changeDir := flag.Bool("c", false, "Change the directory you are selecting walpapers from")
	randomize := flag.Bool("r", true, "Randomize wallpaper to change")
	u, err := user.Current()
	if err != nil {
		log.Fatalf("%sHow the H**K are you not logged in as a user?", errPrefix)
	}
	homeDir := u.HomeDir
	stateFile := flag.String("m", path.Join(homeDir, ".local/papeChanger/state"), "Use a custom location to store the current walpaper set")
	flag.Parse()

	currentWalpaper, err := os.ReadFile(*stateFile)
	if err != nil {
		log.Fatalf("%sCan't read the file", errPrefix)
	}
	pathParts := strings.Split(string(currentWalpaper), string(os.PathSeparator))
	currentDirParts := pathParts[0 : len(pathParts)-1]
	if *changeDir {
		megaDir := string(os.PathSeparator) + filepath.Join(currentDirParts[0:len(currentDirParts)-1]...)
		var files []fs.DirEntry
		files, err = os.ReadDir(megaDir)
		if err != nil {
			log.Fatalf("%sYou've moved your walpapers around and I can't find them now: %e", errPrefix, err)
		}
		dirList := []string{}
		for _, file := range files {
			if file.IsDir() {
				dirList = append(dirList, file.Name())
			}
		}
		var chooseFunc func([]string) (string, error)
		switch runtime.GOOS {
		case "darwin":
			chooseFunc = macos.Chooser
		case "linux":
			chooseFunc = linux.Chooser
		default:
			chooseFunc = func([]string) (string, error) {
				log.Fatalf("%sYour os isn't supported (yet)", errPrefix)
				return "", nil
			}
		}
		var chosen string
		chosen, err = chooseFunc(dirList)
		if err != nil {
			log.Fatalf("%sFailed to choose walpaper directory: %e", errPrefix, err)
		}
		currentDirParts[len(currentDirParts)-1] = string(chosen)
	}
	walpaperFolder := string(os.PathSeparator) + filepath.Join(currentDirParts...)
	papers, err := os.ReadDir(walpaperFolder)
	if err != nil {
		log.Fatalf("%sUnable to get list of individual walpapers: %e", errPrefix, err)
	}
	var fullPath []string
	if *randomize {
		index := rand.Int() % len(papers)
		fullPath = append(currentDirParts, papers[index].Name())
	} else {
		//todo
	}

	var changeWalpaperFunc func(string) error
	switch runtime.GOOS {
	case "darwin":
		changeWalpaperFunc = macos.SetPape
	case "linux":
		changeWalpaperFunc = linux.SetPape
	default:
		changeWalpaperFunc = func(string) error {
			return fmt.Errorf("OS not supported (yet)")
		}
	}
	newWalpaper := string(os.PathSeparator) + filepath.Join(fullPath...)
	err = changeWalpaperFunc(newWalpaper)
	if err != nil {
		log.Fatalf("%sunable to change walpaper: %e", errPrefix, err)
	}
	f, err := os.Create(*stateFile)
	if err != nil {
		log.Fatalf("%sCreation of state file failed: %e", errPrefix, err)
	}
	f.Write([]byte(newWalpaper))
}
