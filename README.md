---

# Project SmartMeterSystem

...

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

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

4. **Set up Go environment variables** (if required) :  
   Ensure `GOPATH` and `GOROOT` are configured in your system paths.  

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

### Common Issues  
**Error**: `make: *** No rule to make target '<command>'. Stop`  
**Solution**: Explicitly specify the Makefile path:  
```bash
make -f Makefile <command>
```  
This ensures compatibility on systems where the default `make` configuration might not recognize the target .  

---