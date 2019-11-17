function getMousePick(event, state) {

    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left
    const y = event.clientY - rect.top

    let normalX = (2 * x) / event.target.width;
    let normalY = 1 - (2 * y) / event.target.height;
    let normalZ = 1;

    let rayNds = vec3.fromValues(normalX, normalY, normalZ);


    let rayClip = vec4.fromValues(rayNds[0], rayNds[1], -1.0, 1.0);

    //console.log(rayClip);

    //console.log(state.projectionMatrix);

    let inverseProjection = mat4.create();
    mat4.invert(inverseProjection, state.projectionMatrix);


    vec4.transformMat4(rayClip, rayClip, inverseProjection);

    let rayEye = vec4.fromValues(rayClip[0], rayClip[1], -1.0, 0.0);
    let inverseViewMatrix = mat4.create();

    mat4.invert(inverseViewMatrix, state.viewMatrix);

    vec4.transformMat4(rayEye, rayEye, inverseViewMatrix);

    let rayWor = vec3.fromValues(rayEye[0], rayEye[1], rayEye[2]);

    vec3.normalize(rayWor, rayWor);

    console.log(rayWor);

    return rayWor;
    
    //let rayWor = vec3.fromValues(rayClip[0], rayClip[1], rayClip[2]);

    //vec3.normalize(rayWor, rayWor);

}