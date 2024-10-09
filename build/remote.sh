#!/bin/bash

# Determine OS type
OS="$(uname)"
ZIP_URL="https://raw.githubusercontent.com/MJAZ93/anto/main/build/mac.zip"
ZIP_FILE="mac.zip"
EXTRACTED_FOLDER="mac"

# Function to exit script on error
exit_on_error() {
    echo "$1"
    exit 1
}

# Download the GitHub zip
echo "Downloading repository zip..."
if [[ "$OS" == "Darwin" || "$OS" == "Linux" ]]; then
    curl -L -o "$ZIP_FILE" "$ZIP_URL" || exit_on_error "Download failed"
elif [[ "$OS" == "MINGW64_NT"* || "$OS" == "MSYS_NT"* ]]; then
    powershell -Command "Invoke-WebRequest -Uri '$ZIP_URL' -OutFile '$ZIP_FILE'" || exit_on_error "Download failed"
else
    exit_on_error "Unsupported OS"
fi

# Check if the zip file is valid
if [[ ! -f "$ZIP_FILE" || $(file --mime-type -b "$ZIP_FILE") != "application/zip" ]]; then
    exit_on_error "Downloaded file is not a valid zip archive"
fi

# Extract the zip to a temporary folder
echo "Extracting zip..."
mkdir -p "$EXTRACTED_FOLDER"
unzip -o "$ZIP_FILE" -d "$EXTRACTED_FOLDER" || exit_on_error "Unzip failed"

# Move the contents to the current folder
echo "Moving files to the root folder..."
mv "$EXTRACTED_FOLDER"/* . || exit_on_error "Failed to move files"

# Example command execution
echo "Running installation commands..."
chmod +x install.sh || exit_on_error "Failed to make install.sh executable"
./install.sh || exit_on_error "Failed to run install.sh"

# Clean up
echo "Cleaning up..."
rm -rf "$ZIP_FILE" "$EXTRACTED_FOLDER"

echo "Installation complete!"
