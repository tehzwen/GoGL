var moveSpeed = 2;

function startGame(state) {
    document.addEventListener("contextmenu", function (e) {
        e.preventDefault();
    }, false);

    document.addEventListener('mousemove', (event) => {
        //handle right click
        if (event.buttons == 2) {
            state.mouse['camMove'] = true;
            state.mouse.rateX = event.movementX;
        }
    });

    document.addEventListener('mouseup', (event) => {
        state.mouse['camMove'] = false;
        state.mouse.rateX = 0;
    })

    document.addEventListener('keypress', (event) => {
        switch (event.code) {
            case "KeyW":
                state.keyboard[event.key] = true;
                break;

            case "KeyS":
                state.keyboard[event.key] = true;
                break;

            case "KeyA":
                state.keyboard[event.key] = true;
                break;

            case "KeyD":
                state.keyboard[event.key] = true;
                break;

            default:
                break;
        }
    });

    document.addEventListener('keyup', (event) => {
        switch (event.code) {
            case "KeyW":
                state.keyboard[event.key] = false;
                break;

            case "KeyS":
                state.keyboard[event.key] = false;
                break;

            case "KeyA":
                state.keyboard[event.key] = false;
                break;

            case "KeyD":
                state.keyboard[event.key] = false;
                break;

            case "KeyZ":
                state.lightIndices.map((index) => {
                    let light = state.objects[index];
                    light.strength += 0.5;
                })
                break;

            case "KeyX":
                state.lightIndices.map((index) => {
                    let light = state.objects[index];
                    light.strength -= 0.5;
                })
                break;

            default:
                break;
        }
    });
}

function printForwardVector(state, important = null) {
    let vMatrix = state.viewMatrix;
    if (important) {
        console.warn(vMatrix[2], vMatrix[6], vMatrix[10])
    } else {
        console.log(vMatrix[2], vMatrix[6], vMatrix[10])
    }

}

function moveForward(state) {
    let inverseView = mat4.create(), forwardVector = vec3.fromValues(0, 0, 0), cameraCenterVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    mat4.invert(inverseView, state.viewMatrix);
    //forward vector from the viewmatrix
    forwardVector = vec3.fromValues(inverseView[2], inverseView[6], -inverseView[10]);

    vec3.normalize(forwardVector, forwardVector);
    vec3.scale(forwardVector, forwardVector, (state.deltaTime * moveSpeed));

    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);
    cameraCenterVector = vec3.fromValues(state.camera.center[0], state.camera.center[1], state.camera.center[2]);

    vec3.add(cameraPositionVector, cameraPositionVector, forwardVector);
    vec3.add(cameraCenterVector, cameraPositionVector, forwardVector);

    state.camera.position = [cameraPositionVector[0], cameraPositionVector[1], cameraPositionVector[2]];
    state.camera.center = [cameraCenterVector[0], cameraCenterVector[1], cameraCenterVector[2]];
}

function moveBackward(state) {
    let inverseView = mat4.create(), forwardVector = vec3.fromValues(0, 0, 0), cameraCenterVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    mat4.invert(inverseView, state.viewMatrix);
    //forward vector from the viewmatrix
    forwardVector = vec3.fromValues(-inverseView[2], -inverseView[6], inverseView[10]);
    vec3.normalize(forwardVector, forwardVector);

    vec3.normalize(forwardVector, forwardVector);
    vec3.scale(forwardVector, forwardVector, (state.deltaTime * moveSpeed));

    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);
    cameraCenterVector = vec3.fromValues(state.camera.center[0], state.camera.center[1], state.camera.center[2]);

    vec3.add(cameraPositionVector, cameraPositionVector, forwardVector);

    state.camera.position = [cameraPositionVector[0], cameraPositionVector[1], cameraPositionVector[2]];
    state.camera.center = [cameraCenterVector[0], cameraCenterVector[1], cameraCenterVector[2]];
}

function moveLeft(state) {
    let forwardVector = vec3.fromValues(0, 0, 0), sidewaysVector = vec3.fromValues(0, 0, 0), cameraCenterVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    sidewaysVector = vec3.fromValues(0, 0, 0);
    forwardVector = vec3.fromValues(state.viewMatrix[2], state.viewMatrix[6], state.viewMatrix[10]);
    vec3.cross(sidewaysVector, forwardVector, state.camera.up);
    vec3.normalize(sidewaysVector, sidewaysVector);
    vec3.scale(sidewaysVector, sidewaysVector, (state.deltaTime * moveSpeed));

    cameraCenterVector = vec3.fromValues(state.camera.center[0], state.camera.center[1], state.camera.center[2]);
    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);

    vec3.add(cameraCenterVector, cameraCenterVector, sidewaysVector);
    vec3.add(cameraPositionVector, cameraPositionVector, sidewaysVector);

    state.camera.center = [cameraCenterVector[0], cameraCenterVector[1], cameraCenterVector[2]];
    state.camera.position = [cameraPositionVector[0], cameraPositionVector[1], cameraPositionVector[2]];
}

function moveRight(state) {
    let forwardVector = vec3.fromValues(0, 0, 0), sidewaysVector = vec3.fromValues(0, 0, 0), cameraCenterVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    sidewaysVector = vec3.fromValues(0, 0, 0);
    forwardVector = vec3.fromValues(-state.viewMatrix[2], -state.viewMatrix[6], -state.viewMatrix[10]);
    vec3.cross(sidewaysVector, forwardVector, state.camera.up);
    vec3.normalize(sidewaysVector, sidewaysVector);
    vec3.scale(sidewaysVector, sidewaysVector, (state.deltaTime * moveSpeed));

    cameraCenterVector = vec3.fromValues(state.camera.center[0], state.camera.center[1], state.camera.center[2]);
    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);

    vec3.add(cameraCenterVector, cameraCenterVector, sidewaysVector);
    vec3.add(cameraPositionVector, cameraPositionVector, sidewaysVector);

    state.camera.center = [cameraCenterVector[0], cameraCenterVector[1], cameraCenterVector[2]];
    state.camera.position = [cameraPositionVector[0], cameraPositionVector[1], cameraPositionVector[2]];
}