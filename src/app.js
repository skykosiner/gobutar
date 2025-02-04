/**
    *@returns {boolean}
*/
function getCurrentTheme() {
    let darkmode = false;

    const darkModePreference = localStorage.getItem("dark-mode");
    if (darkModePreference === null) {
        darkmode = window.matchMedia("(prefers-color-scheme: dark)").matches;
    } else {
        darkmode = darkModePreference === "true";
    }

    return darkmode
}

/**
    *@param darkmode {boolean}
    *@param toggle {HTMLElement}
*/
function setDarkMode(darkmode, toggle) {
    if (darkmode) {
        document.body.classList.remove("light");
        document.body.classList.add("dark");
        toggle.innerHTML = "☀️"
    } else {
        document.body.classList.remove("dark");
        document.body.classList.add("light");
        toggle.innerHTML = "🌙"
    }
}

window.onload = () => {
    const toggle = document.getElementById("theme-toggle")

    const darkmode = getCurrentTheme();
    setDarkMode(darkmode, toggle);

    toggle.onclick = () => {
        const darkmode = getCurrentTheme();
        const newTheme = !darkmode;
        //@ts-ignore
        localStorage.setItem("dark-mode", newTheme.toString())
        setDarkMode(!darkmode, toggle)
    }
}


const percentageLines = document.querySelectorAll("#percentage-line");
percentageLines.forEach(line => {
    const price = parseFloat(line.getAttribute("data-price"));
    const saved = parseFloat(line.getAttribute("data-saved"));

    const savedPercentage = price > 0 ? (saved / price) * 100 : 0;
    const fillElement = line.querySelector("#percentage-fill");

    //@ts-ignore
    fillElement.style.width = `${savedPercentage}%`;
});

document.querySelectorAll("#check").forEach(checkbox => {
    checkbox.addEventListener("change", function () {
        const id = this.dataset.id;
        this.closest("tr").style.backgroundColor = this.checked ? "#F1F1F1" : "#FFFFFF";
        document.getElementById(id).style.opacity = this.checked ? "1" : "0";
    });
});

function newTransaction() {
    const newTransactionForm = document.getElementById("new-transaction")
    if (window.getComputedStyle(newTransactionForm, null).getPropertyValue("display") === "none") {
        newTransactionForm.style.display = "block"
    } else {
        newTransactionForm.style.display = "none"
    }
}
