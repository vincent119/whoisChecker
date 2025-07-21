# 🔍 WhoisChecker

一個高效的 WHOIS 域名查詢工具，支援並發查詢、多種輸出格式和智能重試機制。

## ✨ 功能特色

- 🚀 **並發查詢**: 使用 Goroutines + Channels 實現高效並發查詢
- 📊 **多種輸出格式**: 支援表格、JSON 和 CSV 輸出格式
- 🔄 **智能重試**: 對失敗的查詢自動重試，提高成功率
- ⏰ **查詢間隔**: 可調整的查詢間隔以避免伺服器頻率限制
- 📝 **詳細記錄**: 支援詳細模式，顯示查詢進度和錯誤信息
- 📄 **批量查詢**: 支援從檔案讀取域名清單進行批量查詢
- 🎨 **美觀表格**: 使用 go-pretty 庫提供美觀的表格輸出

## 🚀 快速開始

### 安裝

```bash
git clone https://github.com/vincent119/whoisChecker.git
cd whoisChecker
go mod tidy
go build -o whoisChecker ./cmd/main.go
```

### 基本使用

```bash
# 查詢單個域名
./whoisChecker whois -d google.com

# 從檔案批量查詢
./whoisChecker whois --domainfile domains.txt

# 使用並發查詢（5 個工作者）
./whoisChecker whois --domainfile domains.txt --workers 5

# 查詢域名並顯示 ASN 信息
./whoisChecker whois -d example.com --asn

# 查詢域名的 IPv4 和 IPv6 地址
./whoisChecker whois -d example.com --t4 --t6

# 組合查詢：WHOIS + ASN + IP 地址
./whoisChecker whois -d example.com --asn --t4 --t6
```

## 📖 使用指南

### 命令行選項

```bash
Flags:
  -A, --asn                 顯示 ASN 信息
      --csv                 輸出 CSV 格式
      --delay int           查詢間隔毫秒數 (default 1000)
  -d, --domain string       單一網域查詢
      --domainfile string   從檔案讀取網域清單
  -h, --help                顯示幫助信息
      --json                輸出 JSON 格式
      --retry int           查詢失敗時的重試次數 (default 2)
      --t4                  查詢並顯示 IPv4 地址
      --t6                  查詢並顯示 IPv6 地址
  -v, --verbose             詳細輸出模式
  -w, --workers int         並發查詢數量 (default 3)
```

### 輸出格式

#### 1. 表格格式（默認）

```bash
./whoisChecker whois -d example.com
```

```text
+-------------+----------------------+------------+-----------+
| 域名        | 註冊商               | 到期時間   | 剩餘天數  |
+-------------+----------------------+------------+-----------+
| example.com | RESERVED-Internet... | 2025-08-13 | 22        |
+-------------+----------------------+------------+-----------+
```

#### 2. JSON 格式

```bash
./whoisChecker whois -d example.com --json
```

```json
[
  {
    "domain": "example.com",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "expires": "2025-08-13",
    "days_left": "22",
    "status": "success"
  }
]
```

#### 3. CSV 格式

```bash
./whoisChecker whois -d example.com --csv
```

```csv
域名,註冊商,到期時間,剩餘天數,狀態
example.com,RESERVED-Internet Assigned Numbers Authority,2025-08-13,22,success
```

### 進階查詢選項

#### ASN 查詢

ASN (Autonomous System Number) 是網際網路中用於識別自治系統的唯一編號。

```bash
# 查詢域名的 ASN 信息
./whoisChecker whois -d example.com --asn

# 批量查詢多個域名的 ASN
./whoisChecker whois --domainfile domains.txt --asn --json
```

#### IP 地址查詢

可以同時查詢域名的 IPv4 和 IPv6 地址：

```bash
# 只查詢 IPv4 地址
./whoisChecker whois -d example.com --t4

# 只查詢 IPv6 地址
./whoisChecker whois -d example.com --t6

# 同時查詢 IPv4 和 IPv6
./whoisChecker whois -d example.com --t4 --t6
```

#### 完整信息查詢

```bash
# 查詢所有可用信息
./whoisChecker whois -d example.com --asn --t4 --t6 --json
```

JSON 輸出範例：

```json
[
  {
    "domain": "example.com",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "expires": "2025-08-13",
    "days_left": "22",
    "status": "success",
    "asn": "AS15133",
    "ipv4": "93.184.216.34",
    "ipv6": "2606:2800:220:1:248:1893:25c8:1946",
    "country": "US"
  }
]
```

### 批量查詢

創建域名清單檔案 `domains.txt`：

```text
google.com
github.com
stackoverflow.com
example.com
```

執行批量查詢：

```bash
# 基本批量查詢
./whoisChecker whois --domainfile domains.txt

# 使用 5 個並發工作者，查詢間隔 500ms
./whoisChecker whois --domainfile domains.txt --workers 5 --delay 500

# 詳細模式，顯示查詢進度
./whoisChecker whois --domainfile domains.txt --verbose

# 輸出為 JSON 格式
./whoisChecker whois --domainfile domains.txt --json
```

### 進階設定

#### 調整並發和延遲

```bash
# 使用 10 個並發工作者，查詢間隔 200ms
./whoisChecker whois --domainfile domains.txt --workers 10 --delay 200

# 降低並發數量，增加延遲，適合大量查詢
./whoisChecker whois --domainfile domains.txt --workers 2 --delay 2000
```

#### 重試機制

```bash
# 失敗時重試 5 次
./whoisChecker whois --domainfile domains.txt --retry 5

# 不重試，快速失敗
./whoisChecker whois --domainfile domains.txt --retry 0
```

#### 詳細輸出

```bash
# 啟用詳細模式，查看查詢進度和錯誤信息
./whoisChecker whois --domainfile domains.txt --verbose
```

## 🛠️ 開發

### 依賴安裝

```bash
go get github.com/spf13/cobra@latest
go get github.com/jedib0t/go-pretty/v6@latest
go get github.com/likexian/whois@latest
go get github.com/likexian/whois-parser@latest
```

### 從原始碼編譯

```bash
# 克隆專案
git clone https://github.com/your-username/whoisChecker.git
cd whoisChecker

# 下載依賴
go mod tidy

# 編譯當前平台版本
go build -o whoisChecker ./cmd/main.go

# 執行測試
go test ./...
```

### 跨平台編譯

使用提供的編譯腳本為多個平台建置執行檔：

#### 使用編譯腳本 (macOS/Linux)

```bash
# 賦予腳本執行權限
chmod +x build.sh

# 編譯所有支援的平台
./build.sh
```

這將在 `build/` 目錄中建立以下平台的執行檔：

- **macOS**: Intel (amd64) 和 Apple Silicon (arm64)
- **Windows**: 64 位元 (amd64) 和 32 位元 (386)
- **Linux**: 64 位元 (amd64)、32 位元 (386)、ARM64 和 ARM

#### 使用批次腳本 (Windows)

```cmd
# 執行 Windows 編譯腳本
build.bat
```

#### 手動跨平台編譯

您也可以手動為特定平台進行編譯：

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o whoisChecker_darwin_amd64 ./cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o whoisChecker_darwin_arm64 ./cmd/main.go

# Windows 64位元
GOOS=windows GOARCH=amd64 go build -o whoisChecker_windows_amd64.exe ./cmd/main.go

# Windows 32位元
GOOS=windows GOARCH=386 go build -o whoisChecker_windows_386.exe ./cmd/main.go

# Linux 64位元
GOOS=linux GOARCH=amd64 go build -o whoisChecker_linux_amd64 ./cmd/main.go

# Linux ARM64 (Ubuntu on ARM)
GOOS=linux GOARCH=arm64 go build -o whoisChecker_linux_arm64 ./cmd/main.go

# Linux ARM (Raspberry Pi)
GOOS=linux GOARCH=arm go build -o whoisChecker_linux_arm ./cmd/main.go
```

#### 平台特定安裝

**macOS:**

```bash
# 下載適合您 Mac 的執行檔
# Intel Mac:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_darwin_amd64
# Apple Silicon Mac:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_darwin_arm64

chmod +x whoisChecker
sudo mv whoisChecker /usr/local/bin/
```

**Windows:**

```cmd
# 從 releases 下載 whoisChecker_windows_amd64.exe
# 或使用 PowerShell:
Invoke-WebRequest -Uri "https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_windows_amd64.exe" -OutFile "whoisChecker.exe"
```

**Ubuntu/Linux:**

```bash
# x86_64 系統:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_linux_amd64
# ARM64 系統:
curl -L -o whoisChecker https://github.com/your-username/whoisChecker/releases/latest/download/whoisChecker_linux_arm64

chmod +x whoisChecker
sudo mv whoisChecker /usr/local/bin/
```

### 專案結構

```
whoisChecker/
├── cmd/
│   └── main.go          # 應用程式進入點
├── internal/
│   ├── root.go          # 根命令配置
│   └── whois.go         # WHOIS 查詢實作
├── build.sh             # Unix 編譯腳本
├── build.bat            # Windows 編譯腳本
├── go.mod               # Go 模組檔案
├── go.sum               # Go 校驗檔案
├── README.md            # 說明文件
└── LICENSE              # 授權檔案
```

## 🏗️ 技術架構

### 核心組件

- **Worker Pool Pattern**: 使用固定數量的 Goroutines 處理查詢任務
- **Channel Communication**: 通過 channels 在 goroutines 間傳遞數據
- **Retry Mechanism**: 自動重試失敗的查詢，指數退避延遲
- **Rate Limiting**: 可配置的查詢間隔，避免觸發 API 限制

### 依賴庫

- `github.com/spf13/cobra`: 命令行界面框架
- `github.com/likexian/whois`: WHOIS 查詢庫
- `github.com/likexian/whois-parser`: WHOIS 結果解析庫
- `github.com/jedib0t/go-pretty/v6/table`: 美化表格輸出

## 📊 性能優化

### 並發設定建議

- **小規模查詢（< 10 個域名）**: `--workers 2-3`
- **中等規模查詢（10-100 個域名）**: `--workers 5-10`
- **大規模查詢（> 100 個域名）**: `--workers 3-5`, `--delay 1000-2000`

### 避免限制

某些 WHOIS 伺服器對查詢頻率有限制，建議：

- 使用較少的並發工作者（2-3 個）
- 增加查詢間隔（1000-2000ms）
- 啟用重試機制
- 使用詳細模式監控查詢狀態

## 🔧 故障排除

### 常見問題

#### 1. 查詢失敗

```bash
# 使用詳細模式查看錯誤信息
./whoisChecker whois -d problem-domain.com --verbose

# 增加重試次數
./whoisChecker whois -d problem-domain.com --retry 5
```

#### 2. 網路限制

```bash
# 降低並發數量，增加延遲
./whoisChecker whois --domainfile domains.txt --workers 2 --delay 3000
```

#### 3. 解析錯誤

某些域名的 WHOIS 格式可能不標準，程式會盡量解析並在 JSON/CSV 輸出中提供錯誤詳情。

### 調試模式

```bash
# 啟用詳細輸出查看完整的查詢過程
./whoisChecker whois --domainfile domains.txt --verbose --json
```

## 📝 範例

### 監控域名到期時間

```bash
# 查詢即將到期的域名（輸出 CSV 用於進一步處理）
./whoisChecker whois --domainfile my-domains.txt --csv > domain-status.csv
```

### 批量域名健康檢查

```bash
# 並發查詢大量域名，詳細記錄
./whoisChecker whois --domainfile large-domain-list.txt \
  --workers 5 \
  --delay 1000 \
  --retry 3 \
  --verbose \
  --json > domain-report.json 2> query-log.txt
```

## 🤝 貢獻

歡迎提交議題（Issue）和拉取請求（Pull Request）！

## 📄 授權

MIT 授權

## 🔗 相關連結

- [WHOIS 協議說明](https://tools.ietf.org/html/rfc3912)
- [Go Concurrency Patterns](https://blog.golang.org/pipelines)
- [Cobra CLI 框架](https://github.com/spf13/cobra)
