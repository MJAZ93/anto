# PowerShell script for Windows

# Variables
$OS = "windows"
$ZIP_URL = "https://raw.githubusercontent.com/MJAZ93/anto/main/build/windows.zip"
$ZIP_FILE = "windows.zip"
$EXTRACTED_FOLDER = "."

# Function to exit script on error
function exit_on_error {
    param (
        [string]$message
    )
    Write-Host $message
    exit 1
}

# Download the GitHub zip
Write-Host "Downloading repository zip..."
Invoke-WebRequest -Uri $ZIP_URL -OutFile $ZIP_FILE -ErrorAction Stop
if (!$?) { exit_on_error "Download failed" }

# Check if the zip file is valid
if (-not (Test-Path $ZIP_FILE) -or (Get-Item $ZIP_FILE).Extension -ne ".zip") {
    exit_on_error "Downloaded file is not a valid zip archive"
}

# Extract the zip to a temporary folder
Write-Host "Extracting zip..."
if (-not (Test-Path $EXTRACTED_FOLDER)) {
    New-Item -ItemType Directory -Force -Path $EXTRACTED_FOLDER
}
Add-Type -AssemblyName 'System.IO.Compression.FileSystem'
[System.IO.Compression.ZipFile]::ExtractToDirectory($ZIP_FILE, $EXTRACTED_FOLDER) -ErrorAction Stop
if (!$?) { exit_on_error "Unzip failed" }

# Move the contents to the current folder
Write-Host "Moving files to the root folder..."
$sourceFolder = Join-Path $EXTRACTED_FOLDER "windows"
Get-ChildItem $sourceFolder -Force | Move-Item -Destination "." -ErrorAction Stop
if (!$?) { exit_on_error "Failed to move files" }

# Example command execution
Write-Host "Running installation commands..."
$installScript = "./install.ps1"
if (-not (Test-Path $installScript)) {
exit_on_error "install.ps1 script not found"
}
& $installScript
if (!$?) { exit_on_error "Failed to run install.ps1" }

# Clean up
Write-Host "Cleaning up..."
Remove-Item -Recurse -Force "windows.zip", "install.ps1", "windows", "remote.ps1"
