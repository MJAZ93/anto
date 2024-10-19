#!/bin/bash

# Remove the __MACOSX folder if it exists
if [ -d "__MACOSX" ]; then
    echo "Removing __MACOSX..."
    rm -rf "__MACOSX"
fi

cd .anto || { echo "Failed to enter directory"; exit 1; }

# Unix-like systems - run the 'anto' binary with 'init' parameter
echo "Running anto with 'init' parameter on Unix-like system..."

chmod +x anto
./anto init

echo "Installation complete!"
