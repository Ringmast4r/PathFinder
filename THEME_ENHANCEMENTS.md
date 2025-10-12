# PathFinder Theme System Documentation

**Version:** 3.0.0 - DEATH STAR EDITION
**Last Updated:** 2025-10-12
**Current Theme Count:** 8 Themes

---

## Current Themes (Implemented)

### 1. Matrix Theme (Green)
- **Colors:** Green on black
- **Aesthetic:** Classic hacker terminal
- **Key:** Press `1` or cycle with `F1`

### 2. Rainbow Theme (Purple/Magenta)
- **Colors:** Purple, magenta, vibrant
- **Aesthetic:** Purple gradient vibes
- **Key:** Press `2` or cycle with `F1`
- **Note:** Currently purple-themed, future enhancement for full rainbow gradient

### 3. Cyber Theme (Cyan/Blue)
- **Colors:** Cyan, blue, electric
- **Aesthetic:** Cyberpunk/tech
- **Key:** Press `3` or cycle with `F1`

### 4. Blood Theme (Red) - DEFAULT
- **Colors:** Dark red, crimson
- **Aesthetic:** Aggressive, bold
- **Key:** Press `4` or cycle with `F1`
- **Default:** Loads on startup

### 5. Skittles Theme (Random)
- **Colors:** 6 random vibrant colors from 12-color palette
- **Aesthetic:** Chaotic, fun, each box different color
- **Key:** Press `5` or cycle with `F1`
- **Special:** F2 regenerates random colors

### 6. Dark Theme (Light Mode)
- **Colors:** White background, dark text
- **Aesthetic:** Light mode for bright environments
- **Key:** Press `6` or cycle with `F1`

### 7. Purple Theme
- **Colors:** Purple gradient, lavender text
- **Aesthetic:** Elegant purple aesthetic
- **Key:** Press `7` or cycle with `F1`

### 8. Amber Theme
- **Colors:** Orange/amber terminal colors
- **Aesthetic:** Classic amber monitor
- **Key:** Press `8` or cycle with `F1`

---

## Theme Architecture

### Theme Structure (`main.go:40-51`)

```go
type Theme struct {
    Name       string        // Theme name (displayed in UI)
    Background tcell.Color   // Background color
    Text       tcell.Color   // Default text color
    Primary    tcell.Color   // Primary highlight color
    Success    tcell.Color   // Success messages (200s)
    Warning    tcell.Color   // Warnings (redirects, 3xx)
    Danger     tcell.Color   // Errors (403, 401, 5xx)
    Info       tcell.Color   // Information text
    Border     tcell.Color   // Box borders
    Globe      tcell.Color   // Globe rendering color
}
```

### How Themes Work

1. **Static Themes (Matrix, Cyber, Blood, Dark, Purple, Amber)**
   - Predefined color palettes
   - Consistent colors across sessions
   - Apply to all UI elements

2. **Skittles Theme (Random)**
   - Generates 6 random colors from curated 12-color palette
   - Each box gets a different color from the set
   - Colors cached to prevent flickering
   - F2 key regenerates new random colors
   - Palette: Bright Red, Orange, Yellow, Green, Cyan, Blue, Violet, Magenta, Pink, Hot Pink, Lime, Red-Orange

3. **Theme Switching**
   - F1: Cycle through all 8 themes sequentially
   - 1-8: Jump directly to specific theme
   - F2: Activate Skittles + regenerate colors
   - Instant switching, no restart required

---

## Future Theme Enhancements (ROADMAP)

### Skittles Theme - Globe Enhancement

**Goal:** When Skittles theme is active, globe should have 16 random colored characters.

**Requirements:**
- 16 random positions within globe bounds
- 16 different vibrant colors
- Colors cached (no flicker)
- Only appears in Skittles theme

**Implementation Plan:**
1. Generate 16 random positions when switching to Skittles
2. Generate 16 random colors from vibrant palette
3. Store in `skittlesGlobeColors [16]tcell.Color` and `skittlesGlobePos [16][2]int`
4. In `renderGlobe()`, overlay colored characters if Skittles active

**Status:** ⏳ Planned (struct fields already exist)

---

### Rainbow Theme - Gradient Enhancement

**Goal:** True rainbow gradient across all UI elements (not just purple).

**Requirements:**
- Smooth rainbow gradient: Red → Orange → Yellow → Green → Cyan → Blue → Purple
- Apply to: box borders, text, globe
- Gradient flows left-to-right across screen
- Different from Skittles (ordered vs random)

**Implementation Plan:**

1. **Create Gradient Helper Function:**
   ```go
   func getRainbowColor(position float64) tcell.Color {
       // position: 0.0 (left) to 1.0 (right)
       // Returns RGB color from rainbow spectrum

       // Rainbow stops:
       // 0.00 - Red
       // 0.17 - Orange
       // 0.33 - Yellow
       // 0.50 - Green
       // 0.67 - Cyan
       // 0.83 - Blue
       // 1.00 - Purple
   }
   ```

2. **Apply to Box Borders:**
   - Calculate position along border perimeter
   - Get color from gradient function
   - Apply to each border character

3. **Apply to Globe:**
   - Use X position / width for gradient
   - Colorize globe characters horizontally

4. **Apply to Text:**
   - Calculate character position in string
   - Get gradient color
   - Apply to each character

**Status:** ⏳ Planned

---

## Theme Customization Guide

### Adding a New Theme

1. **Define Theme (`main.go` around line 55):**
   ```go
   ThemeYourName = Theme{
       Name:       "YOUR_NAME",
       Background: tcell.ColorBlack,
       Text:       tcell.ColorWhite,
       Primary:    tcell.NewRGBColor(R, G, B),
       Success:    tcell.NewRGBColor(R, G, B),
       Warning:    tcell.NewRGBColor(R, G, B),
       Danger:     tcell.NewRGBColor(R, G, B),
       Info:       tcell.NewRGBColor(R, G, B),
       Border:     tcell.NewRGBColor(R, G, B),
       Globe:      tcell.NewRGBColor(R, G, B),
   }
   ```

2. **Add to Theme Cycle (`main.go` in `cycleTheme()`):**
   ```go
   case "AMBER":
       CurrentTheme = ThemeYourName  // Add after last theme
   ```

3. **Add Number Key Shortcut:**
   ```go
   case '9':  // Use next available number
       CurrentTheme = ThemeYourName
   ```

4. **Update Controls Text:**
   ```
   "Theme: %s | F1: Cycle | 1-9 themes"
   ```

---

## Color Palette Reference

### RGB Values for Common Colors

**Reds:**
- Bright Red: `255, 0, 0`
- Dark Red: `180, 20, 20`
- Crimson: `220, 20, 60`

**Greens:**
- Matrix Green: `0, 255, 0`
- Dark Green: `0, 180, 0`
- Lime: `50, 205, 50`

**Blues:**
- Cyan: `0, 255, 255`
- Electric Blue: `0, 100, 255`
- Dark Cyan: `0, 200, 255`

**Purples:**
- Magenta: `255, 0, 255`
- Purple: `180, 100, 255`
- Violet: `138, 43, 226`

**Others:**
- Orange: `255, 165, 0`
- Yellow: `255, 255, 0`
- Amber: `255, 180, 0`
- White: `255, 255, 255`
- Black: `0, 0, 0`

---

## Theme Usage Statistics

**Most Popular (User Preference):**
1. Blood (Red) - DEFAULT
2. Matrix (Green) - Classic
3. Cyber (Cyan) - Modern
4. Skittles (Random) - Fun
5. Purple - Aesthetic

**Best for Different Scenarios:**
- **Pentesting:** Blood, Matrix, Cyber (dark backgrounds)
- **Presentations:** Dark (light mode), Purple (professional)
- **Fun/Casual:** Skittles (random colors)
- **Nostalgia:** Amber (classic terminal), Matrix (hacker)

---

## Implementation Priority

### ✅ Completed (v3.0.0)
- [x] 8 theme system fully functional
- [x] F1 theme cycling
- [x] Direct theme selection (1-8)
- [x] Skittles random color generation
- [x] F2 Skittles regeneration
- [x] Theme persistence across UI elements
- [x] Globe theme coloring
- [x] Splash screen theme support

### ⏳ Future Enhancements
- [ ] Skittles globe (16 random colored chars)
- [ ] Rainbow gradient implementation
- [ ] Theme customization via config file
- [ ] User-defined themes
- [ ] Theme preview mode
- [ ] Animated theme transitions

---

## Testing Checklist

### Current Themes (All Working ✅)
- [x] Matrix theme (green on black)
- [x] Rainbow theme (purple/magenta)
- [x] Cyber theme (cyan/blue)
- [x] Blood theme (red, default)
- [x] Skittles theme (6 random colors)
- [x] Dark theme (light mode)
- [x] Purple theme (purple gradient)
- [x] Amber theme (orange/amber)

### Theme Switching
- [x] F1 cycles through all 8 themes
- [x] Numbers 1-8 jump to specific themes
- [x] F2 activates Skittles + regenerates
- [x] Themes persist in globe mode
- [x] Themes apply to all UI elements
- [x] No flickering or glitches

---

## Notes

- **Skittles vs Rainbow:**
  - **Skittles** = Random chaotic colors (different each box)
  - **Rainbow** = Ordered spectrum gradient (flowing)

- **Performance:**
  - Static themes: Zero overhead
  - Skittles: Color generation on theme switch only
  - Rainbow (future): Gradient calc per-frame (negligible)

- **Accessibility:**
  - Dark theme for bright environments
  - High contrast in all themes
  - Colorblind-friendly options (Matrix green, Amber)

---

**Theme System Status: ✅ Production Ready**

*8 themes implemented and tested. Future enhancements planned for Skittles globe and Rainbow gradient.*
