# PathFinder - Changelog

## Version 1.0.1 (2025-10-12) - Adaptive Splash Screen

### 🎨 Splash Screen Improvements
- **Fully Adaptive Logo Rendering**: Splash screen now automatically scales to any terminal size
  - 7 logo size variants from full ASCII art (110 chars) down to minimal "//" (2 chars)
  - Dynamic logo selection based on terminal width
  - Diagonal slash "/" effect preserved across all sizes
  - Zero wrapping guaranteed - always displays cleanly
- **Immediate Screen Clearing**: TUI takes over full screen instantly on startup
- **Smart Centering**: Logos center when space allows, with wrapping protection
- **Logo Variants**:
  - Full ASCII art for wide screens (110+ chars)
  - Compact ASCII art for medium screens (39+ chars)
  - Single-line text with slashes (24+ chars)
  - Diagonal slash variants (17+ chars, 13+ chars)
  - Minimal brand name (14+ chars)
  - Icon-only "//" for tiny terminals

### 🐛 Bug Fixes
- Fixed splash screen starting at command prompt position instead of full screen
- Fixed logo wrapping at normal CMD window sizes
- Removed overly conservative width checking that prevented logos from displaying
- Corrected boundary checking to allow logos that fit exactly

---

## Version 1.0.0 (2025-10-12) - Initial Public Release

### 🎨 Visual Enhancements
- **10 Color Themes**: Added White and Neon themes (total: Matrix, Rainbow, Cyber, Blood, Skittles, Dark, Purple, Amber, White, Neon)
- **BFS Pathfinding Maze Animation**: Real-time animated breadth-first search visualization
  - Green 'S' start → Red 'X' target
  - Blue dots show exploration phase (searching for paths)
  - Yellow stars draw the solution path
  - Animation speed matches actual scan speed (realistic behavior)
  - Auto-loops in idle mode, syncs with scan progress during active scans
- **Splash Screen**: Removed "DEATH STAR EDITION" subtitle for cleaner branding

### ⌨️ Keyboard Controls (Industry Standard Remapping)
- **F1**: Help screen (INDUSTRY STANDARD - changed from theme cycling)
- **F2**: Privacy toggle / Hide network info (moved from F7 for easier access)
- **F3**: Regenerate Skittles colors (dedicated key)
- **F4**: Config menu (unchanged)
- **F5**: Export report (unchanged)
- **F6**: Reset pathfinding maze
- **` (backtick)**: Cycle themes (promoted from F1)
- **\\(backslash)**: Globe visualization (EASTER EGG - hidden, not documented)
- **?**: Help screen (alternative to F1)
- **1-9, 0**: Direct theme selection

### 🔒 Privacy Features
- **Network Info Privacy Toggle (F7)**: Hide local network information for screenshots, recordings, or OpSec scenarios
  - When hidden, pathfinding maze moves up to fill the space
  - Seamless visual transition

### 🎯 Animation Behavior
- **Idle Mode**: Smooth 12-15 second pathfinding demonstration that auto-loops
- **Scan Mode**: Animation speed realistically matches scan speed
  - Fast scans (3 seconds) → Fast pathfinding animation (3 seconds)
  - Slow scans (60 seconds) → Slow pathfinding animation (60 seconds)
  - Blue exploration phase: 0-80% of scan progress
  - Yellow path drawing: 80-100% of scan progress

### 📊 Status Display
- Changed "Searching..." to "Pathfinding..." for maze animation status
- Shows both maze progress and scan progress during active scans
- "Complete! Press F6 to reset" when pathfinding finishes

### 🐛 Bug Fixes
- Fixed maze animation sync with scan completion percentage
- Optimized frame counts for smooth animation (10,000-15,000 frames)
- Fixed maze end marker from 'E' to 'X' with proper red coloring
- Corrected idle animation speed for better visibility

### 📝 Documentation
- Updated README.md with F6, F7 keyboard controls
- Updated ROADMAP.md with completed features
- Comprehensive TECHNICAL_NOTES.md updates
- Added detailed BFS animation documentation

---

**Built by ringmast4r**
