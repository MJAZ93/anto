#!/bin/bash

# Determine OS type
OS="$(uname)"
ZIP_FILE=".anto.zip"
EXTRACTED_FOLDER=".anto"

# Download and unzip the anto.zip (assuming it's already downloaded, or you can add curl/wget here)
echo "Unzipping $ZIP_FILE..."
unzip "$ZIP_FILE" || { echo "Unzip failed"; exit 1; }

# Remove the __MACOSX folder if it exists
if [ -d "$EXTRACTED_FOLDER/__MACOSX" ]; then
    echo "Removing __MACOSX..."
    rm -rf "$EXTRACTED_FOLDER/__MACOSX"
fi

# Navigate to the extracted folder
cd "$EXTRACTED_FOLDER" || { echo "Failed to enter directory"; exit 1; }

# Check if the system is Windows or Unix-like and run the appropriate binary
if [[ "$OS" == "MINGW64_NT"* || "$OS" == "MSYS_NT"* || "$OS" == "CYGWIN_NT"* ]]; then
    # Windows - run anto.exe with the 'init' parameter
    echo "Running anto.exe with 'init' parameter on Windows..."
    ./anto.exe init
else
    # Unix-like systems - run the 'anto' binary with 'init' parameter
    echo "Running anto with 'init' parameter on Unix-like system..."
    chmod +x anto
    ./anto init
fi

echo "Installation complete!"
