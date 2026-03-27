package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/derekmu/g3n/app"
	"github.com/derekmu/g3n/camera"
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/gui"
	"github.com/derekmu/g3n/light"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/renderer"
	"github.com/derekmu/g3n/texture"
	"github.com/derekmu/g3n/util/stats"
)

var up = math32.Vector3{Y: 1}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	ap := app.NewApplication(720, 720, "Demo")
	scene := core.NewNode()

	geom := geometry.NewPlane(10, 10)

	physicalMat := material.NewPhysical()
	physicalMat.SetMetallicFactor(0.5)
	physicalMat.SetRoughnessFactor(0.5)
	colorMap := image.NewRGBA(image.Rect(0, 0, 255, 255))
	normalMap := image.NewRGBA(image.Rect(0, 0, 255, 255))
	for y := range 255 {
		for x := range 255 {
			colorMap.Set(x, y, color.RGBA{
				R: uint8(x),
				G: uint8(y),
				B: uint8(255 - x),
				A: 255,
			})
			v := math32.Vector3{
				X: math32.Sin(float32(x+y/3)/255*math32.Pi*13) / 5,
				Y: math32.Cos(float32(y+x/7)/255*math32.Pi*7) / 5,
				Z: 1,
			}
			v.Normalize()
			normalMap.Set(x, y, color.RGBA{
				R: uint8(v.X*127.5 + 127.5),
				G: uint8(v.Y*127.5 + 127.5),
				B: uint8(v.Z*127.5 + 127.5),
			})
		}
	}
	physicalMat.SetBaseColorMap(texture.NewTexture2DFromRGBA(colorMap))
	physicalMat.SetNormalMap(texture.NewTexture2DFromRGBA(normalMap))
	meshPhysical := graphic.NewMesh(geom, physicalMat)
	meshPhysical.SetName("PHYSICAL")
	meshPhysical.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshPhysical.SetPositionY(-10.1)
	scene.Add(meshPhysical)

	mat := material.NewBlinnPhong(math32.Color3{R: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshRed := graphic.NewMesh(geom, mat)
	meshRed.SetName("RED")
	meshRed.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshRed.SetPositionY(-10.0)
	scene.Add(meshRed)

	mat = material.NewStandard(math32.Color3{B: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshBlue := graphic.NewMesh(geom, mat)
	meshBlue.SetName("BLUE")
	meshBlue.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshBlue.SetPositionY(-9.999)
	scene.Add(meshBlue)

	cam := camera.New(1)
	cam.LookAt(meshPhysical.Position(), up)
	scene.Add(cam)

	dlight := light.NewDirectional(math32.Color3{R: 1, G: 1, B: 1})
	scene.Add(dlight)

	// test font rendering
	panel := gui.NewPanel(0, 0)
	panel.SetPaddings(gui.RectBounds{Top: 5, Right: 10, Bottom: 5, Left: 10})
	panel.SetPanelColor(math32.Color3{R: 0.5, G: 0.5, B: 0.5})
	scene.Add(panel)

	label := panel.AddLabel("This is a demo", true, gui.AlignCenterCenter)
	label.SetFontSize(40)
	label.SetColor(math32.Color4{G: 1, A: 1.0})
	label.FitToText()

	stat := stats.NewStats(ap.Gls())
	statTable := stats.NewStatsTable()
	statTable.SetPosition(0, float32(panel.Height()))
	scene.Add(statTable)

	ap.Subscribe(func(ev core.WindowEvent) bool {
		switch ev := ev.(type) {
		case core.WindowSizeEvent:
			ap.Gls().Viewport(0, 0, int32(ev.Width), int32(ev.Height))

			cam.SetAspect(float32(ev.Width) / float32(ev.Height))
		default:
			return false
		}
		return true
	})

	var st time.Duration
	ap.Run(func(r *renderer.Renderer, dt time.Duration) {
		meshRed.RotateZ(float32(dt) / float32(time.Second) / 13 * math32.Pi)
		meshBlue.RotateZ(float32(dt) / float32(time.Second) / 7 * math32.Pi)

		st += dt
		v := math32.Vector3{X: 5, Y: 5}
		v.ApplyQuaternion(math32.QuaternionFromAxisAngle(up, float32(st)/float32(time.Second)/17*math32.Pi))
		cam.SetPositionVec(v)
		cam.LookAt(meshPhysical.Position(), up)

		v = math32.Vector3{X: 5, Y: 5}
		v.ApplyQuaternion(math32.QuaternionFromAxisAngle(up, float32(st)/float32(time.Second)/5*math32.Pi))
		dlight.SetPositionVec(v)
		dlight.LookAt(meshPhysical.Position(), up)

		if stat.Update(time.Second) {
			statTable.Update(stat)
		}

		ap.Gls().ClearColor(0.1, 0.1, 0.1, 1.0)
		ap.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		_ = r.Render(scene, cam)
	})
}
