package main
import (
	"image/color"
	"golang.org/x/image/colornames"
	"math"
	"github.com/go-gl/mathgl/mgl64"
)
type ray struct {
	origin mgl64.Vec3
	direction mgl64.Vec3
}
type light interface {
}
type intersection struct {
	 material_color color.RGBA
	 hit_location mgl64.Vec3
	 normal mgl64.Vec3
}
type point_light struct {
	color color.RGBA
	location mgl64.Vec3
	attenuation float64
}
type geometry interface {
	intersects(ray ray) (bool, intersection)
}
type sphere struct {
	colour   color.RGBA
	location mgl64.Vec3
	radius   float64
}
func (s sphere) intersects(r ray) (bool, intersection) {
	var oToLoc = s.location.Sub(r.origin)
	var midpointD = oToLoc.Dot(r.direction)
	if(midpointD <= 0) {
		return false, intersection{s.colour, mgl64.Vec3{0,0,0}, mgl64.Vec3{0,0,0}}
	}
	var distanceFromLoc = math.Sqrt(oToLoc.Dot(oToLoc) - midpointD*midpointD)
	if(distanceFromLoc <= 0) {
		return false, intersection{s.colour, mgl64.Vec3{0,0,0}, mgl64.Vec3{0,0,0}}
	}
	if(distanceFromLoc > s.radius) {
		return false, intersection{s.colour, mgl64.Vec3{0,0,0}, mgl64.Vec3{0,0,0}}
	}
	var distance = midpointD - math.Sqrt(s.radius*s.radius - distanceFromLoc*distanceFromLoc)
	var location = r.direction.Mul(distance)
	var norm = (location.Sub(s.location)).Normalize()
	return true, intersection{s.colour, location, norm}
}
func lightingPass(hit intersection) color.RGBA {
	// refactor this to be in a proper place, e.g. recursive trace() function
	var lightDirection = light_sources.location.Sub(hit.hit_location)
	var lightDistance = lightDirection.Len()
	var inShadow = false
	for obj := 0; obj < len(scene); obj++ {
		hit, location := scene[obj].intersects(ray{hit.hit_location,lightDirection.Normalize()})
		if(hit && location.hit_location.Len() < lightDistance && location.hit_location.Len() != 0) {
			inShadow = true
		}
	}
	if(inShadow) {
		return scaleColor(ambient_coefficient,multiplyColour(hit.material_color,ambient_light))
	} else {
		var lightDotNormal = math.Max(0.0, lightDirection.Normalize().Dot(hit.normal.Normalize()))
		var specular = scaleColor(light_sources.attenuation * lightDotNormal ,multiplyColour(hit.material_color,light_sources.color))
		var ambient = scaleColor(ambient_coefficient,multiplyColour(hit.material_color,ambient_light))
		return addColor(ambient,specular)
	}
}
func addColor(lhs color.RGBA, rhs color.RGBA) color.RGBA {
	var r = uint8(lhs.R + rhs.R)
	var g = uint8(lhs.G + rhs.G)
	var b = uint8(lhs.B + rhs.B)
	var a = uint8(lhs.A + rhs.A)
	return color.RGBA{R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a)}
}
func scaleColor(cheatshade float64, colour color.RGBA ) color.RGBA {
	return color.RGBA{uint8(cheatshade * float64(colour.R)),
		uint8(cheatshade * float64(colour.G)),
		uint8(cheatshade * float64(colour.B)),
			colour.A}
}
func multiplyColour(lhs color.RGBA, rhs color.RGBA) color.RGBA {
	var r = float64(lhs.R) * float64(rhs.R) / 255
	var g = float64(lhs.G) * float64(rhs.G) / 255
	var b = float64(lhs.B) * float64(rhs.B) / 255
	var a = float64(lhs.A) * float64(rhs.A) / 255
	return color.RGBA{R: uint8(r),
					  G: uint8(g),
					  B: uint8(b),
					  A: uint8(a)}
}
var camera = mgl64.Vec3{0,0,0}
const fov = float64(90) // should actually use this not the hard coded ones. Ooops
const frustrumWidth = 16
const frustrumHeight = 9
const frustrumDistance = 8
var scene = []geometry{ sphere{colour: colornames.Darkviolet, location:mgl64.Vec3{6,-2,24}, radius: 8},
 						 sphere{colour: colornames.Darkgreen, location:mgl64.Vec3{-7,4,16}, radius: 3}}
var light_sources = point_light{colornames.White, mgl64.Vec3{16,0, 16}, 1.0}
var ambient_light = colornames.White
var ambient_coefficient = 0.2
func transformCoordinate(x int, y int, width int, height int) mgl64.Vec3 {
	var pixelWidth float64 = float64(frustrumWidth) / float64(width)
	var pixelHeight float64 = float64(frustrumHeight) / float64(height)
	var pixelX = float64(x) * pixelWidth - float64(frustrumWidth)/2
	var pixelY = float64(y) * pixelHeight - float64(frustrumHeight)/2
	return mgl64.Vec3{pixelX,pixelY,frustrumDistance}
}
func DrawToBuffer(buffer[]uint8,width int ,height int, stride int) {
	for x := 0; x < width; x++ {
		for y:=0; y < height; y++ {
			pixel := transformCoordinate(x,y,width,height)
			ray := ray{camera, pixel.Normalize()}
			var finalColour color.RGBA
			var smallestDistance float64 = -1
			var collision intersection
			for obj := 0; obj < len(scene); obj++ {
				hit,location := scene[obj].intersects(ray)
				if(hit) {
					distance := location.hit_location.Len()
					if(smallestDistance == -1) {
						smallestDistance = distance
						collision = location
					}
					if(distance < smallestDistance) {
						smallestDistance = distance
						collision = location
					}
				}
			}
			if(smallestDistance != -1 ) {
				finalColour = lightingPass(collision)
			} else {
				finalColour = colornames.Black
			}
			var base = (x + y*width) * stride
			buffer[base] = finalColour.R
			buffer[base+1] = finalColour.G
			buffer[base+2] = finalColour.B
			buffer[base+3] = finalColour.A
		}
	}
}