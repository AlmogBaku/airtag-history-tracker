# Airtag history tracker

This tiny app, tacks your Airtag locations and save them to a CSV file.

## How does it work?

Fortunately, the "Find my" app saves the recent location into a local file. This app continuously reads the file and
saves the location to a CSV file, and keep your Mac awake during this time.

That means, you need to have the "Find my" app installed and running on your phone.

## How to install?

At the moment, we don't release a compiled binary. You need to compile it yourself using Go.
Fortunately, it's very easy.

1. Install Go using [Homebrew](https://brew.sh/): `brew install go`
2. Compile and install this app using:
    ```console
   go install github.com/AlmogBaku/airtag-history-tracker@latest
    ``` 
3. That's it! You can run the app using `airtag-history-tracker` command.

## How to use?

Run the app using `airtag-history-tracker` command. We'll create the CSV files in your current directory.

To stop the app, press `Ctrl+C`.

To collect the locations of a specific device, add the `--device` flag with the device name. You can find the device
namein the "Find my" app.