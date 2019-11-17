class Light {
    constructor(glContext, name, meshDetails, parent = null, ambient, diffuse, specular, n, alpha, colour, strength) {
        this.gl = glContext;
        this.name = name;
        this.parent = parent;
        this.type = "light";
        this.loaded = false;

        this.material = { ambient, diffuse, specular, n, alpha };
        this.model = {
            normals: meshDetails.normals,
            vertices: meshDetails.vertices,
            uvs: meshDetails.uvs,
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
        };
        this.modelMatrix = mat4.create();
        this.colour = vec3.fromValues(colour[0], colour[1], colour[2]);
        this.strength = strength;

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
        this.centroid = calculateCentroid(this.model.vertices, this.lightingShader);
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
                position: initPositionAttribute(this.gl, this.programInfo, positions),
                normal: initNormalAttribute(this.gl, this.programInfo, normals),
                uv: initTextureCoords(this.gl, this.programInfo, textureCoords),
            },
            numVertices: positions.length
        }

        this.loaded = true;
        console.log(this.name + " loaded successfully!");
    }

    lightingShader() {
        //console.log(this.model.vertices)

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
                sampler: this.gl.getUniformLocation(shaderProgram, 'uTexture'),
                samplerExists: this.gl.getUniformLocation(shaderProgram, "samplerExists")
            },
        };

        shaderValuesErrorCheck(programInfo);
        this.programInfo = programInfo;
        this.initBuffers();

    }
}