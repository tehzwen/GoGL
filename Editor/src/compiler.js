var exec = require('child_process').exec;

window.onload = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const sceneFile = urlParams.get('scene');    
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
        console.warn(code);

        runCommand = exec("../Renderer/GoGL " + sceneFile, (err, stdout, stderr) => {
            console.log(stdout);
        })
        runCommand.on('exit', (eCode) => {
            console.warn(eCode);
            window.location.href = "main.html"
        })
    });
}
