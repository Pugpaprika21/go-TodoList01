function closeTodo(todoId) {
    if (todoId) {
        $.ajax({
            type: "PATCH",
            url: "/todo/update/" + todoId,
            data: "data",
            dataType: "json",
            success: function(response) {
                if (response.status) {
                    Swal.fire(response.message, '', 'success').then(resp => {
                        window.location.href = "/todo/index";
                    })
                }
            }
        });
    }
}

function deleteTodo(todoId) {
    const swalWithBootstrapButtons = Swal.mixin({
        customClass: {
            confirmButton: 'btn btn-success',
            cancelButton: 'btn btn-danger'
        },
        buttonsStyling: false
    })

    if (todoId) {
        swalWithBootstrapButtons.fire({
            title: 'ต้องการลบรายการ ' + todoId + ' หรือไม่',
            text: "",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'ดำเนินการ',
            cancelButtonText: 'ยกเลิก',
            reverseButtons: true
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    type: "DELETE",
                    url: "/todo/delete/" + todoId,
                    dataType: "json",
                    success: function(response) {
                        if (response.status) {
                            swalWithBootstrapButtons.fire(
                                'ลบข้อมูลสำเร็จ!',
                                '',
                                'success'
                            ).then(resp => {
                                window.location.href = "/todo/index";
                            });
                        }
                    }
                });

            } else if (result.dismiss === Swal.DismissReason.cancel) {
                swalWithBootstrapButtons.fire(
                    'ยกเลิกดำเนินการ',
                    '',
                    'error'
                )
            }
        })
    }
}
$(function() {
    const url = new URL(window.location.href);
    const msgAction = url.searchParams.get("msg");
    if (msgAction == "created") {
        Swal.fire('เพิ่มข้อมูล task สำเร็จ', '', 'success');
    }

    $("#btn-create-todo").click(function(e) {
        e.preventDefault();
        let title = $("#title").val();
        let description = $("#description").val();
        if (title == "") {
            Swal.fire('กรอกหัวข้อ', '', 'warning');
            $("#title").focus();
            return;
        }

        if (description == "") {
            Swal.fire('กรอกรายละเอียด', '', 'warning');
            $("#description").focus();
            return;
        }

        $("#submit-form-todo").submit();
    });
});