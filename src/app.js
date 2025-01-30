document.querySelectorAll("#check").forEach(checkbox => {
    checkbox.addEventListener("change", function () {
        const id = this.dataset.id;
        this.closest("tr").style.backgroundColor = this.checked ? "#F1F1F1" : "#FFFFFF";
        document.getElementById(id).style.opacity = this.checked ? "1" : "0";
    });
});

function newTransaction() {
    const newTransactionForm = document.getElementById("new-transaction")
    if (window.getComputedStyle(newTransactionForm,null).getPropertyValue("display") === "none") {
        newTransactionForm.style.display = "block"
    } else {
        newTransactionForm.style.display = "none"
    }
}
