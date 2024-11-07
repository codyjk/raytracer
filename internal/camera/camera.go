package camera

import (
	"fmt"
	"io"
	"math"
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/hittable"
	"raytracer/internal/interval"
	"raytracer/internal/ray"
	"raytracer/internal/util"
	"raytracer/internal/vector"
)

// Config holds all camera configuration parameters
type Config struct {
	AspectRatio     float64 // Ratio of image width over height
	ImageWidth      int     // Rendered image width in pixel count
	SamplesPerPixel int     // Count of random samples for each pixel
	MaxDepth        int     // Maximum number of ray bounces into scene

	VerticalFOV float64       // Vertical view angle (field of view)
	LookFrom    vector.Point3 // Point camera is looking from
	LookAt      vector.Point3 // Point camera is looking at
	VUp         vector.Vec3   // Camera-relative "up" direction

	DefocusAngle float64 // Variation angle of rays through each pixel
	FocusDist    float64 // Distance from camera lookFrom point to plane of perfect focus
}

// DefaultConfig returns a Config with reasonable default values
func DefaultConfig() Config {
	return Config{
		AspectRatio:     16.0 / 9.0,
		ImageWidth:      400,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		VerticalFOV:     20.0,
		LookFrom:        vector.NewPoint3(-2, 2, 1),
		LookAt:          vector.NewPoint3(0, 0, -1),
		VUp:             vector.NewVec3(0, 1, 0),
		DefocusAngle:    0.0,
		FocusDist:       10.0,
	}
}

// Camera represents a virtual camera for ray tracing
type Camera struct {
	config           Config
	imageHeight      int           // Rendered image height
	pixelSampleScale float64       // Color scale factor for a sum of pixel samples
	center           vector.Point3 // Camera center
	pixel00Location  vector.Point3 // Offset of pixel 0,0
	pixelDeltaU      vector.Vec3   // Offset to pixel to the right
	pixelDeltaV      vector.Vec3   // Offset to pixel below
	basis            struct {
		// Camera frame basis vectors
		u, v, w vector.Vec3
	}
	defocusDiskU vector.Vec3 // Defocus disk horiztonal radius
	defocusDiskV vector.Vec3 // Defocus disk vertical radius
}

// New creates a new Camera with the given configuration
func New(cfg Config) (*Camera, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid camera configuration: %w", err)
	}

	cam := &Camera{
		config:           cfg,
		pixelSampleScale: 1.0 / float64(cfg.SamplesPerPixel),
	}

	if err := cam.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize camera: %w", err)
	}

	return cam, nil
}

func validateConfig(cfg Config) error {
	if cfg.AspectRatio <= 0 {
		return fmt.Errorf("aspect ratio must be positive, got %v", cfg.AspectRatio)
	}
	if cfg.ImageWidth <= 0 {
		return fmt.Errorf("image width must be positive, got %v", cfg.ImageWidth)
	}
	if cfg.SamplesPerPixel <= 0 {
		return fmt.Errorf("samples per pixel must be positive, got %v", cfg.SamplesPerPixel)
	}
	if cfg.MaxDepth <= 0 {
		return fmt.Errorf("max depth must be positive, got %v", cfg.MaxDepth)
	}
	return nil
}

func (c *Camera) initialize() error {
	// Calculate image height maintaining minimum of 1
	c.imageHeight = int(float64(c.config.ImageWidth) / c.config.AspectRatio)
	if c.imageHeight < 1 {
		c.imageHeight = 1
	}

	c.center = c.config.LookFrom

	// Calculate viewport dimensions
	theta := util.DegreesToRadians(c.config.VerticalFOV)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h * c.config.FocusDist
	viewportWidth := viewportHeight * (float64(c.config.ImageWidth) / float64(c.imageHeight))

	// Calculate camera basis vectors
	c.basis.w = c.config.LookFrom.Sub(c.config.LookAt).Unit()
	c.basis.u = vector.Cross(c.config.VUp, c.basis.w).Unit()
	c.basis.v = vector.Cross(c.basis.w, c.basis.u)

	// Calculate viewport vectors
	viewportU := c.basis.u.Scale(viewportWidth)
	// Invert V to match image coordinate system (top-left origin)
	viewportV := c.basis.v.Scale(-1.0 * viewportHeight)

	// Calculate pixel delta vectors
	c.pixelDeltaU = viewportU.Div(float64(c.config.ImageWidth))
	c.pixelDeltaV = viewportV.Div(float64(c.imageHeight))

	// Calculate upper left pixel location
	viewportUpperLeft := c.center.
		Sub(c.basis.w.Scale(c.config.FocusDist)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))
	c.pixel00Location = viewportUpperLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Scale(0.5))

	// Calculate the camera defocus disk basis vectors.
	defocusRadius := c.config.FocusDist * math.Tan(util.DegreesToRadians(c.config.DefocusAngle/2))
	c.defocusDiskU = c.basis.u.Scale(defocusRadius)
	c.defocusDiskV = c.basis.v.Scale(defocusRadius)

	return nil
}

// Render renders the scene to the provided writer
func (c *Camera) Render(out io.Writer, log io.Writer, world hittable.Hittable) error {
	if _, err := fmt.Fprintf(out, "P3\n%d %d\n255\n", c.config.ImageWidth, c.imageHeight); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for j := 0; j < c.imageHeight; j++ {
		fmt.Fprintf(log, "\rScanlines remaining: %d", c.imageHeight-j)
		if err := c.renderScanline(j, out, world); err != nil {
			return fmt.Errorf("failed to render scanline %d: %w", j, err)
		}
	}

	fmt.Fprintln(log, "\nDone.")
	return nil
}

func (c *Camera) renderScanline(j int, out io.Writer, world hittable.Hittable) error {
	for i := 0; i < c.config.ImageWidth; i++ {
		pixelColor := c.samplePixel(i, j, world)
		if err := color.WriteColor(out, pixelColor.Scale(c.pixelSampleScale)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Camera) samplePixel(i, j int, world hittable.Hittable) color.Color {
	pixelColor := color.NewColor(0, 0, 0)
	for sample := 0; sample < c.config.SamplesPerPixel; sample++ {
		r := c.getRay(i, j)
		pixelColor = pixelColor.Add(c.traceRay(r, c.config.MaxDepth, world))
	}
	return pixelColor
}

func (c *Camera) getRay(i, j int) ray.Ray {
	pixelCenter := c.pixel00Location.
		Add(c.pixelDeltaU.Scale(float64(i))).
		Add(c.pixelDeltaV.Scale(float64(j)))

	// Add random offset within pixel for anti-aliasing
	offset := samplePixelSquare()
	pixelSample := pixelCenter.
		Add(c.pixelDeltaU.Scale(offset.X())).
		Add(c.pixelDeltaV.Scale(offset.Y()))

	rayOrigin := c.center
	if c.config.DefocusAngle > 0 {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	return ray.NewRay(rayOrigin, rayDirection)
}

func (c *Camera) traceRay(r ray.Ray, depth int, world hittable.Hittable) color.Color {
	if depth <= 0 {
		return color.NewColor(0, 0, 0)
	}

	var rec core.HitRecord
	if world.Hit(r, interval.NewInterval(0.001, math.Inf(1)), &rec) {
		var scattered ray.Ray
		var attenuation color.Color
		if rec.Material().Scatter(r, &rec, &attenuation, &scattered) {
			return attenuation.Mul(c.traceRay(scattered, depth-1, world))
		}
		return color.NewColor(0, 0, 0)
	}

	// Render sky gradient
	unitDirection := r.Direction().Unit()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return color.NewColor(1.0, 1.0, 1.0).Scale(1.0 - t).Add(color.NewColor(0.5, 0.7, 1.0).Scale(t))
}

func (c Camera) defocusDiskSample() vector.Point3 {
	// Returns a random point in the camera defocus disk.
	p := vector.RandomInUnitDisk()
	return c.center.
		Add(c.defocusDiskU.Scale(p.X())).
		Add(c.defocusDiskV.Scale(p.Y()))
}

// samplePixelSquare returns a random point in the [-0.5,0.5] square
func samplePixelSquare() vector.Vec3 {
	return vector.NewVec3(
		util.RandomFloat()-0.5,
		util.RandomFloat()-0.5,
		0,
	)
}
