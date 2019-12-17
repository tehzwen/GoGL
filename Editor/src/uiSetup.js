import state from './main.js';

setup();

var saveFile = "testsave.json"

function setup() {
    //listeners for header button presses
    document.getElementById('saveButton').addEventListener('click', () => {
        if (!document.location.href.includes("editor")) {
            createSceneFile(state, "./statefiles/" + saveFile);
        }
    })

    document.getElementById('launchButton').addEventListener('click', () => {
        document.location.href = "compiler.html?scene=" + __dirname + "/statefiles/" + saveFile
    })
}

function createSceneGui(state) {
    //get objects first
    let sideNav = document.getElementById("objectsNav");
    sideNav.innerHTML = "";

    state.objects.map((object) => {
        let objectElement = document.createElement("div");
        let objectName = document.createElement("h5");
        objectName.classList = "object-link";
        objectName.innerHTML = object.name;
        objectName.addEventListener('click', () => {
            displayObjectValues(object);
        });

        objectElement.appendChild(objectName);
        sideNav.appendChild(objectElement);
    });

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
    sideNav.appendChild(objectElement);

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
function displayObjectValues(object) {
    let selectedObjectDiv = document.getElementById("selectedObject");
    selectedObjectDiv.innerHTML = "";

    let positionalInputDiv = document.createElement("div");
    positionalInputDiv.classList = "input-group";

    let prependDivX = document.createElement("div");
    prependDivX.classList = "input-group-prepend";

    //X move input handler
    let objectPositionX = document.createElement("input");
    objectPositionX.type = "number";
    objectPositionX.addEventListener('input', (event) => {
        object.translate([event.target.value - object.model.position[0], 0, 0]);
    })
    objectPositionX.id = object.name + "-positionX";
    objectPositionX.classList = "form-control";
    objectPositionX.value = parseFloat(object.model.position[0]).toFixed(1);

    //Y move input handler
    let objectPositionY = document.createElement("input");
    objectPositionY.type = "number";
    objectPositionY.addEventListener('input', (event) => {
        object.translate([0, event.target.value - object.model.position[1], 0]);
    })
    objectPositionY.id = object.name + "-positionY";
    objectPositionY.classList = "form-control";
    objectPositionY.value = parseFloat(object.model.position[1]).toFixed(1);

    //Z move input handler
    let objectPositionZ = document.createElement("input");
    objectPositionZ.type = "number";
    objectPositionZ.addEventListener('input', (event) => {
        object.translate([0, 0, event.target.value - object.model.position[2]])
    })
    objectPositionZ.id = object.name + "-positionZ";
    objectPositionZ.classList = "form-control";
    objectPositionZ.value = parseFloat(object.model.position[2]).toFixed(1);

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


    if (object.type !== "light") {
        let diffuseTitle = document.createElement("h4");
        diffuseTitle.innerHTML = "Diffuse Color";
        selectedObjectDiv.appendChild(diffuseTitle);
        selectedObjectDiv.appendChild(diffuseColorPicker);
    } else {
        //for light, we want to change its color not diffuse material
    }
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
    document.getElementById('webglError').appendChild(errorTag);

    // Print to the console as well
    console.error(tag + ": " + errorStr);
}

let UI = {
    createSceneGui,
    shaderValuesErrorCheck
}

export default UI;