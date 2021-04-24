[![build](https://github.com/adityak74/google-drive-upload-git-action/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/adityak74/google-drive-upload-git-action/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/adityak74/google-drive-upload-git-action)](https://goreportcard.com/report/github.com/adityak74/google-drive-upload-git-action)

# google-drive-upload-git-action
Github action that uploads files to Google Drive.
**This only works with a Google Service Account!**

Thanks to [Team Tumbleweed](https://github.com/team-tumbleweed) for developing the initial version of this actions package.

To make a GSA go to the [Credentials Dashboard](https://console.cloud.google.com/apis/credentials). You will need to download the **.json key** and base64 encode it. You will use this string as the `credentials` input. To convert the *json* file to base64 without having to use an online tool (which is insecure), use this command:

`echo -n $(cat credentials.json)| base64 -w 0`

On mac this the base64 by default opts for -w as 0, you can skip and just use base64 without any params.

You will also need to **share the drive with the servie account.** To do this, just share the folder like you would normally with a friend, except you share it with the service account email address. Additionally you will need to give the service account acccess to the google drive API. 
Go to `https://console.developers.google.com/apis/api/drive.googleapis.com/overview?project={PROJECT_ID}`. Where `{PROJECT_ID}` is the id of your GCP project. Find more info about that [here.](https://support.google.com/googleapi/answer/7014113?hl=en)

# Inputs

## ``filename``
Required: **YES**.  

The name of the file you want to upload.

## ``name``
Required: **NO**

The name you want the file to have in Google Drive. If this input is not provided, it will use the `filename`.

## ``overwrite``
Required: **NO**

If you want to overwrite the filename with existing file, it will use the `filename`.

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
on: [push]

jobs:
  my_job:
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
        uses: adityak74/google-drive-upload-git-action@main
        with:
          credentials: ${{ secrets.credentials }}
          filename: "archive.zip"
          folderId: ${{ secrets.folderId }}
          name: "documentation.zip" # optional string
          overwrite: "true" # optional boolean
          
```
