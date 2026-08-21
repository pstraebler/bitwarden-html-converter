[![Build and Release](https://github.com/pstraebler/bitwarden-html-converter/actions/workflows/build.yml/badge.svg)](https://github.com/pstraebler/bitwarden-html-converter/actions/workflows/build.yml)

# Bitwarden HTML Converter

Lightweight cross-platform application to convert Bitwarden JSON exports into a printable HTML file. Fully offline and secure.

## Features

- ✅ Simple and intuitive graphical interface
- ✅ Bitwarden JSON → HTML conversion
- ✅ Styled table optimized for printing
- ✅ Single binary with no dependencies
- ✅ Windows, macOS and Linux support
- ✅ Display all item types (logins, cards, identities, notes)

## Installation

### Download precompiled binaries

Go to the [releases page](../../releases) and download the binary for your system:

- **Windows**: `bitwarden-html-converter-windows-amd64.exe`
- **macOS Intel**: `bitwarden-html-converter-macos-amd64.zip`
- **macOS Apple Silicon**: `bitwarden-html-converter-macos-arm64.zip`
- **Linux**: `bitwarden-html-converter-linux-amd64`

### Linux

After downloading, make the file executable:

```bash
chmod +x bitwarden-html-converter-*
```

On macOS, extract the archive and double-click `Bitwarden HTML Converter.app`.

## Usage

1. Launch the application (double-click or via terminal)
2. Click "Select JSON File" and choose your Bitwarden export
3. Click "HTML Destination" and choose where to save the file
4. Click "Convert"
5. Open the generated HTML file in your browser to print it

### Windows note

On Windows, the executable may be blocked by SmartScreen and show a warning such as "Windows protected your PC" or "This app may be dangerous".

To launch it anyway:

1. Double-click the `.exe`
2. Click **More info** in the warning dialog
3. Click **Run anyway**

If needed, you can also right-click the executable, open **Properties**, and check **Unblock** before launching it.

### macOS note

The macOS package is ad hoc signed, but it is not notarized by Apple. Gatekeeper may therefore still require one manual approval on the first launch.

#### Option 1: Finder

1. Open **System Settings**
2. Go to **Privacy & Security**
3. Scroll down and look for a message saying that the app was blocked
4. Click **Open Anyway**
5. Confirm the warning to allow the application to run

After that, try launching the app again.

#### Option 2: terminal

```bash
xattr -dr com.apple.quarantine "Bitwarden HTML Converter.app"
```

The Intel build uses the same procedure.

## Export from Bitwarden

To get a JSON file from Bitwarden:

1. Open Bitwarden (application or browser extension)
2. Go to **Settings** → **Export Vault**
3. Choose **JSON** format (unencrypted)
4. Download the file

⚠️ **Warning**: The JSON file contains your passwords in plain text. Delete it after use.

## Build from source

### Prerequisites

- Go 1.21 or higher
- GCC (for Fyne CGO compilation)

#### Linux
```bash
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

#### macOS
```bash
xcode-select --install
```

#### Windows
Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or use [MSYS2](https://www.msys2.org/)

### Compilation

```bash
git clone https://github.com/pstraebler/bitwarden-html-converter
cd bitwarden-html-converter
go mod download
go mod tidy
go build -o bitwarden-html-converter
```

## Project structure

```
bitwarden-html-converter/
├── main.go              # Graphical interface
├── converter.go         # Conversion logic
├── go.mod              # Go dependencies
├── .github/
│   └── workflows/
│       └── build.yml   # Automated multi-OS build
└── README.md
```

## Technologies used

- **Go**: Programming language
- **Fyne**: Cross-platform GUI framework
- **GitHub Actions**: Automated build and release

## License

MIT

## Security

- This application does not collect any data
- No network connections are made
- Files are processed locally
- Source code is open and auditable

## Contributors

Contributions welcome! Feel free to open an issue or pull request.
