var fs = require('fs');

/**
 * @param  {} gl WebGL2 Context
 * @param  {string} vsSource Vertex shader GLSL source code
 * @param  {string} fsSource Fragment shader GLSL source code
 * @returns {} A shader program object. This is `null` on failure
 */
function initShaderProgram(gl, vsSource, fsSource) {
    // Use our custom function to load and compile the shader objects
    const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource);
    const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource);

    // Create the shader program by attaching and linking the shader objects
    const shaderProgram = gl.createProgram();
    gl.attachShader(shaderProgram, vertexShader);
    gl.attachShader(shaderProgram, fragmentShader);
    gl.linkProgram(shaderProgram);

    // If creating the shader program failed, alert

    if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS)) {
        alert('Unable to link the shader program' + gl.getProgramInfoLog(shaderProgram));
        return null;
    }

    return shaderProgram;
}

/**
 * Loads a shader from source into a shader object. This should later be linked into a program.
 * @param  {} gl WebGL2 context
 * @param  {} type Type of shader. Typically either VERTEX_SHADER or FRAGMENT_SHADER
 * @param  {string} source GLSL source code
 */
function loadShader(gl, type, source) {
    // Create a new shader object
    const shader = gl.createShader(type);

    // Send the source to the shader object
    gl.shaderSource(shader, source);

    // Compile the shader program
    gl.compileShader(shader);

    // See if it compiled successfully
    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
        // Fail with an error message
        var typeStr = '';
        if (type === gl.VERTEX_SHADER) {
            typeStr = 'VERTEX';
        } else if (type === gl.FRAGMENT_SHADER) {
            typeStr = 'FRAGMENT';
        }
        console.error('An error occurred compiling the shader: ' + typeStr, gl.getShaderInfoLog(shader));
        gl.deleteShader(shader);
        return null;
    }

    return shader;
}

/**
 * 
 * @param {array of x,y,z vertices} vertices 
 */
function calculateCentroid(vertices, cb) {
    //console.log(vertices);

    var center = vec3.fromValues(0.0, 0.0, 0.0);
    for (let t = 0; t < vertices.length; t += 3) {
        vec3.add(center, center, vec3.fromValues(vertices[t], vertices[t + 1], vertices[t + 2]));
    }
    vec3.scale(center, center, 1 / (vertices.length / 3));

    if (cb) {
        cb();

        return center;
    } else {
        return center;
    }
}

function initPositionAttribute(gl, programInfo, positionArray) {

    // Create a buffer for the positions.
    const positionBuffer = gl.createBuffer();

    // Select the buffer as the one to apply buffer
    // operations to from here out.
    gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);

    // Now pass the list of positions into WebGL to build the
    // shape. We do this by creating a Float32Array from the
    // JavaScript array, then use it to fill the current buffer.
    gl.bufferData(
        gl.ARRAY_BUFFER, // The kind of buffer this is
        positionArray, // The data in an Array object
        gl.STATIC_DRAW // We are not going to change this data, so it is static
    );

    // Tell WebGL how to pull out the positions from the position
    // buffer into the vertexPosition attribute.
    {
        const numComponents = 3; // pull out 3 values per iteration, ie vec3
        const type = gl.FLOAT; // the data in the buffer is 32bit floats
        const normalize = false; // don't normalize between 0 and 1
        const stride = 0; // how many bytes to get from one set of values to the next
        // Set stride to 0 to use type and numComponents above
        const offset = 0; // how many bytes inside the buffer to start from


        // Set the information WebGL needs to read the buffer properly
        gl.vertexAttribPointer(
            programInfo.attribLocations.vertexPosition,
            numComponents,
            type,
            normalize,
            stride,
            offset
        );
        // Tell WebGL to use this attribute
        gl.enableVertexAttribArray(
            programInfo.attribLocations.vertexPosition);
    }

    return positionBuffer;
}

function initNormalAttribute(gl, programInfo, normalArray) {

    // Create a buffer for the positions.
    const normalBuffer = gl.createBuffer();

    // Select the buffer as the one to apply buffer
    // operations to from here out.
    gl.bindBuffer(gl.ARRAY_BUFFER, normalBuffer);

    // Now pass the list of positions into WebGL to build the
    // shape. We do this by creating a Float32Array from the
    // JavaScript array, then use it to fill the current buffer.
    gl.bufferData(
        gl.ARRAY_BUFFER, // The kind of buffer this is
        normalArray, // The data in an Array object
        gl.STATIC_DRAW // We are not going to change this data, so it is static
    );

    // Tell WebGL how to pull out the positions from the position
    // buffer into the vertexPosition attribute.
    {
        const numComponents = 3; // pull out 4 values per iteration, ie vec3
        const type = gl.FLOAT; // the data in the buffer is 32bit floats
        const normalize = false; // don't normalize between 0 and 1
        const stride = 0; // how many bytes to get from one set of values to the next
        // Set stride to 0 to use type and numComponents above
        const offset = 0; // how many bytes inside the buffer to start from

        // Set the information WebGL needs to read the buffer properly
        gl.vertexAttribPointer(
            programInfo.attribLocations.vertexNormal,
            numComponents,
            type,
            normalize,
            stride,
            offset
        );
        // Tell WebGL to use this attribute
        gl.enableVertexAttribArray(
            programInfo.attribLocations.vertexNormal);
    }

    return normalBuffer;
}

function initTextureCoords(gl, programInfo, textureCoords) {
    if (textureCoords != null && textureCoords.length > 0) {
        // Create a buffer for the positions.
        const textureCoordBuffer = gl.createBuffer();

        // Select the buffer as the one to apply buffer
        // operations to from here out.
        gl.bindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer);

        // Now pass the list of positions into WebGL to build the
        // shape. We do this by creating a Float32Array from the
        // JavaScript array, then use it to fill the current buffer.
        gl.bufferData(
            gl.ARRAY_BUFFER, // The kind of buffer this is
            textureCoords, // The data in an Array object
            gl.STATIC_DRAW // We are not going to change this data, so it is static
        );

        // Tell WebGL how to pull out the positions from the position
        // buffer into the vertexPosition attribute.
        {
            const numComponents = 2;
            const type = gl.FLOAT; // the data in the buffer is 32bit floats
            const normalize = false; // don't normalize between 0 and 1
            const stride = 0; // how many bytes to get from one set of values to the next
            // Set stride to 0 to use type and numComponents above
            const offset = 0; // how many bytes inside the buffer to start from

            // Set the information WebGL needs to read the buffer properly
            gl.vertexAttribPointer(
                programInfo.attribLocations.vertexUV,
                numComponents,
                type,
                normalize,
                stride,
                offset
            );
            // Tell WebGL to use this attribute
            gl.enableVertexAttribArray(
                programInfo.attribLocations.vertexUV);
        }

        return textureCoordBuffer;
    }
}

function initBitangentBuffer(gl, programInfo, bitangents) {
    if (bitangents != null && bitangents.length > 0) {
        // Create a buffer for the positions.
        const bitangentBuffer = gl.createBuffer();

        // Select the buffer as the one to apply buffer
        // operations to from here out.
        gl.bindBuffer(gl.ARRAY_BUFFER, bitangentBuffer);

        // Now pass the list of positions into WebGL to build the
        // shape. We do this by creating a Float32Array from the
        // JavaScript array, then use it to fill the current buffer.
        gl.bufferData(
            gl.ARRAY_BUFFER, // The kind of buffer this is
            bitangents, // The data in an Array object
            gl.STATIC_DRAW // We are not going to change this data, so it is static
        );

        // Tell WebGL how to pull out the positions from the position
        // buffer into the vertexPosition attribute.
        {
            const numComponents = 3;
            const type = gl.FLOAT; // the data in the buffer is 32bit floats
            const normalize = false; // don't normalize between 0 and 1
            const stride = 0; // how many bytes to get from one set of values to the next
            // Set stride to 0 to use type and numComponents above
            const offset = 0; // how many bytes inside the buffer to start from

            // Set the information WebGL needs to read the buffer properly
            gl.vertexAttribPointer(
                programInfo.attribLocations.vertexBitangent,
                numComponents,
                type,
                normalize,
                stride,
                offset
            );
            // Tell WebGL to use this attribute
            gl.enableVertexAttribArray(
                programInfo.attribLocations.vertexBitangent);
        }

        return bitangentBuffer;
    }
}

function initIndexBuffer(gl, elementArray) {

    // Create a buffer for the positions.
    const indexBuffer = gl.createBuffer();

    // Select the buffer as the one to apply buffer
    // operations to from here out.
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);

    // Now pass the list of positions into WebGL to build the
    // shape. We do this by creating a Float32Array from the
    // JavaScript array, then use it to fill the current buffer.
    gl.bufferData(
        gl.ELEMENT_ARRAY_BUFFER, // The kind of buffer this is
        elementArray, // The data in an Array object
        gl.STATIC_DRAW // We are not going to change this data, so it is static
    );

    return indexBuffer;
}

function loadJSONFile(cb, filePath) {
    fetch(filePath)
        .then((data) => {
            return data.json();
        })
        .then((jData) => {
            cb(jData);
        })
        .catch((err) => {
            console.error(err);
        })
}

function getTextures(gl, imgPath) {
    if (imgPath) {
        var texture = gl.createTexture();
        gl.bindTexture(gl.TEXTURE_2D, texture);
        gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE,
            new Uint8Array([255, 0, 0, 255])); // red

        const image = new Image();

        image.onload = function () {
            gl.bindTexture(gl.TEXTURE_2D, texture);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
            gl.texImage2D(
                gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA,
                gl.UNSIGNED_BYTE,
                image
            );
        }

        image.src = imgPath;
        return texture;
    }
}

function parseOBJFileToJSON(objFileURL, object, cb) {
    if (objFileURL) {
        fetch(objFileURL)
            .then((data) => {
                return data.text();
            })
            .then((text) => {
                let mesh = OBJLoader.prototype.parse(text);
                let parentObj;
                //iterate through objects
                for (let j = 0; j < mesh.length; j++) {
                    //iterate through the materials of the mesh
                    for (let i = 0; i < mesh[j].materials.length; i++) {
                        let vertices = mesh[j].geometry.vertices.slice(mesh[j].materials[i].groupStart * 3, mesh[j].materials[i].groupEnd * 3);
                        let uvs = mesh[j].geometry.uvs.slice(mesh[j].materials[i].groupStart * 2, mesh[j].materials[i].groupEnd * 2);
                        let normals = mesh[j].geometry.normals.slice(mesh[j].materials[i].groupStart * 3, mesh[j].materials[i].groupEnd * 3);
                        let newObject = JSON.parse(JSON.stringify(object));

                        if (j > 0) {
                            newObject.name = object.name + i;
                            newObject.parent = object.name;
                            newObject.parentTransform = object.position;
                        }

                        if (i > 0) {
                            newObject.name = newObject.name + i;
                            newObject.parent = object.name;
                            newObject.position = [0, 0, 0];
                            newObject.scale = [1, 1, 1];

                            newObject.parentTransform = parentObj.position;

                            newObject.type = "mesh";
                            let geometry = {
                                vertices,
                                uvs,
                                normals
                            }
                            parseMTL(mesh[j].materials[i].mtllib, mesh[j].materials[i].name, newObject, geometry, cb);
                        } else {
                            parentObj = newObject;
                            let geometry = {
                                vertices,
                                uvs,
                                normals
                            }
                            parseMTL(mesh[j].materials[i].mtllib, mesh[j].materials[i].name, newObject, geometry, cb);
                        }
                    }
                }
            })
            .catch((err) => {
                console.error(err);
            })
    }
}

/**
 * 
 * @param {hex value of color} hex 
 */
function hexToRGB(hex) {
    let r = hex.substring(1, 3);
    let g = hex.substring(3, 5);
    let b = hex.substring(5, 7);
    r = parseInt(r, 16);
    g = parseInt(g, 16);
    b = parseInt(b, 16);
    return [r / 255, g / 255, b / 255];
}

function parseSceneFile(file, state, cb) {
    state.pointLights = [];
    state.objects = [];
    state.directionalLights = [];

    fetch(file)
        .then((data) => {
            return data.json();
        })
        .then((jData) => {
            state.level = jData[0];
            state.numberOfObjectsToLoad = jData[0].objects.length;
            cb();
        })
        .catch((err) => {
            console.error(err);
        })
}

function createSceneFile(state, filename) {
    let totalState = [
        {
            objects: [],
            pointLights: [],
            directionalLights: [],
            settings: {

            }
        }];

    //objects first
    state.objects.forEach((object) => {
        if (object.type === "mesh") {
            if (!object.parent) {
                totalState[0].objects.push({
                    name: object.name ? object.name : null,
                    material: object.material ? object.material : null,
                    type: object.type ? object.type : null, //might change this to be an int value for speed
                    position: object.model.position ? [object.model.position[0], object.model.position[1], object.model.position[2]] : null,
                    scale: object.model.scale ? [object.model.scale[0], object.model.scale[1], object.model.scale[2]] : null,
                    diffuseTexture: object.model.diffuseTexture ? object.model.diffuseTexture : null,
                    normalTexture: object.model.normalTexture ? object.model.normalTexture : null,
                    rotation: object.model.rotation ? object.model.rotation : null,
                    parent: object.parent ? object.parent : null,
                    model: object.modelName ? object.modelName : null
                });
            }
        } else {
            totalState[0].objects.push({
                name: object.name ? object.name : null,
                material: object.material ? object.material : null,
                type: object.type ? object.type : null, //might change this to be an int value for speed
                position: object.model.position ? [object.model.position[0], object.model.position[1], object.model.position[2]] : null,
                scale: object.model.scale ? [object.model.scale[0], object.model.scale[1], object.model.scale[2]] : null,
                diffuseTexture: object.model.diffuseTexture ? object.model.diffuseTexture : null,
                normalTexture: object.model.normalTexture ? object.model.normalTexture : null,
                rotation: object.model.rotation ? object.model.rotation : null,
                parent: object.parent ? object.parent : null,
                model: object.modelName ? object.modelName : null
            });
        }
    });

    console.log(state)

    state.pointLights.forEach((light) => {
        totalState[0].pointLights.push({
            colour: [light.colour[0], light.colour[1], light.colour[2]],
            position: light.position,
            strength: light.strength,
            quadratic: light.quadratic,
            linear: light.linear,
            constant: light.constant
        })
    })

    state.directionalLights.forEach((light) => {
        totalState[0].directionalLights.push({
            colour: [light.colour[0], light.colour[1], light.colour[2]],
            position: light.position,
            direction: light.direction
        })
    })
    //write the savefile 
    fs.writeFile(filename, JSON.stringify(totalState), 'utf-8', () => {
        console.log("Writing complete!")
    })
}

function initShaderUniforms(gl, shaderProgram, uniforms, attribs) {
    let programInfo = {
        attribLocations: {},
        uniformLocations: {}
    };

    //map and check attribs
    attribs.map((attrib) => {
        programInfo.attribLocations[attrib] = gl.getAttribLocation(shaderProgram, attrib);
    })

    //map and check uniforms
    uniforms.map((uniform) => {
        programInfo.uniformLocations[uniform] = gl.getUniformLocation(shaderProgram, uniform);
    })

    programInfo.program = shaderProgram;

    return programInfo;
}

function calculateBitangents(vertices, uvs) {

    let tangents = [], bitangents = [];

    for (let i = 0; i < vertices.length / 3; i += 3) {

        let v0 = getVertexRowN(vertices, i);
        let v1 = getVertexRowN(vertices, i + 1);
        let v2 = getVertexRowN(vertices, i + 2);

        let uv0 = getUVRowN(uvs, i);
        let uv1 = getUVRowN(uvs, i + 1);
        let uv2 = getUVRowN(uvs, i + 2);

        let deltaPos1 = vec3.fromValues(0, 0, 0);
        let deltaPos2 = vec3.fromValues(0, 0, 0);

        vec3.sub(deltaPos1, v1, v0);
        vec3.sub(deltaPos2, v2, v0);

        let deltaUV1 = vec2.fromValues(0, 0);
        let deltaUV2 = vec2.fromValues(0, 0);

        vec2.sub(deltaUV1, uv1, uv0);
        vec2.sub(deltaUV2, uv2, uv0);

        let r = 1.0 / (deltaUV1[0] * deltaUV2[1] - deltaUV1[1] * deltaUV2[0]);

        //calculate the tangent
        let tangent = vec3.fromValues(0, 0, 0);
        let tempTangent1 = vec3.fromValues(0, 0, 0);
        let tempTangent2 = vec3.fromValues(0, 0, 0);

        vec3.scale(tempTangent1, deltaPos1, deltaUV2[1]);
        vec3.scale(tempTangent2, deltaPos2, deltaUV1[1]);

        vec3.subtract(tangent, tempTangent1, tempTangent2);
        vec3.scale(tangent, tangent, r);

        //calculate the bitangent
        let bitangent = vec3.fromValues(0, 0, 0);
        let tempBitangent1 = vec3.fromValues(0, 0, 0);
        let tempBitangent2 = vec3.fromValues(0, 0, 0);

        vec3.scale(tempBitangent1, deltaPos2, deltaUV1[0]);
        vec3.scale(tempBitangent2, deltaPos1, deltaUV2[0]);

        vec3.subtract(bitangent, tempBitangent1, tempBitangent2);
        vec3.scale(bitangent, bitangent, r);
        //push the same tangent and bitangent for all three vertices

        for (let j = 0; j < 3; j++) {
            bitangents.push(bitangent[0]);
            bitangents.push(bitangent[1]);
            bitangents.push(bitangent[2]);

            tangents.push(tangent[0]);
            tangents.push(tangent[1]);
            tangents.push(tangent[2]);
        }
    }
    return { tangents, bitangents };
}

function getVertexRowN(vertices, n) {
    let vertex = vec3.fromValues(vertices[n * 3], vertices[(n * 3) + 1], vertices[(n * 3) + 2]);
    return vertex;
}

function getUVRowN(uvs, n) {
    let uv = vec2.fromValues(uvs[n * 2], uvs[(n * 2) + 1]);
    return uv;
}

function toRadians(angle) {
    return angle * (Math.PI / 180);
}

function initDepthMap(gl, width, height) {
    const texture = gl.createTexture();
    gl.bindTexture(gl.TEXTURE_2D, texture);
    gl.texImage2D(
        gl.TEXTURE_2D,      // target
        0,                  // mip level
        gl.DEPTH_COMPONENT, // internal format
        width,   // width
        height,   // height
        0,                  // border
        gl.DEPTH_COMPONENT, // format
        gl.UNSIGNED_INT,    // type
        null);              // data
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

    const depthMapFBO = gl.createFramebuffer();
    gl.bindFramebuffer(gl.FRAMEBUFFER, depthMapFBO);
    gl.framebufferTexture2D(
        gl.FRAMEBUFFER,       // target
        gl.DEPTH_ATTACHMENT,  // attachment point
        gl.TEXTURE_2D,        // texture target
        texture,         // texture
        0);                   // mip level

    return { depthMapFBO, texture };
}

function initTangentBuffer(gl, programInfo, tangents) {
    if (tangents != null && tangents.length > 0) {
        // Create a buffer for the positions.
        const tangentBuffer = gl.createBuffer();

        // Select the buffer as the one to apply buffer
        // operations to from here out.
        gl.bindBuffer(gl.ARRAY_BUFFER, tangentBuffer);

        // Now pass the list of positions into WebGL to build the
        // shape. We do this by creating a Float32Array from the
        // JavaScript array, then use it to fill the current buffer.
        gl.bufferData(
            gl.ARRAY_BUFFER, // The kind of buffer this is
            tangents, // The data in an Array object
            gl.STATIC_DRAW // We are not going to change this data, so it is static
        );

        // Tell WebGL how to pull out the positions from the position
        // buffer into the vertexPosition attribute.
        {
            const numComponents = 3;
            const type = gl.FLOAT; // the data in the buffer is 32bit floats
            const normalize = false; // don't normalize between 0 and 1
            const stride = 0; // how many bytes to get from one set of values to the next
            // Set stride to 0 to use type and numComponents above
            const offset = 0; // how many bytes inside the buffer to start from

            // Set the information WebGL needs to read the buffer properly
            gl.vertexAttribPointer(
                programInfo.attribLocations.vertexTangent,
                numComponents,
                type,
                normalize,
                stride,
                offset
            );
            // Tell WebGL to use this attribute
            gl.enableVertexAttribArray(
                programInfo.attribLocations.vertexTangent);
        }

        // TODO: Create and populate a buffer for the UV coordinates

        return tangentBuffer;
    }
}

function vectorDistance(vec1, vec2) {
    let xDiff = vec2[0] - vec1[0];
    let yDiff = vec2[1] - vec1[1];
    let zDiff = vec2[2] - vec1[2];

    return Math.sqrt(xDiff + yDiff + zDiff);
}