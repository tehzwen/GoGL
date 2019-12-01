class Cube {
    constructor(glContext, object) {
        this.state = {};
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = object.type;
        this.loaded = false;

        this.material = object.material;
        this.model = {
            vertices: [
                [0.0, 0.0, 0.0],
                [0.0, 0.5, 0.0],
                [0.5, 0.5, 0.0],
                [0.5, 0.0, 0.0],

                [0.0, 0.0, 0.5],
                [0.0, 0.5, 0.5],
                [0.5, 0.5, 0.5],
                [0.5, 0.0, 0.5],

                [0.0, 0.5, 0.5],
                [0.0, 0.5, 0.0],
                [0.5, 0.5, 0.0],
                [0.5, 0.5, 0.5],

                [0.0, 0.0, 0.5],
                [0.5, 0.0, 0.5],
                [0.5, 0.0, 0.0],
                [0.0, 0.0, 0.0],

                [0.5, 0.0, 0.5],
                [0.5, 0.0, 0.0],
                [0.5, 0.5, 0.5],
                [0.5, 0.5, 0.0],

                [0.0, 0.0, 0.5],
                [0.0, 0.0, 0.0],
                [0.0, 0.5, 0.5],
                [0.0, 0.5, 0.0]
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
            bitangents: [
                0, -1, 0,
                0, -1, 0,
                0, -1, 0,
                0, -1, 0, // Front

                0, -1, 0,
                0, -1, 0,
                0, -1, 0,
                0, -1, 0, // Back

                0, -1, 0,
                0, -1, 0,
                0, -1, 0,
                0, -1, 0, // Right

                0, -1, 0,
                0, -1, 0,
                0, -1, 0,
                0, -1, 0, // Left

                0, 0, 1,
                0, 0, 1,
                0, 0, 1,
                0, 0, 1, // Top

                0, 0, -1,
                0, 0, -1,
                0, 0, -1,
                0, 0, -1, // Bot
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
        let xVal = this.model.scale[0];
        let yVal = this.model.scale[1];
        let zVal = this.model.scale[2];

        xVal *= scaleVec[0];
        yVal *= scaleVec[1];
        zVal *= scaleVec[2];

        //need to scale bounding box
        this.boundingBox = scaleBoundingBox(this.boundingBox, scaleVec);

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
            shaderProgram = initShaderProgram(this.gl, shaders.flatNoTexture.vert, shaders.flatNoTexture.frag);
            programInfo = {
                // The actual shader program
                program: shaderProgram,
                attribLocations: setupAttributes(this.gl, shaders.flatNoTexture.attributes, shaderProgram),
                uniformLocations: setupUniforms(this.gl, shaders.flatNoTexture.uniforms, shaderProgram),
            };
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
        
        if (this.material.shaderType === 1 || this.material.shaderType === 0) {
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
        this.lightingShader();
        this.centroid = calculateCentroid(this.model.vertices.flat());
        this.boundingBox = getBoundingBox(this.model.vertices);
        this.initBuffers();
    }
}