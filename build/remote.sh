#!/bin/bash

# Determine OS type
OS="$(uname)"
ZIP_URL="https://github.com/MJAZ93/anto/main/build/.anto.zip"
ZIP_FILE="mac.zip"
EXTRACTED_FOLDER="."

# Download the GitHub zip
echo "Downloading repository zip..."
if [[ "$OS" == "Darwin" || "$OS" == "Linux" ]]; then
    curl -L -o "$ZIP_FILE" "$ZIP_URL"
elif [[ "$OS" == "MINGW64_NT"* || "$OS" == "MSYS_NT"* ]]; then
    powershell -Command "Invoke-WebRequest -Uri '$ZIP_URL' -OutFile '$ZIP_FILE'"
else
    echo "Unsupported OS"
    exit 1
fi

# Extract the zip
echo "Extracting zip..."
unzip "$ZIP_FILE" || { echo "Unzip failed"; exit 1; }

# Navigate to the extracted folder and run commands
cd "$EXTRACTED_FOLDER" || { echo "Failed to enter directory"; exit 1; }

# Example command execution
# ./install.sh (You can replace this with your own commands)
echo "Running installation commands..."
chmod +x install.sh
./install.sh

# Go back and clean up the zip file
cd ..
echo "Cleaning up..."
rm -rf "$ZIP_FILE" "$EXTRACTED_FOLDER"

echo "Installation complete!"
