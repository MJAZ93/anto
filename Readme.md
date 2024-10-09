# Anto
Git project organization and quality made simple

![Project Logo](path/to/logo.png)

## Overview

This project aims to reduce the headache of code review, project integration and code quality. 

We believe code review must focus on the functionality, not on the commit messages stability, project structure and files simple rules. 

It provides a easy way to validate commit messages, project structure (files and folders) and file content using [VSK/MSK](#vskmskfiles) files + [git hooks](https://git-scm.com/book/ms/v2/Customizing-Git-Git-Hooks).

Believe, it's a game changer! :-)

[screenshot-gif]

## Features

### Commit validation (.anto/commit.msk)
The validation commit works by defining rules (regex and max lines), you only need to configure these rules in the *.anto/commit.msk* file:
```
#Desciption of the commit, this will show when a the validation fails
/*
Commit should follow the Conventional Commits structure:

[optional scope]: <description>

Types allowed:
- feat: A new feature
- fix: A bug fix

Commit message should not contain the word commit.

Must have a maximum of 300 characters.
*/

#maximum number of words
l 300 <

# Regex, at least one must match
+ feat:*
+ fix:*

# Regex, no one must match
- commit*
```
More about the .msk files, see [VSK/MSK](#vskmskfiles) section.

### Project structure validation (.anto/validation.vsk)
The project structure validation works by defining files and folder rules, you only need to configure these rules in the *.anto/structure.vsk* file:

```
# a folder defined between [], inside the [] you can put a folder name or a regex that matches files in one directory
[app]
    #a tab separates a childen from a father
    [src]
        [main]
            [java]
                [anto]
                    [feature]
                        # a file defined between {}, inside the {} you can put a file name or a regex that matches files in one directory
                        {*utils.6}
                        [grand]
                    [ui]
                    [utils]
        [test]
            [feature]
                    {*utils.6}
                    [grand]
                    [ui]
                    [utils]
    [build]
        # When you want a file to ignore all folders and files in one file (for example node_modules and build) you set **
        [**]
[commit]
    {commit.*}
```

Anto will generate *validation.vsk* file for your project using this command:

```bash
.anto create-validation
```

### File content validation (.anto/{projectName/*/*.vsk})
The file structure validation works by defining file rule for each file in your project that you want to validate, 
to achieve this, you will have to create folders and files (adding the .msk extension) that match your project structure inside .anto directory.
Hopefully we created a script that will create your project folder and file structure inside .anto, and you will only need to create the rules using this command:

```bash
.anto create-structure
```

Example of an .msk of a *Activity.kt:

```   
# Add the message with the rules, this message will be displayed when a file does not meet the rules:
/*
The MainActivity.kt should only contain navigation code, should not contain any @Composable
...
must have maximum 300 lines to ensure we dont have big classes
*/

# Number of classes
l 300 <

# Add the regexes that the file content must meet with. (+ regex). All must compile.
# Ensure the name is Activity
+ *Activity

# Add the regexes that the file content must meet with. (- regex). All must not compile.
- @Composable*
```

## Installation

### Prerequisites

Before you begin, ensure you have met the following requirements:

- You have git installed in your project [git].
- You have removed all your temp folders and files (eg. node_modules, build, etc).
- You have [other necessary configuration or tool].

### Fast Instalation

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
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/MJAZ93/anto/main/build/remote-windows.ps1" -OutFile "$env:TEMP\remote-windows.ps1"; & "$env:TEMP\remote-windows.ps1" 
```

### Step-by-Step Installation

1. *Download the zip:*
   https://raw.githubusercontent.com/MJAZ93/anto/main/build/mac.zip
2. *Extract and run the install.sh* Or:
3. *Copy the .anto to the root of your git project*
4. *Open the .anto and run the following commands:*
5. *Run the following commands:*
   ```bash
   .anto init
   ```
   *Or*
   Create the validation file (structure.vsk)
   ```bash
   .anto create-validation
   ```
   
   Create the files to validate all the project files (*.msk)
   ```bash
   .anto create-strucure
   ```
   
   Add git commit-msg hook, the validation rules live in commit.msk
   ```bash
   .anto add-precommit
   ```

### How it works

#### Folder structure validation
#### File content validation
#### Commit validation


### Additional validation

You can add additional validation to the commit-msg hook, in other words, you can add additional steps to the validation
any script validation inside the validations.sh for mac and linux and validation.ps1, like this:

```bash
#!/bin/sh
./anto validate

cd ..
#Other scripts (tests, linters)
lint
gradle test
```
### <a name="vskmskfiles"></a>Vsk and Msk files

[link to the documentation]