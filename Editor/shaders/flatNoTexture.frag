#version 300 es
#define MAX_LIGHTS 128
precision highp float;

in vec3 oFragPosition;
in vec3 oNormal;
in vec3 normalInterp;

uniform int numLights;
uniform vec3 diffuseVal;
uniform vec3 ambientVal;
uniform vec3 uLightPositions[MAX_LIGHTS];
uniform vec3 uLightColours[MAX_LIGHTS];
uniform float uLightStrengths[MAX_LIGHTS];

out vec4 fragColor;

void main() {
    vec3 normal = normalize(normalInterp);
    vec3 ambient = vec3(0,0,0);
    vec3 diffuse = vec3(0,0,0);

    for (int i = 0; i < numLights; i++) {

        vec3 lightDirection = normalize(uLightPositions[i] - oFragPosition);
        float diff = max(dot(lightDirection, normal), 0.0);
        vec3 reflectDir = reflect(-lightDirection, normal);
        float lightDistance = length(uLightPositions[i] - oFragPosition);
        float attenuation = 1.0 / (lightDistance * lightDistance);
        attenuation *= uLightStrengths[i];

        //ambient
        ambient += (ambientVal * uLightColours[i]) * diffuseVal;
        //diffuse
        diffuse += diffuseVal * uLightColours[i] * diff;

        //ambient *= attenuation; causes much darker scene
        diffuse *= attenuation;
    }

    fragColor = vec4(diffuse + ambient, 1.0);
}