package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	Version            = "3.0.0"
	UserAgent          = "PathFinder/3.0 (Security Research)"
	DefaultConcurrency = 50
	DefaultTimeout     = 10
	MinRedirectCount   = 3
	MaxRedirects       = 10
)

// ===========================================================================
// THEME SYSTEM
// ===========================================================================

type Theme struct {
	Name       string
	Background tcell.Color
	Text       tcell.Color
	Primary    tcell.Color
	Success    tcell.Color
	Warning    tcell.Color
	Danger     tcell.Color
	Info       tcell.Color
	Border     tcell.Color
	Globe      tcell.Color
}

var (
	ThemeMatrix = Theme{
		Name:       "MATRIX",
		Background: tcell.ColorBlack,
		Text:       tcell.ColorWhite,
		Primary:    tcell.NewRGBColor(0, 255, 0),
		Success:    tcell.NewRGBColor(0, 255, 0),
		Warning:    tcell.ColorYellow,
		Danger:     tcell.NewRGBColor(255, 0, 0),
		Info:       tcell.NewRGBColor(0, 255, 255),
		Border:     tcell.NewRGBColor(0, 180, 0),
		Globe:      tcell.NewRGBColor(0, 255, 0),
	}

	ThemeRainbow = Theme{
		Name:       "RAINBOW",
		Background: tcell.ColorBlack,
		Text:       tcell.ColorWhite,
		Primary:    tcell.NewRGBColor(255, 0, 255),
		Success:    tcell.NewRGBColor(0, 255, 0),
		Warning:    tcell.NewRGBColor(255, 255, 0),
		Danger:     tcell.NewRGBColor(255, 0, 0),
		Info:       tcell.NewRGBColor(0, 255, 255),
		Border:     tcell.NewRGBColor(255, 0, 255),
		Globe:      tcell.NewRGBColor(255, 0, 255),
	}

	ThemeCyber = Theme{
		Name:       "CYBER",
		Background: tcell.ColorBlack,
		Text:       tcell.ColorWhite,
		Primary:    tcell.NewRGBColor(0, 255, 255),
		Success:    tcell.NewRGBColor(0, 255, 0),
		Warning:    tcell.NewRGBColor(255, 255, 0),
		Danger:     tcell.NewRGBColor(255, 0, 0),
		Info:       tcell.NewRGBColor(0, 255, 255),
		Border:     tcell.NewRGBColor(0, 200, 255),
		Globe:      tcell.NewRGBColor(0, 255, 255),
	}

	ThemeBlood = Theme{
		Name:       "BLOOD",
		Background: tcell.ColorBlack,
		Text:       tcell.ColorWhite,
		Primary:    tcell.NewRGBColor(255, 0, 0),
		Success:    tcell.NewRGBColor(0, 255, 0),
		Warning:    tcell.NewRGBColor(255, 255, 0),
		Danger:     tcell.NewRGBColor(255, 0, 0),
		Info:       tcell.NewRGBColor(255, 100, 100),
		Border:     tcell.NewRGBColor(200, 0, 0),
		Globe:      tcell.NewRGBColor(255, 0, 0),
	}

	ThemeSkittles = Theme{
		Name:       "SKITTLES",
		Background: tcell.ColorBlack,
		Text:       tcell.NewRGBColor(255, 255, 255),
		Primary:    tcell.NewRGBColor(255, 0, 255),
		Success:    tcell.NewRGBColor(0, 255, 0),
		Warning:    tcell.NewRGBColor(255, 255, 0),
		Danger:     tcell.NewRGBColor(255, 0, 0),
		Info:       tcell.NewRGBColor(0, 255, 255),
		Border:     tcell.NewRGBColor(255, 105, 180),
		Globe:      tcell.NewRGBColor(255, 0, 255),
	}

	ThemeDark = Theme{
		Name:       "DARK",
		Background: tcell.NewRGBColor(255, 255, 255),
		Text:       tcell.NewRGBColor(0, 0, 0),
		Primary:    tcell.NewRGBColor(30, 70, 170),
		Success:    tcell.NewRGBColor(0, 120, 30),
		Warning:    tcell.NewRGBColor(180, 100, 0),
		Danger:     tcell.NewRGBColor(180, 20, 20),
		Info:       tcell.NewRGBColor(0, 90, 150),
		Border:     tcell.NewRGBColor(100, 100, 110),
		Globe:      tcell.NewRGBColor(30, 70, 170),
	}

	ThemePurple = Theme{
		Name:       "PURPLE",
		Background: tcell.ColorBlack,
		Text:       tcell.NewRGBColor(220, 180, 255),
		Primary:    tcell.NewRGBColor(180, 100, 255),
		Success:    tcell.NewRGBColor(150, 255, 150),
		Warning:    tcell.NewRGBColor(255, 180, 100),
		Danger:     tcell.NewRGBColor(255, 100, 150),
		Info:       tcell.NewRGBColor(200, 150, 255),
		Border:     tcell.NewRGBColor(150, 80, 200),
		Globe:      tcell.NewRGBColor(180, 100, 255),
	}

	ThemeAmber = Theme{
		Name:       "AMBER",
		Background: tcell.ColorBlack,
		Text:       tcell.NewRGBColor(255, 180, 0),
		Primary:    tcell.NewRGBColor(255, 200, 50),
		Success:    tcell.NewRGBColor(255, 220, 100),
		Warning:    tcell.NewRGBColor(255, 150, 0),
		Danger:     tcell.NewRGBColor(200, 100, 0),
		Info:       tcell.NewRGBColor(255, 190, 70),
		Border:     tcell.NewRGBColor(200, 140, 0),
		Globe:      tcell.NewRGBColor(255, 180, 0),
	}

	CurrentTheme = ThemeBlood
)

// ===========================================================================
// DATA STRUCTURES
// ===========================================================================

type RedirectStep struct {
	URL    string
	Status int
}

type ScanResult struct {
	OriginalPath    string
	OriginalURL     string
	FinalStatus     int
	FinalURL        string
	RedirectChain   []RedirectStep
	ContentLength   int
	ContentHash     string
	IsDirect200     bool
	ResponseTime    time.Duration
	Timestamp       time.Time
}

type LiveStats struct {
	mu                sync.RWMutex
	TotalRequests     int64
	CompletedRequests int64
	Direct200s        int64
	Redirects         int64
	Errors            int64
	Protected         int64
	CurrentSpeed      float64
	StartTime         time.Time
	LastUpdate        time.Time
}

type Statistics struct {
	mu              sync.Mutex
	TotalScanned    int
	Direct200s      []*ScanResult
	Redirects       []*ScanResult
	RedirectTargets map[string]int
	ContentHashes   map[string][]*ScanResult
	OtherCodes      []*ScanResult
}

type WildcardBaseline struct {
	Hash   string
	Length int
	Status int
}

type Config struct {
	StatusCodes    []int
	FilterStatuses []int
	FilterSizes    []int
	MatchRegex     *regexp.Regexp
	Extensions     []string
	CustomHeaders  map[string]string
	Cookie         string
	Method         string
	RateLimit      int
	Delay          time.Duration
	Recursive      bool
	RecursionDepth int
	OutputFile     string
	OutputFormat   string
	Theme          string
}

type Scanner struct {
	BaseURL          string
	Concurrency      int
	Timeout          time.Duration
	Verbose          bool
	Client           *http.Client
	Stats            *Statistics
	LiveStats        *LiveStats
	WildcardBaseline *WildcardBaseline
	Config           *Config
	visitedPaths     map[string]bool
	pathMutex        sync.Mutex
	recursionQueue   chan string
	recursionWg      sync.WaitGroup
	rateLimiter      <-chan time.Time
	lastResults      []*ScanResult
	resultsMutex     sync.Mutex
}

// ===========================================================================
// GLOBE RENDERING
// ===========================================================================

type Globe struct {
	Width       int
	Height      int
	Rotation    float64
	AspectRatio float64
	Radius      float64
	EarthMap    []string
	MapWidth    int
	MapHeight   int
}

func getEarthBitmap() []string {
	return []string{
		"                                                                                                                        ",
		"                                                                                                                        ",
		"                                                                                                                        ",
		"                             # ####### #################                                    #                           ",
		"                       #    #   ### #################            ###                                                    ",
		"                      ###  ## ####       ############ #                        ##         ########        #####         ",
		"                  ## ###   #  ### ##      ###########                         #    #### ################   ###          ",
		"      ######## ###### #### # #  #  ###     #########              #######        # ## ##################################",
		" ### ###########################    ####   #####      #          ####### ###############################################",
		"      ########################       ##    ####                #### ####################################################",
		"      ### # #################      ##        #                ##### # ##########################################  ##    ",
		"                ##############     #####                   #     #  #######################################      ##     ",
		"                 ################ #######                # #   ###########################################      ##      ",
		"                  ########################                 ################################################             ",
		"                    ###################  ##                ################################################             ",
		"                   ################### #                    ##########  ####  ############################              ",
		"                   ##################                    ##### ##  ###    ### ##########################                ",
		"                   #################                     ###       # ######## ######################  #    #            ",
		"                    ###############                       #  ###       ##############################  #  #             ",
		"                     #############                        ######        #############################                   ",
		"                       ######## #                        ############################################                   ",
		"                      # ####     #                      ##################### #######################                   ",
		"                       # ###      #                    ################# ######    #################                    ",
		"                         ###  #   #                    ################## ######     ####  #####                        ",
		"                          #####   # #                  ################## #####      ###    ####                        ",
		"                             ####                      ################### ###       ##      ####   #                   ",
		"                               #    #                  ####################           #      # ##                       ",
		"                                #  #####                #####################         #      # #     ##                 ",
		"                                   ######                #### ###############          #      #    #                    ",
		"                                   ########                     ############                 ##   ##                    ",
		"                                  #########                     ###########                   #  ####                   ",
		"                                  #############                 ##########                    ##### #     ##            ",
		"                                 ################                ########                                  ## #         ",
		"                                  ###############                #########                         ## #    # #          ",
		"                                   #############                 #########                                              ",
		"                                   ############                  #########  #                         # ##  #           ",
		"                                     ##########                 #########  ##                        ########           ",
		"                                     ##########                  #######   ##                      ###########     #    ",
		"                                     ########                    #######   #                      #############         ",
		"                                     #######                     ######                           ##############        ",
		"                                     #######                      #####                            #############        ",
		"                                     ######                       ####                             ###   ######         ",
		"                                    #####                                                                  ####       # ",
		"                                    #####                                                                              #",
		"                                    ###                                                                      #        # ",
		"                                    ###                                                                             ##  ",
		"                                    ##                                                                                  ",
		"                                   ##                                                                                   ",
		"                                    ##                                                                                  ",
		"                                                                                                                        ",
		"                                                                                                                        ",
		"                                                                                                                        ",
		"                                       #                                                                                ",
		"                                      #                                #  ##########   ########################         ",
		"                                   #####                 ########################## #################################   ",
		"                  # ## #   #############              #############################################################     ",
		"        ## #########################             ##################################################################     ",
		"           ######################## #  #  ##     #################################################################      ",
		"    ##################################################################################################################  ",
		"########################################################################################################################",
	}
}

func NewGlobe(width, height int) *Globe {
	earthMap := getEarthBitmap()
	effectiveHeight := float64(height) * 2.0
	radius := math.Min(float64(width)/2.5, effectiveHeight/2.5)

	return &Globe{
		Width:       width,
		Height:      height,
		Rotation:    0,
		AspectRatio: 2.0,
		Radius:      radius,
		EarthMap:    earthMap,
		MapWidth:    len(earthMap[0]),
		MapHeight:   len(earthMap),
	}
}

func (g *Globe) sampleEarthAt(lat, lon float64) rune {
	latNorm := (lat + 90) / 180
	lonNorm := (lon + 180) / 360

	y := int(latNorm * float64(g.MapHeight-1))
	x := int(lonNorm * float64(g.MapWidth-1))

	if y < 0 {
		y = 0
	}
	if y >= g.MapHeight {
		y = g.MapHeight - 1
	}
	if x < 0 {
		x = 0
	}
	if x >= g.MapWidth {
		x = g.MapWidth - 1
	}

	return rune(g.EarthMap[y][x])
}

func densityToASCII(density float64) rune {
	if density > 1.0 {
		return '@'
	} else if density > 0.8 {
		return '#'
	} else if density > 0.6 {
		return '%'
	} else if density > 0.4 {
		return 'o'
	} else if density > 0.3 {
		return '='
	} else if density > 0.2 {
		return '+'
	} else if density > 0.15 {
		return '-'
	} else if density > 0.1 {
		return '.'
	} else if density > 0.05 {
		return '`'
	}
	return ' '
}

func (g *Globe) Render() [][]rune {
	if g.Width <= 0 || g.Height <= 0 {
		return [][]rune{[]rune{' '}}
	}

	screen := make([][]rune, g.Height)
	for i := range screen {
		screen[i] = make([]rune, g.Width)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}

	density := make([][]float64, g.Height)
	for i := range density {
		density[i] = make([]float64, g.Width)
	}

	centerX, centerY := g.Width/2, g.Height/2
	effectiveRadius := g.Radius

	// Forward rendering: for each screen pixel, calculate what earth coordinate it represents
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			dx := float64(x - centerX)
			dy := float64(y-centerY) * g.AspectRatio
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= effectiveRadius {
				// Normalize to unit sphere
				nx := dx / effectiveRadius
				ny := dy / effectiveRadius

				// Calculate z coordinate (depth into sphere)
				nz_squared := 1 - nx*nx - ny*ny
				if nz_squared >= 0 {
					nz := math.Sqrt(nz_squared)

					// Convert to latitude/longitude
					lat := math.Asin(ny) * 180 / math.Pi
					lon := math.Atan2(nx, nz)*180/math.Pi + g.Rotation*180/math.Pi

					// Normalize longitude to -180 to 180
					for lon < -180 {
						lon += 360
					}
					for lon > 180 {
						lon -= 360
					}

					// Sample earth bitmap
					earthChar := g.sampleEarthAt(lat, lon)
					if earthChar != ' ' {
						baseDensity := 1.0
						switch earthChar {
						case '#':
							baseDensity = 1.0
						case '.':
							baseDensity = 0.6
						default:
							baseDensity = 0.8
						}

						density[y][x] += baseDensity

						// Anti-aliasing: slightly brighten neighboring pixels
						for dy := -1; dy <= 1; dy++ {
							for dx := -1; dx <= 1; dx++ {
								nx2, ny2 := x+dx, y+dy
								if nx2 >= 0 && nx2 < g.Width && ny2 >= 0 && ny2 < g.Height {
									density[ny2][nx2] += 0.05
								}
							}
						}
					}
				}
			}

			// Draw edge/outline of sphere
			if distance > effectiveRadius-0.5 && distance < effectiveRadius+0.5 {
				density[y][x] += 0.2
			}
		}
	}

	// Convert density to characters
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			screen[y][x] = densityToASCII(density[y][x])
		}
	}

	return screen
}

// ===========================================================================
// TUI DASHBOARD
// ===========================================================================

type TUI struct {
	screen              tcell.Screen
	width               int
	height              int
	scanner             *Scanner
	globe               *Globe
	showGlobe           bool
	running             bool
	mutex               sync.RWMutex
	skittlesBoxColors   [6]tcell.Color  // Store random colors for Skittles theme
	skittlesGlobeColors [16]tcell.Color // 16 random colors for Skittles globe
	skittlesGlobePos    [16][2]int      // Random positions for colored globe chars
	lastTheme           string
	showSplash          bool
	splashProgress      float64 // 0.0 to 1.0 for animation
	progressBarFrame    int     // Animation frame for progress bar
	inputText           string  // User input text for URLs/domains
	inputActive         bool    // Whether input field is active
	showConfigMenu      bool    // Whether config menu is visible
	configMenuSelected  int     // Currently selected menu item
	configEditMode      bool    // Whether editing a config value
	configEditText      string  // Temporary text while editing
	showHelpScreen      bool    // Whether help screen is visible
	resultsScrollOffset int     // Scroll offset for live results
}

func NewTUI(scanner *Scanner) (*TUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
	screen.Clear()

	width, height := screen.Size()

	return &TUI{
		screen:         screen,
		width:          width,
		height:         height,
		scanner:        scanner,
		globe:          NewGlobe(60, 25),
		showGlobe:      false,
		running:        true,
		showSplash:     true,
		splashProgress: 0.0,
	}, nil
}

func (tui *TUI) drawText(x, y int, text string, style tcell.Style) {
	for i, r := range text {
		if x+i < tui.width {
			tui.screen.SetContent(x+i, y, r, nil, style)
		}
	}
}

func (tui *TUI) drawBox(x, y, width, height int, title string, style tcell.Style) {
	// Top border
	tui.screen.SetContent(x, y, '╔', nil, style)
	for i := 1; i < width-1; i++ {
		tui.screen.SetContent(x+i, y, '═', nil, style)
	}
	tui.screen.SetContent(x+width-1, y, '╗', nil, style)

	// Title
	if title != "" {
		titleText := fmt.Sprintf(" %s ", title)
		titleX := x + (width-len(titleText))/2
		tui.drawText(titleX, y, titleText, style.Bold(true))
	}

	// Sides
	for i := 1; i < height-1; i++ {
		tui.screen.SetContent(x, y+i, '║', nil, style)
		tui.screen.SetContent(x+width-1, y+i, '║', nil, style)
	}

	// Bottom border
	tui.screen.SetContent(x, y+height-1, '╚', nil, style)
	for i := 1; i < width-1; i++ {
		tui.screen.SetContent(x+i, y+height-1, '═', nil, style)
	}
	tui.screen.SetContent(x+width-1, y+height-1, '╝', nil, style)
}

func (tui *TUI) generateSkittlesColors() {
	// Bold, vibrant color palette with high contrast
	colorPalette := []tcell.Color{
		tcell.NewRGBColor(255, 0, 0),     // Bright Red
		tcell.NewRGBColor(255, 165, 0),   // Orange
		tcell.NewRGBColor(255, 255, 0),   // Yellow
		tcell.NewRGBColor(0, 255, 0),     // Bright Green
		tcell.NewRGBColor(0, 255, 255),   // Cyan
		tcell.NewRGBColor(0, 100, 255),   // Bright Blue
		tcell.NewRGBColor(138, 43, 226),  // Blue Violet
		tcell.NewRGBColor(255, 0, 255),   // Magenta
		tcell.NewRGBColor(255, 20, 147),  // Deep Pink
		tcell.NewRGBColor(255, 105, 180), // Hot Pink
		tcell.NewRGBColor(50, 205, 50),   // Lime Green
		tcell.NewRGBColor(255, 69, 0),    // Red Orange
	}

	// Shuffle and pick 6 random colors from the palette
	used := make(map[int]bool)
	for i := 0; i < 6; i++ {
		var idx int
		for {
			idx = rand.Intn(len(colorPalette))
			if !used[idx] {
				used[idx] = true
				break
			}
		}
		tui.skittlesBoxColors[i] = colorPalette[idx]
	}
}

func (tui *TUI) Render() {
	tui.screen.Clear()

	if tui.showSplash {
		tui.renderSplashScreen()
	} else if tui.showGlobe {
		tui.renderGlobe()
	} else if tui.showHelpScreen {
		tui.renderDashboard()
		tui.renderHelpScreen()
	} else {
		tui.renderDashboard()
		if tui.showConfigMenu {
			tui.renderConfigMenu()
		}
	}

	tui.screen.Show()
}

func (tui *TUI) renderSplashScreen() {
	// ASCII art for /PATHFINDER - using standard ASCII characters
	logo := []string{
		"",
		"      //   ########     ###    ########  ##     ## ######## #### ##    ## ########  ######## ########  ",
		"     //    ##     ##   ## ##      ##     ##     ## ##        ##  ###   ## ##     ## ##       ##     ## ",
		"    //     ##     ##  ##   ##     ##     ##     ## ##        ##  ####  ## ##     ## ##       ##     ## ",
		"   //      ########  ##     ##    ##     ######### ######    ##  ## ## ## ##     ## ######   ########  ",
		"  //       ##        #########    ##     ##     ## ##        ##  ##  #### ##     ## ##       ##   ##   ",
		" //        ##        ##     ##    ##     ##     ## ##        ##  ##   ### ##     ## ##       ##    ##  ",
		"//         ##        ##     ##    ##     ##     ## ##       #### ##    ## ########  ######## ##     ## ",
		"",
	}

	subtitle := "DEATH STAR EDITION"
	tagline := "Web Path Discovery & Reconnaissance Tool"
	author := "by Ringmast4r"
	versionText := "v" + Version

	// Calculate center position
	logoHeight := len(logo)
	startY := (tui.height - logoHeight - 6) / 2
	if startY < 1 {
		startY = 1
	}

	// Calculate total characters in logo for trickle effect
	totalChars := 0
	for _, line := range logo {
		totalChars += len(line)
	}

	// Trickle effect phase: 0.0 to 1.0 reveals all characters
	// Hold phase: 1.0 to 2.5 holds for 3 seconds at 20fps (60 frames)
	if tui.splashProgress <= 1.0 {
		// Trickle phase - reveal characters gradually
		visibleChars := int(tui.splashProgress * float64(totalChars))

		charCount := 0
		for i, line := range logo {
			y := startY + i
			x := (tui.width - len(line)) / 2
			if x < 0 {
				x = 0
			}

			for j, ch := range line {
				charCount++
				if charCount <= visibleChars {
					style := tcell.StyleDefault.Foreground(CurrentTheme.Primary).Bold(true)
					tui.screen.SetContent(x+j, y, ch, nil, style)
				}
			}
		}

		// Show subtitle/author once logo is complete
		if tui.splashProgress >= 0.95 {
			subtitleY := startY + logoHeight + 2
			subtitleX := (tui.width - len(subtitle)) / 2
			tui.drawText(subtitleX, subtitleY, subtitle, tcell.StyleDefault.Foreground(CurrentTheme.Success).Bold(true))

			taglineY := subtitleY + 1
			taglineX := (tui.width - len(tagline)) / 2
			tui.drawText(taglineX, taglineY, tagline, tcell.StyleDefault.Foreground(CurrentTheme.Info))

			authorY := taglineY + 2
			authorX := (tui.width - len(author)) / 2
			tui.drawText(authorX, authorY, author, tcell.StyleDefault.Foreground(CurrentTheme.Text))

			versionY := authorY + 1
			versionX := (tui.width - len(versionText)) / 2
			tui.drawText(versionX, versionY, versionText, tcell.StyleDefault.Foreground(CurrentTheme.Text).Dim(true))
		}

		// Slow trickle - takes about 1.5 seconds to complete at 20fps
		tui.splashProgress += 0.015
	} else {
		// Hold phase - show complete splash for 3-4 seconds
		// Draw complete logo
		for i, line := range logo {
			y := startY + i
			x := (tui.width - len(line)) / 2
			if x < 0 {
				x = 0
			}
			for j, ch := range line {
				style := tcell.StyleDefault.Foreground(CurrentTheme.Primary).Bold(true)
				tui.screen.SetContent(x+j, y, ch, nil, style)
			}
		}

		// Draw all text
		subtitleY := startY + logoHeight + 2
		subtitleX := (tui.width - len(subtitle)) / 2
		tui.drawText(subtitleX, subtitleY, subtitle, tcell.StyleDefault.Foreground(CurrentTheme.Success).Bold(true))

		taglineY := subtitleY + 1
		taglineX := (tui.width - len(tagline)) / 2
		tui.drawText(taglineX, taglineY, tagline, tcell.StyleDefault.Foreground(CurrentTheme.Info))

		authorY := taglineY + 2
		authorX := (tui.width - len(author)) / 2
		tui.drawText(authorX, authorY, author, tcell.StyleDefault.Foreground(CurrentTheme.Text))

		versionY := authorY + 1
		versionX := (tui.width - len(versionText)) / 2
		tui.drawText(versionX, versionY, versionText, tcell.StyleDefault.Foreground(CurrentTheme.Text).Dim(true))

		// Increment through hold period
		// At 20fps (50ms/frame), 3.5 seconds = 70 frames
		// Progress from 1.0 to 2.5 = 1.5 units / 70 frames = 0.0214 per frame
		tui.splashProgress += 0.0214

		// After 3.5 seconds of hold, transition to dashboard
		if tui.splashProgress >= 2.5 {
			tui.showSplash = false
		}
	}
}

func (tui *TUI) renderGlobe() {
	// Clear screen
	tui.screen.Clear()

	// Render globe
	globeScreen := tui.globe.Render()
	globeStartX := (tui.width - tui.globe.Width) / 2
	globeStartY := (tui.height - tui.globe.Height) / 2

	globeStyle := tcell.StyleDefault.Foreground(CurrentTheme.Globe)
	for y, row := range globeScreen {
		for x, char := range row {
			if char != ' ' {
				tui.screen.SetContent(globeStartX+x, globeStartY+y, char, nil, globeStyle)
			}
		}
	}

	// Globe title
	title := "[ PATHFINDER GLOBE MODE ]"
	titleX := (tui.width - len(title)) / 2
	tui.drawText(titleX, globeStartY-2, title, tcell.StyleDefault.Foreground(CurrentTheme.Primary).Bold(true))

	// Instructions
	instructions := "Press F3 to exit Globe Mode | Press Q to quit"
	instX := (tui.width - len(instructions)) / 2
	tui.drawText(instX, globeStartY+tui.globe.Height+2, instructions, tcell.StyleDefault.Foreground(CurrentTheme.Info))

	// Update rotation - complete rotation in 30 seconds
	// At 50ms per frame (20 fps), increment is: 2*pi / (30 * 20) = 0.01047 radians
	// Negative for west-to-east rotation (matching SecKC-MHN-Globe)
	tui.globe.Rotation -= 0.01047
}

func (tui *TUI) renderDashboard() {
	boxStyle := tcell.StyleDefault.Foreground(CurrentTheme.Border)
	textStyle := tcell.StyleDefault.Foreground(CurrentTheme.Text)

	// Generate random colors for Skittles theme only when switching to it
	if CurrentTheme.Name == "SKITTLES" && tui.lastTheme != "SKITTLES" {
		tui.generateSkittlesColors()
		tui.lastTheme = "SKITTLES"
	} else if CurrentTheme.Name != "SKITTLES" {
		tui.lastTheme = CurrentTheme.Name
	}

	// For Skittles theme, use stored random colors for each box
	var titleBoxStyle, targetBoxStyle, progressBoxStyle, statsBoxStyle, resultsBoxStyle, controlsBoxStyle tcell.Style

	if CurrentTheme.Name == "SKITTLES" {
		titleBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[0])
		targetBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[1])
		progressBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[2])
		statsBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[3])
		resultsBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[4])
		controlsBoxStyle = tcell.StyleDefault.Foreground(tui.skittlesBoxColors[5])
	} else {
		titleBoxStyle = boxStyle
		targetBoxStyle = boxStyle
		progressBoxStyle = boxStyle
		statsBoxStyle = boxStyle
		resultsBoxStyle = boxStyle
		controlsBoxStyle = boxStyle
	}

	// Title box
	titleWidth := tui.width - 4
	titleText := fmt.Sprintf("PATHFINDER v%s - Web Path Discovery Tool", Version)
	tui.drawBox(2, 0, titleWidth, 3, titleText, titleBoxStyle)

	// Add flashing author name in red - centered with gentle pulse
	authorName := "by ringmast4r"
	authorX := 2 + (titleWidth-len(authorName))/2
	// Slower, gentler pulse: range from 200-255 instead of 128-255
	flashIntensity := int((math.Sin(float64(tui.progressBarFrame)/25.0) + 1.0) * 27.5)
	flashColor := tcell.NewRGBColor(int32(flashIntensity+200), 0, 0)
	tui.drawText(authorX, 1, authorName, tcell.StyleDefault.Foreground(flashColor).Bold(true))

	// Add scanning status indicator next to author name if scan is running
	completed := atomic.LoadInt64(&tui.scanner.LiveStats.CompletedRequests)
	total := tui.scanner.LiveStats.TotalRequests
	if completed > 0 && completed < total {
		scanStatus := " ⚡ SCANNING"
		statusX := authorX + len(authorName) + 2
		tui.drawText(statusX, 1, scanStatus, tcell.StyleDefault.Foreground(CurrentTheme.Primary).Bold(true))
	}

	// Input field box - NOW AT TOP for prominence
	var inputBoxStyle tcell.Style
	if CurrentTheme.Name == "SKITTLES" {
		inputBoxStyle = titleBoxStyle
	} else {
		inputBoxStyle = boxStyle
	}

	inputBoxTitle := "INPUT URL/DOMAIN"
	if tui.inputActive {
		inputBoxTitle = "INPUT URL/DOMAIN (ACTIVE - Press Enter to submit)"
	}
	tui.drawBox(2, 3, titleWidth/2-1, 4, inputBoxTitle, inputBoxStyle)

	// Draw input text with cursor if active
	inputDisplay := tui.inputText
	var inputTextStyle tcell.Style
	if tui.inputActive {
		inputDisplay += "_" // Show cursor
		inputTextStyle = textStyle
	} else if inputDisplay == "" {
		inputDisplay = "Press Enter to activate..."
		// Gentle pulsing white glow effect for placeholder text
		glowIntensity := int((math.Sin(float64(tui.progressBarFrame)/20.0) + 1.0) * 60.0)
		glowColor := tcell.NewRGBColor(int32(glowIntensity+150), int32(glowIntensity+150), int32(glowIntensity+150))
		inputTextStyle = tcell.StyleDefault.Foreground(glowColor).Italic(true)
	} else {
		inputTextStyle = textStyle
	}
	tui.drawText(4, 5, truncateString(inputDisplay, titleWidth/2-5), inputTextStyle)

	// Scan Config box - combines URL and all scan settings
	tui.drawBox(2, 7, titleWidth/2-1, 8, "SCAN CONFIG", targetBoxStyle)
	targetText := fmt.Sprintf("URL: %s", truncateString(tui.scanner.BaseURL, titleWidth/2-10))
	tui.drawText(4, 8, targetText, textStyle)

	methodText := fmt.Sprintf("Method: %s  Timeout: %ds", tui.scanner.Config.Method, int(tui.scanner.Timeout.Seconds()))
	tui.drawText(4, 9, methodText, tcell.StyleDefault.Foreground(CurrentTheme.Info))

	timeoutHelp := "(max wait per request)"
	tui.drawText(4, 10, timeoutHelp, tcell.StyleDefault.Foreground(CurrentTheme.Text).Dim(true).Italic(true))

	concText := fmt.Sprintf("Concurrency: %d", tui.scanner.Concurrency)
	tui.drawText(4, 11, concText, textStyle)

	rpsText := "Rate Limit: Unlimited"
	if tui.scanner.Config.RateLimit > 0 {
		rpsText = fmt.Sprintf("Rate Limit: %d req/s", tui.scanner.Config.RateLimit)
	}
	tui.drawText(4, 12, rpsText, tcell.StyleDefault.Foreground(CurrentTheme.Success))

	wordlistText := "Wordlist: wordlist.txt"
	tui.drawText(4, 13, wordlistText, tcell.StyleDefault.Foreground(CurrentTheme.Text).Dim(true))

	// Progress box with animated bar
	tui.drawBox(titleWidth/2+2, 3, titleWidth/2+1, 5, "PROGRESS", progressBoxStyle)
	percentage := float64(0)
	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	// Progress bar - calculate exact available width
	boxInnerWidth := (titleWidth/2 + 1) - 4  // Box width minus borders and padding
	percentText := fmt.Sprintf(" ] %.0f%%", percentage)
	barWidth := boxInnerWidth - 12 - len(percentText)  // Minus "[ SCANNING " prefix and percent text
	if barWidth < 15 {
		barWidth = 15
	}
	filledWidth := int(float64(barWidth) * percentage / 100)

	// Build progress bar with simple ASCII characters for compatibility
	progressBar := "[ SCANNING "

	// Add filled blocks
	for i := 0; i < filledWidth && i < barWidth; i++ {
		progressBar += "#"
	}

	// Add transition blocks (3 of them)
	transitionCount := 0
	if filledWidth > 0 && filledWidth < barWidth {
		maxTransition := 3
		remaining := barWidth - filledWidth
		if remaining < maxTransition {
			maxTransition = remaining
		}
		for i := 0; i < maxTransition; i++ {
			progressBar += "="
			transitionCount++
		}
	}

	// Add empty blocks for remaining space
	emptyCount := barWidth - filledWidth - transitionCount
	for i := 0; i < emptyCount; i++ {
		progressBar += "-"
	}

	progressBar += percentText

	// Draw progress bar with color based on percentage
	barColor := CurrentTheme.Info
	if percentage >= 100 {
		barColor = CurrentTheme.Success
	} else if percentage >= 50 {
		barColor = CurrentTheme.Warning
	}
	tui.drawText(titleWidth/2+4, 4, progressBar, tcell.StyleDefault.Foreground(barColor))

	// Draw target acquired indicator if scan is running
	if completed > 0 && completed < total {
		targetMsg := "⚡ Target acquired →"
		tui.drawText(titleWidth/2+4, 6, targetMsg, tcell.StyleDefault.Foreground(CurrentTheme.Primary))
	}

	// Stats below bar
	speed := tui.scanner.LiveStats.CurrentSpeed
	elapsed := time.Since(tui.scanner.LiveStats.StartTime)
	statsText := fmt.Sprintf("%d/%d | %.0f req/s | %s", completed, total, speed, formatDuration(elapsed))
	tui.drawText(titleWidth/2+4, 5, statsText, textStyle)

	// Increment animation frame
	tui.progressBarFrame++

	// Statistics box
	tui.drawBox(2, 15, titleWidth/2-1, 7, "STATISTICS", statsBoxStyle)
	direct200s := atomic.LoadInt64(&tui.scanner.LiveStats.Direct200s)
	redirects := atomic.LoadInt64(&tui.scanner.LiveStats.Redirects)
	protected := atomic.LoadInt64(&tui.scanner.LiveStats.Protected)
	errors := atomic.LoadInt64(&tui.scanner.LiveStats.Errors)

	stat1 := fmt.Sprintf("%-20s %6d", "Direct 200s:", direct200s)
	tui.drawText(4, 16, stat1, tcell.StyleDefault.Foreground(CurrentTheme.Success))

	stat2 := fmt.Sprintf("%-20s %6d", "Redirects:", redirects)
	tui.drawText(4, 17, stat2, tcell.StyleDefault.Foreground(CurrentTheme.Warning))

	stat3 := fmt.Sprintf("%-20s %6d", "Protected:", protected)
	tui.drawText(4, 18, stat3, tcell.StyleDefault.Foreground(CurrentTheme.Danger))

	stat4 := fmt.Sprintf("%-20s %6d", "Errors:", errors)
	tui.drawText(4, 19, stat4, tcell.StyleDefault.Foreground(CurrentTheme.Danger))

	speedText := fmt.Sprintf("%-20s %6.0f req/s", "Speed:", tui.scanner.LiveStats.CurrentSpeed)
	tui.drawText(4, 20, speedText, tcell.StyleDefault.Foreground(CurrentTheme.Info))

	// Local Network Info box
	var ipBoxStyle tcell.Style
	if CurrentTheme.Name == "SKITTLES" {
		ipBoxStyle = targetBoxStyle
	} else {
		ipBoxStyle = boxStyle
	}
	tui.drawBox(2, 22, titleWidth/2-1, 9, "LOCAL NETWORK INFO", ipBoxStyle)

	netInfo := getLocalNetworkInfo()
	ifaceText := fmt.Sprintf("Interface: %s", truncateString(netInfo.Interface, titleWidth/2-16))
	ipv4Text := fmt.Sprintf("IPv4:      %s", truncateString(netInfo.IPv4, titleWidth/2-16))
	subnetText := fmt.Sprintf("Subnet:    %s", truncateString(netInfo.Subnet, titleWidth/2-16))
	ipv6Text := fmt.Sprintf("IPv6:      %s", truncateString(netInfo.IPv6, titleWidth/2-16))
	macText := fmt.Sprintf("MAC:       %s", truncateString(netInfo.MAC, titleWidth/2-16))
	gatewayText := fmt.Sprintf("Gateway:   %s", truncateString(netInfo.Gateway, titleWidth/2-16))

	tui.drawText(4, 23, ifaceText, tcell.StyleDefault.Foreground(CurrentTheme.Info))
	tui.drawText(4, 24, ipv4Text, tcell.StyleDefault.Foreground(CurrentTheme.Info))
	tui.drawText(4, 25, subnetText, tcell.StyleDefault.Foreground(CurrentTheme.Info))
	tui.drawText(4, 26, ipv6Text, tcell.StyleDefault.Foreground(CurrentTheme.Info).Dim(true))
	tui.drawText(4, 27, macText, tcell.StyleDefault.Foreground(CurrentTheme.Info))
	tui.drawText(4, 28, gatewayText, tcell.StyleDefault.Foreground(CurrentTheme.Info))

	// Live results box - extend down to just above controls
	controlsY := tui.height - 3
	resultsHeight := controlsY - 8  // From Y=8 to controls box
	tui.drawBox(titleWidth/2+2, 8, titleWidth/2+1, resultsHeight, "LIVE RESULTS", resultsBoxStyle)

	tui.scanner.resultsMutex.Lock()
	results := tui.scanner.lastResults
	tui.scanner.resultsMutex.Unlock()

	// Apply scroll offset
	maxVisible := resultsHeight - 2
	totalResults := len(results)

	// Clamp scroll offset
	if tui.resultsScrollOffset < 0 {
		tui.resultsScrollOffset = 0
	}
	if totalResults > maxVisible && tui.resultsScrollOffset > totalResults-maxVisible {
		tui.resultsScrollOffset = totalResults - maxVisible
	}

	// If not scrolling, auto-scroll to bottom (most recent)
	startIdx := tui.resultsScrollOffset
	if tui.resultsScrollOffset == 0 && totalResults > maxVisible {
		startIdx = totalResults - maxVisible
	}

	for i := startIdx; i < totalResults && i-startIdx < maxVisible; i++ {
		result := results[i]
		var label string
		color := CurrentTheme.Text

		if result.IsDirect200 {
			label = "HIT:"
			color = CurrentTheme.Success
		} else if len(result.RedirectChain) > 0 {
			label = "REDIRECT:"
			color = CurrentTheme.Warning
		} else if result.FinalStatus == 403 || result.FinalStatus == 401 {
			label = "PROTECTED:"
			color = CurrentTheme.Danger
		} else if result.FinalStatus >= 500 {
			label = "ERROR:"
			color = CurrentTheme.Danger
		} else {
			label = "FAIL:"
			color = CurrentTheme.Danger
		}

		path := truncateString(result.OriginalPath, titleWidth/2-25)
		line := fmt.Sprintf("%-10s [%d] %s", label, result.FinalStatus, path)
		tui.drawText(titleWidth/2+4, 9+i-startIdx, line, tcell.StyleDefault.Foreground(color))
	}

	// Show scroll indicator if there are more results
	if totalResults > maxVisible {
		scrollIndicator := fmt.Sprintf("(%d/%d)", startIdx+maxVisible, totalResults)
		tui.drawText(titleWidth/2+titleWidth/2-len(scrollIndicator)-3, 8, scrollIndicator, tcell.StyleDefault.Foreground(CurrentTheme.Info).Dim(true))
	}

	// Controls box
	controlsY = tui.height - 3
	tui.drawBox(2, controlsY, titleWidth, 3, "", controlsBoxStyle)
	controls := fmt.Sprintf("Theme: %s | F1: Cycle | F3: Globe | F4: Config | F5: Export | ?: Help | Q: Quit", CurrentTheme.Name)
	controlsX := (titleWidth - len(controls)) / 2
	tui.drawText(2+controlsX, controlsY+1, controls, tcell.StyleDefault.Foreground(CurrentTheme.Info))
}

func (tui *TUI) exportExecutiveSummary() {
	// Generate professional pentest report executive summary
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("PathFinder_Report_%s.txt", timestamp)

	// Build the report
	var report strings.Builder

	report.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	report.WriteString("                      PATHFINDER - EXECUTIVE SUMMARY REPORT                    \n")
	report.WriteString("                           Web Path Discovery Assessment                       \n")
	report.WriteString("═══════════════════════════════════════════════════════════════════════════════\n\n")

	// Report Metadata
	report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
	report.WriteString("│ SCAN METADATA                                                               │\n")
	report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")

	report.WriteString(fmt.Sprintf("Report Generated:    %s\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	report.WriteString(fmt.Sprintf("Tool:                PathFinder v%s\n", Version))
	report.WriteString(fmt.Sprintf("Target URL:          %s\n", tui.scanner.BaseURL))
	report.WriteString(fmt.Sprintf("Scan Started:        %s\n", tui.scanner.LiveStats.StartTime.Format("2006-01-02 15:04:05 MST")))

	elapsed := time.Since(tui.scanner.LiveStats.StartTime)
	report.WriteString(fmt.Sprintf("Scan Duration:       %s\n", formatDuration(elapsed)))
	report.WriteString(fmt.Sprintf("Scan Method:         %s\n", tui.scanner.Config.Method))
	report.WriteString(fmt.Sprintf("Concurrency:         %d workers\n", tui.scanner.Concurrency))
	report.WriteString(fmt.Sprintf("Timeout:             %d seconds\n", int(tui.scanner.Timeout.Seconds())))

	if tui.scanner.Config.RateLimit > 0 {
		report.WriteString(fmt.Sprintf("Rate Limit:          %d req/s\n", tui.scanner.Config.RateLimit))
	} else {
		report.WriteString("Rate Limit:          Unlimited\n")
	}

	report.WriteString("\n")

	// Executive Summary Statistics
	report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
	report.WriteString("│ EXECUTIVE SUMMARY                                                           │\n")
	report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")

	completed := atomic.LoadInt64(&tui.scanner.LiveStats.CompletedRequests)
	direct200s := atomic.LoadInt64(&tui.scanner.LiveStats.Direct200s)
	redirects := atomic.LoadInt64(&tui.scanner.LiveStats.Redirects)
	protected := atomic.LoadInt64(&tui.scanner.LiveStats.Protected)
	errors := atomic.LoadInt64(&tui.scanner.LiveStats.Errors)
	avgSpeed := tui.scanner.LiveStats.CurrentSpeed

	report.WriteString(fmt.Sprintf("Total Requests:      %d\n", completed))
	report.WriteString(fmt.Sprintf("Average Speed:       %.0f req/s\n\n", avgSpeed))

	report.WriteString("FINDINGS BREAKDOWN:\n")
	report.WriteString(fmt.Sprintf("  [✓] Direct 200s:     %d paths (Confirmed accessible resources)\n", direct200s))
	report.WriteString(fmt.Sprintf("  [→] Redirects:       %d paths (Redirection chains detected)\n", redirects))
	report.WriteString(fmt.Sprintf("  [✗] Protected:       %d paths (Authentication/Authorization required)\n", protected))
	report.WriteString(fmt.Sprintf("  [!] Errors:          %d paths (Network/timeout failures)\n\n", errors))

	// Risk Assessment
	report.WriteString("RISK ASSESSMENT:\n")
	if direct200s > 0 {
		report.WriteString(fmt.Sprintf("  • %d accessible paths discovered - Review for sensitive information disclosure\n", direct200s))
	}
	if protected > 0 {
		report.WriteString(fmt.Sprintf("  • %d protected resources found - Verify authentication mechanisms\n", protected))
	}
	if redirects > 0 {
		report.WriteString(fmt.Sprintf("  • %d redirect chains detected - Analyze for open redirects or misconfigurations\n", redirects))
	}
	report.WriteString("\n")

	// Detailed Findings - Direct 200s
	if len(tui.scanner.Stats.Direct200s) > 0 {
		report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
		report.WriteString("│ DETAILED FINDINGS - DIRECT 200s (HIGH PRIORITY)                             │\n")
		report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")
		report.WriteString("The following paths returned HTTP 200 status without redirects, indicating\n")
		report.WriteString("directly accessible resources that should be reviewed for sensitivity.\n\n")

		// Sort by timestamp
		sortedDirect := make([]*ScanResult, len(tui.scanner.Stats.Direct200s))
		copy(sortedDirect, tui.scanner.Stats.Direct200s)
		sort.Slice(sortedDirect, func(i, j int) bool {
			return sortedDirect[i].Timestamp.Before(sortedDirect[j].Timestamp)
		})

		for i, result := range sortedDirect {
			report.WriteString(fmt.Sprintf("[%d] PATH: %s\n", i+1, result.OriginalPath))
			report.WriteString(fmt.Sprintf("    URL:         %s\n", result.OriginalURL))
			report.WriteString(fmt.Sprintf("    Status:      %d (OK)\n", result.FinalStatus))
			report.WriteString(fmt.Sprintf("    Size:        %s\n", formatSize(result.ContentLength)))
			report.WriteString(fmt.Sprintf("    Hash:        %s\n", result.ContentHash[:16]))
			report.WriteString(fmt.Sprintf("    Response:    %dms\n", result.ResponseTime.Milliseconds()))
			report.WriteString(fmt.Sprintf("    Discovered:  %s\n", result.Timestamp.Format("2006-01-02 15:04:05")))
			report.WriteString("\n")
		}
	}

	// Detailed Findings - Redirects
	if len(tui.scanner.Stats.Redirects) > 0 {
		report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
		report.WriteString("│ DETAILED FINDINGS - REDIRECTS (MEDIUM PRIORITY)                             │\n")
		report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")
		report.WriteString("The following paths triggered redirect chains. Review for open redirect\n")
		report.WriteString("vulnerabilities or unexpected redirect behavior.\n\n")

		sortedRedirects := make([]*ScanResult, len(tui.scanner.Stats.Redirects))
		copy(sortedRedirects, tui.scanner.Stats.Redirects)
		sort.Slice(sortedRedirects, func(i, j int) bool {
			return sortedRedirects[i].Timestamp.Before(sortedRedirects[j].Timestamp)
		})

		for i, result := range sortedRedirects {
			report.WriteString(fmt.Sprintf("[%d] PATH: %s\n", i+1, result.OriginalPath))
			report.WriteString(fmt.Sprintf("    Original:    %s\n", result.OriginalURL))
			report.WriteString(fmt.Sprintf("    Final URL:   %s\n", result.FinalURL))
			report.WriteString(fmt.Sprintf("    Status:      %d\n", result.FinalStatus))
			report.WriteString(fmt.Sprintf("    Hops:        %d redirect(s)\n", len(result.RedirectChain)))

			if len(result.RedirectChain) > 0 {
				report.WriteString("    Chain:\n")
				for j, step := range result.RedirectChain {
					report.WriteString(fmt.Sprintf("      %d. [%d] %s\n", j+1, step.Status, step.URL))
				}
			}

			report.WriteString(fmt.Sprintf("    Discovered:  %s\n", result.Timestamp.Format("2006-01-02 15:04:05")))
			report.WriteString("\n")
		}
	}

	// Protected Resources
	if len(tui.scanner.Stats.OtherCodes) > 0 {
		report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
		report.WriteString("│ DETAILED FINDINGS - PROTECTED/ERROR RESPONSES                               │\n")
		report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")

		sortedOther := make([]*ScanResult, len(tui.scanner.Stats.OtherCodes))
		copy(sortedOther, tui.scanner.Stats.OtherCodes)
		sort.Slice(sortedOther, func(i, j int) bool {
			return sortedOther[i].Timestamp.Before(sortedOther[j].Timestamp)
		})

		for i, result := range sortedOther {
			statusLabel := "UNKNOWN"
			if result.FinalStatus == 401 {
				statusLabel = "UNAUTHORIZED"
			} else if result.FinalStatus == 403 {
				statusLabel = "FORBIDDEN"
			} else if result.FinalStatus >= 500 {
				statusLabel = "SERVER ERROR"
			}

			report.WriteString(fmt.Sprintf("[%d] PATH: %s\n", i+1, result.OriginalPath))
			report.WriteString(fmt.Sprintf("    URL:         %s\n", result.OriginalURL))
			report.WriteString(fmt.Sprintf("    Status:      %d (%s)\n", result.FinalStatus, statusLabel))
			report.WriteString(fmt.Sprintf("    Discovered:  %s\n", result.Timestamp.Format("2006-01-02 15:04:05")))
			report.WriteString("\n")
		}
	}

	// Recommendations
	report.WriteString("┌─────────────────────────────────────────────────────────────────────────────┐\n")
	report.WriteString("│ RECOMMENDATIONS                                                              │\n")
	report.WriteString("└─────────────────────────────────────────────────────────────────────────────┘\n\n")

	report.WriteString("1. Review all discovered Direct 200 paths for:\n")
	report.WriteString("   - Sensitive information disclosure (config files, backups, credentials)\n")
	report.WriteString("   - Admin panels or debug interfaces\n")
	report.WriteString("   - Unintended public access to resources\n\n")

	report.WriteString("2. Analyze redirect chains for:\n")
	report.WriteString("   - Open redirect vulnerabilities\n")
	report.WriteString("   - Misconfigured redirects leaking internal paths\n")
	report.WriteString("   - Authentication bypass opportunities\n\n")

	report.WriteString("3. Verify protected resources:\n")
	report.WriteString("   - Confirm authentication mechanisms are properly enforced\n")
	report.WriteString("   - Test for authorization bypass vulnerabilities\n")
	report.WriteString("   - Verify 401/403 responses don't leak sensitive information\n\n")

	report.WriteString("4. General security considerations:\n")
	report.WriteString("   - Implement proper access controls on all discovered paths\n")
	report.WriteString("   - Remove or restrict access to unnecessary endpoints\n")
	report.WriteString("   - Ensure error messages don't reveal system information\n")
	report.WriteString("   - Monitor and log access to sensitive resources\n\n")

	// Footer
	report.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	report.WriteString("                              END OF REPORT                                    \n")
	report.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")

	// Write to file
	err := os.WriteFile(filename, []byte(report.String()), 0644)
	if err != nil {
		// Silently fail - could show error in UI later
		return
	}

	// Success - file was created
	// Could show success message in UI
}

func (tui *TUI) renderConfigMenu() {
	// Draw semi-transparent overlay effect by drawing a box
	menuWidth := 60
	menuHeight := 14
	menuX := (tui.width - menuWidth) / 2
	menuY := (tui.height - menuHeight) / 2

	// Fill background with solid color to make menu readable
	bgStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text)
	for y := menuY; y < menuY+menuHeight; y++ {
		for x := menuX; x < menuX+menuWidth; x++ {
			tui.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}

	// Draw config menu box
	boxStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Primary).Bold(true)
	tui.drawBox(menuX, menuY, menuWidth, menuHeight, "CONFIG MENU", boxStyle)

	// Config options
	textStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text)
	selectedStyle := tcell.StyleDefault.Background(CurrentTheme.Success).Foreground(CurrentTheme.Background).Bold(true)

	options := []string{
		fmt.Sprintf("Concurrency:  %d", tui.scanner.Concurrency),
		fmt.Sprintf("Rate Limit:   %d req/s  (0 = unlimited)", tui.scanner.Config.RateLimit),
		fmt.Sprintf("Timeout:      %d seconds", int(tui.scanner.Timeout.Seconds())),
		fmt.Sprintf("Method:       %s", tui.scanner.Config.Method),
	}

	startY := menuY + 2
	for i, option := range options {
		style := textStyle
		if i == tui.configMenuSelected {
			style = selectedStyle
			option = option + " ◀ ▶"
		}
		tui.drawText(menuX+2, startY+i*2, option, style)
	}

	// Instructions
	instrY := menuY + menuHeight - 3
	instr := "↑/↓: Navigate | ◀/▶: Change Value | F4/Esc: Close"
	tui.drawText(menuX+2, instrY, instr, tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Info))
}

func (tui *TUI) renderHelpScreen() {
	// Full screen help overlay
	helpWidth := tui.width - 8
	helpHeight := tui.height - 4
	helpX := 4
	helpY := 2

	if helpWidth > 120 {
		helpWidth = 120
		helpX = (tui.width - 120) / 2
	}

	// Fill background
	bgStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text)
	for y := helpY; y < helpY+helpHeight; y++ {
		for x := helpX; x < helpX+helpWidth; x++ {
			tui.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}

	// Draw help box
	boxStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Primary).Bold(true)
	tui.drawBox(helpX, helpY, helpWidth, helpHeight, "PATHFINDER HELP GUIDE", boxStyle)

	// Help content
	textStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text)
	titleStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Primary).Bold(true)
	labelStyle := tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Success).Bold(true)

	line := helpY + 2
	col := helpX + 2

	// Overview
	tui.drawText(col, line, "OVERVIEW", titleStyle)
	line += 1
	tui.drawText(col, line, "PathFinder is a web path discovery tool that scans target URLs for existing", textStyle)
	line += 1
	tui.drawText(col, line, "paths and directories using a wordlist. It tracks response codes, redirects,", textStyle)
	line += 1
	tui.drawText(col, line, "and provides real-time feedback during scanning.", textStyle)
	line += 2

	// Configuration Settings
	tui.drawText(col, line, "SCAN CONFIGURATION", titleStyle)
	line += 1

	tui.drawText(col, line, "Method:", labelStyle)
	tui.drawText(col+18, line, "HTTP method to use (GET, POST, HEAD, PUT, DELETE, PATCH)", textStyle)
	line += 1
	tui.drawText(col+18, line, "Default: GET | Use GET for most scans", tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text).Dim(true))
	line += 2

	tui.drawText(col, line, "Timeout:", labelStyle)
	tui.drawText(col+18, line, "Max wait time per HTTP request before giving up", textStyle)
	line += 1
	tui.drawText(col+18, line, "Default: 10s | Increase for slow servers, decrease for fast networks", tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text).Dim(true))
	line += 2

	tui.drawText(col, line, "Concurrency:", labelStyle)
	tui.drawText(col+18, line, "Number of simultaneous requests to send", textStyle)
	line += 1
	tui.drawText(col+18, line, "Default: 50 | Higher = faster but more aggressive", tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text).Dim(true))
	line += 2

	tui.drawText(col, line, "Rate Limit:", labelStyle)
	tui.drawText(col+18, line, "Maximum requests per second (0 = unlimited)", textStyle)
	line += 1
	tui.drawText(col+18, line, "Use to avoid overwhelming target or triggering rate limits", tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text).Dim(true))
	line += 2

	tui.drawText(col, line, "Wordlist:", labelStyle)
	tui.drawText(col+18, line, "File containing paths to scan (one per line)", textStyle)
	line += 1
	tui.drawText(col+18, line, "Default: wordlist.txt in current directory", tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text).Dim(true))
	line += 2

	// Statistics Explained
	tui.drawText(col, line, "STATISTICS", titleStyle)
	line += 1

	tui.drawText(col, line, "Direct 200s:", labelStyle)
	tui.drawText(col+18, line, "Paths that returned HTTP 200 without any redirects", textStyle)
	line += 1
	tui.drawText(col, line, "Redirects:", labelStyle)
	tui.drawText(col+18, line, "Paths that redirected to another URL (3xx status codes)", textStyle)
	line += 1
	tui.drawText(col, line, "Protected:", labelStyle)
	tui.drawText(col+18, line, "Paths requiring authentication (401/403 status codes)", textStyle)
	line += 1
	tui.drawText(col, line, "Errors:", labelStyle)
	tui.drawText(col+18, line, "Failed requests due to network errors or timeouts", textStyle)
	line += 1
	tui.drawText(col, line, "Speed:", labelStyle)
	tui.drawText(col+18, line, "Current scanning speed in requests per second", textStyle)
	line += 2

	// Keyboard Shortcuts
	tui.drawText(col, line, "KEYBOARD SHORTCUTS", titleStyle)
	line += 1
	tui.drawText(col, line, "F1:", labelStyle)
	tui.drawText(col+12, line, "Cycle through color themes", textStyle)
	line += 1
	tui.drawText(col, line, "F3:", labelStyle)
	tui.drawText(col+12, line, "Toggle Globe Mode (animated Earth visualization)", textStyle)
	line += 1
	tui.drawText(col, line, "F4:", labelStyle)
	tui.drawText(col+12, line, "Open configuration menu to adjust scan settings", textStyle)
	line += 1
	tui.drawText(col, line, "F5:", labelStyle)
	tui.drawText(col+12, line, "Export executive summary report for pentest documentation", textStyle)
	line += 1
	tui.drawText(col, line, "Enter:", labelStyle)
	tui.drawText(col+12, line, "Activate input field to enter a new target URL", textStyle)
	line += 1
	tui.drawText(col, line, "↑/↓:", labelStyle)
	tui.drawText(col+12, line, "Scroll through live results (Up/Down arrow keys)", textStyle)
	line += 1
	tui.drawText(col, line, "?:", labelStyle)
	tui.drawText(col+12, line, "Toggle this help screen", textStyle)
	line += 1
	tui.drawText(col, line, "Q:", labelStyle)
	tui.drawText(col+12, line, "Quit PathFinder", textStyle)
	line += 1
	tui.drawText(col, line, "1-8:", labelStyle)
	tui.drawText(col+12, line, "Switch to specific theme (Matrix, Rainbow, Cyber, Blood, etc.)", textStyle)
	line += 2

	// Close instruction
	closeText := "Press ? or Esc to close this help screen"
	closeX := helpX + (helpWidth-len(closeText))/2
	tui.drawText(closeX, helpY+helpHeight-2, closeText, tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Warning).Bold(true))
}

func (tui *TUI) HandleInput() {
	for tui.running {
		ev := tui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Handle config menu navigation first
			if tui.showConfigMenu {
				switch ev.Key() {
				case tcell.KeyUp:
					tui.configMenuSelected--
					if tui.configMenuSelected < 0 {
						tui.configMenuSelected = 3 // 4 options, so max index is 3
					}
					tui.Render()
					continue
				case tcell.KeyDown:
					tui.configMenuSelected++
					if tui.configMenuSelected > 3 {
						tui.configMenuSelected = 0
					}
					tui.Render()
					continue
				case tcell.KeyLeft:
					// Decrement selected option
					switch tui.configMenuSelected {
					case 0: // Concurrency
						if tui.scanner.Concurrency > 1 {
							tui.scanner.Concurrency -= 5
							if tui.scanner.Concurrency < 1 {
								tui.scanner.Concurrency = 1
							}
						}
					case 1: // Rate Limit
						if tui.scanner.Config.RateLimit > 0 {
							tui.scanner.Config.RateLimit -= 10
							if tui.scanner.Config.RateLimit < 0 {
								tui.scanner.Config.RateLimit = 0
							}
						}
					case 2: // Timeout
						timeoutSec := int(tui.scanner.Timeout.Seconds())
						if timeoutSec > 1 {
							timeoutSec -= 1
							if timeoutSec < 1 {
								timeoutSec = 1
							}
							tui.scanner.Timeout = time.Duration(timeoutSec) * time.Second
						}
					case 3: // Method
						methods := []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
						currentIdx := 0
						for i, m := range methods {
							if m == tui.scanner.Config.Method {
								currentIdx = i
								break
							}
						}
						currentIdx--
						if currentIdx < 0 {
							currentIdx = len(methods) - 1
						}
						tui.scanner.Config.Method = methods[currentIdx]
					}
					tui.Render()
					continue
				case tcell.KeyRight:
					// Increment selected option
					switch tui.configMenuSelected {
					case 0: // Concurrency
						tui.scanner.Concurrency += 5
						if tui.scanner.Concurrency > 500 {
							tui.scanner.Concurrency = 500
						}
					case 1: // Rate Limit
						tui.scanner.Config.RateLimit += 10
						if tui.scanner.Config.RateLimit > 1000 {
							tui.scanner.Config.RateLimit = 1000
						}
					case 2: // Timeout
						timeoutSec := int(tui.scanner.Timeout.Seconds())
						timeoutSec += 1
						if timeoutSec > 300 {
							timeoutSec = 300
						}
						tui.scanner.Timeout = time.Duration(timeoutSec) * time.Second
					case 3: // Method
						methods := []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
						currentIdx := 0
						for i, m := range methods {
							if m == tui.scanner.Config.Method {
								currentIdx = i
								break
							}
						}
						currentIdx++
						if currentIdx >= len(methods) {
							currentIdx = 0
						}
						tui.scanner.Config.Method = methods[currentIdx]
					}
					tui.Render()
					continue
				case tcell.KeyEscape:
					tui.showConfigMenu = false
					tui.Render()
					continue
				case tcell.KeyF4:
					tui.showConfigMenu = false
					tui.Render()
					continue
				}
			}

			switch ev.Key() {
			case tcell.KeyUp:
				// Scroll results up
				if !tui.showConfigMenu && !tui.inputActive {
					tui.resultsScrollOffset -= 1
					if tui.resultsScrollOffset < 0 {
						tui.resultsScrollOffset = 0
					}
				}
			case tcell.KeyDown:
				// Scroll results down
				if !tui.showConfigMenu && !tui.inputActive {
					tui.resultsScrollOffset += 1
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				if tui.showHelpScreen {
					tui.showHelpScreen = false
					tui.Render()
					continue
				}
				if tui.showConfigMenu {
					tui.showConfigMenu = false
					tui.Render()
					continue
				}
				tui.running = false
				return
			case tcell.KeyF1:
				// F1 - Cycle theme
				tui.cycleTheme('0')
			case tcell.KeyF2:
				// F2 - Activate Skittles theme with random colors
				CurrentTheme = ThemeSkittles
				tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
				// Reset lastTheme to force color regeneration
				tui.lastTheme = ""
			case tcell.KeyF3:
				// F3 - Toggle Globe Mode
				tui.showGlobe = !tui.showGlobe
			case tcell.KeyF4:
				// F4 - Toggle Config Menu
				tui.showConfigMenu = !tui.showConfigMenu
			case tcell.KeyF5:
				// F5 - Export Executive Summary Report
				tui.exportExecutiveSummary()
			case tcell.KeyEnter:
				// Enter - Toggle input field active/inactive, or submit if active with text
				if !tui.showGlobe {
					if tui.inputActive && len(tui.inputText) > 0 {
						// Submit the input and start new scan
						tui.submitInput()
					} else {
						// Toggle input field
						tui.inputActive = !tui.inputActive
					}
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				// Backspace - Delete last character from input
				if tui.inputActive && len(tui.inputText) > 0 {
					tui.inputText = tui.inputText[:len(tui.inputText)-1]
				}
			case tcell.KeyRune:
				// If input is active, add typed characters to input text
				if tui.inputActive {
					tui.inputText += string(ev.Rune())
				} else {
					// Only handle theme shortcuts when input is not active
					switch ev.Rune() {
					case 'q', 'Q':
						tui.running = false
						return
					case '?':
						// Toggle help screen
						tui.showHelpScreen = !tui.showHelpScreen
					case ' ':
						// Spacebar - Skip splash screen if showing
						if tui.showSplash {
							tui.showSplash = false
						}
					case '1':
						CurrentTheme = ThemeMatrix
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '2':
						CurrentTheme = ThemeRainbow
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '3':
						CurrentTheme = ThemeCyber
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '4':
						CurrentTheme = ThemeBlood
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '5':
						CurrentTheme = ThemeSkittles
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '6':
						CurrentTheme = ThemeDark
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '7':
						CurrentTheme = ThemePurple
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					case '8':
						CurrentTheme = ThemeAmber
						tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
					}
				}
			}
		case *tcell.EventResize:
			tui.width, tui.height = tui.screen.Size()
			tui.screen.Sync()
		}
		tui.Render()
	}
}

func (tui *TUI) cycleTheme(key rune) {
	switch key {
	case '1':
		CurrentTheme = ThemeMatrix
	case '2':
		CurrentTheme = ThemeRainbow
	case '3':
		CurrentTheme = ThemeCyber
	case '4':
		CurrentTheme = ThemeBlood
	case '5':
		CurrentTheme = ThemeSkittles
	case '6':
		CurrentTheme = ThemeDark
	case '7':
		CurrentTheme = ThemePurple
	case '8':
		CurrentTheme = ThemeAmber
	default:
		// Cycle through themes
		switch CurrentTheme.Name {
		case "MATRIX":
			CurrentTheme = ThemeRainbow
		case "RAINBOW":
			CurrentTheme = ThemeCyber
		case "CYBER":
			CurrentTheme = ThemeBlood
		case "BLOOD":
			CurrentTheme = ThemeSkittles
		case "SKITTLES":
			CurrentTheme = ThemeDark
		case "DARK":
			CurrentTheme = ThemePurple
		case "PURPLE":
			CurrentTheme = ThemeAmber
		case "AMBER":
			CurrentTheme = ThemeMatrix
		}
	}
	tui.screen.SetStyle(tcell.StyleDefault.Background(CurrentTheme.Background).Foreground(CurrentTheme.Text))
}

func (tui *TUI) Run() {
	go tui.HandleInput()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for tui.running {
		<-ticker.C
		tui.Render()
	}
}

func (tui *TUI) Stop() {
	tui.running = false
	tui.screen.Fini()
}

func (tui *TUI) submitInput() {
	input := strings.TrimSpace(tui.inputText)
	if input == "" {
		return
	}

	// Normalize URL - add https:// if no scheme present
	var targetURL string
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		targetURL = input
	} else {
		// Assume https for plain domains
		targetURL = "https://" + input
	}

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		// Invalid URL - silently return (could show error message in future)
		return
	}

	// Update scanner's BaseURL
	tui.scanner.BaseURL = strings.TrimRight(targetURL, "/")

	// Reset statistics
	tui.scanner.Stats = &Statistics{
		RedirectTargets: make(map[string]int),
		ContentHashes:   make(map[string][]*ScanResult),
	}
	tui.scanner.LiveStats = &LiveStats{
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}
	tui.scanner.lastResults = make([]*ScanResult, 0, 50)

	// Load wordlist - try default
	paths, err := LoadWordlist("wordlist.txt")
	if err != nil {
		// Can't start scan without wordlist
		return
	}

	// Clear input and deactivate
	tui.inputText = ""
	tui.inputActive = false

	// Start new scan in background
	go func() {
		tui.scanner.ScanAll(paths, tui)
	}()
}

// ===========================================================================
// SCANNER
// ===========================================================================

func NewScanner(baseURL string, concurrency int, timeout int, verbose bool, config *Config) *Scanner {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        concurrency * 2,
		MaxIdleConnsPerHost: concurrency,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var rateLimiter <-chan time.Time
	if config.RateLimit > 0 {
		interval := time.Second / time.Duration(config.RateLimit)
		rateLimiter = time.Tick(interval)
	}

	recursionQueue := make(chan string, 10000)

	return &Scanner{
		BaseURL:        strings.TrimRight(baseURL, "/"),
		Concurrency:    concurrency,
		Timeout:        time.Duration(timeout) * time.Second,
		Verbose:        verbose,
		Client:         client,
		Config:         config,
		visitedPaths:   make(map[string]bool),
		recursionQueue: recursionQueue,
		rateLimiter:    rateLimiter,
		lastResults:    make([]*ScanResult, 0, 50),
		LiveStats: &LiveStats{
			StartTime:  time.Now(),
			LastUpdate: time.Now(),
		},
		Stats: &Statistics{
			RedirectTargets: make(map[string]int),
			ContentHashes:   make(map[string][]*ScanResult),
		},
	}
}

func (s *Scanner) FetchWithRedirectTracking(targetURL string) (*ScanResult, error) {
	if s.rateLimiter != nil {
		<-s.rateLimiter
	}

	if s.Config.Delay > 0 {
		time.Sleep(s.Config.Delay)
	}

	var redirectChain []RedirectStep
	currentURL := targetURL
	startTime := time.Now()

	method := s.Config.Method
	if method == "" {
		method = "GET"
	}

	for i := 0; i < MaxRedirects; i++ {
		req, err := http.NewRequest(method, currentURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", UserAgent)

		for key, value := range s.Config.CustomHeaders {
			req.Header.Set(key, value)
		}

		if s.Config.Cookie != "" {
			req.Header.Set("Cookie", s.Config.Cookie)
		}

		resp, err := s.Client.Do(req)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		status := resp.StatusCode

		if status >= 300 && status < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				break
			}

			redirectChain = append(redirectChain, RedirectStep{
				URL:    currentURL,
				Status: status,
			})

			baseURL, _ := url.Parse(currentURL)
			nextURL, err := baseURL.Parse(location)
			if err != nil {
				break
			}
			currentURL = nextURL.String()
			continue
		}

		contentHash := fmt.Sprintf("%x", md5.Sum(body))
		isDirect := len(redirectChain) == 0 && status == 200
		responseTime := time.Since(startTime)

		result := &ScanResult{
			OriginalPath:  strings.TrimPrefix(targetURL, s.BaseURL),
			OriginalURL:   targetURL,
			FinalStatus:   status,
			FinalURL:      currentURL,
			RedirectChain: redirectChain,
			ContentLength: len(body),
			ContentHash:   contentHash,
			IsDirect200:   isDirect,
			ResponseTime:  responseTime,
			Timestamp:     time.Now(),
		}

		return result, nil
	}

	return nil, fmt.Errorf("max redirects exceeded")
}

func (s *Scanner) DetectWildcard() *WildcardBaseline {
	randomPaths := []string{
		randomString(32),
		"this-path-never-exists-" + randomString(16),
		"__test__" + fmt.Sprintf("%d", rand.Intn(900000)+100000),
	}

	var baselines []*ScanResult
	for _, randPath := range randomPaths {
		testURL := fmt.Sprintf("%s/%s", s.BaseURL, randPath)
		result, err := s.FetchWithRedirectTracking(testURL)
		if err == nil && result.FinalStatus == 200 {
			baselines = append(baselines, result)
		}
	}

	if len(baselines) > 0 {
		firstHash := baselines[0].ContentHash
		allSame := true
		for _, b := range baselines {
			if b.ContentHash != firstHash {
				allSame = false
				break
			}
		}

		if allSame {
			return &WildcardBaseline{
				Hash:   firstHash,
				Length: baselines[0].ContentLength,
				Status: 200,
			}
		}
	}

	return nil
}

func (s *Scanner) IsWildcardResponse(result *ScanResult) bool {
	if s.WildcardBaseline == nil {
		return false
	}
	return result.FinalStatus == 200 && result.ContentHash == s.WildcardBaseline.Hash
}

func (s *Scanner) ShouldFilterResult(result *ScanResult) bool {
	if len(s.Config.StatusCodes) > 0 {
		found := false
		for _, code := range s.Config.StatusCodes {
			if result.FinalStatus == code {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	if len(s.Config.FilterStatuses) > 0 {
		for _, code := range s.Config.FilterStatuses {
			if result.FinalStatus == code {
				return true
			}
		}
	}

	if len(s.Config.FilterSizes) > 0 {
		for _, size := range s.Config.FilterSizes {
			if result.ContentLength == size {
				return true
			}
		}
	}

	return false
}

func (s *Scanner) ScanPath(path string) (*ScanResult, error) {
	targetURL := s.BaseURL + "/" + strings.TrimPrefix(path, "/")
	result, err := s.FetchWithRedirectTracking(targetURL)

	if err != nil {
		atomic.AddInt64(&s.LiveStats.Errors, 1)
		return nil, err
	}

	if s.IsWildcardResponse(result) {
		return nil, nil
	}

	if s.ShouldFilterResult(result) {
		return nil, nil
	}

	// Update live stats
	s.Stats.mu.Lock()
	s.Stats.TotalScanned++

	if result.IsDirect200 {
		s.Stats.Direct200s = append(s.Stats.Direct200s, result)
		atomic.AddInt64(&s.LiveStats.Direct200s, 1)
	} else if len(result.RedirectChain) > 0 {
		s.Stats.Redirects = append(s.Stats.Redirects, result)
		s.Stats.RedirectTargets[result.FinalURL]++
		atomic.AddInt64(&s.LiveStats.Redirects, 1)
	}

	if result.FinalStatus == 401 || result.FinalStatus == 403 {
		atomic.AddInt64(&s.LiveStats.Protected, 1)
	}

	if result.FinalStatus == 401 || result.FinalStatus == 403 || result.FinalStatus == 500 {
		s.Stats.OtherCodes = append(s.Stats.OtherCodes, result)
	}

	s.Stats.ContentHashes[result.ContentHash] = append(s.Stats.ContentHashes[result.ContentHash], result)
	s.Stats.mu.Unlock()

	// Add to live display buffer
	s.AddLiveResult(result)

	return result, nil
}

func (s *Scanner) AddLiveResult(result *ScanResult) {
	s.resultsMutex.Lock()
	defer s.resultsMutex.Unlock()

	s.lastResults = append(s.lastResults, result)
	if len(s.lastResults) > 100 {
		s.lastResults = s.lastResults[len(s.lastResults)-100:]
	}
}

func (s *Scanner) ScanAll(paths []string, tui *TUI) []*ScanResult {
	if len(s.Config.Extensions) > 0 {
		paths = GeneratePathsWithExtensions(paths, s.Config.Extensions)
	}

	totalPaths := len(paths)
	s.LiveStats.TotalRequests = int64(totalPaths)

	// Wildcard detection
	s.WildcardBaseline = s.DetectWildcard()

	// Start speed calculator
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			completed := atomic.LoadInt64(&s.LiveStats.CompletedRequests)
			elapsed := time.Since(s.LiveStats.StartTime).Seconds()
			if elapsed > 0 {
				s.LiveStats.mu.Lock()
				s.LiveStats.CurrentSpeed = float64(completed) / elapsed
				s.LiveStats.mu.Unlock()
			}

			if completed >= int64(totalPaths) {
				return
			}
		}
	}()

	// Scanning
	sem := make(chan struct{}, s.Concurrency)
	var wg sync.WaitGroup

	var results []*ScanResult
	var resultsMutex sync.Mutex

	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result, _ := s.ScanPath(path)
			if result != nil {
				resultsMutex.Lock()
				results = append(results, result)
				resultsMutex.Unlock()
			}

			atomic.AddInt64(&s.LiveStats.CompletedRequests, 1)
		}(p)
	}

	wg.Wait()

	return results
}

func (s *Scanner) AnalyzeResults() string {
	output := "\n" + strings.Repeat("=", 80) + "\n"
	output += "SCAN SUMMARY\n"
	output += strings.Repeat("=", 80) + "\n\n"

	output += fmt.Sprintf("Total paths scanned: %d\n", s.Stats.TotalScanned)
	output += fmt.Sprintf("Direct 200s found: %d\n", len(s.Stats.Direct200s))
	output += fmt.Sprintf("Redirects found: %d\n\n", len(s.Stats.Redirects))

	if len(s.Stats.Direct200s) > 0 {
		output += strings.Repeat("=", 80) + "\n"
		output += "[✓] DIRECT 200s (Actual hosted content - NO redirects)\n"
		output += strings.Repeat("=", 80) + "\n"

		sort.Slice(s.Stats.Direct200s, func(i, j int) bool {
			return s.Stats.Direct200s[i].OriginalPath < s.Stats.Direct200s[j].OriginalPath
		})

		for _, result := range s.Stats.Direct200s {
			output += fmt.Sprintf("  %s\n", result.OriginalURL)
			output += fmt.Sprintf("     Length: %s | Hash: %s... | Time: %dms\n",
				formatSize(result.ContentLength), result.ContentHash[:12], result.ResponseTime.Milliseconds())
		}
	}

	return output
}

// ===========================================================================
// HELPER FUNCTIONS
// ===========================================================================

type LocalNetInfo struct {
	IPv4      string
	IPv6      string
	MAC       string
	Interface string
	Gateway   string
	Subnet    string
}

func getLocalNetworkInfo() LocalNetInfo {
	info := LocalNetInfo{
		IPv4:      "N/A",
		IPv6:      "N/A",
		MAC:       "N/A",
		Interface: "N/A",
		Gateway:   "N/A",
		Subnet:    "N/A",
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return info
	}

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var ipnet *net.IPNet

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ipnet = v
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Check if IPv4
			if ip.To4() != nil {
				if info.IPv4 == "N/A" {
					info.IPv4 = ip.String()
					info.Interface = iface.Name
					if len(iface.HardwareAddr) > 0 {
						info.MAC = iface.HardwareAddr.String()
					}
					if ipnet != nil {
						info.Subnet = ipnet.String()
					}
				}
			} else {
				// IPv6
				if info.IPv6 == "N/A" {
					info.IPv6 = ip.String()
				}
			}
		}
	}

	return info
}

func GeneratePathsWithExtensions(basePaths []string, extensions []string) []string {
	if len(extensions) == 0 {
		return basePaths
	}

	var allPaths []string
	for _, p := range basePaths {
		allPaths = append(allPaths, p)
		for _, ext := range extensions {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			allPaths = append(allPaths, p+ext)
		}
	}
	return allPaths
}

func LoadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paths []string
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			paths = append(paths, line)
		}
	}

	return paths, nil
}

func ExportToJSON(filename string, results []*ScanResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func ExportToCSV(filename string, results []*ScanResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Path", "URL", "Status", "Final URL", "Redirects", "Length", "Hash", "Direct200", "Time(ms)"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		redirectCount := fmt.Sprintf("%d", len(result.RedirectChain))
		direct200 := "No"
		if result.IsDirect200 {
			direct200 = "Yes"
		}

		row := []string{
			result.OriginalPath,
			result.OriginalURL,
			strconv.Itoa(result.FinalStatus),
			result.FinalURL,
			redirectCount,
			strconv.Itoa(result.ContentLength),
			result.ContentHash[:12],
			direct200,
			strconv.FormatInt(result.ResponseTime.Milliseconds(), 10),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func formatSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func parseIntList(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []int
	for _, p := range parts {
		if val, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			result = append(result, val)
		}
	}
	return result
}

func parseStringList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ===========================================================================
// MAIN
// ===========================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	target := flag.String("target", "", "Target base URL")
	wordlist := flag.String("wordlist", "wordlist.txt", "Wordlist file")
	_ = wordlist // Keep flag for future use, loaded dynamically when scan starts
	concurrency := flag.Int("concurrency", DefaultConcurrency, "Concurrent requests")
	timeout := flag.Int("timeout", DefaultTimeout, "Timeout in seconds")
	verbose := flag.Bool("verbose", false, "Verbose output")
	statusCodes := flag.String("mc", "", "Match status codes")
	filterStatuses := flag.String("fc", "", "Filter status codes")
	filterSizes := flag.String("fs", "", "Filter content sizes")
	extensions := flag.String("x", "", "File extensions")
	headers := flag.String("H", "", "Custom header")
	cookie := flag.String("cookie", "", "Cookie data")
	method := flag.String("X", "GET", "HTTP method")
	rateLimit := flag.Int("rate", 0, "Max requests/sec")
	delay := flag.Int("delay", 0, "Delay between requests (ms)")
	recursive := flag.Bool("r", false, "Recursive scanning")
	recursionDepth := flag.Int("depth", 3, "Recursion depth")
	outputFile := flag.String("o", "", "Output file")
	outputFormat := flag.String("of", "text", "Output format")
	theme := flag.String("theme", "matrix", "Color theme: matrix, rainbow, cyber, blood")

	flag.Parse()

	// Set theme
	switch strings.ToLower(*theme) {
	case "rainbow":
		CurrentTheme = ThemeRainbow
	case "cyber":
		CurrentTheme = ThemeCyber
	case "blood":
		CurrentTheme = ThemeBlood
	case "skittles":
		CurrentTheme = ThemeSkittles
	case "dark", "white":
		CurrentTheme = ThemeDark
	case "purple":
		CurrentTheme = ThemePurple
	case "amber":
		CurrentTheme = ThemeAmber
	default:
		CurrentTheme = ThemeMatrix
	}

	// If no target provided, use placeholder
	if *target == "" {
		*target = "https://enter-url-to-scan"
	}

	parsedURL, err := url.Parse(*target)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		fmt.Printf("Error: Invalid URL: %s\n", *target)
		os.Exit(1)
	}

	customHeaders := make(map[string]string)
	if *headers != "" {
		parts := strings.SplitN(*headers, ":", 2)
		if len(parts) == 2 {
			customHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	config := &Config{
		StatusCodes:    parseIntList(*statusCodes),
		FilterStatuses: parseIntList(*filterStatuses),
		FilterSizes:    parseIntList(*filterSizes),
		Extensions:     parseStringList(*extensions),
		CustomHeaders:  customHeaders,
		Cookie:         *cookie,
		Method:         *method,
		RateLimit:      *rateLimit,
		Delay:          time.Duration(*delay) * time.Millisecond,
		Recursive:      *recursive,
		RecursionDepth: *recursionDepth,
		OutputFile:     *outputFile,
		OutputFormat:   *outputFormat,
		Theme:          *theme,
	}

	// Don't load wordlist here - it will be loaded when user starts a scan
	scanner := NewScanner(*target, *concurrency, *timeout, *verbose, config)

	// Create TUI
	tui, err := NewTUI(scanner)
	if err != nil {
		fmt.Printf("Error creating TUI: %v\n", err)
		os.Exit(1)
	}
	defer tui.Stop()

	// Don't auto-start scanning - wait for user to enter URL
	// Scanning will be triggered via the input field

	// Run TUI
	tui.Run()

	// After TUI exits, print summary
	fmt.Println(scanner.AnalyzeResults())

	if *outputFile != "" {
		allResults := append([]*ScanResult{}, scanner.Stats.Direct200s...)
		allResults = append(allResults, scanner.Stats.Redirects...)
		allResults = append(allResults, scanner.Stats.OtherCodes...)

		switch *outputFormat {
		case "json":
			if err := ExportToJSON(*outputFile, allResults); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("[OK] Exported %d results\n", len(allResults))
			}
		case "csv":
			if err := ExportToCSV(*outputFile, allResults); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("[OK] Exported %d results\n", len(allResults))
			}
		}
	}
}
