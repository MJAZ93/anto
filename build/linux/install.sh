#!/bin/bash

# Determine OS type
OS="$(uname)"
ZIP_FILE=".anto.zip"
EXTRACTED_FOLDER=".anto"

# Download and unzip the anto.zip (assuming it's already downloaded, or you can add curl/wget here)
echo "Unzipping $ZIP_FILE..."
unzip "$ZIP_FILE" || { echo "Unzip failed"; exit 1; }

# Navigate to the extracted folder
cd "$EXTRACTED_FOLDER" || { echo "Failed to enter directory"; exit 1; }

# Unix-like systems - run the 'anto' binary with 'init' parameter
echo "Running anto with 'init' parameter on Unix-like system..."
chmod +x anto
./anto init

echo "Installation complete!"
