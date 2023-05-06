// Mathis Van Eetvelde
// 2021-present

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
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
	encodedInput     = "encoded"
)

type Args struct {
	filename string
	name     string
	folderId string
	creds    string
}

func run() error {
	args, err := parseArguments()
	if err != nil {
		return err
	}

	// fetching a JWT config with credentials and the right scope
	conf, err := google.JWTConfigFromJSON([]byte(args.creds), scope)
	if err != nil {
		return fmt.Errorf("fetching JWT credentials failed with error: %v", err)
	}

	// instantiating a new drive service
	ctx := context.Background()
	svc, err := drive.New(conf.Client(ctx))
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(args.filename)
	if err != nil {
		return fmt.Errorf("opening file with filename: %v failed with error: %v", args.filename, err)
	}

	// decide name of file in GDrive
	if args.name == "" {
		args.name = file.Name()
	}

	f := &drive.File{
		Name:    args.name,
		Parents: []string{args.folderId},
	}

	_, err = svc.Files.Create(f).Media(file).SupportsAllDrives(true).Do()
	if err != nil {
		return fmt.Errorf("creating file: %+v failed with error: %v", f, err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		githubactions.Fatalf(err.Error())
	}
}

func missingInput(inputName string) error {
	return fmt.Errorf("missing input '%v'", inputName)
}

func incorrectInput(inputName string, reason string) error {
	if reason == "" {
		return fmt.Errorf("incorrect input '%v'", inputName)
	} else {
		return fmt.Errorf("incorrect input '%v' reason: %v", inputName, reason)
	}
}

func parseArguments() (*Args, error) {
	// get filename argument from action input
	filename := githubactions.GetInput(filenameInput)
	if filename == "" {
		return nil, missingInput(filenameInput)
	}

	// get name argument from action input
	name := githubactions.GetInput(nameInput)

	// get folderId argument from action input
	folderId := githubactions.GetInput(folderIdInput)
	if folderId == "" {
		return nil, missingInput(folderIdInput)
	}

	// get base64 encoded credentials argument from action input
	credentialsStr := githubactions.GetInput(credentialsInput)
	if credentialsStr == "" {
		return nil, missingInput(credentialsInput)
	}
	// add base64 encoded credentials argument to mask
	githubactions.AddMask(credentialsStr)

	// get encoded boolean input
	var encoded bool
	encodedStr := githubactions.GetInput(encodedInput)
	if encodedStr == "" || encodedStr == "true" {
		encoded = true
	} else if encodedStr == "false" {
		encoded = false
	} else {
		return nil, incorrectInput(encodedInput, "encoded needs to be either empty, `false` or `true`.")
	}

	// decode if encoded is true
	var credentials string
	if encoded {
		// decode credentials to []byte
		credentialsByte, err := base64.StdEncoding.DecodeString(credentialsStr)
		if err != nil {
			return nil, incorrectInput(credentials, fmt.Sprintf("base64 decoding of 'credentials' failed with error: %v", err))
		}
		credentials = string(credentialsByte)
	} else {
		credentials = credentialsStr
	}

	creds := strings.TrimSuffix(string(credentials), "\n")

	// add decoded credentials argument to mask
	githubactions.AddMask(creds)

	return &Args{filename: filename, name: name, folderId: folderId, creds: creds}, nil
}
