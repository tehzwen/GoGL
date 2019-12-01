const fs = require('fs');

var shaders = {};

fs.readdir(__dirname + "/shaders/", (err, files) => {
    if (err) {
        console.error(err);
    } else {
        numShaders = files.length;
        var loaded = 0;
        let requests = files.map((file) => {
            tempFile = file.split('.');

            if (tempFile[1] === "vert" || tempFile[1] === "frag") {
                //console.warn(file);
                return new Promise((resolve) => {
                    readShader(__dirname + "/shaders/" + file, tempFile[0], tempFile[1], resolve);
                })
            }

        });
        Promise.all(requests).then(() => console.log("Completed loading shaders..."))
    }
});

function getUniformsFromShaders(shaders) {

    let shaderKeys = Object.keys(shaders);

    for (let i = 0; i < shaderKeys.length; i++) {
        //iterate through the vert shader first
        let vertSplit = shaders[shaderKeys[i]].vert.split("\n");
        var attributeRE = new RegExp('\\bin\\b');

        vertSplit.map((vLine) => {
            if (vLine.indexOf('uniform') !== -1) {
                let spaceSplit = vLine.split(" ");
                if (spaceSplit[2].indexOf('[') !== -1) {
                    shaders[shaderKeys[i]].uniforms.push(spaceSplit[2].split('[')[0]);
                } else {
                    shaders[shaderKeys[i]].uniforms.push(spaceSplit[2].slice(0, -1));
                }
            } else if (attributeRE.exec(vLine)) {
                shaders[shaderKeys[i]].attributes.push(vLine.split(" ")[2].slice(0, -1));
            }
        })
        //iterate through the fragShader
        let fragSplit = shaders[shaderKeys[i]].frag.split("\n");
        fragSplit.map((fLine) => {
            if (fLine.indexOf('uniform') !== -1) {
                //check if the uniform is an array
                let spaceSplit = fLine.split(" ");
                if (spaceSplit[2].indexOf('[') !== -1) {
                    shaders[shaderKeys[i]].uniforms.push(spaceSplit[2].split('[')[0]);
                } else {
                    shaders[shaderKeys[i]].uniforms.push(spaceSplit[2].slice(0, -1));
                }
            }
        })
    }
}

function setupAttributes(gl, shaderAttribs, shaderProgram) {
    let attribs = {};
    shaderAttribs.map((attr) => {
        attribs[attr] = gl.getAttribLocation(shaderProgram, attr)
    })

    return attribs;
} 

function setupUniforms(gl, shaderUniforms, shaderProgram) {
    let uniforms = {};

    shaderUniforms.map((uni) => {
        uniforms[uni] = gl.getUniformLocation(shaderProgram, uni)
    })

    return uniforms;
}

function readShader(file, shaderName, type, cb) {
    fetch(file)
        .then((res) => {
            return res.text();
        })
        .then((data) => {
            shaders[shaderName] = {...shaders[shaderName], [type]: data, uniforms:[], attributes: []}
            cb();
        })
        .catch((err) => {
            console.error(err);
        })
}