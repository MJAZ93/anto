
### Step-by-Step Installation

1. Download the zip:
   mac - https://raw.githubusercontent.com/MJAZ93/anto/main/build/mac.zip
   windows - https://raw.githubusercontent.com/MJAZ93/anto/main/build/windows.zip
   linux - https://raw.githubusercontent.com/MJAZ93/anto/main/build/linux.zip
2. Extract and run install script:

*mac and linux*:
```bash
   ./install.sh
```

*windows*:
```powershell
   .\install.ps1
```

#### Or

3. Copy the `.anto` folder to the root of your Git project.
4. Open the `.anto` folder and run the following commands:
5. Initialize Anto with:
   *mac and linux*:
   ```bash
   ./anto init
   ```
   *windows*:
   ```powershell
   .\anto.ps1 init
   ```

##### Or

1. Create the validation file (`structure.vsk`):
   *mac and linux*:
   ```bash
   ./anto create-validation
   ```
*windows*:
   ```powershell
    .\anto.ps1 create-validation
   ```

2. Create the `.msk` files for validating project files:

*mac and linux*:
   ```bash
   ./anto create-structure
   ```

*windows*:
   ```powershell
   .\anto.ps1 create-structure
   ```  

3. Add the Git `commit-msg` hook (validation rules live in `commit.msk`):  
   *mac and linux*:
   ```bash
   ./anto add-precommit
   ```
*windows*:
   ```powershell
    .\anto.ps1 add-precommit
   ```  
