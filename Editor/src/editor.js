const fs = require('fs')
const gameFilesPath = __dirname + "/../Renderer/game/"
const exec = require('child_process').exec;

window.onload = () => {
    state = {};
    //set the size of the parent div to the client height
    document.getElementById('editorParent').style.height = document.body.clientHeight * 1.5 + "vh";

    var editor = ace.edit("editorWindow", {
        theme: "ace/theme/tomorrow_night",
        mode: "ace/mode/golang",
        maxLines: 200,
        minLines: 30
    });
    editor.commands.addCommand({
        name: 'saveCommand',
        bindKey: { win: "Ctrl-S", mac: "Command-S", linux: "Ctrl-S" },
        exec: function () {
            saveCurrentFile(state, editor)
        }
    })
    editor.setFontSize(16);
    editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true
    })

    state.editor = editor;

    document.getElementById('saveButton').addEventListener('click', () => {
        saveCurrentFile(state, editor);
    })

    document.getElementById('launchButton').addEventListener('click', () => {
        if (document.location.search) {
            document.location.href = "compiler.html" + document.location.search;
        } else {
            document.location.href = "compiler.html";
        }
    })

    fs.readdir(gameFilesPath, {withFileTypes: true}, (err, items) => {
        if (!err) {
            const files = items
            .filter(file => file.isFile())
            .map(file => file.name)
            //check each item if its a directory or not

            let fileNav = document.getElementById("fileTabs");

            setEditorTextFromFile(gameFilesPath + files[0]);
            state.currentFile = gameFilesPath + files[0];

            files.forEach((item) => {
                let fileTab = document.createElement('li');
                let fileButton = document.createElement('button');
                fileButton.classList = "btn btn-outline-primary btn-block btn-lg"
                fileButton.id = item;
                fileButton.innerHTML = item;
                fileButton.style = `text-transform: initial`;

                fileButton.addEventListener('click', (e) => {
                    setEditorTextFromFile(gameFilesPath + e.target.id);
                    state.currentFile = gameFilesPath + e.target.id;
                })
                fileTab.appendChild(fileButton);
                fileNav.appendChild(fileTab);
            })
        } else {
            console.error(err);
        }
    })
}

function setEditorTextFromFile(filePath) {
    fs.readFile(filePath, 'utf-8', (err, data) => {
        if (!err) {
            state.editor.setValue(data, -1);
        } else {
            console.error(err);
        }
    })
}

function saveCurrentFile(state) {
    //I will want to run go.fmt on this file first to make sure it has no errors
    fs.writeFile(state.currentFile, state.editor.getValue(), (err) => {
        if (err) {
            console.error(err);
        } else {
            checkFileSyntax(state);
            console.log("writing of ", state.currentFile, " successful!");
        }
    })
}

function checkFileSyntax(state) {
    file = state.currentFile.trim()
    let importCommand = "goimports -w " + state.currentFile.trim();

    importsCommandExec = exec(importCommand, (err, stdout, stderr) => {
        if (err) {
            console.error(err);
        }
    })

    let command = "make -C ../Renderer/";
    syntaxCommand = exec(command, (err, stdout, stderr) => {
        if (err) {
            console.error(err);
        }
    })

    syntaxCommand.on('exit', (code) => {
        if (code === 0) {
            console.log("GOOD")
            setEditorTextFromFile(state.currentFile);
        }
    })
}