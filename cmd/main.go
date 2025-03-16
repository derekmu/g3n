package main

import (
	"github.com/derekmu/g3n/app"
	"github.com/derekmu/g3n/camera"
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/light"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/renderer"
	"log"
	"time"
)

var origin = math32.Vector3{}
var up = math32.Vector3{Y: 1}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	ap := app.NewApplication(1280, 720, "Demo")
	scene := core.NewNode()

	cam := camera.New(1)
	cam.SetPosition(0, 0, -10)
	cam.LookAt(&origin, &up)
	scene.Add(cam)

	dlight := light.NewDirectional(math32.Color{R: 1, G: 1, B: 1}, 1.0)
	dlight.SetPosition(0, 0, -10)
	dlight.LookAt(&origin, &up)
	scene.Add(dlight)

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

	geom := geometry.NewPlane(10, 10)
	mat := material.NewStandard(math32.Color{G: 1})
	mesh := graphic.NewMesh(geom, mat)
	mesh.RotateOnAxis(math32.Vector3{X: 1}, math32.Pi)
	mesh.SetPositionZ(1.1)
	scene.Add(mesh)

	geom = geometry.NewPlane(10, 10)
	mat = material.NewStandard(math32.Color{R: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshRed := graphic.NewMesh(geom, mat)
	meshRed.SetName("RED")
	meshRed.SetPositionZ(1)
	meshRed.RotateOnAxis(math32.Vector3{X: 1}, math32.Pi)
	scene.Add(meshRed)

	geom = geometry.NewPlane(10, 10)
	mat = material.NewStandard(math32.Color{B: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshBlue := graphic.NewMesh(geom, mat)
	meshBlue.SetName("BLUE")
	meshBlue.SetPositionZ(0.9)
	meshBlue.RotateOnAxis(math32.Vector3{X: 1}, math32.Pi)
	scene.Add(meshBlue)

	ap.Run(func(r *renderer.Renderer, dt time.Duration) {
		meshRed.RotateZ(float32(dt) / float32(time.Second) / 13 * math32.Pi)
		meshBlue.RotateZ(float32(dt) / float32(time.Second) / 7 * math32.Pi)

		ap.Gls().ClearColor(0.1, 0.1, 0.1, 1.0)
		ap.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		_ = r.Render(scene, cam)
	})
}
