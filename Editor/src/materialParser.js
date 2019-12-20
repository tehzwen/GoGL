function parseMTL(mtlFile, mtlName, object, geometry, cb) {
    fetch("./models/" + mtlFile)
        .then((res) => {
            return res.text()
        })
        .then((text) => {
            let foundIndex = null;
            let lineArray = text.split("\n");
            let material = {
                diffuse: null,
                ambient: null,
                specular: null,
                n: null,
                shaderType: null,
                alpha: null,
                diffuseMap: null
            }

            //find the index it occurs on
            for (let i = 0; i < lineArray.length; i++) {
                if (lineArray[i].indexOf("newmtl") !== -1 && lineArray[i].indexOf(mtlName) !== -1) {
                    foundIndex = i;
                    break;
                }
            }

            //now loop the array starting at that index until a blank line is found and we ought to have our info
            for (let j = foundIndex; j < lineArray.length; j++) {
                if (lineArray[j] === "") {
                    break;
                } else {
                    //check for n value
                    let firstChars = lineArray[j].slice(0, 2);

                    if (firstChars === "Ns") {
                        //split on white space to get value
                        let whiteSpaceSplit = lineArray[j].split(" ");
                        material.n = parseFloat(whiteSpaceSplit[1]);
                    } else if (firstChars === "Ka") {
                        let whiteSpaceSplit = lineArray[j].split(" ");
                        material.ambient = [parseFloat(whiteSpaceSplit[1]), parseFloat(whiteSpaceSplit[2]), parseFloat(whiteSpaceSplit[3])];
                    } else if (firstChars === "Kd") {
                        let whiteSpaceSplit = lineArray[j].split(" ");
                        material.diffuse = [parseFloat(whiteSpaceSplit[1]), parseFloat(whiteSpaceSplit[2]), parseFloat(whiteSpaceSplit[3])];
                    } else if (firstChars === "Ks") {
                        let whiteSpaceSplit = lineArray[j].split(" ");
                        material.specular = [parseFloat(whiteSpaceSplit[1]), parseFloat(whiteSpaceSplit[2]), parseFloat(whiteSpaceSplit[3])];
                    } else if (firstChars === "d ") {
                        let whiteSpaceSplit = lineArray[j].split(" ");
                        material.alpha = parseFloat(whiteSpaceSplit[1]);
                    } else if (lineArray[j].indexOf("map_Kd") !== -1) {
                        //now we have to split it on // for the file TODO: make sure this works on linux kekw
                        let fileSplit = lineArray[j].split(" ");
                        material.diffuseMap = fileSplit[fileSplit.length - 1]
                    }
                }
            }

            if (material.diffuseMap) {
                material.shaderType = 3;
            } else {
                material.shaderType = 1;
            }
            
            object.mtl = material;
            cb(geometry, object);
        })
        .catch((err) => {
            if (err.message === "Failed to fetch") {
                cb(geometry, object);
            } else {
                console.error(err);
            }
            
        })
}