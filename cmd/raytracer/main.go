package main

import (
	"fmt"
	"os"
	"raytracer/internal/color"
)

func main() {
	// Image dimensions
	imageWidth := 256
	imageHeight := 256

	// Header for PPM format
	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	for j := 0; j < imageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d", imageHeight-j)
		for i := 0; i < imageWidth; i++ {
			pixelColor := color.NewColor(
				float64(i)/float64(imageWidth-1),
				float64(j)/float64(imageHeight-1),
				0.0,
			)
			color.WriteColor(os.Stdout, pixelColor)
		}
	}

	fmt.Fprintln(os.Stderr, "\nDone.")
}
