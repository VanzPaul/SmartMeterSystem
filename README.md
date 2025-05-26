# Project SmartMeterSystem

...

## Getting Started

### Windows-Specific Setup  
For **Windows users**, follow these steps to install dependencies via [Chocolatey](https://chocolatey.org/) :  

1. **Install Chocolatey** (if not already installed):  
   Open an **administrative Command Prompt** and run:  
   ```bash
   @"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
   ```  

2. **Install Go, Make, and Docker**:  
   ```bash
   choco install golang -y          # Installs Go (ensure GOPATH is configured) 
   choco install make -y           # Installs GNU Make
   choco install docker -y         # Installs Docker CLI
   ```  

3. **Verify installations**:  
   ```bash
   go version
   make --version
   docker --version
   ```  

4. **Clone the repository**:  
   Run the command 
   ```bash
   git clone https://github.com/VanzPaul/SmartMeterSystem
   ```
   and switch to `dev-ui` branch.  

5. **Install VsCode Extensions**:
   It is recommended to install the following extensions for vscode:
   - Dcoker
   - Go
   - Makefile 
   -


6. **Build and Run the prject**:
   In the terminal, run `make build` and `make run`. 
   
   To run install and run the mongo-db database while in development, run the command
   ```bash
   make docker-run
   ```
   Head to the docker extension in vscode and stop the container `smartmetersystem-app` (this is the container for the application) while leaving the `mongo:latest` (this is the container for mongodb) conatiner run.


   **Common Issues**  
   Error: `make: *** No rule to make target '<command>'. Stop`  
   Solution: Explicitly specify the Makefile path:  
   ```bash
   make -f Makefile <command>
   ```  

7. **Using the application**:
   In your browser, head to `localhost:8080/home` to access the web application.

---

## MakeFile

Run build make command with tests:
```bash
make all
```

Build the application:
```bash
make build
```

Run the application:
```bash
make run
```

Create DB container:
```bash
make docker-run
```

Shutdown DB Container:
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

---

## Docs

cmd