//
// PACKAGES
//
package main

//
// IMPORTS
//
import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
	// "github.com/wrfly/ecp"
)

/*
 * CONSTANTS
 */
const (
	AppDesc            AppMetaData = "File change detection tool"
	AppName            AppMetaData = "FileDelta"
	AppVersion         AppMetaData = "0.1.0"
	CLIName            AppMetaData = "filedelta"
	CommandCheck                   = "check"
	CommandStore                   = "store"
	ErrorArgumentError             = 10
	ErrorCacheMissed               = 1
)

// const AppHelp AppMetaData = ``

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
	cmd      string
	debug    = false
	hashFile string
	lastArg  string
)

/*
 * FUNCTIONS
 */
func init() {
	for _, a := range os.Args[1:] {
		switch a {
		case CommandStore, CommandCheck:
			lastArg = a
		case "-d", "--debug":
			debug = true
		case "-v", "--version":
			fmt.Println(AppLabel)
			os.Exit(0)
		case "-h", "--help":
			fmt.Println(AppHelp)
			os.Exit(0)
		default:
			switch lastArg {
			case CommandStore:
				cmd = CommandStore
				hashFile = a
			case CommandCheck:
				cmd = CommandCheck
				hashFile = a
			default:
				lastArg = ""
				hashFile = a
			}
		}
	}
}

/*
 * MAIN ENTRYPOINT
 */
func main() {
	// log.Printf("%s", AppLabel)

	if len(hashFile) > 0 {
		cacheFile := url.QueryEscape(hashFile)
		f, err := os.Open(hashFile)
		if err != nil {
			log.Fatal(err)
		}

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		f.Close()
		theHash := h.Sum(nil)

		db, err := scribble.New(".", nil)
		if err != nil {
			fmt.Println("Error", err)
		}

		switch cmd {
		case CommandStore:
			hexHash := fmt.Sprintf("%x", theHash)
			if debug {
				fmt.Printf("cache: %s\n", hashFile)
				fmt.Printf("cache: %s\n", cacheFile)
				fmt.Printf("cache: %s\n", hexHash)
			}
			fmt.Printf("%x *%s\n", theHash, hashFile)
			if err := db.Write("cache", cacheFile, hexHash); err != nil {
				fmt.Println("Error", err)
			}
		case CommandCheck:
			hexHash := fmt.Sprintf("%x", theHash)

			var cachedHash string
			if err := db.Read("cache", cacheFile, &cachedHash); err != nil {
				fmt.Println("Error", err)
			}

			if debug {
				line := fmt.Sprintf("%x *%s\n", theHash, hashFile)
				fmt.Printf("check: %s", line)
				fmt.Printf("cache: Cache File = %s\n", cacheFile)
				fmt.Printf("check:  File Hash = %x\n", theHash)
				fmt.Printf("check: Cache Hash = %s\n", cachedHash)
			}
			if hexHash == string(cachedHash) {
				fmt.Printf("%s: OK\n", hashFile)
			} else {
				fmt.Printf("%s: ERROR\n", hashFile)
				os.Exit(ErrorCacheMissed)
			}
		default:
			fmt.Printf("%x *%s\n", theHash, hashFile)
		}
	} else {
		// fmt.Fprintln(os.Stderr, "You need to specify a file or option")
		fmt.Println(AppHelp)
		os.Exit(ErrorArgumentError)
	}
}
