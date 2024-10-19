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
Invoke-WebRequest -Uri $ZIP_URL -OutFile $ZIP_FILE
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

try {
    Add-Type -AssemblyName 'System.IO.Compression.FileSystem'
    [System.IO.Compression.ZipFile]::ExtractToDirectory($ZIP_FILE, $EXTRACTED_FOLDER)
} catch {
    exit_on_error "Unzip failed: $($_.Exception.Message)"
}

# Move the contents to the current folder
Write-Host "Moving files to the root folder..."
$sourceFolder = Join-Path $EXTRACTED_FOLDER "windows"
try {
    Get-ChildItem $sourceFolder -Force | Move-Item -Destination "."
} catch {
    exit_on_error "Failed to move files: $($_.Exception.Message)"
}

# Example command execution
Write-Host "Running anto with 'init' parameter on Windows..."

# Check if the __MACOSX folder exists and remove it if it does
if (Test-Path "__MACOSX") {
    Write-Host "Removing __MACOSX..."
    Remove-Item -Recurse -Force "__MACOSX"
}

# Navigate to the extracted folder
Set-Location ".anto" -ErrorAction Stop

# Unblock the 'anto.exe' file to prevent execution issues
Unblock-File -Path "./anto.exe"

# Windows - Run the 'anto' binary with 'init' parameter
Write-Host "Running anto with 'init' parameter on Windows..."
Start-Process -NoNewWindow -Wait "./anto.exe" "init" -ErrorAction Stop

Write-Host "Installation complete!"

# Move back to the parent directory
Set-Location ..

# Clean up
Write-Host "Cleaning up..."
Remove-Item -Recurse -Force "windows.zip", "install.ps1", "windows", "anto.exe"
