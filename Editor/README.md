# WebGLEngine
Simple WebGL rendering engine used for simple playgrounding with webgl

***This is a desktop version of the engine***

To run:

run command ```npm i``` in the root of this project
run command ```npm start``` to start the app

## Design
### <u>Objects</u>
The main paradigm of the engine is that of objects. Using different classes, the engine performs different actions based off of the type of class the object is but the objects are all similarly rendered in the scene. 

For example, if an object is loaded in that is of type Cube, the engine will automatically create a Cube object, assign in the values that are for the cube and add it into the scene.

Each object can have its own vertex and fragment shader depending on how you'd like to render it in the scene. For the sake of the current project, each object is using a single vertex and single fragment shader that I've written.

All objects are stored in a JSON file under /statefiles/*file*.json and are loaded in by the scene from the file. By looking at the "alienScene.json" example scene file you can see a structure like 
```
{
    "objects" : [
            {
                "name" : "alien",
                "material" : {"diffuse" : [0.2, 0.2, 0.2], "ambient": [0.001, 0.001, 0.001], "specular": [0.8, 0.8, 0.8], "n":12, "alpha":1.0},
                "type": "mesh",
                "model" : "./models/alien.obj",
                "parent" : null,
                "position": [1.0, 0.0, 1.0],
                "scale": [1.0, 1.0, 1.0]
            },
            {
                "name" : "wall0",
                "material" : {"diffuse" : [0.2, 0.2, 0.2], "ambient": [0.1, 0.1, 0.1], "specular": [0.8, 0.8, 0.8], "n":10, "alpha":1.0},
                "type": "cube",
                "texture": "./materials/plywood.jpg",
                "textureNorm": "./materials/checkerNorm.jpg",
                "parent" : null,
                "position": [2.0, -1.0, -2],
                "scale": [1.0, 4.0, 10.0]
            },
    ]
}
```

From the example above you can see 2 types of objects, one being a mesh (Model class) and the other being a cube (Cube class). All values are loaded into the object class they belong to and then rendered from there. It's important to note that the texture files belong in `/materials/materialName ` and the model files belong in `/models/modelName`. 

The meshes used in the Model class rely on the **three-object-loader.js** in `/lib/` which is a modified model loader that was used by three.js. 

### <u>Scene</u>

The scene consists of a main state object that keeps a reference to all current settings, objects and rendering options for the entire scene. As of right now the state currently holds information on the camera, keyboard, mouse, objectTable, lightIndices, objectCount and objects. 

### <u>Lights</u>

Lights are shown using a simple lightbulb mesh I found so you can see where it appears in the scene. They are loaded from the JSON scene file as mentioned above and look something like this:

```
{
    "objects" : [
        {
            "name" : "apple",
            "material" : {"diffuse" : [0.8, 0.2, 0.2], "ambient": [0.001, 0.001, 0.001], "specular": [0.8, 0.8, 0.8], "n":50, "alpha":1.0},
            "type": "mesh",
            "model" : "./models/apple.obj",
            "parent" : null,
            "position": [0.0, 0.0, 1.0],
            "scale": [2.0, 2.0, 2.0]
        },
        {
            "name": "light0",
            "type": "light",
            "model" : "./models/lightbulb.obj",
            "position" : [-1.45, 0.8, -1.5],
            "material" : {"diffuse" : [1.0, 1.0, 1.0], "ambient": [0.1, 0.1, 0.1], "specular": [0.3, 0.3, 0.3], "n":1000, "alpha":1.0},
            "colour" : [1.0, 1.0, 1.0],
            "parent" : null,
            "strength" : 0.75
        },
        {
            "name": "light1",
            "type": "light",
            "model" : "./models/lightbulb.obj",
            "position" : [1.9, 0.8, 2.5],
            "material" : {"diffuse" : [1.0, 1.0, 1.0], "ambient": [0.1, 0.1, 0.1], "specular": [0.3, 0.3, 0.3], "n":1000, "alpha":1.0},
            "colour" : [1.0, 1.0, 1.0],
            "parent" : null,
            "strength" : 0.75
        }
    ]
}

```
Similar to the example above, this will load the apple object, place it in the scene and load two lights of the same values in different positions of the scene.

### <u>Game</u>

For now, the engine has no actual "game" and rather has a simple example of where a game might be called that allows the camera to move about the scene with WASD keys and look horizontally with right click + mouse move. 

```
function render(now) {
        stats.begin();
        now *= 0.001; // convert to seconds
        const deltaTime = now - then;
        then = now;

        state.deltaTime = deltaTime;

        //wait until the scene is completely loaded to render it
        if (state.numberOfObjectsToLoad <= state.objects.length) {
            if (!state.gameStarted) {
                startGame(state);
                state.gameStarted = true;
            }
```

This example calls a start game function that I've defined in `/src/myGame.js` and sets up the controls for us to use. 

## Examples

### Non-textured obj model, lighting, textured primitives
![Alt text](./samples/bunny.png?raw=true)

### Multiple obj models, lighting, normal bump mapping on primitives (plane & wall)
![Alt text](./samples/alienscene.png?raw=true)

### Multiple obj models, one textured, the other flat shading. 
![Alt text](./samples/crate.png?raw=true)

### Camera movement, object rotation around a parent, gui example.
![Alt text](./samples/sample.gif?raw=true)
