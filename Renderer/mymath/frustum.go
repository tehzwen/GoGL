package mymath

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CPlane struct {
	normal   mgl32.Vec3
	distance float32
}

type Frustum struct {
	Planes [6]CPlane
}

func (f Frustum) SphereIntersection(vecCenter mgl32.Vec3, radius float32) bool {
	for i := 0; i < 6; i++ {
		if vecCenter.Dot(f.Planes[i].normal)+f.Planes[i].distance+radius <= 0 {
			return false
		}
	}
	return true
}

func ConstructFrustrum(view mgl32.Mat4, proj mgl32.Mat4) Frustum {
	VP := proj.Mul4(view)

	//get all the planes sides in an array of vec4
	frustum := Frustum{}
	//near
	frustum.Planes[0].normal[0] = VP[3] + VP[2]
	frustum.Planes[0].normal[1] = VP[7] + VP[6]
	frustum.Planes[0].normal[2] = VP[11] + VP[10]
	frustum.Planes[0].distance = VP[15] + VP[14]

	//far
	frustum.Planes[1].normal[0] = VP[3] - VP[2]
	frustum.Planes[1].normal[1] = VP[7] - VP[6]
	frustum.Planes[1].normal[2] = VP[11] - VP[10]
	frustum.Planes[1].distance = VP[15] - VP[14]

	//left
	frustum.Planes[2].normal[0] = VP[3] + VP[0]
	frustum.Planes[2].normal[1] = VP[7] + VP[4]
	frustum.Planes[2].normal[2] = VP[11] + VP[8]
	frustum.Planes[2].distance = VP[15] + VP[12]

	//right
	frustum.Planes[3].normal[0] = VP[3] - VP[0]
	frustum.Planes[3].normal[1] = VP[7] - VP[4]
	frustum.Planes[3].normal[2] = VP[11] - VP[8]
	frustum.Planes[3].distance = VP[15] - VP[12]

	//up
	frustum.Planes[4].normal[0] = VP[3] - VP[1]
	frustum.Planes[4].normal[1] = VP[7] - VP[5]
	frustum.Planes[4].normal[2] = VP[11] - VP[9]
	frustum.Planes[4].distance = VP[15] - VP[13]

	//down
	frustum.Planes[5].normal[0] = VP[3] + VP[1]
	frustum.Planes[5].normal[1] = VP[7] + VP[5]
	frustum.Planes[5].normal[2] = VP[11] + VP[9]
	frustum.Planes[5].distance = VP[15] + VP[13]

	for i := 0; i < 6; i++ {
		frustum.Planes[i].normal = frustum.Planes[i].normal.Normalize()
	}

	return frustum
}
