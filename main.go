package main

import (
	"fmt"
	"math/rand"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		screenWidth  = 1000
		screenHeight = 850
	)

	rl.InitWindow(screenWidth, screenHeight, "mtx(3) - Call of the Night Manual")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Initialize Starfield
	stars := make([]Star, 100)
	for i := range stars {
		stars[i] = Star{
			x:     rand.Float32() * screenWidth,
			y:     rand.Float32() * screenHeight,
			speed: 0.2 + rand.Float32()*0.8,
			size:  0.5 + rand.Float32()*1.5,
		}
	}

	for !rl.WindowShouldClose() {
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

		// 1. Draw Background Glow
		rl.DrawCircleGradient(200, 200, 600, rl.NewColor(26, 11, 46, 100), rl.Blank)

		// 2. Draw Starfield
		for _, s := range stars {
			rl.DrawCircle(int32(s.x), int32(s.y), s.size, rl.Fade(rl.White, 0.6))
		}

		// 3. Draw Main Glass Container
		containerRect := rl.NewRectangle(100, 50, 800, 750)
		rl.DrawRectangleRec(containerRect, rl.NewColor(18, 7, 33, 220)) // Glass fill
		rl.DrawRectangleLinesEx(containerRect, 1, rl.Fade(cotnNeonMagenta, 0.5))

		// Neon Accent Line (Top)
		rl.DrawLineEx(rl.NewVector2(100, 50), rl.NewVector2(900, 50), 2, cotnNeonMagenta)

		// 4. Header
		rl.DrawText("MTX(3)", 120, 70, 20, cotnNeonMagenta)
		rl.DrawText("C11 THREADS MANUAL", 400, 70, 16, cotnDimText)
		rl.DrawText("MTX(3)", 820, 70, 20, cotnNeonMagenta)

		// 5. Documentation Content
		yOff := int32(130)

		// Section: NAME
		drawSectionHeader("NAME", 120, yOff)
		yOff += 35
		rl.DrawText("mtx_init, mtx_lock, mtx_unlock - mutex primitives", 120, yOff, 18, cotnBrightText)
		yOff += 60

		// Section: SYNOPSIS
		drawSectionHeader("SYNOPSIS", 120, yOff)
		yOff += 35
		// Code Block Background
		rl.DrawRectangle(120, yOff, 760, 140, rl.NewColor(5, 2, 10, 255))
		rl.DrawRectangleLines(120, yOff, 760, 140, rl.Fade(cotnElectricBlue, 0.3))
		
		rl.DrawText("#include <threads.h>", 140, yOff+15, 18, cotnDimText)
		rl.DrawText("int mtx_init(mtx_t *mtx, int type);", 140, yOff+45, 18, cotnElectricBlue)
		rl.DrawText("int mtx_lock(mtx_t *mtx);", 140, yOff+75, 18, cotnElectricBlue)
		rl.DrawText("int mtx_unlock(mtx_t *mtx);", 140, yOff+105, 18, cotnElectricBlue)
		yOff += 170

		// Section: DESCRIPTION
		drawSectionHeader("DESCRIPTION", 120, yOff)
		yOff += 35
		desc := "In the quiet of the execution cycle, mutexes serve as gatekeepers.\nA C11 mutex object ensures that only one thread can wander\nthrough a critical section at any given time."
		rl.DrawText(desc, 120, yOff, 18, cotnBrightText)
		yOff += 100

		// Types
		rl.DrawText("Types:", 120, yOff, 18, cotnElectricBlue)
		yOff += 25
		rl.DrawText("- mtx_plain: Simple non-recursive lock.", 140, yOff, 18, cotnDimText)
		yOff += 22
		rl.DrawText("- mtx_recursive: Allows re-entry by the same thread.", 140, yOff, 18, cotnDimText)

		// 6. Footer
		rl.DrawLine(120, 760, 880, 760, 1, rl.Fade(cotnDimText, 0.3))
		rl.DrawText("C11 STANDARD - TOKYO NIGHT EDITION", 350, 775, 14, cotnDimText)

		// 7. Post-Processing: Scanlines
		for i := int32(0); i < screenHeight; i += 3 {
			rl.DrawLine(0, i, screenWidth, i, rl.NewColor(0, 0, 0, 40))
		}

		rl.EndDrawing()
	}
}

func drawSectionHeader(text string, x, y int32) {
	// Glow line
	rl.DrawLine(x, y+10, x+30, y+10, 2, cotnNeonMagenta)
	rl.DrawText(text, x+40, y, 20, cotnNeonMagenta)
}
