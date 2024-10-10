
# VSK and MSK Files: Structure and Rules

VSK and MSK files form the core of Anto’s validation system. These files allow you to define rules for validating the structure and content of your Git projects. Here’s a breakdown of each file type and its structure.

## VSK Files: Project Structure Validation

VSK (Validation Structure Key) files are used to define the structure of your project—both the folders and files it must contain. The syntax follows a tree-like structure where you define directories and files using brackets `[]` and curly braces `{}` respectively.

### Structure:

- **Folders**: Defined by square brackets `[ ]`. You can specify either an explicit folder name or use regular expressions (regex) to match multiple folder names.
  
- **Files**: Defined by curly braces `{ }`. Similarly, you can specify an exact file name or use regex to match multiple file names in a directory.

### Example:

```plaintext
[app]
    [src]
        [main]
            [java]
                [anto]
                    [feature]
                        {*utils.6}   # Matches files that end with 'utils.6'
                    [ui]            # Matches 'ui' directory
                    [utils]         # Matches 'utils' directory
        [test]
            [feature]
                {*utils.6}
    [build]
        [**]   # Matches all files and folders recursively (useful for ignoring entire directories)
[commit]
    {commit.*}  # Matches files like 'commit.log', 'commit.txt', etc.
```

### Rules:

1. **Folder Definitions**: You define folders using `[folder_name]`, and any folders inside are indented with a tab.
2. **File Definitions**: You define files using `{file_name}` and can use regex to match file patterns. For instance, `{commit.*}` matches files like `commit.log`, `commit.txt`, etc.
3. **Wildcard Rules**: Use `**` to match all files and subdirectories under a directory.

This structure ensures that your project always adheres to a predefined folder layout, with required files present in the correct locations.

### Commands:

To generate the VSK file for your project, use the command:
```bash
.anto create-validation
```

## MSK Files: File Content Validation

MSK (Mask) files define the rules for the content within individual files. These rules allow you to specify what must or must not be present in the file, ensuring consistency in code style, size, and content.

### Structure:

- **Comments**: Use `/* ... */` to define the error message that will be shown if the validation fails.
- **Rules**:
  - `l <number>`: Defines a limit (usually for the number of lines in the file).
  - `+ <regex>`: Specifies content that **must** be present in the file (one or more).
  - `- <regex>`: Specifies content that **must not** be present in the file (none).

### Example:

```plaintext
/*
MainActivity.kt must adhere to the following rules:
- It should contain only navigation code.
- It must not contain any @Composable annotations.
- The file should not exceed 300 lines to maintain readability.
*/

# Line count rule
l 300 <   # The file must have fewer than 300 lines

# Must include this regex pattern (one or more)
+ *Activity   # The file name must contain 'Activity'

# Must exclude this regex pattern (none allowed)
- @Composable*   # The file must not contain any @Composable annotations
```

### Rules:

1. **Line Count**: You can set a maximum number of lines a file can have using `l <number> <`. This ensures files don’t become too large and difficult to maintain.
2. **Required Content**: Use `+ <regex>` to specify patterns that must appear in the file. All these patterns must compile and be found in the file.
3. **Excluded Content**: Use `- <regex>` to define patterns that must not appear in the file. These patterns should not compile or exist in the file.
4. **Error Message**: In the comments section, you define a message that will be displayed if the file violates any of the rules.

### Commands:

You can automatically generate the MSK file structure for your project using the command:
```bash
.anto create-structure
```

## How They Work Together

- **VSK Files** ensure that your project’s folder and file structure is correct and consistent across different environments.
- **MSK Files** guarantee that the contents of your files follow predefined guidelines, keeping your code clean and maintainable.

Together, these files help enforce structure and content rules across your entire codebase, making code reviews more focused on functionality and less on trivial format or structural issues.
