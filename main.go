package main

import (
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Call of the Night Palette
var (
	cotnMidnight     = rl.NewColor(8, 2, 13, 255)
	cotnDeepPurple   = rl.NewColor(26, 11, 46, 255)
	cotnNeonMagenta  = rl.NewColor(255, 46, 151, 255)
	cotnElectricBlue = rl.NewColor(0, 212, 255, 255)
	cotnBrightText   = rl.NewColor(224, 209, 240, 255)
	cotnDimText      = rl.NewColor(139, 124, 163, 255)
)

type Star struct {
	x, y  float32
	speed float32
	size  float32
}

func getManPage(page string) string {
	// Execute man command with col -b to strip backspace formatting
	cmd := exec.Command("man", page)
	// Force plain text output by setting TERM to dumb
	cmd.Env = append(os.Environ(), "TERM=dumb")
	out, err := cmd.Output()
	if err != nil {
		return "Error: Could not find manual page for '" + page + "'"
	}

	// Regex to clean up any remaining weirdness
	content := string(out)
	// Remove backspace overstrikes (e.g., 'm\bman' -> 'man')
	re := regexp.MustCompile(".\x08")
	content = re.ReplaceAllString(content, "")

	return content
}

func main() {
	// Get page from args or default to mtx_init
	pageName := "mtx_init"
	if len(os.Args) > 1 {
		pageName = os.Args[1]
	}

	manContent := getManPage(pageName)
	manLines := strings.Split(manContent, "\n")

	rand.Seed(time.Now().UnixNano())

	const (
		screenWidth  = 1000
		screenHeight = 900
	)

	rl.InitWindow(screenWidth, screenHeight, "COTN Docs - "+pageName)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Initialize Starfield
	stars := make([]Star, 100)
	for i := range stars {
		stars[i] = Star{
			x:     rand.Float32() * screenWidth,
			y:     rand.Float32() * screenHeight,
			speed: 0.1 + rand.Float32()*0.5,
			size:  0.5 + rand.Float32()*1.5,
		}
	}

	scrollY := float32(0)

	for !rl.WindowShouldClose() {
		// Input handling for scrolling
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			scrollY += wheel * 30
		}
		if rl.IsKeyDown(rl.KeyDown) {
			scrollY -= 5
		}
		if rl.IsKeyDown(rl.KeyUp) {
			scrollY += 5
		}

		// Clamp scrolling
		maxScroll := -float32(len(manLines)*22) + (screenHeight - 200)
		if scrollY > 0 {
			scrollY = 0
		}
		if scrollY < maxScroll && maxScroll < 0 {
			scrollY = maxScroll
		}

		// Update Stars
		for i := range stars {
			stars[i].y += stars[i].speed
			if stars[i].y > screenHeight {
				stars[i].y = 0
				stars[i].x = rand.Float32() * screenWidth
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(cotnMidnight)

		// 1. Starfield
		for _, s := range stars {
			rl.DrawCircle(int32(s.x), int32(s.y), s.size, rl.Fade(rl.White, 0.4))
		}

		// 2. Main Glass Container
		containerRect := rl.NewRectangle(80, 40, 840, 820)
		rl.DrawRectangleRec(containerRect, rl.NewColor(18, 7, 33, 230))
		rl.DrawRectangleLinesEx(containerRect, 1, rl.Fade(cotnNeonMagenta, 0.4))

		// Scissor mode to keep text inside the "glass"
		rl.BeginScissorMode(85, 45, 830, 810)

		// 3. Header (Static-ish but moves slightly with scroll if desired, here static)
		rl.DrawText(strings.ToUpper(pageName)+"(X)", 110, 70, 20, cotnNeonMagenta)
		rl.DrawLineEx(rl.NewVector2(110, 95), rl.NewVector2(140, 95), 2, cotnNeonMagenta)

		// 4. Render Man Page Content
		for i, line := range manLines {
			posY := int32(120 + scrollY + float32(i*22))

			// Simple logic to highlight Section Headers (all caps lines starting with no whitespace)
			isHeader := len(line) > 0 && !strings.HasPrefix(line, " ") && line == strings.ToUpper(line)

			if posY > -20 && posY < screenHeight {
				if isHeader {
					rl.DrawText(line, 110, posY, 20, cotnElectricBlue)
				} else {
					rl.DrawText(line, 110, posY, 16, cotnBrightText)
				}
			}
		}
		rl.EndScissorMode()

		// 5. HUD/Frame Elements
		rl.DrawLineEx(rl.NewVector2(80, 40), rl.NewVector2(920, 40), 2, cotnNeonMagenta)

		// Bottom Status
		rl.DrawRectangle(80, 830, 840, 30, rl.NewColor(5, 2, 10, 255))
		rl.DrawText("UP/DOWN OR WHEEL TO WALK THROUGH THE NIGHT", 320, 840, 12, cotnDimText)

		// 6. Post-Processing: Scanlines
		for i := int32(0); i < screenHeight; i += 3 {
			rl.DrawLine(0, i, screenWidth, i, rl.NewColor(0, 0, 0, 35))
		}

		rl.EndDrawing()
	}
}
