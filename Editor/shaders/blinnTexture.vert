#version 300 es
in vec3 vertexPosition;
in vec3 vertexNormal;
in vec2 vertexUV;
in vec3 vertexBitangent;

uniform mat4 uProjectionMatrix;
uniform mat4 uViewMatrix;
uniform mat4 uModelMatrix;
uniform mat4 normalMatrix;

out vec3 oFragPosition;
out vec3 oNormal;
out vec3 normalInterp;
out vec2 oUV;
out vec3 oVertBitang;

void main() {
    gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(vertexPosition, 1.0);

    oFragPosition = (uModelMatrix * vec4(vertexPosition, 1.0)).xyz;
    oNormal = normalize((uModelMatrix * vec4(vertexNormal, 1.0)).xyz);
    normalInterp = vec3(normalMatrix * vec4(vertexNormal, 0.0));
    oUV = vertexUV;
    oVertBitang = vertexBitangent;
}
