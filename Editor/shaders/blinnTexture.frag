#version 300 es
#define MAX_LIGHTS 128
precision highp float;

in vec3 oFragPosition;
in vec3 oNormal;
in vec3 normalInterp;
in vec2 oUV;
in vec3 oVertBitang;

uniform vec3 uCameraPosition;
uniform int numLights;
uniform vec3 diffuseVal;
uniform vec3 ambientVal;
uniform vec3 specularVal;
uniform float nVal;
uniform sampler2D uTexture;
uniform int samplerExists;
uniform int uTextureNormExists;
uniform sampler2D uTextureNorm;
uniform vec3 uLightPositions[MAX_LIGHTS];
uniform vec3 uLightColours[MAX_LIGHTS];
uniform float uLightStrengths[MAX_LIGHTS];

out vec4 fragColor;

void main() {
    vec3 normal = vec3(0);
    vec3 regularNormal = normalize(normalInterp);
    vec3 ambient = vec3(0,0,0);
    vec3 diffuse = vec3(0,0,0);
    vec3 specular = vec3(0,0,0);

    for (int i = 0; i < numLights; i++) {
        if (uTextureNormExists == 1) {
            normal = texture(uTextureNorm, oUV).xyz;
            normal = 2.0 * normal - 1.0;
            normal = normal * vec3(5.0, 5.0, 5.0);
            vec3 biTangent = cross(oNormal, oVertBitang);
            mat3 nMatrix = mat3(oVertBitang, biTangent, oNormal);
            normal = normalize(nMatrix * normal);
        }

        vec3 nCameraPosition = normalize(uCameraPosition); // Normalize the camera position
        vec3 V = normalize(nCameraPosition - oFragPosition);

        vec3 lightDirection = normalize(uLightPositions[i] - oFragPosition);
        float diff = max(dot(lightDirection, regularNormal + normal), 0.0);
        vec3 reflectDir = reflect(-lightDirection, normal);
        float spec = pow(max(dot(V, reflectDir), 0.0), nVal);
        float lightDistance = length(uLightPositions[i] - oFragPosition);
        float attenuation = 1.0 / (lightDistance * lightDistance);
        attenuation *= uLightStrengths[i];

        //ambient
        ambient += (ambientVal * uLightColours[i]);
        //diffuse
        diffuse += diffuseVal + uLightColours[i] * diff;

        if (diff > 0.0f)
        {
            specular += specularVal * uLightColours[i] * spec;
            specular *= attenuation;
        }

        ambient *= attenuation; //causes much darker scene
        diffuse *= attenuation;
    }

    vec4 textureColor = texture(uTexture, oUV);

    if (samplerExists == 1) {
        fragColor = vec4((ambient + diffuse + specular) * textureColor.rgb, 1.0);
    } else {
        fragColor = vec4(ambient + diffuse + specular, 1.0);
    }
    
}
