import myMath from "./mymath/index.js";
import UI from "./uiSetup.js";
import { Cube, PointLight, Plane, Model, DirectionalLight } from "./objects/index.js";

var currentlyRendered = 0;
var state = {
    saveFile: "testsave.json"
};

if (window.location.pathname.indexOf("main.html") !== -1) {
    window.onload = () => {
        const { ipcRenderer } = require('electron')
        ipcRenderer.on('sceneOpen', (event, arg) => {
            if (arg.data) {
                console.log(arg.data.filePaths[0]);
            }
        })

        state.width = window.innerWidth;
        state.height = window.innerHeight;
        state.renderedText = document.getElementById("renderedNumText");
        state.renderedText.style.color = "white";
        state.loadingBar = {
            parent: document.getElementById("loadingBar"),
            child: document.getElementById("loadingBarProgress")
        }

        //add event listener to canvas resize
        window.addEventListener("resize", (e) => {
            state.width = window.innerWidth;
            state.height = window.innerHeight;
            state.gl.viewport(0, 0, state.width, state.height);
        })

        parseSceneFile("./statefiles/" + state.saveFile, state, main);
    }
}

/**
 * 
 * @param {string - type of object to be added to the scene} type 
 * @param {string - url of the model being added to the game} url 
 * @purpose **WIP** Adds a new object to the scene from using the gui to add said object //move to helpers
 */
function addObject(type, url = null) {
    if (type === "Cube") {
        let testCube = new Cube(state.gl, "Cube", null, [0.1, 0.1, 0.1], [0.0, 0.0, 0.0], [0.0, 0.0, 0.0], 10, 1.0);
        testCube.vertShader = state.vertShaderSample;
        testCube.fragShader = state.fragShaderSample;
        testCube.setup();
        addObjectToScene(state, testCube);
        UI.createSceneGui(state);
    }
}

function createModalFromMesh(mesh, object) {
    if (object.type === "mesh") {
        let tempMesh = new Model(state.gl, object, mesh);
        let testVal = tempMesh.setup();
        testVal.then((val) => {
            addObjectToScene(state, val);
        })
    } else {
        let tempLight = new Light(state.gl, object, mesh);
        tempLight.setup();
        addObjectToScene(state, tempLight);
    }
}

function main() {
    const canvas = document.querySelector("#glCanvas");
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    // Initialize the WebGL2 context
    var gl = canvas.getContext("webgl2");

    // Only continue if WebGL2 is available and working
    if (gl === null) {
        printError('WebGL 2 not supported by your browser',
            'Check to see you are using a <a href="https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API#WebGL_2_2" class="alert-link">modern browser</a>.');
        return;
    }

    state = {
        ...state,
        gl,
        canvas: canvas,
        objectCount: 0,
        objectTable: {},
        lightIndices: [],
        initialRender: false,
        keyboard: {},
        mouse: { sensitivity: 0.2 },
        gameStarted: false,
        camera: {
            name: 'camera',
            position: vec3.fromValues(6, 6, 0),
            front: vec3.fromValues(0.0, 0.0, 1.0),
            up: vec3.fromValues(0.0, 1.0, 0.0),
            pitch: 0,
            yaw: 90, //works when the camera front is 0, 0, 1 to start
            roll: 0
        },
        samplerExists: 0,
        samplerNormExists: 0
    };

    //iterate through the level's objects and add them
    state.level.objects.map((object) => {
        if (object.type === "mesh") {
            parseOBJFileToJSON("./models/" + object.model, object, createModalFromMesh);
        } else if (object.type === "cube") {
            let tempCube = new Cube(gl, object);
            tempCube.setup();
            addObjectToScene(state, tempCube);
        } else if (object.type === "plane") {
            let tempPlane = new Plane(gl, object);
            tempPlane.setup();
            addObjectToScene(state, tempPlane);
        }
    })

    state.level.pointLights.forEach((light) => {
        let tempPointLight = new PointLight(gl, light);
        state.pointLights.push(tempPointLight);
        UI.createSceneGui(state);
    })

    state.level.directionalLights.forEach((light) => {
        let tempDirLight = new DirectionalLight(gl, light);
        state.directionalLights.push(tempDirLight);
        UI.createSceneGui(state);
    })

    //setup mouse click listener
    /*
    canvas.addEventListener('click', (event) => {
        getMousePick(event, state);
    }) */

    startRendering(gl, state);
}

/**
 * 
 * @param {object - object containing scene values} state 
 * @param {object - the object to be added to the scene} object 
 * @purpose - Helper function for adding a new object to the scene and refreshing the GUI
 */
function addObjectToScene(state, object) {
    if (object.type === "light") {
        state.lightIndices.push(state.objectCount);
        state.numLights++;
    }
    //check if its a child to a mesh, if so we need to increase the amount of objects we are waiting to load
    if (object.type === "mesh" && object.parent) {
        state.numberOfObjectsToLoad++;
    }

    object.name = object.name;
    state.objects.push(object);
    state.objectTable[object.name] = state.objectCount;
    state.objectCount++;
    UI.createSceneGui(state);
}

/**
 * 
 * @param {gl context} gl 
 * @param {object - object containing scene values} state 
 * @purpose - Calls the drawscene per frame
 */
function startRendering(gl, state) {
    // A variable for keeping track of time between frames
    var then = 0.0;

    // This function is called when we want to render a frame to the canvas
    function render(now) {
        now *= 0.001; // convert to seconds
        const deltaTime = now - then;
        then = now;

        state.deltaTime = deltaTime;

        //wait until the scene is completely loaded to render it


        if (state.numberOfObjectsToLoad <= state.objects.length) {
            if (!state.initialRender) {
                drawScene(gl, deltaTime, state);
                state.initialRender = true;
            }

            if (!state.gameStarted) {
                startGame(state);
                state.gameStarted = true;
            }

            if (state.keyboard["w"]) {
                moveForward(state);
            }
            if (state.keyboard["s"]) {
                moveBackward(state);
            }
            if (state.keyboard["a"]) {
                moveLeft(state);
            }
            if (state.keyboard["d"]) {
                moveRight(state);
            }

            if (state.mouse['camMove']) {
                let front = vec3.fromValues(0, 0, 0);
                state.camera.yaw += state.mouse.rateX * state.mouse.sensitivity;
                state.camera.pitch += state.mouse.rateY * state.mouse.sensitivity;

                if (state.camera.pitch > 89) {
                    state.camera.pitch = 89
                }
                if (state.camera.pitch < -89) {
                    state.camera.pitch = -89
                }

                front[0] = Math.cos(toRadians(state.camera.yaw)) * Math.cos(toRadians(state.camera.pitch));
                front[1] = Math.sin(toRadians(state.camera.pitch));
                front[2] = Math.sin(toRadians(state.camera.yaw)) * Math.cos(toRadians(state.camera.pitch));

                vec3.normalize(state.camera.front, front);

                //vec3.rotateY(state.camera.front, state.camera.front, state.camera.position, (-state.mouse.rateX * state.mouse.sensitivity));
            }

            let keyMove = Object.values(state.keyboard).includes(true);
            // Draw our scene
            if (state.mouse.camMove || keyMove || state.render) {
                drawScene(gl, deltaTime, state);
            }

            state.renderedText.innerHTML = "Rendered: " + currentlyRendered;
            currentlyRendered = 0;
        } else {
            drawScene(gl, deltaTime, state);
        }
        // Request another frame when this one is done
        requestAnimationFrame(render);
    }

    // Draw the scene
    requestAnimationFrame(render);
}

/**
 * 
 * @param {gl context} gl 
 * @param {float - time from now-last} deltaTime 
 * @param {object - contains the state for the scene} state 
 * @purpose Iterate through game objects and render the objects aswell as update uniforms
 */
function drawScene(gl, deltaTime, state) {
    console.log("rendering!")
    gl.enable(gl.DEPTH_TEST); // Enable depth testing
    gl.depthFunc(gl.LEQUAL); // Near things obscure far things
    gl.enable(gl.CULL_FACE);
    gl.bindFramebuffer(gl.FRAMEBUFFER, null);
    gl.viewport(0, 0, state.width, state.height);
    gl.clearColor(0.0, 0.0, 0.0, 1.0);
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    var sortedObjects = state.objects.slice().sort((a, b) => {
        return vec3.distance(state.camera.position, a.model.position) >= vec3.distance(state.camera.position, b.model.position) ? -1 : 1;
    });

    sortedObjects.forEach((object) => {
        if (object.loaded) {
            gl.useProgram(object.programInfo.program);
            {

                if (object.material.alpha < 1.0) {
                    gl.enable(gl.BLEND);
                    gl.disable(gl.DEPTH_TEST);
                    gl.blendFunc(gl.ONE_MINUS_CONSTANT_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
                    gl.clearDepth(object.material.alpha);
                } else {
                    gl.disable(gl.BLEND);
                    gl.depthMask(true);
                    gl.enable(gl.DEPTH_TEST);
                    gl.depthFunc(gl.LEQUAL); // Near things obscure far things
                    gl.clearDepth(1.0);
                }

                gl.activeTexture(gl.TEXTURE0);
                gl.uniform1i(object.programInfo.uniformLocations.projectedTexture, 0);
                gl.bindTexture(gl.TEXTURE_2D, state.depthTexture);

                var projectionMatrix = mat4.create();
                var fovy = 60.0 * Math.PI / 180.0; // Vertical field of view in radians
                var aspect = state.canvas.clientWidth / state.canvas.clientHeight; // Aspect ratio of the canvas
                var near = 0.1; // Near clipping plane
                var far = 100.0; // Far clipping plane
                gl.uniform1f(object.programInfo.uniformLocations.near_plane, near);
                gl.uniform1f(object.programInfo.uniformLocations.far_plane, far);

                mat4.perspective(projectionMatrix, fovy, aspect, near, far);
                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uProjectionMatrix, false, projectionMatrix);

                state.projectionMatrix = projectionMatrix;

                var viewMatrix = mat4.create();
                //create camera front value
                let camFront = vec3.fromValues(0, 0, 0);
                vec3.add(camFront, state.camera.position, state.camera.front);
                mat4.lookAt(
                    viewMatrix,
                    state.camera.position,
                    camFront,
                    state.camera.up,
                );

                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uViewMatrix, false, viewMatrix);
                gl.uniform3fv(object.programInfo.uniformLocations.uCameraPosition, state.camera.position);
                state.viewMatrix = viewMatrix;

                //perform frustum culling
                let frustum = new myMath.Frustum(projectionMatrix, viewMatrix);
                if (!frustum.sphereIntersection(object.model.position,
                    Math.pow(vec3.len(vec3.fromValues(object.boundingBox.xMax, object.boundingBox.yMax, object.boundingBox.zMax)), 2))) { //use the squared len value for greater tolerance
                    return;
                }

                currentlyRendered++;

                //TODO Centroid rotation of scaled objects not calculated properly
                //Apply transformations to model matrix
                var modelMatrix = mat4.create();
                var negCentroid = vec3.fromValues(0.0, 0.0, 0.0);
                vec3.negate(negCentroid, object.centroid);
                mat4.translate(modelMatrix, modelMatrix, object.model.position);
                mat4.translate(modelMatrix, modelMatrix, object.centroid);
                mat4.mul(modelMatrix, modelMatrix, object.model.rotation);
                mat4.translate(modelMatrix, modelMatrix, negCentroid);
                mat4.scale(modelMatrix, modelMatrix, object.model.scale);

                if (object.parent) {
                    let parent = getObject(state, object.parent);
                    if (parent.modelMatrix) {
                        mat4.mul(modelMatrix, parent.modelMatrix, modelMatrix);
                    }
                }

                object.modelMatrix = modelMatrix;
                var normalMatrix = mat4.create();
                mat4.invert(normalMatrix, modelMatrix);
                mat4.transpose(normalMatrix, normalMatrix);

                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uModelMatrix, false, modelMatrix);
                gl.uniformMatrix4fv(object.programInfo.uniformLocations.normalMatrix, false, normalMatrix);
                gl.uniform3fv(object.programInfo.uniformLocations.diffuseVal, object.material.diffuse);
                gl.uniform3fv(object.programInfo.uniformLocations.ambientVal, object.material.ambient);
                gl.uniform3fv(object.programInfo.uniformLocations.specularVal, object.material.specular);
                gl.uniform1f(object.programInfo.uniformLocations.alpha, object.material.alpha);
                gl.uniform1f(object.programInfo.uniformLocations.nVal, object.material.n);
                gl.uniform1i(object.programInfo.uniformLocations.numPointLights, state.pointLights.length);
                gl.uniform1i(object.programInfo.uniformLocations.numDirLights, state.directionalLights.length);

                state.pointLights.forEach((pL, index) => {
                    gl.uniform3fv(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].position"), pL.position);
                    gl.uniform3fv(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].color"), pL.colour);
                    gl.uniform1f(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].strength"), pL.strength);
                    gl.uniform1f(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].constant"), pL.constant);
                    gl.uniform1f(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].linear"), pL.linear);
                    gl.uniform1f(gl.getUniformLocation(object.programInfo.program, "pointLights[" + index + "].quadratic"), pL.quadratic);
                })

                state.directionalLights.forEach((dL, index) => {
                    gl.uniform3fv(gl.getUniformLocation(object.programInfo.program, "directionalLights[" + index + "].position"), dL.position);
                    gl.uniform3fv(gl.getUniformLocation(object.programInfo.program, "directionalLights[" + index + "].color"), dL.colour);
                    gl.uniform3fv(gl.getUniformLocation(object.programInfo.program, "directionalLights[" + index + "].direction"), dL.direction);
                })

                {
                    // Bind the buffer we want to draw
                    gl.bindVertexArray(object.buffers.vao);

                    // check for diffuse texture and apply it
                    if (object.model.texture != null) {
                        gl.activeTexture(gl.TEXTURE0 + 1);
                        gl.uniform1i(object.programInfo.uniformLocations.uTexture, 1);
                        gl.bindTexture(gl.TEXTURE_2D, object.model.texture);
                    }

                    //check for normal texture and apply it
                    if (object.model.textureNorm != null) {
                        gl.activeTexture(gl.TEXTURE0 + 2);
                        gl.uniform1i(object.programInfo.uniformLocations.uTextureNorm, 2);
                        gl.bindTexture(gl.TEXTURE_2D, object.model.textureNorm);
                    }

                    // Draw the object
                    const offset = 0; // Number of elements to skip before starting

                    //if its a mesh then we don't use an index buffer and use drawArrays instead of drawElements
                    if (object.type === "mesh") {
                        gl.drawArrays(gl.TRIANGLES, offset, object.buffers.numVertices / 3);
                    } else {
                        gl.drawElements(gl.TRIANGLES, object.buffers.numVertices, gl.UNSIGNED_SHORT, offset);
                    }

                    gl.bindTexture(gl.TEXTURE_2D, null);
                }
            }
        }
    });
    state.render = false;
}

export default state;