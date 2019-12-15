class Light {
    constructor(glContext, object) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "light";
        this.loaded = false;

        this.material = object.material;
        this.model = {
            normals: null,
            vertices: null,
            uvs: null,
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
        };
        this.modelMatrix = mat4.create();
        this.colour = vec3.fromValues(object.colour[0], object.colour[1], object.colour[2]);
        this.strength = object.strength;
        this.modelName = "./models/lightbulb.obj";

        this.lightingShader = this.lightingShader.bind(this);
    }

    scale(scaleVec) {
        let xVal = this.model.scale[0];
        let yVal = this.model.scale[1];
        let zVal = this.model.scale[2];

        xVal *= scaleVec[0];
        yVal *= scaleVec[1];
        zVal *= scaleVec[2];

        this.model.scale = vec3.fromValues(xVal, yVal, zVal);
    }

    setup() {
        //this.centroid = calculateCentroid(this.model.vertices, this.lightingShader);
        parseOBJFileToJSON(this.modelName, this.lightingShader);
    }


    initBuffers() {
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

        this.loaded = true;
        console.log(this.name + " loaded successfully!");
    }

    lightingShader(mesh) {
        this.model.vertices = mesh.vertices;
        this.model.normals = mesh.normals;
        this.model.uvs = mesh.uvs;
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
                    shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers();
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
                    shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
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
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers();
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
                    shaderValuesErrorCheck(programInfo);
                    this.programInfo = programInfo;
                    this.centroid = calculateCentroid(this.model.vertices);
                    this.boundingBox = getBoundingBox(this.model.vertices);
                    this.initBuffers();
                })
                .catch((err) => {
                    console.error(err);
                })
        }
    }
}