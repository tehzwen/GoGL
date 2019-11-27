var exec = require('child_process').exec;

window.onload = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const sceneFile = urlParams.get('scene');    
    const compilationText = document.getElementById("compilationOutput")
    console.log(sceneFile);

    buildCommand = exec("make -C ../Renderer/", function (err, stdout, stderr) {
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
            let currTime = new Date();
            tempText.innerHTML = currTime.toLocaleTimeString() + "<b> Compilation successful!</b>";
            compilationText.appendChild(tempText)
        } else {
            let tempText = document.createElement('p');
            let currTime = new Date();
            tempText.innerHTML = currTime.toLocaleTimeString() + "<b> Compilation failed!</b>";
            compilationText.appendChild(tempText)
        }

        runCommand = exec("../Renderer/GoGL " + sceneFile, (err, stdout, stderr) => {
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
                    tempText.innerHTML = currTime.toLocaleTimeString() + "<b> " +data + "</b>";
                    compilationText.appendChild(tempText)
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
