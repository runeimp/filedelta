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
	"hash"
	"io/ioutil"
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
	AppVersion         AppMetaData = "0.4.0"
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

USAGE: filedelta [OPTIONS] COMMAND FILENAME

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
	cacheFile  string
	cachePath  string
	cmd        string
	debug      = false
	hashFile   string
	hashSHA256 hash.Hash
	homePath   string
)

/*
 * FUNCTIONS
 */
func cacheFilePath(file string) string {
	fileAbs, err := filepath.Abs(file)
	if err == nil {
		file = fileAbs
	}
	return fmt.Sprintf("%s/%s", cachePath, hashSHA256String(file))
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
		err       error
		f         *os.File
		fileBytes []byte
		hashStr   string
	)

	f, err = os.Open(file)
	if err != nil {
		return hashStr, err
	}
	fileBytes, err = ioutil.ReadAll(f)
	hashStr = hashSHA256ByteString(fileBytes)
	f.Close()

	return hashStr, err
}

func hashSHA256Bytes(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}

func hashSHA256ByteString(bytes []byte) string {
	hashBytes := hashSHA256Bytes(bytes)
	return fmt.Sprintf("%x", hashBytes)
}

func hashSHA256String(value string) string {
	// fmt.Printf("hashSHA256String() |      nil | hashStr = %q\n", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	hashBytes := hashSHA256Bytes([]byte(value))
	return fmt.Sprintf("%x", hashBytes)
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

	hashSHA256 = sha256.New()
	homePath, _ = homedir.Dir()
	cachePath = fmt.Sprintf("%s/.local/filedelta/cache", homePath)

	_, err := os.Stat(cachePath)
	if err != nil {
		os.MkdirAll(cachePath, 0700)
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
