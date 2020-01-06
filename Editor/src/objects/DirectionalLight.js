class DirectionalLight {
    constructor(glContext, object) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "directionalLight";
        this.loaded = false;
        this.position = object.position;
        this.colour = vec3.fromValues(object.colour[0], object.colour[1], object.colour[2]);
        this.direction = vec3.fromValues(object.direction[0], object.direction[1], object.direction[2]);
    }

    setup() {
        
    }

    translate(translateVec) {
        vec3.add(this.position, this.position, vec3.fromValues(translateVec[0], translateVec[1], translateVec[2]));
    }
}

export default DirectionalLight;