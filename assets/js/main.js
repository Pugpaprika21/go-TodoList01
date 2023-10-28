function editTodo(todoId) {
    console.log(todoId);
}

function deleteTodo(todoId) {
    console.log(todoId);
}

$(function() {
    $("#table-todo").DataTable({
        scrollY: 500,
    });
    $("#btn-create-todo").click(function(e) {
        e.preventDefault();
        let title = $("#title").val();
        let description = $("#description").val();
        if (title != "" && description != "") {
            $("#submit-form-todo").submit();
        }
    });
});