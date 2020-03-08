const _ = require('lodash')

function setup(newState) {
    //listeners for header button presses
    document.getElementById('saveButton').addEventListener('click', () => {
        if (!document.location.href.includes("editor")) {
            createSceneFile(newState, "./statefiles/" + newState.saveFile);
        }
    })

    document.getElementById('launchButton').addEventListener('click', () => {
        document.location.href = "compiler.html?scene=" + __dirname + "/statefiles/" + newState.saveFile
    })
}

function createSceneGui(state) {
    let loadingProgress = ((state.objects.length / state.numberOfObjectsToLoad) * 100);
    state.loadingBar.child.style.width = loadingProgress + "%";

    if (loadingProgress === 100) {
        state.loadingBar.parent.style.display = "none";
    }

    //get objects first
    let sideNav = document.getElementById("objectsNav");
    sideNav.innerHTML = "";

    state.objects.forEach((object) => {
        if (!object.parent) {
            let objectElement = document.createElement("div");
            let childrenDiv = document.createElement("div");
            childrenDiv.id = object.name + "childDiv";
            childrenDiv.className = "collapse";
            let objectName = document.createElement("h5");
            objectName.className = "object-link";
            objectName.innerHTML = object.name;
            objectName.addEventListener('click', () => {
                //its not open yet
                if (childrenDiv.className === "collapse") {
                    childrenDiv.className = "collapse-active";
                    objectName.className = "object-link"
                    displayObjectValues(object, state);
                } else {
                    childrenDiv.className = "collapse";
                    displayObjectValues(object, state);
                    //displayObjectValues(null);
                }
            });

            objectElement.appendChild(objectName);
            objectElement.appendChild(childrenDiv);
            sideNav.appendChild(objectElement);
        } else {
            //find the parent's children div and throw some data in it
            let parentChildrenDiv = document.getElementById(object.parent + "childDiv");
            let objectName = document.createElement("h6");
            objectName.className = "object-link";
            objectName.innerHTML = `<i>${object.name}</i>`;

            objectName.addEventListener('click', () => {
                displayObjectValues(object, state);
            });
            parentChildrenDiv.appendChild(objectName);
        }
    });

    state.pointLights.forEach((object) => {
        let objectElement = document.createElement("div");
        let childrenDiv = document.createElement("div");
        let objectName = document.createElement("h5");
        objectName.className = "object-link";
        objectName.innerHTML = object.name;

        objectName.addEventListener('click', () => {
            displayPointLightValues(object, state);
        });

        objectElement.appendChild(objectName);
        objectElement.appendChild(childrenDiv);
        sideNav.appendChild(objectElement);
    })

    state.directionalLights.forEach((object) => {
        let objectElement = document.createElement("div");
        let childrenDiv = document.createElement("div");
        let objectName = document.createElement("h5");
        objectName.className = "object-link";
        objectName.innerHTML = object.name;

        objectName.addEventListener('click', () => {
            displayObjectValues(object, state);
        });

        objectElement.appendChild(objectName);
        objectElement.appendChild(childrenDiv);
        sideNav.appendChild(objectElement);
    })
}

function displayPointLightValues(object, state) {
    let selectedObjectDiv = document.getElementById("selectedObject");
    selectedObjectDiv.name = object.name;
    selectedObjectDiv.innerHTML = "";
    let nameInput = document.createElement('input');
    nameInput.value = object.name;
    nameInput.style.textAlign = 'center'
    nameInput.addEventListener('input', (e) => {
        //change the children's parent name aswell
        state.objects.forEach((o) => {
            if (o.parent === object.name) {
                o.parent = e.target.value;
            }
        })
        object.name = e.target.value;
    })
    nameInput.addEventListener('focusout', (e) => {
        createSceneGui(state);
    })

    let positionalInputDiv = displayPositionalValues(object.position, object, state);

    //create input for color of light
    let colorPicker = document.createElement("input");
    colorPicker.type = "color";
    colorPicker.classList = "form-control";
    colorPicker.value = rgbToHex(object.colour);
    colorPicker.addEventListener('change', (event) => {
        let newColor = hexToRGB(event.target.value);
        object.colour = newColor;
        state.render = true;
    });


    //create input for light strength
    let strengthInput = document.createElement("input");
    strengthInput.type = "number";
    strengthInput.value = object.strength;
    strengthInput.style.textAlign = 'center';
    strengthInput.addEventListener("input", (e) => {
        object.strength = parseInt(e.target.value);
        state.render = true;
    })

    //add a delete button to remove this object from the scene
    let deleteDiv = document.createElement("div");
    let deleteButton = document.createElement("button");
    deleteButton.classList = "btn btn-danger";
    deleteButton.innerHTML = "Delete"
    deleteButton.style.marginTop = '15px';
    deleteButton.addEventListener('click', () => {
        let lightIndex = _.findIndex(state.pointLights, { name: object.name })
        state.pointLights.splice(lightIndex, 1);
        state.render = true;
        createSceneGui(state);
        selectedObjectDiv.innerHTML = "";
    });


    selectedObjectDiv.appendChild(nameInput);
    selectedObjectDiv.appendChild(createHeader("Position", "h4"));
    selectedObjectDiv.appendChild(positionalInputDiv);
    selectedObjectDiv.appendChild(createHeader("Color", "h4"));
    selectedObjectDiv.appendChild(colorPicker);
    selectedObjectDiv.appendChild(createHeader("Strength", "h4"));
    selectedObjectDiv.appendChild(strengthInput);
    deleteDiv.appendChild(deleteButton);
    selectedObjectDiv.appendChild(deleteDiv);
}

function createSelectionOptions(optionsArr, selectObj) {
    optionsArr.map((option) => {
        let tempSelect = document.createElement("option");
        tempSelect.innerHTML = option;
        selectObj.appendChild(tempSelect);
    })
}

function handleTypeSelectChange(event) {
    if (event.target.value === "Mesh") {
        let addNav = document.getElementById("addObjectsNav");
        let addButton = addNav.lastChild;

        let fileUpload = document.createElement("input");
        fileUpload.id = "meshUpload";
        fileUpload.type = "file";
        fileUpload.classList = "form-control-file";

        addNav.insertBefore(fileUpload, addButton);
    } else {
        let fileUpload = document.getElementById("meshUpload");
        if (fileUpload) {
            fileUpload.remove();
        }
    }
}

/**
 * Purpose: Function to show controls to manipulate object,
 * does different actions depending on what object we select
 * @param {Game Object to manipulate} object 
 */
function displayObjectValues(object, state) {
    if (!object) {
        let selectedObjectDiv = document.getElementById("selectedObject");
        selectedObjectDiv.innerHTML = "";
        return;
    }

    let selectedObjectDiv = document.getElementById("selectedObject");
    selectedObjectDiv.name = object.name;
    selectedObjectDiv.innerHTML = "";

    let positionalInputDiv = displayPositionalValues(object.model.position, object, state);
    //create input with the name of the object and add listener to onChange set the name of the object accordingly
    let nameInput = document.createElement('input');
    nameInput.value = object.name;
    nameInput.style.textAlign = 'center'
    nameInput.addEventListener('input', (e) => {
        //change the children's parent name aswell
        state.objects.forEach((o) => {
            if (o.parent === object.name) {
                o.parent = e.target.value;
            }
        })
        object.name = e.target.value;
    })
    nameInput.addEventListener('focusout', (e) => {
        createSceneGui(state);
    })

    //add a delete button to remove this object from the scene
    let deleteButton = document.createElement("button");
    deleteButton.classList = "btn btn-danger";
    deleteButton.innerHTML = "Delete"
    deleteButton.style.marginTop = '15px';
    deleteButton.addEventListener('click', () => {
        if (object.type !== "mesh") {
            let deleteIndex = _.findIndex(state.objects, (o) => {
                return o.name === object.name;
            })
            state.objects.splice(deleteIndex, 1);
            state.numberOfObjectsToLoad--;
            state.render = true;
            createSceneGui(state);
            object.delete()
            selectedObjectDiv.innerHTML = "";
        } else {
            let deleteIndexes = _.map(_.keys(_.pickBy(state.objects, { modelName: object.modelName })), Number)
            for (let i = 0; i < deleteIndexes.length; i++) {
                state.objects[deleteIndexes[(deleteIndexes.length - 1) - i]].delete();
                state.objects.splice(deleteIndexes[(deleteIndexes.length - 1) - i], 1);
                state.numberOfObjectsToLoad--;
                state.render = true;
                createSceneGui(state);
                selectedObjectDiv.innerHTML = "";
            }
        }
    })

    selectedObjectDiv.appendChild(nameInput);
    selectedObjectDiv.appendChild(createHeader("Position", "h4"));
    selectedObjectDiv.appendChild(positionalInputDiv);
    selectedObjectDiv.appendChild(createHeader("Rotation", "h4"));
    selectedObjectDiv.appendChild(displayRotationValues(object, state));
    selectedObjectDiv.appendChild(createHeader("Scale", "h4"));
    selectedObjectDiv.appendChild(displayScaleValues(object, state));

    //create material ui elements
    createMaterialUI(state, object, selectedObjectDiv);
    selectedObjectDiv.appendChild(deleteButton);
}

function displayPositionalValues(position, object, state) {
    let positionalInputDiv = document.createElement("div");
    positionalInputDiv.classList = "input-group";

    //X move input handler
    let objectPositionX = document.createElement("input");
    objectPositionX.type = "number";
    objectPositionX.addEventListener('input', (event) => {
        object.translate([event.target.value - position[0], 0, 0]);
        state.render = true;
    })
    objectPositionX.id = object.name + "-positionX";
    objectPositionX.classList = "form-control";
    objectPositionX.value = parseFloat(position[0]).toFixed(1);

    //Y move input handler
    let objectPositionY = document.createElement("input");
    objectPositionY.type = "number";
    objectPositionY.addEventListener('input', (event) => {
        object.translate([0, event.target.value - position[1], 0]);
        state.render = true;
    })
    objectPositionY.id = object.name + "-positionY";
    objectPositionY.classList = "form-control";
    objectPositionY.value = parseFloat(position[1]).toFixed(1);

    //Z move input handler
    let objectPositionZ = document.createElement("input");
    objectPositionZ.type = "number";
    objectPositionZ.addEventListener('input', (event) => {
        object.translate([0, 0, event.target.value - position[2]])
        state.render = true;
    })
    objectPositionZ.id = object.name + "-positionZ";
    objectPositionZ.classList = "form-control";
    objectPositionZ.value = parseFloat(position[2]).toFixed(1);

    let prependDivX = document.createElement("div");
    prependDivX.classList = "input-group-prepend";

    prependDivX.innerHTML = `
        <span class="input-group-text">X</span>
        `;
    let prependDivY = prependDivX.cloneNode(true);
    prependDivY.innerHTML = `
        <span class="input-group-text">Y</span>
        `;

    let prependDivZ = prependDivX.cloneNode(true);
    prependDivZ.innerHTML = `
        <span class="input-group-text">Z</span>
        `;

    //add all the elements in
    positionalInputDiv.appendChild(prependDivX);
    positionalInputDiv.appendChild(objectPositionX);
    positionalInputDiv.appendChild(prependDivY);
    positionalInputDiv.appendChild(objectPositionY);
    positionalInputDiv.appendChild(prependDivZ);
    positionalInputDiv.appendChild(objectPositionZ);

    return positionalInputDiv;
}


function createMaterialUI(state, object, mainDiv) {
    let texturesDiv = document.createElement("div");
    let materialTitle = document.createElement("h3");
    materialTitle.innerHTML = "Material";
    materialTitle.style.marginTop = "25px";

    let diffuseNameTitle = document.createElement("p");
    diffuseNameTitle.innerHTML = object.model.diffuseTexture;
    diffuseNameTitle.classList = "orange-text";
    let normalNameTitle = document.createElement("p");
    normalNameTitle.innerHTML = object.model.normalTexture;
    normalNameTitle.classList = "orange-text";

    let diffuseInput = document.createElement("input");
    diffuseInput.type = "file";
    diffuseInput.id = "diffuseMaterialInput";
    //event listener for adding diffuse texture
    diffuseInput.addEventListener('input', (e) => {
        object.model.texture = getTextures(state.gl, e.target.files[0].name);
        object.model.diffuseTexture = e.target.files[0].name;
        diffuseNameTitle.innerHTML = e.target.files[0].name;
    })

    let normalInput = document.createElement("input");
    normalInput.type = "file";
    normalInput.id = "normalMaterialInput";

    normalInput.addEventListener('input', (e) => {
        object.model.textureNorm = getTextures(state.gl, e.target.files[0].name);
        object.model.normalTexture = e.target.files[0].name;
        normalNameTitle.innerHTML = e.target.files[0].name;
    })

    if (object.material.shaderType === 3) {
        texturesDiv.append(diffuseNameTitle);
        texturesDiv.append(diffuseInput);
    } else if (object.material.shaderType === 4) {
        texturesDiv.append(diffuseNameTitle);
        texturesDiv.append(diffuseInput);
        texturesDiv.append(normalNameTitle);
        texturesDiv.append(normalInput);
    }

    //create select dropdown for materials
    let matTypeDropdown = document.createElement("select")
    matTypeDropdown.addEventListener('change', (e) => {
        //check if the material has a texture applied to it already or not
        let numShaderType = parseInt(e.target.value);
        if (numShaderType === 0 || numShaderType === 1) {
            object.material.shaderType = numShaderType;
            object.reset();
            state.render = true;
            diffuseNameTitle.style.display = 'none';
            normalNameTitle.style.display = 'none';
            diffuseInput.style.display = 'none';
            normalInput.style.display = 'none';
        } else if (numShaderType === 3) {
            //create input showing diffuse texture
            object.material.shaderType = numShaderType;
            object.reset();
            state.render = true;
            diffuseNameTitle.style.display = 'inline';
            normalNameTitle.style.display = 'none';
            diffuseInput.style.display = 'inline';
            normalInput.style.display = 'none';
        } else if (numShaderType === 4) {
            object.material.shaderType = numShaderType;
            object.reset();
            state.render = true;
            diffuseNameTitle.style.display = 'inline';
            normalNameTitle.style.display = 'inline';
            diffuseInput.style.display = 'inline';
            normalInput.style.display = 'inline';
        }
    })
    let options = [{ text: "2D", value: 0 }, { text: "Flat Blinn", value: 1 }, { text: "Texture Blinn", value: 3 }, { text: "Normal & Diffuse", value: 4 }];
    options.forEach((opt) => {
        let tempOption = document.createElement("option");
        tempOption.value = opt.value;
        tempOption.innerHTML = opt.text;
        matTypeDropdown.appendChild(tempOption);
    });
    matTypeDropdown.value = object.material.shaderType;

    //diffuse color picker
    let diffuseTitle = document.createElement("h4");
    diffuseTitle.style.marginTop = "15px";
    diffuseTitle.innerHTML = "Diffuse Color";
    let diffuseColorPicker = document.createElement("input");
    diffuseColorPicker.type = "color";
    diffuseColorPicker.classList = "form-control";
    diffuseColorPicker.value = rgbToHex(object.material.diffuse);
    diffuseColorPicker.addEventListener('change', (event) => {
        let newColor = hexToRGB(event.target.value);
        object.material.diffuse = newColor;
        state.render = true;
    });

    mainDiv.appendChild(materialTitle);
    mainDiv.appendChild(matTypeDropdown);
    mainDiv.appendChild(diffuseTitle);
    mainDiv.appendChild(diffuseColorPicker);
    mainDiv.appendChild(texturesDiv);
}


/**
 * Purpose: Function creates a ui row for editing the scale of 
 * a selected scene object
 * @param {Scene object to be edited/displayed} object 
 * @param {State object containing all objects} state 
 */
function displayScaleValues(object, state) {
    let rotationalInputDiv = document.createElement("div");
    rotationalInputDiv.classList = "input-group";

    let prependDivX = document.createElement("div");
    prependDivX.classList = "input-group-prepend";

    prependDivX.innerHTML = `
        <span class="input-group-text">X</span>
        `;
    let prependDivY = prependDivX.cloneNode(true);
    prependDivY.innerHTML = `
        <span class="input-group-text">Y</span>
        `;

    let prependDivZ = prependDivX.cloneNode(true);
    prependDivZ.innerHTML = `
        <span class="input-group-text">Z</span>
        `;

    //X scale
    let objectScaleX = document.createElement("input");
    objectScaleX.type = "number";
    objectScaleX.classList = "form-control";
    objectScaleX.value = 0;
    objectScaleX.addEventListener('input', (event) => {
        if (event.target.value > 0) {
            //grow
            object.scale([1.5, 1.0, 1.0]);
        } else {
            //shrink
            object.scale([0.5, 1.0, 1.0]);
        }
        objectScaleX.value = 0;
        state.render = true;
    })

    //Y scale
    let objectScaleY = document.createElement("input");
    objectScaleY.type = "number";
    objectScaleY.classList = "form-control";
    objectScaleY.value = 0;
    objectScaleY.addEventListener('input', (event) => {
        if (event.target.value > 0) {
            //grow
            object.scale([1.0, 1.5, 1.0]);
        } else {
            //shrink
            object.scale([1.0, 0.5, 1.0]);
        }
        objectScaleX.value = 0;
        state.render = true;
    })

    //Z scale
    let objectScaleZ = document.createElement("input");
    objectScaleZ.type = "number";
    objectScaleZ.classList = "form-control";
    objectScaleZ.value = 0;
    objectScaleZ.addEventListener('input', (event) => {
        if (event.target.value > 0) {
            //grow
            object.scale([1.0, 1.0, 1.5]);
        } else {
            //shrink
            object.scale([1.0, 1.0, 0.5]);
        }
        objectScaleX.value = 0;
        state.render = true;
    })

    rotationalInputDiv.appendChild(prependDivX);
    rotationalInputDiv.appendChild(objectScaleX);
    rotationalInputDiv.appendChild(prependDivY);
    rotationalInputDiv.appendChild(objectScaleY);
    rotationalInputDiv.appendChild(prependDivZ);
    rotationalInputDiv.appendChild(objectScaleZ);

    return rotationalInputDiv;
}


/**
 * Purpose: Function creates a ui row for editing the rotation of 
 * a selected scene object
 * @param {Scene object to be edited/displayed} object 
 * @param {State object containing all objects} state 
 */
function displayRotationValues(object, state) {
    let rotationalInputDiv = document.createElement("div");
    rotationalInputDiv.classList = "input-group";

    let prependDivX = document.createElement("div");
    prependDivX.classList = "input-group-prepend";

    prependDivX.innerHTML = `
        <span class="input-group-text">X</span>
        `;
    let prependDivY = prependDivX.cloneNode(true);
    prependDivY.innerHTML = `
        <span class="input-group-text">Y</span>
        `;

    let prependDivZ = prependDivX.cloneNode(true);
    prependDivZ.innerHTML = `
        <span class="input-group-text">Z</span>
        `;

    //X rotation
    let objectRotationX = document.createElement("input");
    objectRotationX.type = "number";
    objectRotationX.classList = "form-control";
    objectRotationX.value = 0;
    objectRotationX.addEventListener('input', (event) => {
        mat4.rotateX(object.model.rotation, object.model.rotation, 1)
        state.render = true;
    })

    //Y rotation
    let objectRotationY = document.createElement("input");
    objectRotationY.type = "number";
    objectRotationY.classList = "form-control";
    objectRotationY.value = 0;
    objectRotationY.addEventListener('input', (event) => {
        mat4.rotateY(object.model.rotation, object.model.rotation, 1)
        state.render = true;
    })

    //Z rotation
    let objectRotationZ = document.createElement("input");
    objectRotationZ.type = "number";
    objectRotationZ.classList = "form-control";
    objectRotationZ.value = 0;
    objectRotationZ.addEventListener('input', (event) => {
        mat4.rotateZ(object.model.rotation, object.model.rotation, 1)
        state.render = true;
    })

    rotationalInputDiv.appendChild(prependDivX);
    rotationalInputDiv.appendChild(objectRotationX);
    rotationalInputDiv.appendChild(prependDivY);
    rotationalInputDiv.appendChild(objectRotationY);
    rotationalInputDiv.appendChild(prependDivZ);
    rotationalInputDiv.appendChild(objectRotationZ);

    return rotationalInputDiv;
}

function createHeader(text, size) {
    let tempTitle = document.createElement(size);
    tempTitle.innerHTML = text;
    return tempTitle;
}


function shaderValuesErrorCheck(programInfo) {
    let missing = [];
    //do attrib check
    Object.keys(programInfo.attribLocations).map((attrib) => {
        if (programInfo.attribLocations[attrib] === -1) {
            missing.push(attrib);
        }
    });
    //do uniform check
    Object.keys(programInfo.uniformLocations).map((attrib) => {
        if (!programInfo.uniformLocations[attrib]) {
            missing.push(attrib);
        }
    });

    if (missing.length > 0) {
        printError('Shader Location Error', 'One or more of the uniform and attribute variables in the shaders could not be located or is not being used : ' + missing);
    }
}

/**
 * A custom error function. The tag with id `webglError` must be present
 * @param  {string} tag Main description
 * @param  {string} errorStr Detailed description
 */
function printError(tag, errorStr) {
    // Create a HTML tag to display to the user
    var errorTag = document.createElement('div');
    errorTag.classList = 'alert alert-danger';
    errorTag.innerHTML = '<strong>' + tag + '</strong><p>' + errorStr + '</p>';

    // Insert the tag into the HMTL document
    document.getElementById('webglError').innerHTML = errorTag;

    // Print to the console as well
    console.error(tag + ": " + errorStr);
}

let UI = {
    createSceneGui,
    shaderValuesErrorCheck,
    setup
}

export default UI;