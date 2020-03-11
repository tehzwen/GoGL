import UI from "../uiSetup.js";

class Model {
    constructor(glContext, object, meshDetails) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "mesh";
        this.loaded = false;
        this.modelName = object.model;
        this.parentTransform = object.parentTransform ? object.parentTransform : null;
        this.initialTransform = { position: object.position, scale: object.scale, rotation: object.rotation };
        this.material = object.mtl ? object.mtl : object.material;
        this.collide = object.collide;
        this.model = {
            normals: meshDetails.normals,
            vertices: meshDetails.vertices,
            uvs: meshDetails.uvs,
            position: vec3.fromValues(0, 0, 0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1, 1, 1),
            diffuseTexture: object.diffuseTexture ? object.diffuseTexture : null,
            texture: object.diffuseTexture ? getTextures(glContext, object.diffuseTexture ? object.diffuseTexture : "default.png") : getTextures(glContext, object.mtl.diffuseMap, true)
        };
        this.modelMatrix = mat4.create();
        this.lightingShader = this.lightingShader.bind(this);
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

    setup() {
        return new Promise((res, rej) => {
            this.lightingShader(res);
        })
    }


    initBuffers(resp) {
        //create vertices, normal and indicies arrays
        const positions = new Float32Array(this.model.vertices);
        const normals = new Float32Array(this.model.normals);
        const textureCoords = new Float32Array(this.model.uvs);

        var vertexArrayObject = this.gl.createVertexArray();
        this.gl.bindVertexArray(vertexArrayObject);

        this.buffers = {
            vao: vertexArrayObject,
            attributes: {
                position: this.programInfo.attribLocations.vertexPosition != null ? initPositionAttribute(this.gl, this.programInfo, positions) : null,
                normal: this.programInfo.attribLocations.vertexNormal != null ? initNormalAttribute(this.gl, this.programInfo, normals) : null,
                uv: this.programInfo.attribLocations.vertexUV != null ? initTextureCoords(this.gl, this.programInfo, textureCoords) : null,
            },
            numVertices: positions.length
        }

        if (this.parent) {
            //transform the bounding box to the parents location/scale
            this.boundingBox = translateBoundingBox(this.boundingBox, this.parentTransform);
        } else {
            this.scale(this.initialTransform.scale);
            this.translate(this.initialTransform.position);

            if (this.initialTransform.rotation) {
                this.model.rotation = this.initialTransform.rotation;
            }
        }


        this.loaded = true;
        console.log(this.name + " loaded successfully!");

        resp(this);
    }

    lightingShader(res) {
        //iterate through the objects
        let shaderProgram;
        let programInfo;

        //plain flat shading
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
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers(res);
                    return new Promise((res, rej) => {
                        res({ shaderLoad: true })
                    });
                })
                .catch((err) => {
                    console.error(err);
                })
        }
        //blinn phong with no textures
        else if (this.material.shaderType === 1) {
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
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers(res);
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
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers(res);
                })
                .catch((err) => {
                    console.error(err);
                })
        } else if (this.material.shaderType === 4) {
            fetch('./shaders/blinnNormalAndDiffuseTexture.json')
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
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers(res);
                })
                .catch((err) => {
                    console.error(err);
                })
        }
    }

    reset() {
        this.gl.disableVertexAttribArray(this.buffers.vao);
        this.gl.deleteProgram(this.programInfo.program);
        this.model.position = vec3.fromValues(0.0, 0.0, 0.0);
        this.model.rotation = mat4.create();
        this.model.scale = vec3.fromValues(1.0, 1.0, 1.0);
        this.setup();
        this.gl.enableVertexAttribArray(this.buffers.vao);
    }

    delete() {
        Object.keys(this.buffers.attributes).forEach((key) => {
            this.gl.deleteBuffer(this.buffers.attributes[key]);
        })
    }
}

export default Model;