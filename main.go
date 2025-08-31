// Copyright (c) 2024, Samuel Title. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"log"
	"time"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/text/text"
	"cogentcore.org/core/xyz"
	"cogentcore.org/core/xyz/xyzcore"

	"cogentcore.org/core/math32"
)

// SimpleAnim handles animation for our 3D scene
type SimpleAnim struct {
	// Whether animation is running
	On bool

	// Animation speed
	Speed float32 `min:"0.01" step:"0.01"`

	// Current angle
	Angle float32 `edit:"-"`

	// Animation ticker
	Ticker *time.Ticker `display:"-"`

	// Scene editor reference
	SceneEditor *xyzcore.SceneEditor

	// Animated objects
	Cube   *xyz.Solid
	Sphere *xyz.Solid

	// Original positions
	CubePosOrig   math32.Vector3
	SpherePosOrig math32.Vector3
}

// Start initializes the animation
func (a *SimpleAnim) Start(se *xyzcore.SceneEditor, on bool) {
	a.SceneEditor = se
	a.On = on
	a.Speed = 0.05
	a.GetObjects()
	a.Ticker = time.NewTicker(time.Second / 30) // 30 fps
	go a.Animate()
}

// GetObjects finds the objects to animate
func (a *SimpleAnim) GetObjects() {
	sc := a.SceneEditor.SceneXYZ()

	cubeObj := sc.ChildByName("animated-cube", 0)
	if cubeObj == nil {
		log.Println("Couldn't find cube to animate")
		return
	}
	a.Cube = cubeObj.(*xyz.Solid)
	a.CubePosOrig = a.Cube.Pose.Pos

	sphereObj := sc.ChildByName("animated-sphere", 0)
	if sphereObj == nil {
		log.Println("Couldn't find sphere to animate")
		return
	}
	a.Sphere = sphereObj.(*xyz.Solid)
	a.SpherePosOrig = a.Sphere.Pose.Pos
}

// Animate runs the animation loop
func (a *SimpleAnim) Animate() {
	for {
		if a.Ticker == nil || a.SceneEditor.This == nil {
			return
		}
		<-a.Ticker.C // wait for tick
		if !a.On || a.SceneEditor.This == nil || a.Cube == nil || a.Sphere == nil {
			continue
		}

		// Calculate new positions
		radius := float32(0.5)

		// Move cube in a circle
		dx := radius * math32.Cos(a.Angle)
		dz := radius * math32.Sin(a.Angle)
		cubePos := a.CubePosOrig
		cubePos.X += dx
		cubePos.Z += dz
		a.Cube.SetPosePos(cubePos)

		// Move sphere in opposite direction
		spherePos := a.SpherePosOrig
		spherePos.X -= dx * 0.5
		spherePos.Z -= dz * 0.5
		a.Sphere.SetPosePos(spherePos)

		// Rotate cube
		a.Cube.Pose.SetAxisRotation(0, 1, 0, a.Angle*180/math32.Pi)

		// Update scene
		a.SceneEditor.SceneWidget().UpdateWidget()
		a.Angle += a.Speed
	}
}

func main() {
	// Create animation controller
	anim := &SimpleAnim{}

	// Create main body
	b := core.NewBody("Simple XYZ Demo")

	// Add title
	core.NewText(b).SetText(`Simple <b>XYZ</b> <i>3D</i> Demo`).
		SetType(core.TextHeadlineSmall).
		Styler(func(s *styles.Style) {
			s.Text.Align = text.Center
		})

	// Add animation control button
	animButton := core.NewButton(b).SetText("Start Animation")
	animButton.OnClick(func(e events.Event) {
		anim.On = !anim.On
		if anim.On {
			animButton.SetText("Stop Animation")
		} else {
			animButton.SetText("Start Animation")
		}
	})

	// Create scene editor
	se := xyzcore.NewSceneEditor(b)
	se.UpdateWidget()
	sw := se.SceneWidget()
	sc := se.SceneXYZ()
	sw.SelectionMode = xyzcore.Manipulable

	// Set up camera
	sc.Camera.Pose.Pos.Set(0, 3, 8)
	sc.Camera.LookAt(math32.Vector3{}, math32.Vec3(0, 1, 0))

	// Add lighting
	xyz.NewAmbient(sc, "ambient", 0.3, xyz.DirectSun)
	xyz.NewDirectional(sc, "directional", 1, xyz.DirectSun).Pos.Set(0, 2, 1)

	// Set background color
	se.Styler(func(s *styles.Style) {
		sc.Background = colors.Scheme.Select.Container
	})

	// Create floor
	floorMesh := xyz.NewPlane(sc, "floor-plane", 10, 10)
	floor := xyz.NewSolid(sc).SetMesh(floorMesh).
		SetColor(colors.Tan).SetPos(0, -1, 0)
	floor.SetName("floor")

	// Create 3D text
	text3D := xyz.NewText2D(sc).SetText("XYZ 3D Demo")
	text3D.Styles.Text.Align = text.Center
	text3D.Pose.Scale.SetScalar(0.2)
	text3D.SetPos(0, 2, 0)

	// Create animated cube
	cubeMesh := xyz.NewBox(sc, "cube-mesh", 1, 1, 1)
	cube := xyz.NewSolid(sc).SetMesh(cubeMesh).
		SetColor(colors.Blue).SetShiny(20).SetPos(-1.5, 0, 0)
	cube.SetName("animated-cube")

	// Create animated sphere
	sphereMesh := xyz.NewSphere(sc, "sphere-mesh", 0.5, 32)
	sphere := xyz.NewSolid(sc).SetMesh(sphereMesh).
		SetColor(colors.Orange).SetPos(1.5, 0, 0)
	sphere.SetName("animated-sphere")

	// Create cylinder
	cylinderMesh := xyz.NewCylinder(sc, "cylinder-mesh", 1.5, 0.3, 32, 1, true, true)
	cylinder := xyz.NewSolid(sc).SetMesh(cylinderMesh).
		SetColor(colors.Green).SetPos(0, 0, -2)
	cylinder.Pose.SetAxisRotation(1, 0, 0, 90)

	// Create semi-transparent torus
	torusMesh := xyz.NewTorus(sc, "torus-mesh", 0.7, 0.1, 32)
	torus := xyz.NewSolid(sc).SetMesh(torusMesh).
		SetColor(color.RGBA{255, 0, 255, 150}).SetPos(0, 1.5, 0)
	torus.Pose.SetAxisRotation(1, 0, 0, 45)

	// Create lines
	linesMesh := xyz.NewLines(sc, "lines", []math32.Vector3{
		{X: -2, Y: -0.5, Z: 2},
		{X: 0, Y: 1, Z: 2},
		{X: 2, Y: -0.5, Z: 2},
	}, math32.Vec2(0.1, 0.05), xyz.CloseLines)
	xyz.NewSolid(sc).SetMesh(linesMesh).SetColor(colors.Yellow)

	// Add arrow
	xyz.NewArrow(sc, sc, "arrow", math32.Vec3(-2, 0, 0), math32.Vec3(2, 0, 0),
		0.05, colors.Red, xyz.StartArrow, xyz.EndArrow, 4, 0.5, 8)

	// Start animation but don't run it yet
	anim.Start(se, false)

	// Run the application
	b.RunMainWindow()
}
