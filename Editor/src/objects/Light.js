class Light {
    constructor(glContext, object, meshDetails) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "light";
        this.loaded = false;

        this.material = object.material;
        this.model = {
            normals: meshDetails.normals,
            vertices: meshDetails.vertices,
            uvs: meshDetails.uvs,
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
        var shaderProgram;
        var programInfo;

        shaderProgram = initShaderProgram(this.gl, shaders.flatNoTexture.vert, shaders.flatNoTexture.frag);
        programInfo = {
            // The actual shader program
            program: shaderProgram,
            attribLocations: setupAttributes(this.gl, shaders.flatNoTexture.attributes, shaderProgram),
            uniformLocations: setupUniforms(this.gl, shaders.flatNoTexture.uniforms, shaderProgram),
        };

        shaderValuesErrorCheck(programInfo);
        this.programInfo = programInfo;
        this.initBuffers();

    }
}