// TTW Software Team
// Mathis Van Eetvelde
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
	"strings"

	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	scope            = "https://www.googleapis.com/auth/drive.file"
	filenameInput    = "filename"
	nameInput        = "name"
	folderIdInput    = "folderId"
	credentialsInput = "credentials"
)

func uploadToDrive(svc *drive.Service, filename string, folderId string, driveFile *drive.File, name string) {
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
			Name:    name,
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

	// get filename argument from action input
	filename := githubactions.GetInput(filenameInput)
	if filename == "" {
		missingInput(filenameInput)
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
	// get name argument from action input
	name := githubactions.GetInput(nameInput)

	// get folderId argument from action input
	folderId := githubactions.GetInput(folderIdInput)
	if folderId == "" {
		missingInput(folderIdInput)
	}

	// get base64 encoded credentials argument from action input
	credentials := githubactions.GetInput(credentialsInput)
	if credentials == "" {
		missingInput(credentialsInput)
	}
	// add base64 encoded credentials argument to mask
	githubactions.AddMask(credentials)

	// decode credentials to []byte
	decodedCredentials, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("base64 decoding of 'credentials' failed with error: %v", err))
	}

	creds := strings.TrimSuffix(string(decodedCredentials), "\n")

	// add decoded credentials argument to mask
	githubactions.AddMask(creds)

	// fetching a JWT config with credentials and the right scope
	conf, err := google.JWTConfigFromJSON([]byte(creds), scope)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("fetching JWT credentials failed with error: %v", err))
	}

	// instantiating a new drive service
	ctx := context.Background()
	svc, err := drive.New(conf.Client(ctx))
	if err != nil {
		log.Println(err)
	}

	if name == "" {
		file, err := os.Open(filename)
		if err != nil {
			githubactions.Fatalf(fmt.Sprintf("opening file with filename: %v failed with error: %v", filename, err))
		}
		name = file.Name()
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
			uploadToDrive(svc, filename, folderId, nil, name)
		} else {
			for _, i := range r.Files {
				if filename == i.Name {
					fmt.Printf("Overwriting file: %s (%s)\n", i.Name, i.Id)
					uploadToDrive(svc, filename, folderId, i, "")
				}
			}
		}
	} else {
		uploadToDrive(svc, filename, folderId, nil, name)
	}
}

func missingInput(inputName string) {
	githubactions.Fatalf(fmt.Sprintf("missing input '%v'", inputName))
}
