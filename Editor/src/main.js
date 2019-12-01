var total = 0;
var state = {};
var stats = new Stats();

window.onload = () => {
    var sceneFile = "testsave.json";
    parseSceneFile("./statefiles/" + sceneFile, state, main);
    //console.log(shaders);
    getUniformsFromShaders(shaders);
}

/**
 * 
 * @param {object - contains vertex, normal, uv information for the mesh to be made} mesh 
 * @param {object - the game object that will use the mesh information} object 
 * @purpose - Helper function called as a callback function when the mesh is done loading for the object
 */
function createMesh(mesh, object) {
    if (object.type === "mesh") {
        let testModel = new Model(state.gl, object, mesh);
        testModel.setup();
        testModel.model.position = object.position;
        if (object.scale) {
            testModel.scale(object.scale);
        }
        addObjectToScene(state, testModel);
    } else {
        let testLight = new Light(state.gl, object, mesh);
        testLight.setup();
        testLight.model.position = object.position;
        if (object.scale) {
            testLight.scale(object.scale);
        }

        addObjectToScene(state, testLight);
    }
}

/**
 * 
 * @param {string - type of object to be added to the scene} type 
 * @param {string - url of the model being added to the game} url 
 * @purpose **WIP** Adds a new object to the scene from using the gui to add said object 
 */
function addObject(type, url = null) {
    if (type === "Cube") {
        let testCube = new Cube(state.gl, "Cube", null, [0.1, 0.1, 0.1], [0.0, 0.0, 0.0], [0.0, 0.0, 0.0], 10, 1.0);
        testCube.vertShader = state.vertShaderSample;
        testCube.fragShader = state.fragShaderSample;
        testCube.setup();

        addObjectToScene(state, testCube);
        createSceneGui(state);
    }
}

function main() {
    stats.showPanel(0);
    document.getElementById("fps").appendChild(stats.dom);
    //document.body.appendChild( stats.dom );
    const canvas = document.querySelector("#glCanvas");

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
        keyboard: {},
        mouse: { sensitivity: 0.2 },
        gameStarted: false,
        camera: {
            name: 'camera',
            position: vec3.fromValues(0.5, 0.0, -2.5),
            center: vec3.fromValues(0.5, 0.0, 0.0),
            up: vec3.fromValues(0.0, 1.0, 0.0),
            pitch: 0,
            yaw: 0,
            roll: 0
        },
        samplerExists: 0,
        samplerNormExists: 0
    };

    state.numLights = state.lights.length;

    //iterate through the level's objects and add them
    state.level.objects.map((object) => {
        if (object.type === "mesh" || object.type === "light") {
            parseOBJFileToJSON(object.model, createMesh, object);
        } else if (object.type === "cube") {
            let tempCube = new Cube(gl, object);
            tempCube.setup();
            //tempCube.model.position = vec3.fromValues(object.position[0], object.position[1], object.position[2]);

            if (object.scale) {
                tempCube.scale(object.scale);
            }
            tempCube.translate(object.position);
            addObjectToScene(state, tempCube);
        } else if (object.type === "plane") {
            let tempPlane = new Plane(gl, object);
            tempPlane.setup();

            tempPlane.model.position = vec3.fromValues(object.position[0], object.position[1], object.position[2]);
            if (object.scale) {
                tempPlane.scale(object.scale);
            }
            addObjectToScene(state, tempPlane);
        }
    })

    state.level.lights.map((light) => {
        parseOBJFileToJSON(light.model, createMesh, light);
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
    //console.log(object);
    if (object.type === "light") {
        state.lightIndices.push(state.objectCount);
        state.numLights++;
    }

    object.name = object.name;
    state.objects.push(object);
    state.objectTable[object.name] = state.objectCount;
    state.objectCount++;
    createSceneGui(state);
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
        stats.begin();
        now *= 0.001; // convert to seconds
        const deltaTime = now - then;
        then = now;

        state.deltaTime = deltaTime;

        //wait until the scene is completely loaded to render it
        if (state.numberOfObjectsToLoad <= state.objects.length) {
            if (!state.gameStarted) {
                console.log(state.objects)
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
                vec3.rotateY(state.camera.center, state.camera.center, state.camera.position, (-state.mouse.rateX * deltaTime * state.mouse.sensitivity));
            }
            
            // Draw our scene
            drawScene(gl, deltaTime, state);
        }
        stats.end();
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
    gl.clearColor(0.0, 0.0, 0.0, 1.0);
    //gl.clearColor(1.0, 1.0, 1.0, 1.0);
    gl.enable(gl.DEPTH_TEST); // Enable depth testing
    gl.depthFunc(gl.LEQUAL); // Near things obscure far things
    gl.cullFace(gl.BACK);
    gl.enable(gl.CULL_FACE);
    gl.clearDepth(1.0); // Clear everything
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    let lightPositionArray = [], lightColourArray = [], lightStrengthArray = [];

    for (let i = 0; i < state.lightIndices.length; i++) {
        let light = state.objects[state.lightIndices[i]];
        for (let j = 0; j < 3; j++) {
            lightPositionArray.push(light.model.position[j]);
            lightColourArray.push(light.colour[j]);
        }
        lightStrengthArray.push(light.strength);
    }

    state.objects.map((object) => {
        if (object.loaded) {

            gl.useProgram(object.programInfo.program);
            {

                var projectionMatrix = mat4.create();
                var fovy = 60.0 * Math.PI / 180.0; // Vertical field of view in radians
                var aspect = state.canvas.clientWidth / state.canvas.clientHeight; // Aspect ratio of the canvas
                var near = 0.1; // Near clipping plane
                var far = 100.0; // Far clipping plane

                mat4.perspective(projectionMatrix, fovy, aspect, near, far);

                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uProjectionMatrix, false, projectionMatrix);

                state.projectionMatrix = projectionMatrix;

                var viewMatrix = mat4.create();
                mat4.lookAt(
                    viewMatrix,
                    state.camera.position,
                    state.camera.center,
                    state.camera.up,
                );
                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uViewMatrix, false, viewMatrix);

                gl.uniform3fv(object.programInfo.uniformLocations.uCameraPosition, state.camera.position);

                state.viewMatrix = viewMatrix;

                var modelMatrix = mat4.create();
                var negCentroid = vec3.fromValues(0.0, 0.0, 0.0);
                vec3.negate(negCentroid, object.centroid);

                mat4.translate(modelMatrix, modelMatrix, object.model.position);
                mat4.translate(modelMatrix, modelMatrix, object.centroid);
                mat4.mul(modelMatrix, modelMatrix, object.model.rotation);
                mat4.translate(modelMatrix, modelMatrix, negCentroid);
                mat4.scale(modelMatrix, modelMatrix, object.model.scale);

                object.modelMatrix = modelMatrix;

                var normalMatrix = mat4.create();
                mat4.invert(normalMatrix, modelMatrix);
                mat4.transpose(normalMatrix, normalMatrix);

                gl.uniformMatrix4fv(object.programInfo.uniformLocations.uModelMatrix, false, modelMatrix);
                gl.uniformMatrix4fv(object.programInfo.uniformLocations.normalMatrix, false, normalMatrix);

                gl.uniform3fv(object.programInfo.uniformLocations.diffuseVal, object.material.diffuse);
                gl.uniform3fv(object.programInfo.uniformLocations.ambientVal, object.material.ambient);
                gl.uniform3fv(object.programInfo.uniformLocations.specularVal, object.material.specular);
                gl.uniform1f(object.programInfo.uniformLocations.nVal, object.material.n);

                gl.uniform1i(object.programInfo.uniformLocations.numLights, state.numLights);

                //use this check to wait until the light meshes are loaded properly
                if (lightColourArray.length > 0 && lightPositionArray.length > 0 && lightStrengthArray.length > 0) {
                    gl.uniform3fv(object.programInfo.uniformLocations.uLightPositions, lightPositionArray);
                    gl.uniform3fv(object.programInfo.uniformLocations.uLightColours, lightColourArray);
                    gl.uniform1fv(object.programInfo.uniformLocations.uLightStrengths, lightStrengthArray);
                }

                {
                    // Bind the buffer we want to draw
                    gl.bindVertexArray(object.buffers.vao);

                    //check for diffuse texture and apply it
                    if (object.model.texture != null) {
                        state.samplerExists = 1;
                        gl.activeTexture(gl.TEXTURE0);
                        gl.uniform1i(object.programInfo.uniformLocations.samplerExists, state.samplerExists);
                        gl.uniform1i(object.programInfo.uniformLocations.uTexture, 0);
                        gl.bindTexture(gl.TEXTURE_2D, object.model.texture);

                    } else {
                        gl.activeTexture(gl.TEXTURE0);
                        state.samplerExists = 0;
                        gl.uniform1i(object.programInfo.uniformLocations.samplerExists, state.samplerExists);
                    }

                    //check for normal texture and apply it
                    if (object.model.textureNorm != null) {
                        state.samplerNormExists = 1;
                        gl.activeTexture(gl.TEXTURE1);
                        gl.uniform1i(object.programInfo.uniformLocations.uTextureNormExists, state.samplerNormExists);
                        gl.uniform1i(object.programInfo.uniformLocations.uTextureNorm, 1);
                        gl.bindTexture(gl.TEXTURE_2D, object.model.textureNorm);
                        //console.log("here")
                    } else {
                        gl.activeTexture(gl.TEXTURE1);
                        state.samplerNormExists = 0;
                        gl.uniform1i(object.programInfo.uniformLocations.uTextureNormExists, state.samplerNormExists);
                    }

                    // Draw the object
                    const offset = 0; // Number of elements to skip before starting

                    //if its a mesh then we don't use an index buffer and use drawArrays instead of drawElements
                    if (object.type === "mesh" || object.type === "light") {
                        gl.drawArrays(gl.TRIANGLES, offset, object.buffers.numVertices / 3);
                    } else {
                        gl.drawElements(gl.TRIANGLES, object.buffers.numVertices, gl.UNSIGNED_SHORT, offset);
                    }
                }
            }
        }
    });
}