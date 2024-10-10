
# Anto
Simplifying Git Project Organization and Code Quality

![Project Logo](build/logo.png)

## Overview

Anto is designed to ease the burden of code reviews, project integration, and maintaining code quality.

We believe that code reviews should focus on functionality, rather than commit message formats, project structure, or simple file rules.

Anto provides an easy way to validate commit messages, project structure (files and folders), and file content using [VSK/MSK](./vsk_msk_structure.md) files, combined with [Git hooks](https://git-scm.com/book/ms/v2/Customizing-Git-Git-Hooks).

It's a game changer! :-)

[screenshot-gif]

## Features

### Commit Validation (`.anto/commit.msk`)
Commit validation works by defining rules (regex and max lines) in the `.anto/commit.msk` file:

```
/*
Commits should follow the Conventional Commits structure:

[optional scope]: <description>

Allowed types:
- feat: A new feature
- fix: A bug fix

Commit messages should not contain the word 'commit' and must be a maximum of 300 characters.
*/

l 300 <

+ feat:*
+ fix:*

- commit*
```
For more details about `.msk` files, see the [VSK/MSK](#vskmskfiles) section.

### Project Structure Validation (`.anto/validation.vsk`)
The project structure validation is defined through rules for files and folders in the `.anto/structure.vsk` file:

```plaintext
[app]
    [src]
        [main]
            [java]
                [anto]
                    [feature]
                        {*utils.6} # Matches files that end with 'utils.6'
                        [grand]
                    [ui] # Matches 'ui' directory
                    [utils] # Matches 'utils' directory
        [test]
            [feature]
                {*utils.6}
                [grand]
                [ui]
                [utils]
    [build]
        [**] # Matches all files and folders recursively (useful for ignoring entire directories)
[commit]
    {commit.*} # Matches files like 'commit.log', 'commit.txt', etc.
```

You can generate the `validation.vsk` file for your project with this command:

```bash
.anto create-validation
```

### File Content Validation (`.anto/{projectName/*/*.vsk}`)
File content validation is based on rules defined for specific files within your project. Create directories and files (with the `.msk` extension) that mirror your project structure inside the `.anto` folder.

To automate this, use the following command to create your project folder and file structure:

```bash
.anto create-structure
```

Example of an `.msk` file for `MainActivity.kt`:

```plaintext
/*
The MainActivity.kt file should only contain navigation code and must not include any @Composable annotations.

The file should be a maximum of 300 lines to prevent overly large classes.
*/

l 300 < # must have 300 lines max

+ *Activity # the class name must contain Activity

- @Composable* # the activity shouldnot contain any composable
```

## Installation

### Prerequisites

Before getting started, ensure you have:

- Git installed [git].
- Removed temporary folders and files (e.g., node_modules, build).
- [Other necessary configurations or tools].

### Fast Installation

In the root of your .git project run the following command:
#### Mac

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/MJAZ93/anto/main/build/remote-mac.sh)"
```

#### Linux

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/MJAZ93/anto/main/build/remote-linux.sh)"
```

#### Windows

```bash
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/MJAZ93/anto/main/build/remote-windows.ps1" -OutFile "$env:TEMP
emote-windows.ps1"; & "$env:TEMP
emote-windows.ps1"
```

### Step-by-Step Installation

1. Download the zip: https://raw.githubusercontent.com/MJAZ93/anto/main/build/mac.zip
2. Extract and run `install.sh`:

```bash
   ./install.sh
```

#### Or

3. Copy the `.anto` folder to the root of your Git project.
4. Open the `.anto` folder and run the following commands:
5. Initialize Anto with:
   ```bash
   ./anto init
   ```
   Or create the validation file (`structure.vsk`):
   ```bash
   ./anto create-validation
   ```
   Create the `.msk` files for validating project files:
   ```bash
   ./anto create-structure
   ```
   Add the Git `commit-msg` hook (validation rules live in `commit.msk`):
   ```bash
   ./anto add-precommit
   ```

### Testing
Just make a commit and **anto** will ensure to make the all validations.

## How It Works

- **Folder Structure Validation**: Ensures the folder structure follows predefined rules.
- **File Content Validation**: Validates specific content inside project files.
- **Commit Validation**: Ensures commit messages follow the predefined rules.
- **Documentation**: Ensures the project has a proper documentation, because it forces the project to describe each file and folder.

## Additional Validation

You can extend the commit-msg hook by adding custom steps. For example, you can include other scripts (tests, linters) in `validations.sh` (for Mac and Linux) or `validation.ps1` (for Windows):

```bash
#!/bin/sh
./anto validate

# goes to the root of the project
cd ..
# Additional scripts (e.g., tests, linters)
lint
gradle test
```

## Licence
Anto is under the [Fair Code Licence](https://faircode.io/).

## Donate - Support Development
To help Anto growth please donate using github sponsors.
