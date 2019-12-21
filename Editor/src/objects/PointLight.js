class PointLight {
    constructor(glContext, object) {
        this.gl = glContext;
        this.name = object.name;
        this.parent = object.parent;
        this.type = "pointLight";
        this.loaded = false;
        this.position = object.position;
        this.colour = vec3.fromValues(object.colour[0], object.colour[1], object.colour[2]);
        this.strength = object.strength;
        this.quadratic = object.quadratic;
        this.linear = object.linear;
        this.constant = object.constant;
    }

    setup() {
        
    }

    translate(translateVec) {
        vec3.add(this.position, this.position, vec3.fromValues(translateVec[0], translateVec[1], translateVec[2]));
    }
}

export default PointLight;