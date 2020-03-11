var exec = require('child_process').exec;
var spawn = require('child_process').spawn;
var path = require('path');

window.onload = () => {
    const urlParams = new URLSearchParams(window.location.search);
    let sceneFile = urlParams.get('scene');
    const compilationText = document.getElementById("compilationOutput")
    compilationText.style.height = screen.height - 200 + 'px';
    console.log(sceneFile);

    if (window.navigator.userAgent.indexOf("Windows") !== -1) {
        sceneFile = sceneFile.split("/").join("\\");
        psCommand = spawn("powershell.exe", ['./src/windowsBuild.ps1 "' + sceneFile + '"']);
        psCommand.stdout.on("data", (data) => {
            let splitData = data.toString().split('\n');
            splitData.pop() //take off the newline character
            splitData.map((data) => {
                let tempText = document.createElement('p');
                if (tempText) {
                    let currTime = new Date();
                    tempText.innerHTML = currTime.toLocaleTimeString() + "<b> " + data + "</b>";
                    tempText.classList = "orange-text";
                    compilationText.appendChild(tempText)
                    compilationText.scrollTop = compilationText.scrollHeight; //causes output to auto scroll down
                }
            });
        })

        psCommand.stderr.on("data", (data) => {
            let splitData = data.toString().split('\n');
            splitData.pop() //take off the newline character
            splitData.map((data) => {
                let tempText = document.createElement('p');
                if (tempText) {
                    let currTime = new Date();
                    tempText.innerHTML = currTime.toLocaleTimeString() + "<b> " + data + "</b>";
                    tempText.style.color = 'red';
                    compilationText.appendChild(tempText)
                    compilationText.scrollTop = compilationText.scrollHeight; //causes output to auto scroll down
                }
            });
        })
        psCommand.on("exit", function () {
            console.log("Build script finished");
        });
        psCommand.stdin.end();
    } else {
        buildCommand = exec("cd ../Renderer/;go build", function (err, stdout, stderr) {
            if (err) {
                // should have err.code here?  
                console.error(err);
            }
            console.log(stdout);
        })

        buildCommand.on('exit', function (code) {
            // exit code is code
            if (code === 0) {
                let tempText = document.createElement('p');
                tempText.classList = "orange-text";
                let currTime = new Date();
                tempText.innerHTML = currTime.toLocaleTimeString() + "<b> Compilation successful!</b>";
                compilationText.appendChild(tempText)
            } else {
                let tempText = document.createElement('p');
                tempText.classList = "orange-text";
                let currTime = new Date();
                tempText.innerHTML = currTime.toLocaleTimeString() + "<b> Compilation failed!</b>";
                compilationText.appendChild(tempText)
            }

            runCommand = exec("cd ../Renderer/ ; ./Renderer " + sceneFile, (err, stdout, stderr) => {
                if (stderr) {
                    console.error(stderr);
                }
            })
            runCommand.stdout.on('data', (data) => {
                let splitData = data.split('\n');
                splitData.pop() //take off the newline character
                splitData.map((data) => {
                    let tempText = document.createElement('p');
                    if (tempText) {
                        let currTime = new Date();
                        tempText.innerHTML = currTime.toLocaleTimeString() + "<b> " + data + "</b>";
                        tempText.classList = "orange-text";
                        compilationText.appendChild(tempText)
                        compilationText.scrollTop = compilationText.scrollHeight; //causes output to auto scroll down
                    }
                });
            })

            runCommand.on('exit', (eCode) => {
                console.warn(eCode);
                if (eCode === 0) {
                    //window.location.href = "main.html"
                }
            })
        });
    }
}
