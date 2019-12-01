class Model {
    constructor(glContext, object, meshDetails) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "mesh";
        this.loaded = false;
        this.modelName = object.model;

        this.material = object.material;
        this.model = {
            normals: meshDetails.normals,
            vertices: meshDetails.vertices,
            uvs: meshDetails.uvs,
            position: vec3.fromValues(0.0, 0.0, 0.0),
            rotation: mat4.create(),
            scale: vec3.fromValues(1.0, 1.0, 1.0),
            texture: object.texture ? getTextures(glContext, object.texture) : null
        };
        this.modelMatrix = mat4.create();

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
        this.buffers;

        if (this.material.shaderType === 1) {
            this.buffers = {
                vao: vertexArrayObject,
                attributes: {
                    position: initPositionAttribute(this.gl, this.programInfo, positions),
                    normal: initNormalAttribute(this.gl, this.programInfo, normals),
                },
                numVertices: positions.length
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
                numVertices: positions.length
            }
        }

        this.loaded = true;
        console.log(this.name + " loaded successfully!");
    }

    lightingShader() {
        //console.log(this.model.vertices)

        var shaderProgram;
        var programInfo;

        if (this.material.shaderType === 0) {
            console.warn("DO FLAT SHADING!")
        } else if (this.material.shaderType === 1) {
            shaderProgram = initShaderProgram(this.gl, shaders.blinnNoTexture.vert, shaders.blinnNoTexture.frag);
            programInfo = {
                // The actual shader program
                program: shaderProgram,
                attribLocations: setupAttributes(this.gl, shaders.blinnNoTexture.attributes, shaderProgram),
                uniformLocations: setupUniforms(this.gl, shaders.blinnNoTexture.uniforms, shaderProgram),
            };
        } else if (this.material.shaderType === 3) {
            //blinn phong with diffusetexture only
        } else if (this.material.shaderType === 4) {
            shaderProgram = initShaderProgram(this.gl, shaders.blinnTexture.vert, shaders.blinnTexture.frag);
            programInfo = {
                // The actual shader program
                program: shaderProgram,
                attribLocations: setupAttributes(this.gl, shaders.blinnTexture.attributes, shaderProgram),
                uniformLocations: setupUniforms(this.gl, shaders.blinnTexture.uniforms, shaderProgram),
            };
        }

        
        shaderValuesErrorCheck(programInfo);
        this.programInfo = programInfo;
        this.initBuffers();

    }
}