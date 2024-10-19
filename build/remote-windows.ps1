# PowerShell script for Windows

# Set variables
$ZIP_FILE = ".windows.zip"
$EXTRACTED_FOLDER = ".windows"

# Function to exit script on error
function exit_on_error {
    param([string]$message)
    Write-Host $message
    exit 1
}

# Unzip the anto.zip
Write-Host "Unzipping $ZIP_FILE..."
if (-Not (Test-Path $ZIP_FILE)) {
    exit_on_error "Zip file not found: $ZIP_FILE"
}

# Extract zip using .NET functionality (equivalent of unzip in bash)
Add-Type -AssemblyName 'System.IO.Compression.FileSystem'
[System.IO.Compression.ZipFile]::ExtractToDirectory($ZIP_FILE, $EXTRACTED_FOLDER) || exit_on_error "Unzip failed"

# Remove the __MACOSX folder if it exists
$macosxFolder = Join-Path $EXTRACTED_FOLDER "__MACOSX"
if (Test-Path $macosxFolder) {
    Write-Host "Removing __MACOSX..."
    Remove-Item -Recurse -Force $macosxFolder || exit_on_error "Failed to remove __MACOSX"
}

# Navigate to the extracted folder
if (-Not (Test-Path $EXTRACTED_FOLDER)) {
    exit_on_error "Failed to enter directory: $EXTRACTED_FOLDER"
}
Set-Location $EXTRACTED_FOLDER

# Windows - run the 'anto.exe' binary with 'init' parameter
$antoExe = "anto.exe"
if (Test-Path $antoExe) {
    Write-Host "Running anto.exe with 'init' parameter..."
    .\anto.exe init || exit_on_error "Failed to run anto.exe with init parameter"
} else {
    exit_on_error "anto.exe not found in $EXTRACTED_FOLDER"
}

Write-Host "Installation complete!"

# Go back and clean up the zip file
Write-Host "Cleaning up..."

# Remove files and folders
Remove-Item -Recurse -Force "mac.zip", "install.sh", "mac", "remote.sh", "__MACOSX"
