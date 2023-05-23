# Contribution user guide

You'll need to [download the Go binary][go-binary] to build and test this application.

You may also need to follow the user guide in the README.md to obtain the
Google service account to interact with the API with your locally built app.

## Check out the repo

```bash
git clone mathisve/gdrive-upload-action
```

## Run the tests

You may want to run the tests before you start development, and make sure all tests are passing.

```bash
go test -v
```

## Build the app locally

You can build the app locally with either `go run main.go` or `go build .`.

## Running the app locally

The `go run` method is very useful during development, as you can make sure
new or altered functionality works end-to-end.

In order for the githubactions.GitInput functions to work, you'll need to
export the needed env variables. E.g if using [a direnv `.envrc` file][direnv]:

```bash
export INPUT_FILENAME=README.md
export INPUT_ENCODED=true
export INPUT_CREDENTIALS=...some-base64-service-account-json...
export INPUT_FOLDERID=...some-gdrive-folder-id...
```

Then, run the app, and examine the output

```bash
go run main.go
::add-mask::...
::add-mask::...
::debug::Creating file README.md in folder 1cPZfHTv4Btz-wazqowfEPurP4Ede_zyv
```

[go-binary]: https://go.dev/dl/
[direnv]: https://direnv.net/
