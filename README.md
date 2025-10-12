# PathFinder v3.0.0 - DEATH STAR EDITION

**The fastest, most visually stunning web path discovery tool with full TUI dashboard.**

PathFinder combines the best features from **gobuster**, **ffuf**, and **feroxbuster** with a beautiful terminal UI, real-time statistics, and unique redirect chain tracking.

![Version](https://img.shields.io/badge/version-3.0.0-red)
![Go](https://img.shields.io/badge/Go-1.16%2B-blue)
![License](https://img.shields.io/badge/license-Educational-yellow)

---

## ✨ What's New in v3.0.0

🎨 **Full TUI Dashboard** - Real-time scanning with live statistics and colored themes
🌍 **Globe Mode** - Rotating Earth visualization (press F3)
🎨 **8 Color Themes** - Matrix, Rainbow, Cyber, Blood, Skittles, Dark, Purple, Amber
⚡ **Splash Screen** - Epic /PATHFINDER ASCII art with trickle animation
🎯 **Function Key Controls** - F1 (Cycle), F3 (Globe), F4 (Config), F5 (Export), ? (Help)
📊 **Live Results** - Watch paths getting discovered in real-time with scrollable results
📝 **Executive Summary Export** - Professional pentest reports with F5
🔧 **Interactive Config Menu** - Adjust settings on-the-fly with F4
❓ **Built-in Help System** - Press ? for comprehensive documentation
🌐 **Local Network Info** - Interface, IP, MAC, Subnet, Gateway display
⌨️ **Dynamic Input** - Enter new URLs directly in the TUI

---

## 🚀 Quick Start

### Windows (Simplest Method)
```cmd
cd PathFinder
SCAN.bat honeypotlogs.com
```

### Advanced Usage
```bash
pathfinder.exe -target https://example.com -wordlist wordlist.txt
```

---

## 🎮 TUI Dashboard Features

When you launch PathFinder, you'll see a **full-screen terminal interface** with:

### 📊 Dashboard Layout
- **Title Bar** - PathFinder version and author (glowing red)
- **Input Field** - Dynamic URL input (Press Enter to activate)
- **Scan Config** - Target URL, method, timeout, concurrency, rate limit, wordlist
- **Progress Box** - Animated progress bar, stats, speed, elapsed time
- **Statistics** - Direct 200s, Redirects, Protected paths, Errors, Speed
- **Local Network Info** - Interface name, IPv4, Subnet, IPv6, MAC, Gateway
- **Live Results** - Scrolling list of discovered paths (Up/Down to scroll)
- **Controls** - Quick reference for all keyboard shortcuts

### 🎨 Color Themes (Press 1-8 or F1)
1. **Matrix** - Classic green-on-black hacker aesthetic
2. **Rainbow** - Purple/magenta vibrant theme
3. **Cyber** - Cyan/blue cyberpunk vibes
4. **Blood** - Red theme (default)
5. **Skittles** - Random bold colors everywhere!
6. **Dark** - Light mode (white background, dark text)
7. **Purple** - Purple gradient aesthetic
8. **Amber** - Amber/orange terminal theme

### ⌨️ Keyboard Controls
- **F1** - Cycle through themes (Matrix → Rainbow → Cyber → Blood → Skittles → Dark → Purple → Amber)
- **F3** - Toggle Globe Mode (spinning Earth visualization)
- **F4** - Open Config Menu (adjust concurrency, rate limit, timeout, method)
- **F5** - Export Executive Summary Report (pentest documentation)
- **?** - Toggle Help Screen (comprehensive guide)
- **Enter** - Activate input field / Submit new URL
- **↑/↓** - Scroll through live results (Up/Down arrows)
- **1-8** - Jump directly to specific theme
- **Q** - Quit and show scan summary
- **Esc** - Close menus / Alternative quit

### 🌍 Globe Mode
Press **F3** to enter Globe Mode:
- Beautiful ASCII Earth rendered in 3D
- Rotates west-to-east (full rotation in 30 seconds)
- Uses real Earth bitmap data
- Theme-colored globe (changes with your active theme)
- Press F3 again to return to dashboard

### 🔧 Interactive Config Menu (F4)
Adjust scan settings in real-time:
- **Concurrency** - Number of simultaneous requests (◀/▶ to adjust)
- **Rate Limit** - Max requests per second (0 = unlimited)
- **Timeout** - Max wait time per request (seconds)
- **Method** - HTTP method (GET, POST, HEAD, PUT, DELETE, PATCH)
- Navigate with ↑/↓, change values with ◀/▶, close with F4 or Esc

### 📝 Executive Summary Export (F5)
Generate professional pentest reports:
- **Scan Metadata** - Tool version, target, timing, configuration
- **Executive Summary** - Statistics, findings breakdown, risk assessment
- **Detailed Findings** - All discoveries with timestamps, sorted chronologically
  - Direct 200s (HIGH PRIORITY) - Accessible resources
  - Redirects (MEDIUM PRIORITY) - Full redirect chains
  - Protected/Error Responses - 401/403/500 status codes
- **Recommendations** - Professional security guidance
- Auto-saves as `PathFinder_Report_YYYY-MM-DD_HH-MM-SS.txt`

### ❓ Help System (Press ?)
Built-in comprehensive documentation:
- **Overview** - What PathFinder does
- **Scan Configuration** - Explanation of all settings (Method, Timeout, Concurrency, Rate Limit)
- **Statistics** - What each metric means
- **Keyboard Shortcuts** - Complete control reference
- Press ? or Esc to close

---

## 🎯 Core Features

### Industry-Leading Scanner
✅ **Fast** - 2000-5000+ req/s with optimized Go concurrency
✅ **Complete redirect chain tracking** - See every redirect step with timestamps
✅ **Wildcard detection** - Automatically filters catch-all responses
✅ **Content fingerprinting** - Identify duplicate content via MD5 hashing
✅ **Direct 200 vs Redirected** - Know what's real vs redirected
✅ **Smart filtering** - By status codes, content size, regex patterns

### Output Features
- **Live TUI Dashboard** - Real-time visual feedback with 8 themes
- **Export to JSON/CSV** - Save results for later analysis
- **Executive Summary Reports** - Professional pentest documentation (F5)
- **Color-coded results** - Direct 200s (green), Redirects (yellow), Protected (red)
- **Detailed summary** - After quitting, see complete analysis
- **Scrollable results** - Navigate through thousands of findings with arrow keys

### Network Intelligence
- **Local Network Info Display** - Know your attack surface
  - Interface name (e.g., eth0, wlan0)
  - IPv4 and Subnet/CIDR
  - IPv6 address
  - MAC address
  - Default Gateway

---

## 📖 Usage Examples

### Basic Scan with TUI
```bash
pathfinder.exe -target https://example.com
```

### Fast Scan (High Concurrency)
```bash
pathfinder.exe -target https://example.com -concurrency 100
```

### With File Extensions
```bash
pathfinder.exe -target https://example.com -x php,html,js,txt
```

### Custom Wordlist
```bash
pathfinder.exe -target https://example.com -wordlist /path/to/big-wordlist.txt
```

### Filter Status Codes
```bash
# Only show 200s and redirects
pathfinder.exe -target https://example.com -mc 200,301,302

# Hide 404s
pathfinder.exe -target https://example.com -fc 404
```

### Rate Limited Scanning
```bash
# Be gentle: 50 req/s max
pathfinder.exe -target https://example.com -rate 50

# Add delay between requests
pathfinder.exe -target https://example.com -delay 100
```

### With Authentication
```bash
# Bearer token
pathfinder.exe -target https://api.example.com -H "Authorization:Bearer YOUR_TOKEN"

# Session cookie
pathfinder.exe -target https://example.com -cookie "session=abc123"
```

### Export Results
```bash
# Export to JSON
pathfinder.exe -target https://example.com -o results.json -of json

# Export to CSV
pathfinder.exe -target https://example.com -o results.csv -of csv
```

### Choose Theme at Launch
```bash
# Start with Skittles theme
pathfinder.exe -target https://example.com -theme skittles

# Start with Matrix theme
pathfinder.exe -target https://example.com -theme matrix
```

---

## 🛠️ Command-Line Options

### Core Options
```
-target <url>           Target base URL (required)
-wordlist <file>        Wordlist file (default: wordlist.txt)
-concurrency <n>        Concurrent requests (default: 50)
-timeout <n>            Timeout in seconds (default: 10)
-verbose                Show errors and debug info
```

### Filtering Options
```
-mc <codes>            Match status codes (200,301,302)
-fc <codes>            Filter status codes (404,403)
-fs <sizes>            Filter content sizes (1234,5678)
```

### Extensions & HTTP
```
-x <exts>              File extensions (php,html,js,txt)
-H <header>            Custom header (Name:Value)
-cookie <data>         Cookie string
-X <method>            HTTP method (default: GET)
```

### Performance Tuning
```
-rate <n>              Max requests/sec (0=unlimited)
-delay <n>             Delay between requests (ms)
```

### Recursion
```
-r                     Enable recursive scanning
-depth <n>             Max recursion depth (default: 3)
```

### Output
```
-o <file>              Output file path
-of <format>           Format: text, json, csv
-theme <name>          Theme: matrix, rainbow, cyber, blood, skittles, dark, purple, amber
```

---

## 🏗️ Building from Source

### Prerequisites
- Go 1.16 or higher
- Windows/Linux/macOS

### Build Script (Windows)
```cmd
build.bat
```
Creates: `pathfinder.exe`, `pathfinder-linux`, `pathfinder-mac`

### Manual Build
```bash
# Windows
go build -ldflags="-s -w" -o pathfinder.exe main.go

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o pathfinder main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o pathfinder main.go
```

---

## 📊 Performance Comparison

| Tool | Speed (req/s) | TUI | Redirect Tracking | Themes | Export Reports |
|------|---------------|-----|-------------------|--------|----------------|
| **PathFinder** | **2000-5000+** | ✅ Full | ✅ Complete chain | ✅ 8 themes | ✅ Pentest reports |
| gobuster | 1000-3000 | ❌ | ❌ | ❌ | ❌ |
| ffuf | 1500-4000 | ❌ | ❌ | ❌ | ❌ |
| feroxbuster | 2000-3000 | ⚠️ Basic | ⚠️ Partial | ❌ | ❌ |

---

## 📂 Project Structure

```
PathFinder/
├── main.go                    # Main application (scanner + TUI)
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── wordlist.txt               # 167 common web paths
├── build.bat                  # Build script (cross-platform)
├── SCAN.bat                   # Quick scan launcher
├── START_PATHFINDER.bat       # Simple launcher
├── README.md                  # This file
├── ROADMAP.md                 # Development roadmap and future features
├── TECHNICAL_NOTES.md         # Detailed implementation notes
└── THEME_ENHANCEMENTS.md      # Theme system documentation
```

---

## 🎯 Real-World Use Cases

### 1. Bug Bounty Recon
```bash
pathfinder.exe -target https://target.com \
  -wordlist SecLists-common.txt \
  -x php,js,txt,bak \
  -concurrency 100 \
  -o bug-bounty-results.json -of json
# Press F5 to export pentest report during scan
```

### 2. API Discovery
```bash
pathfinder.exe -target https://api.target.com \
  -wordlist api-endpoints.txt \
  -H "Authorization:Bearer TOKEN" \
  -mc 200,201,401
```

### 3. Admin Panel Hunting
```bash
pathfinder.exe -target https://target.com \
  -wordlist admin-paths.txt \
  -fc 404 \
  -mc 200,301,302,401,403
```

### 4. Authenticated Scanning
```bash
pathfinder.exe -target https://target.com \
  -cookie "session=abc123; auth=xyz" \
  -wordlist authenticated-paths.txt
```

---

## 🐛 Troubleshooting

### "Error loading wordlist"
- Make sure `wordlist.txt` is in the same directory as `pathfinder.exe`
- Or specify full path: `-wordlist C:\path\to\wordlist.txt`

### Terminal Display Issues
- Use Windows Terminal (recommended) or Command Prompt
- Avoid PowerShell ISE (limited terminal support)
- Make sure terminal supports 256 colors

### Scan Running Too Fast/Slow
```bash
# Too fast? Reduce concurrency and add rate limit
pathfinder.exe -target https://example.com -concurrency 10 -rate 50

# Too slow? Increase concurrency and reduce timeout
pathfinder.exe -target https://example.com -concurrency 200 -timeout 5

# Or use F4 to adjust settings in real-time!
```

### No Results Found
- Try different wordlist (SecLists recommended)
- Check target is reachable: `curl -I https://target.com`
- Use `-verbose` to see errors
- Verify you're not being rate limited

### Can't See All Results
- Use ↑/↓ arrow keys to scroll through live results
- Results are limited to last 100 in TUI (all saved to file on export)

---

## ⚖️ Legal & Ethical Use

⚠️ **WARNING:** Only scan systems you own or have explicit written permission to test.

### Legal Risks
- Unauthorized scanning may violate laws (CFAA in US, equivalents worldwide)
- May breach Terms of Service agreements
- Could be considered hostile network reconnaissance

### Best Practices
1. ✅ **Always** get written authorization before scanning
2. ✅ Start with low concurrency (`-concurrency 10`)
3. ✅ Use rate limiting on production systems (`-rate 50`)
4. ✅ Have emergency contact information ready
5. ✅ Stop immediately if requested

---

## 🔮 Roadmap & Future Features

See `ROADMAP.md` for complete development roadmap.

### Critical Priority
- [ ] Proxy Support (Burp Suite/ZAP integration)
- [ ] Response Body Regex Filtering
- [ ] Advanced Authentication (Basic Auth, Bearer tokens)
- [ ] Resume/Save State
- [ ] Better Smart 404 Detection

### High Priority
- [ ] Recursive Link Extraction
- [ ] Multiple Wordlist Support
- [ ] User-Agent Randomization
- [ ] Response Time Analysis
- [ ] Multi-Target Mode
- [ ] Technology Detection

### Completed ✅
- [x] Real-time TUI Dashboard
- [x] 8 Color Themes
- [x] Globe Visualization
- [x] Executive Summary Export (F5)
- [x] Interactive Config Menu (F4)
- [x] Built-in Help System (?)
- [x] Local Network Info Display
- [x] Results Scrolling (↑/↓)
- [x] Dynamic URL Input
- [x] Wildcard Detection
- [x] Redirect Chain Tracking

---

## 📝 Technical Details

### Built With
- **Language:** Go 1.16+
- **TUI Library:** [tcell/v2](https://github.com/gdamore/tcell) - Terminal cell management
- **HTTP Client:** Native `net/http` with custom transport
- **Concurrency:** Goroutines with semaphore pattern

### Architecture
- **Scanner Engine** - HTTP request handling, redirect tracking, wildcard detection
- **TUI System** - Full-screen dashboard with event-driven rendering
- **Globe Renderer** - 3D sphere projection with Earth bitmap
- **Theme System** - Modular color schemes (8 themes)
- **Statistics Tracker** - Thread-safe atomic counters for live updates
- **Report Generator** - Professional pentest documentation export

### Performance Optimizations
- Connection pooling (MaxIdleConnsPerHost)
- Goroutine semaphore for controlled concurrency
- Content hashing for duplicate detection
- Wildcard baseline caching
- Efficient TUI rendering (50ms refresh)

---

## 🤝 Contributing

Want to improve PathFinder? Check `ROADMAP.md` for planned features or suggest your own!

---

## 📜 License

Educational and authorized security testing only.

---

## 🙏 Acknowledgments

Inspired by:
- **gobuster** by OJ Reeves
- **ffuf** by Joona Hoikkala
- **feroxbuster** by epi052

TUI inspiration:
- **SecKC-MHN-Globe** - Globe visualization concept

---

## 📞 Support

Found a bug? Have a feature request?
- Check `TECHNICAL_NOTES.md` for implementation details
- Check `ROADMAP.md` for planned features
- Open an issue on GitHub

---

**Built with ⚡ in Go | Designed for 🎨 visual impact | Optimized for 🚀 speed**

**Happy hunting! (Responsibly, of course)**
