#!/bin/bash

# Determine OS type
OS="$(uname)"
ZIP_URL="https://raw.githubusercontent.com/MJAZ93/anto/main/build/mac.zip"
ZIP_FILE="mac.zip"
EXTRACTED_FILES=()

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

# Capture the list of files in the zip archive
echo "Listing files in the zip archive..."
EXTRACTED_FILES=($(unzip -Z1 "$ZIP_FILE"))

# Extract the zip to the current directory (same as where the zip is located)
echo "Extracting zip to the current directory..."
unzip -o "$ZIP_FILE" -d "." || exit_on_error "Unzip failed"

# Example command execution
echo "Running installation commands..."
chmod +x install.sh || exit_on_error "Failed to make install.sh executable"
./install.sh || exit_on_error "Failed to run install.sh"

# Delete the extracted files and folders
echo "Cleaning up extracted files..."
for file in "${EXTRACTED_FILES[@]}"; do
    echo "Deleting $file..."
    rm -rf "$file" || exit_on_error "Failed to delete $file"
done

# Delete the zip file
echo "Deleting zip file..."
rm -f "$ZIP_FILE"

echo "Cleanup complete!"
echo "Installation complete!"
