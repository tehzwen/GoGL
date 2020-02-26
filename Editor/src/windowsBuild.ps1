echo "starting build"
$sceneFile=$args[0]
echo $sceneFile

# if ($sceneFile) {
#     cd ../Renderer/; go build; ./Renderer.exe $sceneFile
# } else {
#     cd ../Renderer/; go build; ./Renderer.exe
# }

#cd ../Renderer/; go build; ./Renderer.exe $sceneFile
cd ../Renderer/; go build; ./Renderer.exe $sceneFile