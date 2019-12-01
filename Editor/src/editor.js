const fs = require('fs')
const gameFilesPath = __dirname + "/../Renderer/game/"
const exec = require('child_process').exec;

window.onload = () => {
    state = {};

    var editor = ace.edit("editorWindow", {
        theme: "ace/theme/monokai",
        mode: "ace/mode/golang",
        autoScrollEditorIntoView: true,
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

    fs.readdir(gameFilesPath, (err, items) => {
        if (!err) {
            let fileNav = document.getElementById("fileTabs");

            setEditorTextFromFile(gameFilesPath + items[0], editor);
            state.currentFile = gameFilesPath + items[0];

            items.map((item) => {
                let fileTab = document.createElement('li');
                let fileButton = document.createElement('button');
                fileButton.classList = "btn btn-outline-primary btn-block btn-lg"
                fileButton.id = item;
                fileButton.innerHTML = item;
                fileButton.style = `text-transform: initial`;

                fileButton.addEventListener('click', (e) => {
                    setEditorTextFromFile(gameFilesPath + e.target.id, editor);
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