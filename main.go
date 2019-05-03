//
// PACKAGES
//
package main

/*
 * IMPORTS
 */
import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

/*
 * CONSTANTS
 */
const (
	AppDesc            AppMetaData = "File change detection tool"
	AppName            AppMetaData = "FileDelta"
	AppVersion         AppMetaData = "0.3.0"
	CLIName            AppMetaData = "filedelta"
	CommandCheck                   = "check"
	CommandStore                   = "store"
	ErrorArgumentError             = 10
	ErrorCacheMissed               = 1
)

/*
 * DERIVED CONSTANTS
 */
var (
	AppLabel = AppMetaData(fmt.Sprintf("%s v%s", string(AppName), string(AppVersion)))
	AppHelp  = AppMetaData(fmt.Sprintf("%s\n\n%s", string(AppLabel), string(AppDesc)) + `

COMMANDS
  check   Compares the provided files hash against the one stored
  store   Stores the provided files hash for later comparison

OPTIONS
  -d | --debug     Output debugging info
  -h | --help      Output this help info
  -v | --version   Output app version info
`)
)

/*
 * TYPES
 */
type (
	// AppMetaData defines meta-data about an application
	AppMetaData string
)

/*
 * VARIABLES
 */
var (
	cacheFile string
	homePath  string
	cachePath string
	cmd       string
	debug     = false
	hashFile  string
)

/*
 * FUNCTIONS
 */
func cacheFilePath(file string) string {
	fileAbs, err := filepath.Abs(file)
	if err == nil {
		file = fileAbs
	}
	return fmt.Sprintf("%s/%s", cachePath, url.QueryEscape(file))
}

func cacheHashGet(file string) (string, error) {
	fileHash, err := ioutil.ReadFile(cacheFilePath(file))
	return string(fileHash), err
}

func cacheHashPut(file, hash string) error {
	cacheFile := cacheFilePath(file)
	return ioutil.WriteFile(cacheFile, []byte(hash), 0644)
}

func fileHashGet(file string) (string, error) {
	var (
		err  error
		f    *os.File
		hash string
	)
	// fmt.Printf("fileHashGet() | file = %q\n", file)
	f, err = os.Open(file)
	if err != nil {
		// fmt.Printf("fileHashGet() | err = %q\n", err)
		return hash, err
	}

	h := sha256.New()
	_, err = io.Copy(h, f)
	f.Close()

	hash = fmt.Sprintf("%x", h.Sum(nil))
	// fmt.Printf("fileHashGet() | hash = %q\n", hash)

	return hash, err
}

func init() {
	var skip int
	for i, a := range os.Args {
		if i > 0 && i != skip {
			switch a {
			case CommandStore:
				cmd = CommandStore
				skip = i + 1
				hashFile = os.Args[skip]
			case CommandCheck:
				cmd = CommandCheck
				skip = i + 1
				hashFile = os.Args[skip]
			case "-c", "--cache":
				skip = i + 1
				cachePath = os.Args[skip]
			case "-d", "--debug":
				debug = true
			case "-v", "--version":
				fmt.Println(AppLabel)
				os.Exit(0)
			case "-h", "--help":
				fmt.Println(AppHelp)
				os.Exit(0)
			default:
				hashFile = a
			}
		}
	}

	homePath, _ = homedir.Dir()
	cachePath = fmt.Sprintf("%s/.local/filedelta/cache", homePath)

	_, err := os.Stat(cachePath)
	if err != nil {
		os.MkdirAll(cachePath, 0700)
	}

	if len(hashFile) > 0 {
		cacheFile = url.QueryEscape(hashFile)
	}
}

/*
 * MAIN ENTRYPOINT
 */
func main() {
	if len(hashFile) > 0 {
		hexHash, _ := fileHashGet(hashFile)

		switch cmd {
		case CommandStore:
			if debug {
				fmt.Printf("cache: %s\n", hashFile)
				fmt.Printf("cache: %s\n", cacheFile)
				fmt.Printf("cache: %s\n", hexHash)
			}
			fmt.Printf("%s *%s\n", hexHash, hashFile)
			cacheHashPut(hashFile, hexHash)
		case CommandCheck:
			cachedHash, _ := cacheHashGet(hashFile)
			if debug {
				line := fmt.Sprintf("%s *%s\n", hexHash, hashFile)
				fmt.Printf("check: %s", line)
				fmt.Printf("cache: Cache File = %s\n", cacheFile)
				fmt.Printf("check:  File Hash = %s\n", hexHash)
				fmt.Printf("check: Cache Hash = %s\n", cachedHash)
			}
			if hexHash == string(cachedHash) {
				fmt.Printf("%s: OK\n", hashFile)
			} else {
				fmt.Printf("%s: ERROR\n", hashFile)
				os.Exit(ErrorCacheMissed)
			}
		default:
			fmt.Printf("%s *%s\n", hexHash, hashFile)
		}
	} else {
		fmt.Println(AppHelp)
		os.Exit(ErrorArgumentError)
	}
}
