class Plane {
    constructor(glContext, object) {
        this.state = {};
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "plane";
        this.loaded = false;

        this.material = object.material;
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

        this.boundingBox = scaleBoundingBox(this.boundingBox, scaleVec);
        this.model.scale = vec3.fromValues(xVal, yVal, zVal);
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
                    shaderValuesErrorCheck(programInfo);
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
                    shaderValuesErrorCheck(programInfo);
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
                    shaderValuesErrorCheck(programInfo);
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
                    shaderValuesErrorCheck(programInfo);
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

    setup() {
        this.centroid = calculateCentroid(this.model.vertices.flat());
        this.boundingBox = getBoundingBox(this.model.vertices);
        this.lightingShader();
    }
}