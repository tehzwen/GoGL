// function getObject(state, name) {
//     return state.objects[state.objectTable[name]];
// }

function getObject(state, name) {
    for (let i = 0; i < state.objects.length; i++) {
        if (state.objects[i].name === name) {
            return state.objects[i];
        }
    }
    return null;
}

function intersect(a, b) {
    return (a.xMin <= b.xMax && a.xMax >= b.xMin) &&
        (a.yMin <= b.yMax && a.yMax >= b.yMin) &&
        (a.zMin <= b.zMax && a.zMax >= b.zMin)
}

function getBoundingBox(vertices) {
    let xMin = 0, xMax = 0, yMin = 0, yMax = 0, zMin = 0, zMax = 0;

    for (let i = 0; i < vertices.length / 3; i += 3) {
        if (vertices[i] > xMax) {
            xMax = vertices[i];
        }
        if (vertices[i] < xMin) {
            xMin = vertices[i];
        }
        if (vertices[i + 1] > yMax) {
            yMax = vertices[i + 1];
        }
        if (vertices[i + 1] < yMin) {
            yMin = vertices[i + 1];
        }
        if (vertices[i + 2] > zMax) {
            zMax = vertices[i + 2];
        }
        if (vertices[i + 2] < zMin) {
            zMin = vertices[i + 2];
        }
    }
    
    return { xMin, yMin, zMin, xMax, yMax, zMax };
}

function scaleBoundingBox(boundingBox, scaleVec) {
    let newBox = {
        xMin: boundingBox.xMin,
        yMin: boundingBox.yMin,
        zMin: boundingBox.zMin
    };

    if (boundingBox.xMax === 0 && scaleVec[0] > 1) {
        newBox.xMax = boundingBox.xMax + scaleVec[0]
    } else {
        newBox.xMax = boundingBox.xMax * scaleVec[0];
    }

    if (boundingBox.yMax === 0 && scaleVec[1] > 1) {
        newBox.yMax = boundingBox.yMax + scaleVec[1]
    } else {
        newBox.yMax = boundingBox.yMax * scaleVec[1];
    }

    if (boundingBox.zMax === 0 && scaleVec[2] > 1) {
        newBox.zMax = boundingBox.zMax + scaleVec[2]
    } else {
        newBox.zMax = boundingBox.zMax * scaleVec[2];
    }
    return newBox;
}

function translateBoundingBox(boundingBox, translateVector) {
    let newBox = {};

    newBox.xMin = boundingBox.xMin + translateVector[0];
    newBox.xMax = boundingBox.xMax + translateVector[0];
    newBox.yMin = boundingBox.yMin + translateVector[1];
    newBox.yMax = boundingBox.yMax + translateVector[1];
    newBox.zMin = boundingBox.zMin + translateVector[2];
    newBox.zMax = boundingBox.zMax + translateVector[2];

    return newBox;
}