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

var up = math32.Vector3{Y: 1}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	ap := app.NewApplication(720, 720, "Demo")
	scene := core.NewNode()

	geom := geometry.NewPlane(10, 10)

	greenMat := material.NewPhysical()
	greenMat.SetBaseColorFactor(math32.Color4{G: 0.5, A: 1})
	greenMat.SetMetallicFactor(0.5)
	greenMat.SetRoughnessFactor(0.5)
	meshGreen := graphic.NewMesh(geom, greenMat)
	meshGreen.SetName("GREEN")
	meshGreen.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshGreen.SetPositionY(-10.1)
	scene.Add(meshGreen)

	mat := material.NewBlinnPhong(math32.Color{R: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshRed := graphic.NewMesh(geom, mat)
	meshRed.SetName("RED")
	meshRed.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshRed.SetPositionY(-10.0)
	scene.Add(meshRed)

	mat = material.NewStandard(math32.Color{B: 1})
	mat.SetTransparent(true)
	mat.SetOpacity(0.5)
	meshBlue := graphic.NewMesh(geom, mat)
	meshBlue.SetName("BLUE")
	meshBlue.RotateOnAxis(math32.Vector3{X: 1}, -math32.Pi/2)
	meshBlue.SetPositionY(-9.999)
	scene.Add(meshBlue)

	cam := camera.New(1)
	cam.LookAt(meshGreen.Position(), up)
	scene.Add(cam)

	dlight := light.NewDirectional(math32.Color{R: 1, G: 1, B: 1}, 1.0)
	dlight.SetPosition(0, 10, 0)
	dlight.LookAt(meshGreen.Position(), up)
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

	var st time.Duration
	ap.Run(func(r *renderer.Renderer, dt time.Duration) {
		meshRed.RotateZ(float32(dt) / float32(time.Second) / 13 * math32.Pi)
		meshBlue.RotateZ(float32(dt) / float32(time.Second) / 7 * math32.Pi)

		st += dt
		v := math32.Vector3{X: 5, Y: 5}
		v.ApplyQuaternion(math32.QuaternionFromAxisAngle(up, float32(st)/float32(time.Second)/23*math32.Pi))
		cam.SetPositionVec(v)
		cam.LookAt(meshGreen.Position(), up)

		ap.Gls().ClearColor(0.1, 0.1, 0.1, 1.0)
		ap.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		_ = r.Render(scene, cam)
	})
}
