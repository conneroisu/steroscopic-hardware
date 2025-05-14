package despair

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPNG(t *testing.T) {
	// Setup - create test images
	tempDir := t.TempDir()

	// Create grayscale test image
	grayPath := filepath.Join(tempDir, "gray.png")
	createGrayTestImage(t, grayPath)

	// Create RGBA test image
	rgbaPath := filepath.Join(tempDir, "rgba.png")
	createRGBATestImage(t, rgbaPath)

	// Create generic test image (NRGBA)
	nrgbaPath := filepath.Join(tempDir, "nrgba.png")
	createNRGBATestImage(t, nrgbaPath)

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"Load gray image", grayPath, false},
		{"Load RGBA image", rgbaPath, false},
		{"Load NRGBA image", nrgbaPath, false},
		{"Load nonexistent file", filepath.Join(tempDir, "nonexistent.png"), true},
		{"Load invalid file", filepath.Join(tempDir, "invalid.png"), true},
	}

	// Create invalid PNG file
	invalidPath := filepath.Join(tempDir, "invalid.png")
	if err := os.WriteFile(invalidPath, []byte("not a png"), 0644); err != nil {
		t.Fatalf("Failed to create invalid test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPNG(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPNG() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !tt.wantErr && got == nil {
				t.Error("LoadPNG() returned nil image without error")
			}
			if !tt.wantErr {
				// Verify image dimensions match the test images (10x10)
				if got.Bounds().Dx() != 10 || got.Bounds().Dy() != 10 {
					t.Errorf("LoadPNG() image dimensions = %v, want %v", got.Bounds().Size(), image.Point{10, 10})
				}
			}
		})
	}
}

func TestMustLoadPNG(t *testing.T) {
	// Setup - create test image
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test.png")
	createGrayTestImage(t, testPath)

	// Test successful load
	t.Run("Successful load", func(t *testing.T) {
		img := MustLoadPNG(testPath)
		if img == nil {
			t.Error("MustLoadPNG() returned nil for valid image")
		}
	})

	// Test panic behavior
	t.Run("Should panic on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustLoadPNG() did not panic for nonexistent file")
			}
		}()
		MustLoadPNG(filepath.Join(tempDir, "nonexistent.png"))
	})
}

func TestSavePNG(t *testing.T) {
	// Create a test image
	grayImg := image.NewGray(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			grayImg.SetGray(x, y, color.Gray{uint8(x * y % 256)})
		}
	}

	tempDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		img      image.Image
		wantErr  bool
	}{
		{
			"Save to valid path",
			filepath.Join(tempDir, "valid.png"),
			grayImg,
			false,
		},
		{
			"Save to invalid directory",
			filepath.Join(tempDir, "nonexistent", "test.png"),
			grayImg,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SavePNG(tt.filename, tt.img)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"SavePNG() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !tt.wantErr {
				// Verify the file exists
				_, err := os.Stat(tt.filename)
				if err != nil {
					t.Errorf(
						"SavePNG() didn't create file at %s",
						tt.filename,
					)
				}

				// Verify the file can be loaded
				file, err := os.Open(tt.filename)
				if err != nil {
					t.Errorf(
						"SavePNG() created file that can't be opened: %v",
						err,
					)

					return
				}
				defer file.Close()

				_, err = png.Decode(file)
				if err != nil {
					t.Errorf(
						"SavePNG() created file that can't be decoded as PNG: %v",
						err,
					)
				}
			}
		})
	}
}

func TestMustSavePNG(t *testing.T) {
	// Create a test image
	grayImg := image.NewGray(image.Rect(0, 0, 10, 10))

	tempDir := t.TempDir()
	validPath := filepath.Join(tempDir, "valid.png")
	invalidPath := filepath.Join(tempDir, "nonexistent", "test.png")

	// Test successful save
	t.Run("Successful save", func(t *testing.T) {
		MustSavePNG(validPath, grayImg)
		// Verify the file exists
		_, err := os.Stat(validPath)
		if err != nil {
			t.Errorf(
				"MustSavePNG() didn't create file at %s",
				validPath,
			)
		}
	})

	// Test panic behavior
	t.Run("Should panic on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSavePNG() did not panic for invalid path")
			}
		}()
		MustSavePNG(invalidPath, grayImg)
	})
}

func TestEndToEnd(t *testing.T) {
	// Create a test image
	tempDir := t.TempDir()
	origPath := filepath.Join(tempDir, "original.png")
	createRGBATestImage(t, origPath)

	// Test save and load cycle
	t.Run("End-to-end test", func(t *testing.T) {
		// Load the image
		img, err := LoadPNG(origPath)
		if err != nil {
			t.Fatalf("LoadPNG() error = %v", err)
		}

		// Save the image
		savePath := filepath.Join(tempDir, "saved.png")
		err = SavePNG(savePath, img)
		if err != nil {
			t.Fatalf("SavePNG() error = %v", err)
		}

		// Re-load the image
		reloaded, err := LoadPNG(savePath)
		if err != nil {
			t.Fatalf("LoadPNG() error = %v", err)
		}

		// Compare dimensions
		if !reloaded.Bounds().Eq(img.Bounds()) {
			t.Errorf(
				"End-to-end image dimensions don't match: got %v, want %v",
				reloaded.Bounds(),
				img.Bounds(),
			)
		}

		// Compare some pixel values
		for y := 0; y < 10; y += 3 {
			for x := 0; x < 10; x += 3 {
				if reloaded.GrayAt(x, y).Y != img.GrayAt(x, y).Y {
					t.Errorf(
						"End-to-end pixel at (%d,%d) doesn't match: got %v, want %v",
						x,
						y,
						reloaded.GrayAt(x, y).Y,
						img.GrayAt(x, y).Y,
					)
				}
			}
		}
	})
}

// Helper functions to create test images

func createGrayTestImage(t *testing.T, path string) {
	img := image.NewGray(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.SetGray(x, y, color.Gray{uint8((x + y) % 256)})
		}
	}

	err := saveTestImage(path, img)
	if err != nil {
		t.Fatalf("Failed to create test gray image: %v", err)
	}
}

func createRGBATestImage(t *testing.T, path string) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(x * 25 % 256),
				G: uint8(y * 25 % 256),
				B: uint8((x + y) * 15 % 256),
				A: 255,
			})
		}
	}

	err := saveTestImage(path, img)
	if err != nil {
		t.Fatalf("Failed to create test RGBA image: %v", err)
	}
}

func createNRGBATestImage(t *testing.T, path string) {
	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x * 25 % 256),
				G: uint8(y * 25 % 256),
				B: uint8((x + y) * 15 % 256),
				A: 255,
			})
		}
	}

	err := saveTestImage(path, img)
	if err != nil {
		t.Fatalf("Failed to create test NRGBA image: %v", err)
	}
}

func saveTestImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
