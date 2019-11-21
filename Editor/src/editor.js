window.onload = () => {

    var editor = ace.edit("editorWindow", {
        theme: "ace/theme/monokai",
        mode: "ace/mode/golang",
        autoScrollEditorIntoView: true,
        maxLines: 30,
        minLines: 30
    });
    editor.setFontSize(16);

    editor.on('change', (data) => {
        console.log(editor.getValue())
    })
}