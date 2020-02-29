import UI from "../uiSetup.js"

class Cube {
    constructor(glContext, object) {
        this.state = {};
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = object.type;
        this.loaded = false;
        this.initialTransform = { position: object.position, scale: object.scale, rotation: object.rotation };
        this.material = object.material;
        this.collide = object.collide;
        this.model = {
            vertices: [
                0.0, 0.0, 0.0,
                0.0, 0.5, 0.0,
                0.5, 0.5, 0.0,
                0.5, 0.0, 0.0,

                0.0, 0.0, 0.5,
                0.0, 0.5, 0.5,
                0.5, 0.5, 0.5,
                0.5, 0.0, 0.5,

                0.0, 0.5, 0.5,
                0.0, 0.5, 0.0,
                0.5, 0.5, 0.0,
                0.5, 0.5, 0.5,

                0.0, 0.0, 0.5,
                0.5, 0.0, 0.5,
                0.5, 0.0, 0.0,
                0.0, 0.0, 0.0,

                0.5, 0.0, 0.5,
                0.5, 0.0, 0.0,
                0.5, 0.5, 0.5,
                0.5, 0.5, 0.0,

                0.0, 0.0, 0.5,
                0.0, 0.0, 0.0,
                0.0, 0.5, 0.5,
                0.0, 0.5, 0.0
            ],
            triangles: [
                //front face
                2, 0, 1, 3, 0, 2,
                //backface
                5, 4, 6, 6, 4, 7,
                //top face
                10, 9, 8, 10, 8, 11,
                //bottom face
                13, 12, 14, 14, 12, 15,
                //
                18, 16, 17, 18, 17, 19,

                22, 21, 20, 23, 21, 22,
            ],
            uvs: [
                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0,

                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0,

                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0,

                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0,

                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0,

                0.0, 0.0,
                1.0, 0.0,
                1.0, 1.0,
                0.0, 1.0
            ],
            normals: [
                0.0, 0.0, -1.0,
                0.0, 0.0, -1.0,
                0.0, 0.0, -1.0,
                0.0, 0.0, -1.0,

                0.0, 0.0, 1.0,
                0.0, 0.0, 1.0,
                0.0, 0.0, 1.0,
                0.0, 0.0, 1.0,

                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,
                0.0, 1.0, 0.0,

                0.0, -1.0, 0.0,
                0.0, -1.0, 0.0,
                0.0, -1.0, 0.0,
                0.0, -1.0, 0.0,

                1.0, 0.0, 0.0,
                1.0, 0.0, 0.0,
                1.0, 0.0, 0.0,
                1.0, 0.0, 0.0,

                -1.0, 0.0, 0.0,
                -1.0, 0.0, 0.0,
                -1.0, 0.0, 0.0,
                -1.0, 0.0, 0.0
            ],
            diffuseTexture: object.diffuseTexture ? object.diffuseTexture : null,
            normalTexture: object.normalTexture ? object.normalTexture : null,
            texture: object.diffuseTexture ? getTextures(glContext, object.diffuseTexture) : null,
            textureNorm: object.normalTexture ? getTextures(glContext, object.normalTexture) : null,
            buffers: null,
            modelMatrix: mat4.create(),
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
        };
    }

    scale(scaleVec) {

        //model scale
        let xVal = this.model.scale[0];
        let yVal = this.model.scale[1];
        let zVal = this.model.scale[2];

        //centroid scale
        let cenX = this.centroid[0];
        let cenY = this.centroid[1];
        let cenZ = this.centroid[2];

        cenX *= scaleVec[0];
        cenY *= scaleVec[1];
        cenZ *= scaleVec[2];

        xVal *= scaleVec[0];
        yVal *= scaleVec[1];
        zVal *= scaleVec[2];

        //need to scale bounding box
        this.boundingBox = scaleBoundingBox(this.boundingBox, scaleVec);
        this.centroid = vec3.fromValues(cenX, cenY, cenZ);
        this.model.scale = vec3.fromValues(xVal, yVal, zVal);
    }

    translate(translateVec) {
        vec3.add(this.model.position, this.model.position, vec3.fromValues(translateVec[0], translateVec[1], translateVec[2]));
        this.boundingBox = translateBoundingBox(this.boundingBox, translateVec);
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
        const positions = new Float32Array(this.model.vertices);
        const normals = new Float32Array(this.model.normals);
        const indices = new Uint16Array(this.model.triangles);
        const textureCoords = new Float32Array(this.model.uvs);
        const bitangents = new Float32Array(this.model.bitangents);
        const tangents = new Float32Array(this.model.tangents);

        var vertexArrayObject = this.gl.createVertexArray();
        this.gl.bindVertexArray(vertexArrayObject);

        this.buffers = {
            vao: vertexArrayObject,
            attributes: {
                position: this.programInfo.attribLocations.vertexPosition != null ? initPositionAttribute(this.gl, this.programInfo, positions) : null,
                normal: this.programInfo.attribLocations.vertexNormal != null ? initNormalAttribute(this.gl, this.programInfo, normals) : null,
                uv: this.programInfo.attribLocations.vertexUV != null ? initTextureCoords(this.gl, this.programInfo, textureCoords) : null,
                bitangents: this.programInfo.attribLocations.vertexBitangent != null ? initBitangentBuffer(this.gl, this.programInfo, bitangents) : null,
                tangents: this.programInfo.attribLocations.vertexTangent != null ? initBitangentBuffer(this.gl, this.programInfo, tangents) : null
            },
            indicies: initIndexBuffer(this.gl, indices),
            numVertices: indices.length
        }
        this.loaded = true;
    }

    setup() {
        this.centroid = calculateCentroid(this.model.vertices);
        this.boundingBox = getBoundingBox(this.model.vertices);
        this.scale(this.initialTransform.scale);
        this.translate(this.initialTransform.position);

        if (this.initialTransform.rotation) {
            this.model.rotation = this.initialTransform.rotation;
        }

        this.lightingShader();
    }

    delete() {
        Object.keys(this.buffers.attributes).forEach((key) => {
            this.gl.deleteBuffer(this.buffers.attributes[key]);
        })
        this.gl.deleteBuffer(this.buffers.indicies);
    }
}

export default Cube;