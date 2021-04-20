// TTW Software Team
// 2021-present

// Modified by Aditya Karnam
// 2021
// Added file overwrite support

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	scope = "https://www.googleapis.com/auth/drive.file"
)

func uploadToDrive(svc *drive.Service, filename string, folderId string, driveFile *drive.File) {
	file, err := os.Open(filename)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("opening file with filename: %v failed with error: %v", filename, err))
	}

	if driveFile != nil {
		f := &drive.File{
			Name: file.Name(),
		}
		_, err = svc.Files.Update(driveFile.Id, f).Media(file).Do()
	} else {
		f := &drive.File{
			Name:    file.Name(),
			Parents: []string{folderId},
		}
		_, err = svc.Files.Create(f).Media(file).Do()
	}

	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("creating/updating file failed with error: %v", err))
	} else {
		githubactions.Debugf("Uploaded/Updated file.")
	}
}

func main() {

	// get file argument from action input
	filename := githubactions.GetInput("filename")
	if filename == "" {
		githubactions.Fatalf("missing input 'filename'")
	}

	// get overwrite flag
	var overwriteFlag bool
	overwrite := githubactions.GetInput("overwrite")
	if overwrite == "" {
		githubactions.Warningf("Overwrite is disabled.")
		overwriteFlag = false
	} else {
		overwriteFlag, _ = strconv.ParseBool(overwrite)
	}
	// get folderId argument from action input
	folderId := githubactions.GetInput("folderId")
	if folderId == "" {
		githubactions.Fatalf("missing input 'folderId'")
	}

	// get base64 encoded credentials argument from action input
	credentials := githubactions.GetInput("credentials")
	if credentials == "" {
		githubactions.Fatalf("missing input 'credentials'")
	}
	// add base64 encoded credentials argument to mask
	githubactions.AddMask(credentials)

	// decode credentials to []byte
	decodedCredentials, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("base64 decoding of 'credentials' failed with error: %v", err))
	}
	// add decoded credentials argument to mask
	githubactions.AddMask(string(decodedCredentials))

	// fetching a JWT config with credentials and the right scope
	conf, err := google.JWTConfigFromJSON(decodedCredentials, scope)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("fetching JWT credentials failed with error: %v", err))
	}

	// instantiating a new drive service
	ctx := context.Background()
	svc, err := drive.New(conf.Client(ctx))
	if err != nil {
		log.Println(err)
	}

	if overwriteFlag {
		r, err := svc.Files.List().Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
			fmt.Println("Unable to retrieve files")
		}
		fmt.Println("Files:")
		if len(r.Files) == 0 {
			fmt.Println("No similar files found. Creating a new file")
			uploadToDrive(svc, filename, folderId, nil)
		} else {
			for _, i := range r.Files {
				if filename == i.Name {
					fmt.Printf("Overwriting file: %s (%s)\n", i.Name, i.Id)
					uploadToDrive(svc, filename, folderId, i)
				}
			}
		}
	} else {
		uploadToDrive(svc, filename, folderId, nil)
	}
}
