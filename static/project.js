function deleteDoc(id) {
    $.ajax({
        url: "/deleteDoc",
        method: "POST",
        data : { docId: id}
    });
}