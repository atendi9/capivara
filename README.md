# Capivara

<img src="./logo/logo.png" width=200/>

## ⚡ Quick Install (No Root Required)

You can easily install Capivara globally for your user without needing administrative or root privileges. The scripts below will download the source, build the binary, and place it in your local user directory (`~/.local/bin`).

> **Note:** Ensure you have [Go](https://go.dev/doc/install) and Git installed on your system to build Capivara. To execute tests, you will also need Go and/or [Node.js](https://nodejs.org/pt-br/download) installed on your machine depending on your project.

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

Capivara is a smart test runner wrapper that automatically detects your project's language based on the presence of a `go.mod` (Go) or `package.json` (Node.js) file. It wraps around standard testing commands (`go test` and `node --test`) to provide a cleaner, localized output. By default, the output is in English. 

To run your tests with Capivara, navigate to your project's root directory and run:

**For Go projects:**
```bash
capivara ./...
```
*(Any standard `go test` flags can still be passed alongside it!)*

**For Node.js projects:**
```bash
capivara
```

### 🇧🇷 Portuguese Output

If you prefer to see the test results and error messages in Portuguese, you can use the `--lang=portuguese` flag:

```bash
capivara --lang=portuguese ./...
```

### 🇷🇺 Russian Output

If you prefer to see the test results and error messages in Russian, you can use the `--lang=russian` flag:

```bash
capivara --lang=russian ./...
```

### 🇯🇵 Japanese Output

If you prefer to see the test results and error messages in Japanese, you can use the `--lang=japanese` flag:

```bash
capivara --lang=japanese ./...
```

### 🇨🇳 Chinese Output

If you prefer to see the test results and error messages in Chinese, you can use the `--lang=chinese` flag:

```bash
capivara --lang=chinese ./...
```
