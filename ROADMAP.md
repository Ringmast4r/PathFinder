# PathFinder - Development Roadmap

## 🎯 Current Status
PathFinder v1.0.1 - Fully functional web path discovery tool with TUI, adaptive splash screen, 10 themes, BFS pathfinding animation, live results, and executive summary export.

---

## 🔴 CRITICAL PRIORITY - Must-Have Features

### 1. Proxy Support
**Status:** Not Started
**Importance:** Essential for pentest workflow
**Implementation:**
- Add `--proxy` flag to route traffic through Burp Suite/ZAP
- Support HTTP/HTTPS/SOCKS5 proxies
- Example: `--proxy http://127.0.0.1:8080`
- Modify HTTP client transport to use proxy

### 2. Advanced Authentication
**Status:** Partial (basic header support exists)
**Improvements Needed:**
- Basic Auth: `-u username:password` flag
- Bearer Token: Better token handling
- Digest Auth support
- NTLM Auth support
- Session cookie management

### 3. Response Body Regex Filtering
**Status:** Not Started
**Importance:** Much more powerful than status code filtering
**Implementation:**
- `--match-regex <pattern>` - Only show responses matching pattern
- `--filter-regex <pattern>` - Hide responses matching pattern
- Parse response body content for matches
- Use Go's `regexp` package

### 4. Resume/Save State
**Status:** Not Started
**Importance:** Critical for large scans
**Implementation:**
- Save scan state to JSON file periodically
- `--resume <state-file>` to continue interrupted scan
- Track completed paths to avoid duplicates
- Save on Ctrl+C gracefully

### 5. Better 404 Detection (Smart Soft-404)
**Status:** Basic wildcard detection exists
**Improvements Needed:**
- Multiple baseline comparisons
- Similarity scoring (Levenshtein distance)
- Dynamic page detection (timestamps, random strings)
- Pattern-based 404 identification

---

## 🟡 HIGH PRIORITY - Competitive Features

### 6. Recursive Link Extraction
**Status:** Framework exists (recursive flag present)
**Needs:** Full implementation
- Parse HTML responses for `<a href>`, `<link>`, `<script src>` tags
- Extract and queue new paths automatically
- Respect `--depth` flag for recursion levels
- De-duplicate discovered paths

### 7. Multiple Wordlist Support
**Status:** Single wordlist only
**Enhancement:**
- `--wordlist file1.txt,file2.txt,file3.txt`
- Merge and de-duplicate entries
- Chain wordlists sequentially

### 8. User-Agent Randomization
**Status:** Static UA only
**Implementation:**
- `--random-ua` flag
- Pool of realistic user agents (Chrome, Firefox, Safari, mobile)
- Rotate UA per request or per session
- Custom UA list from file

### 9. Response Time Analysis
**Status:** Response time tracked but not analyzed
**Enhancement:**
- Flag responses slower/faster than threshold
- `--response-time-threshold 5000` (ms)
- Show outliers in results
- Export time analysis in reports

### 10. Multi-Target Mode
**Status:** Single target only
**Enhancement:**
- `--target-list targets.txt`
- Scan multiple URLs from file
- Aggregate results per target
- Parallel target scanning

### 11. Technology Detection
**Status:** Not implemented
**Features:**
- Detect CMS (WordPress, Joomla, Drupal)
- Framework detection (Django, Laravel, Spring)
- Server fingerprinting
- Parse headers, response patterns, file paths

---

## 🟢 LOW PRIORITY - Nice to Have

### 12. Vhost/Subdomain Enumeration Mode
**Status:** Not implemented
- Virtual host discovery
- Subdomain brute forcing mode
- Different scanning approach vs path discovery

### 13. Screenshot Capture
**Status:** Not implemented
- Headless browser integration (chromedp)
- Capture screenshots of discovered pages
- Visual record for reports

### 14. HTTP/2 Support
**Status:** HTTP/1.1 only
- Upgrade to HTTP/2 for better performance
- Multiplexing support

### 15. Request Throttling Patterns
**Status:** Basic rate limiting exists
**Enhancement:**
- Random delays between requests
- Jitter patterns for stealth
- Adaptive throttling based on errors

---

## ✅ COMPLETED FEATURES

- [x] Adaptive splash screen (7 responsive logo variants, zero wrapping guarantee)
- [x] Real-time TUI with live results
- [x] 10 color themes (Matrix, Rainbow, Cyber, Blood, Skittles, Dark, Purple, Amber, White, Neon)
- [x] BFS pathfinding maze animation (synced with scan progress)
- [x] Globe visualization mode (F3)
- [x] Executive summary export (F5)
- [x] Local network information display (with F7 privacy toggle)
- [x] Wildcard detection and filtering
- [x] Redirect chain tracking
- [x] Interactive configuration menu (F4)
- [x] Built-in help screen (?)
- [x] Results scrolling (Up/Down arrows)
- [x] Concurrency control
- [x] Rate limiting
- [x] Multiple HTTP methods (GET, POST, HEAD, PUT, DELETE, PATCH)
- [x] Custom headers support
- [x] Content hash comparison
- [x] Status code filtering
- [x] Content length filtering
- [x] Response time tracking
- [x] JSON/CSV export options
- [x] Theme quick-switch (1-9, 0 keys)
- [x] Maze reset (F6)

---

## 📋 Development Notes

### Top 3 Recommended Next Steps:
1. **Proxy Support** - Most critical for pentest workflow
2. **Response Body Regex** - More powerful filtering
3. **User-Agent Randomization** - Stealth/evasion

### Architecture Considerations:
- Keep TUI responsive during long scans
- Maintain thread-safe operations for concurrent scanning
- Preserve current theme system and visual design
- Ensure backward compatibility with existing flags

### Testing Priorities:
- Proxy integration testing with Burp Suite
- Regex performance with large response bodies
- State save/resume reliability
- Recursive scanning depth limits

---

## 🤝 Contributing

Features marked as "Not Started" are open for implementation. Priority should be given to Critical and High Priority items.

**Last Updated:** 2025-10-12
