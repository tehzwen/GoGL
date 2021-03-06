import UI from "../uiSetup.js";

class Plane {
    constructor(glContext, object) {
        this.state = {};
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "plane";
        this.loaded = false;
        this.reflective = object.reflective;
        this.refractionIndex = object.refractionIndex;
        this.initialTransform = { position: object.position, scale: object.scale, rotation: object.rotation };
        this.material = object.material;
        this.collide = object.collide;
        this.model = {
            vertices: [
                0.0, 0.5, 0.5,
                0.0, 0.5, 0.0,
                0.5, 0.5, 0.0,
                0.5, 0.5, 0.5,
            ],
            triangles: [
                0, 2, 1, 2, 0, 3,
            ],
            uvs: [
                0.0, 0.0,
                5.0, 0.0,
                5.0, 5.0,
                0.0, 5.0,
            ],
            normals: [
                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,
            ],
            bitangents: [
                0, -1, 0,
                0, -1, 0,
                0, -1, 0,
                0, -1, 0, // top
            ],
            diffuseTexture: object.diffuseTexture ? object.diffuseTexture : "default.png",
            normalTexture: object.normalTexture ? object.normalTexture : "defaultNorm.png",
            texture: object.diffuseTexture ? getTextures(glContext, object.diffuseTexture) : getTextures(glContext, "default.png"),
            textureNorm: object.normalTexture ? getTextures(glContext, object.normalTexture) : getTextures(glContext, "defaultNorm.png"),
            buffers: null,
            modelMatrix: mat4.create(),
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
        };
    }

    translate(translateVec) {
        vec3.add(this.model.position, this.model.position, vec3.fromValues(translateVec[0], translateVec[1], translateVec[2]));
        this.boundingBox = translateBoundingBox(this.boundingBox, translateVec);
    }

    scale(scaleVec) {
        let xVal = this.model.scale[0];
        let yVal = this.model.scale[1];
        let zVal = this.model.scale[2];

        xVal *= scaleVec[0];
        yVal *= scaleVec[1];
        zVal *= scaleVec[2];

        this.model.uvs = this.scaleUVs(this.model.uvs, scaleVec);
        this.boundingBox = scaleBoundingBox(this.boundingBox, scaleVec);
        this.model.scale = vec3.fromValues(xVal, yVal, zVal);
    }

    scaleUVs(uvs, scaleVec) {
        let newUVs = [];

        for (let i = 0; i < uvs.length; i++) {
            if (i % 2 === 0) {
                newUVs.push(uvs[i] * (scaleVec[0] / vec3.len(scaleVec)));
            } else {
                newUVs.push(uvs[i] * (scaleVec[2] / vec3.len(scaleVec)));
            }
        }
        return newUVs;
    }

    lightingShader() {
        var shaderProgram;
        var programInfo;

        if (this.material.shaderType === 0) {
            fetch('./shaders/basicShader.json')
                .then((res) => {
                    return res.json();
                })
                .then((data) => {
                    this.fragShader = data.fragShader.join("\n");
                    this.vertShader = data.vertShader.join("\n");
                    shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
                    programInfo = initShaderUniforms(this.gl, shaderProgram, data.uniforms, data.attribs);
                    UI.shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        } else if (this.material.shaderType === 1) {
            fetch('./shaders/blinnNoTexture.json')
                .then((res) => {
                    return res.json();
                })
                .then((data) => {
                    this.fragShader = data.fragShader.join("\n");
                    this.vertShader = data.vertShader.join("\n");
                    shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
                    programInfo = initShaderUniforms(this.gl, shaderProgram, data.uniforms, data.attribs);
                    UI.shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        } else if (this.material.shaderType === 2) {
            fetch('./shaders/basicDepthShader.json')
                .then((res) => {
                    return res.json();
                })
                .then((data) => {
                    this.fragShader = data.fragShader.join("\n");
                    this.vertShader = data.vertShader.join("\n");
                    shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
                    programInfo = initShaderUniforms(this.gl, shaderProgram, data.uniforms, data.attribs);
                    UI.shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        } else if (this.material.shaderType === 3) {
            fetch('./shaders/blinnTexture.json')
                .then((res) => {
                    return res.json();
                })
                .then((data) => {
                    this.fragShader = data.fragShader.join("\n");
                    this.vertShader = data.vertShader.join("\n");
                    shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
                    programInfo = initShaderUniforms(this.gl, shaderProgram, data.uniforms, data.attribs);
                    UI.shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        } else if (this.material.shaderType === 4) {
            let tangentCalc = calculateBitangents(this.model.vertices, this.model.uvs);
            this.model.bitangents = tangentCalc.bitangents;
            this.model.tangents = tangentCalc.tangents;

            fetch('./shaders/blinnDiffuseAndNormal.json')
                .then((res) => {
                    return res.json();
                })
                .then((data) => {
                    this.fragShader = data.fragShader.join("\n");
                    this.vertShader = data.vertShader.join("\n");
                    shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
                    programInfo = initShaderUniforms(this.gl, shaderProgram, data.uniforms, data.attribs);
                    UI.shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        }

    }

    initBuffers() {
        //create vertices, normal and indicies arrays
        const positions = new Float32Array(this.model.vertices.flat());
        const normals = new Float32Array(this.model.normals.flat());
        const indices = new Uint16Array(this.model.triangles);
        const textureCoords = new Float32Array(this.model.uvs);
        const bitangents = new Float32Array(this.model.bitangents);

        var vertexArrayObject = this.gl.createVertexArray();
        this.gl.bindVertexArray(vertexArrayObject);
        this.buffers;

        if (this.material.shaderType === 1) {
            this.buffers = {
                vao: vertexArrayObject,
                attributes: {
                    position: initPositionAttribute(this.gl, this.programInfo, positions),
                    normal: initNormalAttribute(this.gl, this.programInfo, normals),
                },
                indicies: initIndexBuffer(this.gl, indices),
                numVertices: indices.length
            }
        } else if (this.material.shaderType === 3) {
            this.buffers = {
                vao: vertexArrayObject,
                attributes: {
                    position: initPositionAttribute(this.gl, this.programInfo, positions),
                    normal: initNormalAttribute(this.gl, this.programInfo, normals),
                    uv: initTextureCoords(this.gl, this.programInfo, textureCoords),
                },
                indicies: initIndexBuffer(this.gl, indices),
                numVertices: indices.length
            }

            this.loaded = true;
        } else if (this.material.shaderType === 4) {
            this.buffers = {
                vao: vertexArrayObject,
                attributes: {
                    position: initPositionAttribute(this.gl, this.programInfo, positions),
                    normal: initNormalAttribute(this.gl, this.programInfo, normals),
                    uv: initTextureCoords(this.gl, this.programInfo, textureCoords),
                    bitangents: initBitangentBuffer(this.gl, this.programInfo, bitangents)
                },
                indicies: initIndexBuffer(this.gl, indices),
                numVertices: indices.length
            }
        }
        this.loaded = true;
    }

    reset() {
        this.gl.disableVertexAttribArray(this.buffers.vao);
        this.model.position = vec3.fromValues(0.0, 0.0, 0.0);
        this.model.rotation = mat4.create();
        this.model.scale = vec3.fromValues(1.0, 1.0, 1.0);
        this.model.uvs = [
            0.0, 0.0,
            5.0, 0.0,
            5.0, 5.0,
            0.0, 5.0,
        ]
        this.setup();
        this.gl.enableVertexAttribArray(this.buffers.vao);
    }

    setup() {
        this.centroid = calculateCentroid(this.model.vertices.flat());
        this.boundingBox = getBoundingBox(this.model.vertices);
        this.lightingShader();
        this.scale(this.initialTransform.scale);
        this.translate(this.initialTransform.position);

        if (this.initialTransform.rotation) {
            this.model.rotation = this.initialTransform.rotation;
        }
    }

    delete() {
        Object.keys(this.buffers.attributes).forEach((key) => {
            this.gl.deleteBuffer(this.buffers.attributes[key]);
        })
        this.gl.deleteBuffer(this.buffers.indicies);
    }
}

export default Plane;