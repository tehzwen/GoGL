import state from './main.js';

function setup(newState) {
    console.log(newState)
    //listeners for header button presses
    document.getElementById('saveButton').addEventListener('click', () => {
        if (!document.location.href.includes("editor")) {
            createSceneFile(newState, "./statefiles/" + newState.saveFile);
        }
    })

    document.getElementById('launchButton').addEventListener('click', () => {
        //
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
                    displayObjectValues(null);
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

    //camera stuff here TODO
    /*
        let camera = state.camera;
        let objectElement = document.createElement("div");
        let objectName = document.createElement("h5");
        objectName.classList = "object-link";
        objectName.innerHTML = camera.name;
    
        objectName.addEventListener('click', () => {
            let objectModel = {
                model: { ...camera },
                name: camera.name
            }
    
            displayObjectValues(objectModel);
        });
    
        objectElement.appendChild(objectName);
        sideNav.appendChild(objectElement); */

    /*
    let addNav = document.getElementById("addObjectsNav");
    addNav.innerHTML = "";
    let objectTypeSelect = document.createElement("select");
    objectTypeSelect.classList = "form-control";
    objectTypeSelect.addEventListener('change', (event) => {
        handleTypeSelectChange(event);
    })

    createSelectionOptions(["Cube", "Mesh"], objectTypeSelect);

    let addNewButton = document.createElement("button");
    addNewButton.innerHTML = "New Object";
    addNewButton.classList = "btn btn-primary";
    addNewButton.addEventListener('click', () => {
        addObject(objectTypeSelect.value);
    });

    addNav.appendChild(objectTypeSelect);
    addNav.appendChild(addNewButton); */

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
    let position;
    if (object.type !== "directionalLight") {
        position = object.model.position;
    } else {
        position = object.position;
    }

    if (!object) {
        let selectedObjectDiv = document.getElementById("selectedObject");
        selectedObjectDiv.innerHTML = "";
        return;
    }

    let selectedObjectDiv = document.getElementById("selectedObject");
    selectedObjectDiv.innerHTML = "";

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


    //diffuse color picker
    let diffuseColorPicker = document.createElement("input");
    diffuseColorPicker.type = "color";
    diffuseColorPicker.classList = "form-control";
    diffuseColorPicker.value = "#ffffff";
    diffuseColorPicker.addEventListener('change', (event) => {
        let newColor = hexToRGB(event.target.value);
        object.material.diffuse = newColor;
        state.render = true;
    });

    //add all the elements in
    positionalInputDiv.appendChild(prependDivX);
    positionalInputDiv.appendChild(objectPositionX);
    positionalInputDiv.appendChild(prependDivY);
    positionalInputDiv.appendChild(objectPositionY);
    positionalInputDiv.appendChild(prependDivZ);
    positionalInputDiv.appendChild(objectPositionZ);
    selectedObjectDiv.appendChild(createHeader(`<i>${object.name}</i>`, "h3"));
    selectedObjectDiv.appendChild(createHeader("Position", "h4"));
    selectedObjectDiv.appendChild(positionalInputDiv);
    selectedObjectDiv.appendChild(createHeader("Rotation", "h4"));
    selectedObjectDiv.appendChild(displayRotationValues(object, state))


    if (object.type !== "light") {
        let diffuseTitle = document.createElement("h4");
        diffuseTitle.innerHTML = "Diffuse Color";
        selectedObjectDiv.appendChild(diffuseTitle);
        selectedObjectDiv.appendChild(diffuseColorPicker);
    } else {
        //for light, we want to change its color not diffuse material
    }
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