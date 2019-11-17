class Plane {
    constructor(glContext, name, parent = null, ambient, diffuse, specular, n, alpha, texture, textureNorm) {
        this.state = {};
        this.gl = glContext;
        this.name = name;
        this.parent = parent;
        this.type = "primitive";
        this.loaded = false;

        this.material = { ambient, diffuse, specular, n, alpha, textureID: texture };
        this.model = {
            vertices: [
                [0.0, 0.5, 0.5],
                [0.0, 0.5, 0.0],
                [0.5, 0.5, 0.0],
                [0.5, 0.5, 0.5],
            ],
            triangles: [
                0, 1, 2, 0, 2, 3,
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
            texture: texture ? getTextures(glContext, texture) : null,
            textureNorm: textureNorm ? getTextures(glContext, textureNorm) : null,
            buffers: null,
            modelMatrix: mat4.create(),
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
            programInfo: null,
            fragShader: "",
            vertShader: ""
        };
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

    lightingShader() {
        const shaderProgram = initShaderProgram(this.gl, this.vertShader, this.fragShader);
        // Collect all the info needed to use the shader program.
        const programInfo = {
            // The actual shader program
            program: shaderProgram,
            // The attribute locations. WebGL will use there to hook up the buffers to the shader program.
            // NOTE: it may be wise to check if these calls fail by seeing that the returned location is not -1.
            attribLocations: {
                vertexPosition: this.gl.getAttribLocation(shaderProgram, 'aPosition'),
                vertexNormal: this.gl.getAttribLocation(shaderProgram, 'aNormal'),
                vertexUV: this.gl.getAttribLocation(shaderProgram, 'aUV'),
                vertexBitangent: this.gl.getAttribLocation(shaderProgram, 'aVertBitang')
            },
            uniformLocations: {
                projection: this.gl.getUniformLocation(shaderProgram, 'uProjectionMatrix'),
                view: this.gl.getUniformLocation(shaderProgram, 'uViewMatrix'),
                model: this.gl.getUniformLocation(shaderProgram, 'uModelMatrix'),
                normalMatrix: this.gl.getUniformLocation(shaderProgram, 'normalMatrix'),
                diffuseVal: this.gl.getUniformLocation(shaderProgram, 'diffuseVal'),
                ambientVal: this.gl.getUniformLocation(shaderProgram, 'ambientVal'),
                specularVal: this.gl.getUniformLocation(shaderProgram, 'specularVal'),
                nVal: this.gl.getUniformLocation(shaderProgram, 'nVal'),
                cameraPosition: this.gl.getUniformLocation(shaderProgram, 'uCameraPosition'),
                numLights: this.gl.getUniformLocation(shaderProgram, 'numLights'),
                lightPositions: this.gl.getUniformLocation(shaderProgram, 'uLightPositions'),
                lightColours: this.gl.getUniformLocation(shaderProgram, 'uLightColours'),
                lightStrengths: this.gl.getUniformLocation(shaderProgram, 'uLightStrengths'),
                samplerExists: this.gl.getUniformLocation(shaderProgram, "samplerExists"),
                sampler: this.gl.getUniformLocation(shaderProgram, 'uTexture'),
                normalSamplerExists: this.gl.getUniformLocation(shaderProgram, 'uTextureNormExists'),
                normalSampler: this.gl.getUniformLocation(shaderProgram, 'uTextureNorm')
            },
        };

        shaderValuesErrorCheck(programInfo);
        this.programInfo = programInfo;

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

        this.loaded = true;
    }

    setup() {
        this.lightingShader();
        this.centroid = calculateCentroid(this.model.vertices.flat());
        this.initBuffers();
    }
}