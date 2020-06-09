Introduction
============

The project focuses on creating a cross platform, highly customized
rendering engine which includes a front end editor for scene rendering
and game development. Technologies used include OpenGL and Golang in the
rendering engine and WebGL, Javascript and Electron in the editor. The
editor provides scene management as well as a scripting API for
controlling the scene which passes data to the engine to then be run in
the rendering engine.

Related Works and Motivation
============================

Unity
-----

Much of the overall editor design used in the project was inspired a lot
by Unity. Offering a pane for scene objects, a pane for object specific
values and then a large rendering window was very helpful. Unity focuses
on attaching scripts to objects and having them inherit properties using
an object oriented approach.

Three.js
--------

Three.js is a web framework that uses WebGL and abstracts away much of
the difficulties of working with graphics in the browser. After using
Three.js, I felt that the way it used classes to categorize primitive
types (squares, planes, etc.) worked really well and so I followed the
same paradigm in the engine and the editor.

Godot Engine
------------

Godot is an open source game engine that focuses mostly on 2D rendering
applications and does a good job of them. The engine focuses on a
hierarchical scheme for organizing scenes and objects with parents
passing properties to children and so on. For the purpose of handling
transforms, texturing and organization I chose to follow this idea and
organize my scene using parent and children type relationships.

Motivation
----------

Having taken CMPT 370 and learning more about graphics in WebGL, I
became interested in how a rendering engine written in a compiled
language would perform. Golang was of interest as it handles concurrency
and cross platform compilation quite nicely. The editor was something I
originally wanted to write in Golang as well but since I already had
done work on a WebGL engine in the browser I decided to build on top of
that work and create a desktop version using Electron.

Methodology
===========

![[\[fig:1\]]{#fig:1 label="fig:1"}System
Architecture](diagram.png)

As seen in figure 1, the overall system works by sharing data between
the editor and the engine. The editor passes scene related data which
contains definitions for properties of objects in the scene including
transforms, materials and hierarchy. Included in the scene data is
information on lighting, scene and camera settings. The editor also
passes information to the engine in the form of Golang code which
accesses an internal API written in the engine which allows the code to
gain access to objects in the scene and their properties. The libraries
used for the editor are Electron and gl-matrix, the latter being used
for matrix and vector operations.

The editor as described before is run in Electron and WebGL which both
use JavaScript. JavaScript is a loosely typed interpreted scripting
language used mostly for web development. The benefits of such a
language is that data can be passed around very easily without declaring
types or strict rules for each data type. The disadvantages to such a
language is that by not having strict types, errors can occur and can be
quite difficult to debug. JavaScript interacts with the WebGL API to
provide it with attributes, buffers and uniform values. These updates
are done very easily with WebGL and I found that there was ample
documentation for it online.

The rendering engine runs on Golang, a young, strictly typed, compiled
language created by Google. I chose Go for a couple of reasons including
the fact that there weren't many projects that used it and did graphics
rendering and because of Go's concurrency model which makes it quite
easy using 'goroutines'. Similar to the editor, the engine passes
attribute, buffer and uniform values to the OpenGL API. Performing the
OpenGL operations in the engine proved to be somewhat difficult as the
documentation for OpenGL in Go is not very extensive and doesn't show
many examples on sample use. In addition, because the engine needs to
render different types of objects in similar ways I decided to use an
object oriented approach where the engine could call similar methods on
'classes' for rendering purposes. Since Go is not technically an object
oriented language this took a little bit of figuring where I used
interfaces and structs to accomplish similar things.

![[\[fig:2\]]{#fig:2 label="fig:2"}OpenGL
Pipeline](pipeline-v2.png)

In order for OpenGL and WebGL to both work they require shader programs
to be written for them in GLSL (OpenGL Shader Language). The language is
a high level shader language with syntax very similar to the C
programming language with the exception of built-in matrix and vector
operations. The shaders can have input variables, uniform variables and
output variables. The usual setup is to have a vertex shader program
then output values which are handled in the fragment shader program
which then displays the output in the form of pixels to the screen. The
uniforms and attributes listed in the programs require those values to
be updated and set by the programming language that is controlling these
programs so in my case Go and JS. The order in which the programs run
and values are provided to the shader programs is visualized in figure
2, the engine and editor both follow the same order and run their
respective shader programs based off the same design. It is worth noting
that there are many ways to load shader programs my editor loads the
shader programs in the form of JSON stored strings and the engine loads
its shaders in the form of structs that contain string fields.

Theory
======

In this section I will provide an overview of the theory used in the
development of the engine for each feature.

Basic 3D Rendering
------------------

The basics of rendering a scene in 3D for my engine required the setup
of a vertex shader, a fragment shader and the object to be rendered. The
object contains information including the positions, normals, and
tangents of each vertex for the object. This information is passed to
the shader programs in the form of array buffers per vertex. Included
with these values are object specific values for their materials to be
rendered.

In the vertex shader we perform projection to take our 3D scene and
project it for viewing on our 2D screen. We use a commonly referred to
operation called the MVP matrix. $$M * V * P$$ M refers to the model
matrix which contains the transformations for the specific object being
rendered. V refers to the view matrix which is constructed for the
camera or viewing of the scene. The P refers to the projection matrix
which represents the projection of the scene, there are two usual ways
of doing this: orthographic and perspective. In the case of my engine
and editor, I chose to use a perspective projection.

Shading
-------

For the purpose of efficiency, the engine relies on many different types
of shader programs to perform different types of shading, however for
the purpose of the paper I will encapsulate the basics of them all here.

The shaders rely on the Blinn-Phong lighting model with some changes
depending on whether the object in question contains textures to be
rendered, reflection or refraction. Each object to be rendered contains
a material that details the ambient, diffuse and specular colour for the
object as-well as the shininess value, and any textures to be rendered.
$\text{Ambient: }A = Ka * Lk$\
$\text{Diffuse: }D = Kd * Lk * (N • L)$\
$\text{Specular: }S = Ks * Lk * (N • (V + L))^n$\
$\text{Overall: }A * D * S$\
This formula is used for every light in the scene to calculate a total
shading value. It is worth noting that in the case of my shader programs
the following code caused issues with multiple lights. The issues
included inconsistent colours being cast by each light, sometimes
certain lights not being rendered or shaded at all.

    for (int i = 0; i < numLights; i++) {
        vec3 lightDirection = normalize(lights[i].position - oFragPosition);
        vec3 ambient = ambientValue * lights[i].colour;
        float NdotL = max(dot(oNormal, lightDirection), 1.0);
        vec3 diffuse = (diffuseValue * lights[i].colour) * NdotL;
        float NDotH = max(dot(oNormal, H), 0.0);
        float NHPow = pow(NDotH, nVal);
        vec3 specular = (specularVal * lights[i].colour) * NHPow;
        vec4 fragColour = vec4(ambient * diffuse * specular, 1.0);
        result += fragColour;
    }

The problems that occured with that code was solved by creating
functions to do essentially the same math but calling them inside the
iteration similar to the following.

    vec4 CalcLighting(light Pointlight) {
        vec3 lightDirection = normalize(light.position - oFragPosition);
        vec3 ambient = ambientValue * light.colour;
        float NdotL = max(dot(oNormal, lightDirection), 1.0);
        vec3 diffuse = (diffuseValue * light.colour) * NdotL;
        float NDotH = max(dot(oNormal, H), 0.0);
        float NHPow = pow(NDotH, nVal);
        vec3 specular = (specularVal * light.colour) * NHPow;
        vec4 fragColour = vec4(ambient * diffuse * specular, 1.0);
        return fragColour;
    }


    for (int i = 0; i < numLights; i++) {
        result += CalcLighting(lights[i]);
    }

Mesh & Material Loading
-----------------------

mtllib alien.mtl\
o alien\
v 0.066307 -0.710018 0.485289\
v 0.108720 -0.575360 0.535121\
\...\
vn -0.4646 0.6778 0.5698\
vn -0.4060 0.7724 0.4883\
\...\
f 15917//15916 11954//11953 15918//15917\
f 22936//22935 66340//66336 41829//41827\

As part of the different objects that are able to be rendered by the
engine, I wanted to include complex 3D meshes. I first attempted to use
a C library called Assimp for loading mesh information which was fast
but I often ended up with broken models and incomplete vertex data. I
ended up writing my own parser that uses regular expressions to gather
information on each mesh object and its corresponding material values.
More complex meshes contain definitions and values for multiple objects,
to account for this I created a master vertex list and indexes into it
for each material so that each child object has the correct material
associated with it. I found that although my own parser was tedious to
write and loads meshes much slower than Assimp did, the accuracy of my
mesh loader out performed assimp and so I stuck with it.

newmtl DefaultOBJ\
Ns 225.000000\
Ka 1.000000 1.000000 1.000000\
Kd 0.800000 0.800000 0.800000\
Ks 0.500000 0.500000 0.500000\
d 1.000000\
illum 2\
map\_Kd Dish.jpg\
map\_Bump DishNormal.jpg\

The mesh loader currently supports obj files only and these files come
with a material file or .mtl file. These files as seen above contain
information (in order) of the material name, shininess, ambient,
diffuse, specular, alpha, illumination model, diffuse texture and
normal/bump texture. It's important to note that different materials can
contain more or less information than that listed above. I wrote a
simple parser to look for these values and apply them to the indices I
calculate above when loading a mesh in.

Shadows
-------

The main type of shadows present in the engine are that of point-light
shadows which cast shadows in all directions. In order for this to work,
the engine has an array of all point-lights present in the scene and
before rendering any colours we render the scene's depth to a 3D texture
and store that with the point-light. When the coloured objects are then
rendered, we have the information of these textures stored with the
point-lights and compare the depth from the current pixel to the light
in question against the depth of the initial depth render from this
light to determine if the pixel should be in shadow or not.

![[\[fig:1\]]{#fig:1 label="fig:1"}+X Direction
Depth](depth1.jpg)

![[\[fig:2\]]{#fig:2 label="fig:2"}-X Direction
Depth](depth2.jpg)

It's important to note that the performance of multiple lights
constantly rendering depth to each texture is quite costly and as such I
made improvements by not re-rendering certain depths if lights had not
moved or objects had not moved into their sphere of influence. This
technique can be done by rendering depth using the typical vertex and
fragment shader, however my solution sped things up and simplified
things by writing a geometry shader and completing the six directional
render writes in the geometry shader instead of having to iterate six
times.

Reflection & Refraction
-----------------------

![[\[fig:2\]]{#fig:2 label="fig:2"}Reflection &
Refraction](reflections.png)

The technique I used for simple reflections and refraction uses a
cube-map texture. A cube-map is a 3D texture created by six different 2D
textures to give the illusion of a sky-box like surrounding. The scene
is rendered as normal, and at the end the scene is then rendered with
the cube-map as a background that is infinitely far away. The objects in
the scene can access the colours of the cube-map very efficiently in
comparison to something like ray tracing. We can then reflect or refract
the color of the cube-map depending on what the object is specified to
do. Each object would have a specific reflection mode attached to the
material (0 - don't reflect, 1 - reflect and 2 - refract) and also a
refraction index value if the object were of the refraction type.

System Details
==============

Efficiencies
------------

In order to make the editing and running experience as pleasant as
possible several efficiencies were put in place by me to enable low
CPU/GPU usage while using the application to make nice looking graphics.

### Frustum Culling

I used the idea of frustum culling to speed up the overall render speed
and lower the system resources being used by both the editor and the
engine. The basic idea of this is to not render anything that the camera
(or viewing frustum) cannot see.

![[\[fig:1\]]{#fig:1 label="fig:1"}Example Viewing
Frustum](frustum.png)

The planes of the frustum are then used in a dot product with the object
in question to determine whether or not the object is inside the plane
or outside the plane. If the object is within the planes then we render
it, otherwise we return early and do not render the object. This
technique has an upfront cost to perform this check but overall I found
it to be more efficient than simply rendering every object without
checking.

### Lazy Rendering

When developing the editor I was using my small laptop without a
graphics card and ran into the case of where the editor simply took up
too many system resources to keep open for too long. To solve this
problem I came up with the idea of lazy rendering. Essentially we do not
render objects or re-render our canvas unless something has changed in
the scene (ie. adding an object, transforming an object, light changes,
camera movement). This resulted in a low footprint editor that ran nice
and smooth on all my devices.

### Mesh Caching

As the engine started to be able to load meshes I experimented with more
and more complex meshes to see how it would perform. As the complexity
of meshes increase so does the amount of time it takes to load larger
and larger obj files. As a result if I loaded in a large mesh it would
take considerable time to reload the same scene if I made any changes to
it. I decided to take the mesh details including the vertices, normals
and other data and compress it then save it in a binary format. If an
object were to be loaded in, it would check for a cached version of
itself first and load from there instead of the obj file. This sped up
loading complex meshes considerably but also comes at the cost of disk
space and a slower loading speed the first time through as it needs to
write the compressed file.

Debugging
---------

The engine system has error messages and tests in place so that in the
event of an error the editor will catch it and display the error for the
user without crashing the entire system. In addition to this I found
that using a program called RenderDoc was highly helpful as it would
provide detailed information on the values your program was currently
running in the OpenGL pipeline.

Collisions
----------

As part of the internal scene API that I was developing for the engine,
I also included the creation and handling of bounding box collisions.
Included in this was the ability to add on-collision event listeners to
objects as-well as transformations of the bounding boxes so as the
object moved so did the box for its collisions.

Results
=======

Editor
------

![[\[fig:1\]]{#fig:1 label="fig:1"}Scene Editor](editor1.png)

![[\[fig:1\]]{#fig:1 label="fig:1"}Code Editor](editor2.png)

![[\[fig:1\]]{#fig:1 label="fig:1"}Engine Compile &
Execution](editor3.png)

The end result of the editor included a desktop application where I
could easily adjust the current scene and its objects, add code that
interacted with the engine and compile & execute the engine all from
that application.

Engine
------

![[\[fig:1\]]{#fig:1 label="fig:1"}Multiple meshes with
shadows](engine1.png)

![[\[fig:1\]]{#fig:1 label="fig:1"}Sky-box above
mesh](engine2.png)

![[\[fig:1\]]{#fig:1 label="fig:1"}Looking down on mesh with single
point-light](engine3.png)

![[\[fig:1\]]{#fig:1 label="fig:1"}Reflective mesh, textured primitives
and shadows](engine4.png)

Resources & Acknowledgements
============================

Resources
---------


https://www.assimp.org/ Assimp Mesh Importer,\ 
https://learnopengl.com/ Joey de Vries, *Learn OpenGL*,\ 
https://renderdoc.org/ RenderDoc,\ 
https://www.blender.org/ Blender,\

Acknowledgements
----------------

Dana Cobzas, *Supervisor*\
Jon Coulson, *Technical Support*\
Andrew Whittle, *Helpful graphics resource*\
https://github.com/awhittle3
