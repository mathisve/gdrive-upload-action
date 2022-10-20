[![build](https://github.com/team-tumbleweed/gdrive-upload-action/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/team-tumbleweed/gdrive-upload-action/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/team-tumbleweed/gdrive-upload-action)](https://goreportcard.com/report/github.com/team-tumbleweed/gdrive-upload-action)

# gdrive-upload-action
Github action that uploads files to Google Drive.
**This only works with a Google Service Account!**

To make a GSA go to the [Credentials Dashboard](https://console.cloud.google.com/apis/credentials). You will need to download the **.json key** and base64 encode it. You will use this string as the `credentials` input. To convert the *json* file to base64 without having to use an online tool (which is insecure), use this command:

`echo -n $(cat credentials.json)| base64`

You will also need to **share the drive with the service account.** To do this, just share the folder like you would normally with a friend, except you share it with the service account email address. Additionally you will need to give the service account acccess to the google drive API. 
Go to `https://console.developers.google.com/apis/api/drive.googleapis.com/overview?project={PROJECT_ID}`. Where `{PROJECT_ID}` is the id of your GCP project. Find more info about that [here.](https://support.google.com/googleapi/answer/7014113?hl=en)

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

A base64 encoded string with the [GSA credentials](https://stackoverflow.com/questions/46287267/how-can-i-get-the-file-service-account-json-for-google-translate-api/46290808).


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
        uses: team-tumbleweed/gdrive-upload-action@main
        with:
          filename: archive.zip
          name: documentation.zip
          folderId: ${{ secrets.folderId }}
          credentials: ${{ secrets.credentials }}
```
