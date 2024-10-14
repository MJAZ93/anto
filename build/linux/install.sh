#!/bin/bash

# Navigate to the extracted folder
cd .anto || { echo "Failed to enter directory"; exit 1; }

# Unix-like systems - run the 'anto' binary with 'init' parameter
echo "Running anto with 'init' parameter on Unix-like system..."
chmod +x anto
./anto init

echo "Installation complete!"
