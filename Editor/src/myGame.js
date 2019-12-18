var moveSpeed = 5;

function startGame(state) {
    let canvas = document.getElementById('glCanvas');

    canvas.addEventListener("contextmenu", function (e) {
        e.preventDefault();
    }, false);

    canvas.addEventListener('mousemove', (event) => {
        //handle right click
        if (event.buttons == 2) {
            state.mouse['camMove'] = true;
            state.mouse.rateX = event.movementX;
            state.mouse.rateY = -event.movementY;
        }
    });

    canvas.addEventListener('mouseup', (event) => {
        state.mouse['camMove'] = false;
        state.mouse.rateX = 0;
        state.mouse.rateY = 0;
    })

    canvas.addEventListener('mouseout', (event) => {
        state.mouse['camMove'] = false;
        state.mouse.rateX = 0;
        state.mouse.rateY = 0;
    })

    canvas.addEventListener('keypress', (event) => {
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

    canvas.addEventListener('keyup', (event) => {
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

            case "ArrowRight":
                moveTestCubeTestCollision(state, "right");
                break;

            case "ArrowLeft":
                moveTestCubeTestCollision(state, "left");
                break;

            case "ArrowUp":
                moveTestCubeTestCollision(state, "forward");
                break;

            case "ArrowDown":
                moveTestCubeTestCollision(state, "backward");
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
    let inverseView = mat4.create(), forwardVector = vec3.fromValues(0, 0, 0), camerafrontVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    mat4.invert(inverseView, state.viewMatrix);
    //forward vector from the viewmatrix
    forwardVector = vec3.fromValues(inverseView[2], -inverseView[6], -inverseView[10]);

    vec3.normalize(forwardVector, forwardVector);
    vec3.scale(forwardVector, forwardVector, (state.deltaTime * moveSpeed));

    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);
    camerafrontVector = vec3.fromValues(state.camera.front[0], state.camera.front[1], state.camera.front[2]);

    vec3.add(cameraPositionVector, cameraPositionVector, forwardVector);
    vec3.add(camerafrontVector, cameraPositionVector, forwardVector);

    state.camera.position = [cameraPositionVector[0], state.camera.position[1], cameraPositionVector[2]];
}

function moveBackward(state) {
    let inverseView = mat4.create(), forwardVector = vec3.fromValues(0, 0, 0), camerafrontVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    mat4.invert(inverseView, state.viewMatrix);
    //forward vector from the viewmatrix
    forwardVector = vec3.fromValues(-inverseView[2], inverseView[6], inverseView[10]);
    vec3.normalize(forwardVector, forwardVector);

    vec3.normalize(forwardVector, forwardVector);
    vec3.scale(forwardVector, forwardVector, (state.deltaTime * moveSpeed));

    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);
    camerafrontVector = vec3.fromValues(state.camera.front[0], state.camera.front[1], state.camera.front[2]);

    vec3.add(cameraPositionVector, cameraPositionVector, forwardVector);

    state.camera.position = [cameraPositionVector[0], state.camera.position[1], cameraPositionVector[2]];
}

function moveLeft(state) {
    let forwardVector = vec3.fromValues(0, 0, 0), sidewaysVector = vec3.fromValues(0, 0, 0), camerafrontVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    sidewaysVector = vec3.fromValues(0, 0, 0);
    forwardVector = vec3.fromValues(state.viewMatrix[2], state.viewMatrix[6], state.viewMatrix[10]);
    vec3.cross(sidewaysVector, forwardVector, state.camera.up);
    vec3.normalize(sidewaysVector, sidewaysVector);
    vec3.scale(sidewaysVector, sidewaysVector, (state.deltaTime * moveSpeed));

    camerafrontVector = vec3.fromValues(state.camera.front[0], state.camera.front[1], state.camera.front[2]);
    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);

    vec3.add(camerafrontVector, camerafrontVector, sidewaysVector);
    vec3.add(cameraPositionVector, cameraPositionVector, sidewaysVector);

    //state.camera.front = [camerafrontVector[0], camerafrontVector[1], camerafrontVector[2]];
    state.camera.position = [cameraPositionVector[0], state.camera.position[1], cameraPositionVector[2]];
}

function moveRight(state) {
    let forwardVector = vec3.fromValues(0, 0, 0), sidewaysVector = vec3.fromValues(0, 0, 0), camerafrontVector = vec3.fromValues(0, 0, 0), cameraPositionVector = vec3.fromValues(0, 0, 0);

    sidewaysVector = vec3.fromValues(0, 0, 0);
    forwardVector = vec3.fromValues(-state.viewMatrix[2], -state.viewMatrix[6], -state.viewMatrix[10]);
    vec3.cross(sidewaysVector, forwardVector, state.camera.up);
    vec3.normalize(sidewaysVector, sidewaysVector);
    vec3.scale(sidewaysVector, sidewaysVector, (state.deltaTime * moveSpeed));

    camerafrontVector = vec3.fromValues(state.camera.front[0], state.camera.front[1], state.camera.front[2]);
    cameraPositionVector = vec3.fromValues(state.camera.position[0], state.camera.position[1], state.camera.position[2]);

    vec3.add(camerafrontVector, camerafrontVector, sidewaysVector);
    vec3.add(cameraPositionVector, cameraPositionVector, sidewaysVector);

    //state.camera.front = [camerafrontVector[0], camerafrontVector[1], camerafrontVector[2]];
    state.camera.position = [cameraPositionVector[0], state.camera.position[1], cameraPositionVector[2]];
}

function moveTestCubeTestCollision(state, direction) {
    let testCube = getObject(state, "testCube0");

    if (direction === "left") {
        testCube.translate([0.25, 0, 0]);
    } else if (direction === "right") {
        testCube.translate([-0.25, 0, 0]);
    } else if (direction === "forward") {
        testCube.translate([0, 0, 0.25]);
    } else if (direction === "backward") {
        testCube.translate([0, 0, -0.25]);
    }

    let collide = false;
    state.objects.map((obj) => {
        if (obj.name !== "testCube0") {

            if (obj.type === "mesh" && obj.parent == null) {
                //console.error(obj.boundingBox)
                collide = intersect(testCube.boundingBox, obj.boundingBox);
                if (collide) {
                    console.warn("Collided with", obj.name);
                }
            } else {
                collide = intersect(testCube.boundingBox, obj.boundingBox);
                if (collide) {
                    console.warn("Collided with", obj.name);
                }
            }
        }
    })
}