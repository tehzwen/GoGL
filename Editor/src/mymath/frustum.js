class Frustum {
    constructor(projection, view){
        this.planes = []

        //construct the frustum from the view and projection matrices
        let VP = mat4.create();
        mat4.mul(VP, projection, view);
        
        //near
        let nearPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        nearPlane.normal = vec3.fromValues(VP[3] + VP[2], VP[7] + VP[6], VP[11] + VP[10]);
        nearPlane.distance = VP[15] + VP[14];
        this.planes.push(nearPlane);

        //far
        let farPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        farPlane.normal = vec3.fromValues(VP[3] - VP[2], VP[7] - VP[6], VP[11] - VP[10]);
        farPlane.distance = VP[15] - VP[14];
        this.planes.push(farPlane);

        //left
        let leftPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        leftPlane.normal = vec3.fromValues(VP[3] + VP[0], VP[7] + VP[4], VP[11] + VP[8]);
        leftPlane.distance = VP[15] + VP[12];
        this.planes.push(leftPlane);

        //right
        let rightPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        rightPlane.normal = vec3.fromValues(VP[3] - VP[0], VP[7] - VP[4], VP[11] - VP[8]);
        rightPlane.distance = VP[15] - VP[12];
        this.planes.push(rightPlane);

        //up
        let upPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        upPlane.normal = vec3.fromValues(VP[3] - VP[1], VP[7] - VP[5], VP[11] - VP[9]);
        upPlane.distance = VP[15] - VP[13];
        this.planes.push(upPlane);

        //down
        let downPlane = {
            normal: vec3.fromValues(0,0,0),
            distance: null
        }
        downPlane.normal = vec3.fromValues(VP[3] + VP[1], VP[7] + VP[5], VP[11] + VP[9]);
        downPlane.distance = VP[15] + VP[13];
        this.planes.push(downPlane);

        
        //normalize the normals of each plane

        for (let i = 0; i < 6; i++) {
            vec3.normalize(this.planes[i].normal, this.planes[i].normal);
        }
    }

    sphereIntersection(vecCenter, radius) {
        for (let i = 0; i < 6; i++) {
            if (vec3.dot(vecCenter, this.planes[i].normal) + this.planes[i].distance + radius <= 0 ) {
                return false;
            }
        }
        return true;
    }
}

export default Frustum;