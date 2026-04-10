# Capivara

<img src="./logo/logo.png" width=200/>

## ⚡ Quick Install (No Root Required)

You can easily install Capivara globally for your user without needing administrative or root privileges. The scripts below will download the source, build the binary, and place it in your local user directory (`~/.local/bin`).

> **Note:** Ensure you have [Go](https://go.dev/doc/install) and Git installed on your system before running these scripts.

### 🍎 macOS and 🐧 Linux

Open your terminal and run the following command to download and execute the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/atendi9/capivara/main/install.sh | bash
```

### 🪟 Windows
Open PowerShell and run the following command:

```PowerShell
Invoke-RestMethod -Uri https://raw.githubusercontent.com/atendi9/capivara/main/install.ps1 | Invoke-Expression
```

## 🚀 Usage

Capivara acts as a wrapper around the standard `go test` command. By default, the output is in English. 

To run your tests with Capivara, navigate to your Go project and run:

```bash
capivara ./...
```

### 🇧🇷 Portuguese Output

If you prefer to see the test results and error messages in Portuguese, you can use the `--lang=portuguese` flag:

```bash
capivara --lang=portuguese ./...
```

*(Any standard `go test` flags can still be passed alongside it!)*
