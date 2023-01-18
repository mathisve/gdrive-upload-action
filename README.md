[![build](https://github.com/mathisve/gdrive-upload-action/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/mathisve/gdrive-upload-action/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mathisve/gdrive-upload-action)](https://goreportcard.com/report/github.com/mathisve/gdrive-upload-action)

# gdrive-upload-action
### Github action that uploads files to Google Drive.

**Only works with a Google Service Account!**

To make a Google Service Account go to the [Credentials Dashboard](https://console.cloud.google.com/apis/credentials) and press `CREATE CREDENTIALS`, then click on `Service Account`. After creating the SA, go back to the Credentials Dashboard, click on the SA you just created, click on the `KEYS` tabs. Then click on `ADD KEY`, `Create New Key` and select `json`.

### Encoded credentials
1. Encode the credentials

    `cat credentials.json | base64`

2. Create a new secret called `credentials` and copy the output of the previous command into this secret.

3. Use this secret as your credentials input in your workflow.

### Plaintext credentials

1. Create a new secret called `credentials` and copy the output of the previous command into this secret.

2. Use this secret as your credentials input in your workflow.

3. Set the `encoded` input to `false` in your workflow.


You will also need to **share the drive with the service account.** To do this, just share the folder like you would normally with a friend, except you share it with the service account email address. Go to `https://console.developers.google.com/apis/api/drive.googleapis.com/overview?project={PROJECT_ID}`. Where `{PROJECT_ID}` is the id of your GCP project. 

Find more info about your Drives project ID [here.](https://support.google.com/googleapi/answer/7014113?hl=en)


Additionally you need to enable the Google Drive API for your GCP project. Do this by going [here](https://console.cloud.google.com/marketplace/product/google/drive.googleapis.com) and pressing `ENABLE`.

# Inputs

## ``filename``
Required: **YES**.  

The name of the file you want to upload.

## ``name``
Required: **NO**

The name you want the file to have in Google Drive. If this input is not provided, it will use the `filename`.

## ``folderId``
Required: **YES**. 

The [ID of the folder](https://ploi.io/documentation/database/where-do-i-get-google-drive-folder-id) you want to upload to.

## ``credentials``
Required: **YES**.

A string with the [GSA credentials](https://stackoverflow.com/questions/46287267/how-can-i-get-the-file-service-account-json-for-google-translate-api/46290808).
This string should be base64 encoded. If it is not, set the `encoded` input to `false`.

## ``encoded``
Required: **NO**

Whether or not the credentials string is base64 encoded. Defaults to `true`.



# Usage Example

## Simple Workflow
In this example we stored the folderId and credentials as action secrets. This is highly recommended as leaking your credentials key will allow anyone to use your service account.
```yaml
# .github/workflows/main.yml
name: Main
on:
  push:
    branches:
      - 'master'
      - 'main'
      - '!test'

jobs:
  gdrive-upload:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Archive files
        run: |
          sudo apt-get update
          sudo apt-get install zip
          zip -r archive.zip *

      - name: Upload to gdrive
        uses: mathisve/gdrive-upload-action@main
        with:
          filename: archive.zip
          name: documentation.zip
          folderId: ${{ secrets.folderId }}
          credentials: ${{ secrets.credentials }}
          encoded: false
```
