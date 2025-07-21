[中文版](./README_CN.md)

# 🔍 WhoisChecker

An efficient WHOIS domain query tool that supports concurrent queries, multiple output formats, and intelligent retry mechanisms.

## ✨ Features

- 🚀 **Concurrent Queries**: High-efficiency concurrent queries using Goroutines + Channels
- 📊 **Multiple Output Formats**: Supports table, JSON, and CSV output formats
- 🔄 **Intelligent Retry**: Automatic retry for failed queries to improve success rate
- ⏰ **Query Intervals**: Adjustable query intervals to avoid server rate limits
- 📝 **Detailed Logging**: Supports verbose mode to display query progress and error information
- 📄 **Batch Queries**: Supports reading domain lists from files for batch queries
- 🎨 **Beautiful Tables**: Uses go-pretty library for attractive table output
- 🌐 **ASN Information**: Query and display Autonomous System Number information
- 📍 **IP Address Resolution**: Resolve and display IPv4/IPv6 addresses for domains
- 🏳️ **Country Information**: Extract and display country information from WHOIS data

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/your-username/whoisChecker.git
cd whoisChecker
go mod tidy
go build -o whoisChecker ./cmd/main.go
```

### Basic Usage

```bash
# Query a single domain
./whoisChecker whois -d google.com

# Batch query from file
./whoisChecker whois --domainfile domains.txt

# Use concurrent queries (5 workers)
./whoisChecker whois --domainfile domains.txt --workers 5

# Query domain and show ASN information
./whoisChecker whois -d example.com --asn

# Query domain's IPv4 and IPv6 addresses
./whoisChecker whois -d example.com --t4 --t6

# Combined query: WHOIS + ASN + IP addresses
./whoisChecker whois -d example.com --asn --t4 --t6
```

## 📖 Usage Guide

### Command Line Options

| Flag           | Short | Description                                | Default |
| -------------- | ----- | ------------------------------------------ | ------- |
| `--domain`     | `-d`  | Single domain to query                     | -       |
| `--domainfile` | -     | File containing domain list (one per line) | -       |
| `--workers`    | `-w`  | Number of concurrent workers               | 3       |
| `--json`       | -     | Output results in JSON format              | false   |
| `--csv`        | -     | Output results in CSV format               | false   |
| `--verbose`    | `-v`  | Enable verbose output                      | false   |
| `--retry`      | `-r`  | Number of retry attempts                   | 3       |
| `--delay`      | -     | Delay between queries (milliseconds)       | 1000    |
| `--asn`        | -     | Query and display ASN information          | false   |
| `--t4`         | -     | Resolve and display IPv4 addresses         | false   |
| `--t6`         | -     | Resolve and display IPv6 addresses         | false   |

### Output Formats

#### Table Format (Default)

```bash
./whoisChecker whois -d example.com
```

```
┌─────────────┬────────────┬────────────┬───────────┐
│ DOMAIN      │ REGISTRAR  │ EXPIRES    │ DAYS LEFT │
├─────────────┼────────────┼────────────┼───────────┤
│ example.com │ IANA       │ 2024-08-13 │ 120       │
└─────────────┴────────────┴────────────┴───────────┘
```

#### JSON Format

```bash
./whoisChecker whois -d example.com --json
```

```json
[
  {
    "domain": "example.com",
    "registrar": "IANA",
    "expires": "2024-08-13",
    "days_left": "120",
    "status": "clientDeleteProhibited",
    "error": "",
    "asn": "",
    "ipv4": "",
    "ipv6": "",
    "country": "US"
  }
]
```

#### CSV Format

```bash
./whoisChecker whois -d example.com --csv
```

```csv
Domain,Registrar,Expires,Days Left,Status,Country
example.com,IANA,2024-08-13,120,clientDeleteProhibited,US
```

### Advanced Features

#### ASN Query

```bash
# Query domain with ASN information
./whoisChecker whois -d example.com --asn

# Output includes ASN column
┌─────────────┬────────────┬────────────┬───────────┬─────────┐
│ DOMAIN      │ REGISTRAR  │ EXPIRES    │ DAYS LEFT │ ASN     │
├─────────────┼────────────┼────────────┼───────────┼─────────┤
│ example.com │ IANA       │ 2024-08-13 │ 120       │ AS15133 │
└─────────────┴────────────┴────────────┴───────────┴─────────┘
```

#### IP Address Resolution

```bash
# Query IPv4 addresses
./whoisChecker whois -d example.com --t4

# Query IPv6 addresses
./whoisChecker whois -d example.com --t6

# Query both IPv4 and IPv6
./whoisChecker whois -d example.com --t4 --t6
```

#### Batch Processing

Create a file `domains.txt`:

```
google.com
github.com
stackoverflow.com
example.com
```

```bash
# Process all domains with 5 concurrent workers
./whoisChecker whois --domainfile domains.txt --workers 5 --verbose
```

#### Performance Optimization

```bash
# High-performance batch processing
./whoisChecker whois --domainfile large_list.txt \
  --workers 10 \
  --delay 500 \
  --retry 2 \
  --json > results.json
```

## 🔧 Configuration

### Environment Variables

You can set environment variables to configure default behavior:

```bash
export WHOIS_WORKERS=5
export WHOIS_DELAY=500
export WHOIS_RETRY=2
```

### Domain File Format

Create a text file with one domain per line:

```
example.com
google.com
github.com
# Comments are supported
stackoverflow.com
```

## 📊 Performance

### Benchmarks

- **Single Query**: ~1-2 seconds per domain
- **Concurrent Queries**: Up to 10x faster with optimal worker count
- **Memory Usage**: ~50MB for 1000 domains
- **Recommended Settings**: 3-5 workers for best balance

### Rate Limiting

To avoid being blocked by WHOIS servers:

- Use reasonable delay (500-1000ms)
- Limit concurrent workers (3-5)
- Enable retry mechanism
- Monitor verbose output for errors

## 🛠️ Development

### Dependencies

```bash
go get github.com/spf13/cobra@latest
go get github.com/jedib0t/go-pretty/v6@latest
go get github.com/likexian/whois@latest
go get github.com/likexian/whois-parser@latest
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/your-username/whoisChecker.git
cd whoisChecker

# Download dependencies
go mod tidy

# Build for current platform
go build -o whoisChecker ./cmd/main.go

# Run tests
go test ./...
```

### Cross-Platform Compilation

Build binaries for multiple platforms using the provided build scripts:

#### Using build script (macOS/Linux)

```bash
# Make the script executable
chmod +x build.sh

# Build for all supported platforms
./build.sh
```

This will create binaries in the `build/` directory for:

- **macOS**: Intel (amd64) and Apple Silicon (arm64)
- **Windows**: 64-bit (amd64) and 32-bit (386)
- **Linux**: 64-bit (amd64), 32-bit (386), ARM64, and ARM

#### Using batch script (Windows)

```cmd
# Run the Windows build script
build.bat
```

#### Manual cross-compilation

You can also build for specific platforms manually:

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o whoisChecker_darwin_amd64 ./cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o whoisChecker_darwin_arm64 ./cmd/main.go

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o whoisChecker_windows_amd64.exe ./cmd/main.go

# Windows 32-bit
GOOS=windows GOARCH=386 go build -o whoisChecker_windows_386.exe ./cmd/main.go

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o whoisChecker_linux_amd64 ./cmd/main.go

# Linux ARM64 (Ubuntu on ARM)
GOOS=linux GOARCH=arm64 go build -o whoisChecker_linux_arm64 ./cmd/main.go

# Linux ARM (Raspberry Pi)
GOOS=linux GOARCH=arm go build -o whoisChecker_linux_arm ./cmd/main.go
```

#### Platform-specific Installation

**macOS:**

```bash
# Download the appropriate binary for your Mac
# For Intel Macs:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_darwin_amd64
# For Apple Silicon Macs:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_darwin_arm64

chmod +x whoisChecker
sudo mv whoisChecker /usr/local/bin/
```

**Windows:**

```cmd
# Download whoisChecker_windows_amd64.exe from releases
# Or use PowerShell:
Invoke-WebRequest -Uri "https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_windows_amd64.exe" -OutFile "whoisChecker.exe"
```

**Ubuntu/Linux:**

```bash
# For x86_64 systems:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_linux_amd64
# For ARM64 systems:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_linux_arm64

chmod +x whoisChecker
sudo mv whoisChecker /usr/local/bin/
```

### Project Structure

```
whoisChecker/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── root.go          # Root command configuration
│   └── whois.go         # WHOIS query implementation
├── go.mod               # Go module file
├── go.sum               # Go checksum file
├── README.md            # This file
└── LICENSE              # License file
```

## 🐛 Troubleshooting

### Common Issues

1. **Connection Timeouts**

   ```bash
   # Increase delay and retry count
   ./whoisChecker whois -d domain.com --delay 2000 --retry 5
   ```

2. **Rate Limited**

   ```bash
   # Reduce workers and increase delay
   ./whoisChecker whois --domainfile domains.txt --workers 2 --delay 1500
   ```

3. **No WHOIS Data**

   - Some domains may not have public WHOIS data
   - Try with verbose mode to see detailed error messages

   ```bash
   ./whoisChecker whois -d domain.com --verbose
   ```

4. **Memory Issues with Large Files**
   ```bash
   # Process in smaller batches
   split -l 100 large_domains.txt batch_
   for file in batch_*; do
     ./whoisChecker whois --domainfile "$file" --json >> all_results.json
   done
   ```

### Debug Mode

```bash
# Enable verbose output for debugging
./whoisChecker whois -d domain.com --verbose

# Output example:
🔍 Querying example.com...
✅ Query completed for example.com
```

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [jedib0t/go-pretty](https://github.com/jedib0t/go-pretty) - Table formatting
- [likexian/whois](https://github.com/likexian/whois) - WHOIS client
- [likexian/whois-parser](https://github.com/likexian/whois-parser) - WHOIS data parser

## 📞 Support

If you have any questions or issues, please:

1. Check the [troubleshooting](#-troubleshooting) section
2. Search existing [issues](https://github.com/your-username/whoisChecker/issues)
3. Create a new issue if needed

---

⭐ **Star this project if you find it helpful!**
