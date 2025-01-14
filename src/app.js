const sectionContainer = document.getElementById("section");
const sections = Array.from(sectionContainer.getElementsByClassName("section-item"));

sections.sort(function (a, b) {
    const nameA = a.getAttribute("data-name").toLowerCase();
    const nameB = b.getAttribute("data-name").toLowerCase();
    if (nameA < nameB) return -1;
    if (nameA > nameB) return 1;
    return 0;
});

sections.forEach(function (section) {
    sectionContainer.appendChild(section);
});

sections.map(section => {
    const id = section.getAttribute("data-id");
    /** @type {HTMLFormElement} */
    const updateSectionTitle = document.forms[`updateSectionTitle${id}`]

    updateSectionTitle.addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(updateSectionTitle);
        fetch(`/api/section/new-name?id=${id}`, {
            method: "POST",
            body: formData,
        }).then(async resp => {
            if (resp.ok) {
                window.location.href = "/";
            }
        })
    })

    /** @type {HTMLFormElement} */
    const updateSectionColor = document.forms[`updateSectionColor${id}`]

    updateSectionColor.addEventListener("submit", function (event) {
        event.preventDefault();
        fetch(`/api/section/new-color?id=${id}`, {
            method: "POST",
            body: JSON.stringify({
                //@ts-ignore
                newColor: document.getElementById(`newColor${id}`).value
            })
        }).then(async resp => {
            if (resp.ok) {
                window.location.href = "/";
            }
        })
    })

    const editButton = document.getElementById(`editButton${id}`)
    editButton.addEventListener("click", function () {
        const sectionEditor = document.getElementById(`sectionEditor${id}`)
        if (window.getComputedStyle(sectionEditor).display == "none") {
            sectionEditor.style.display = "block"
        } else {
            sectionEditor.style.display = "none"
        }
    });
})

const percentageLines = document.querySelectorAll(".percentage-line");
percentageLines.forEach(line => {
    const price = parseFloat(line.getAttribute("data-price"));
    const saved = parseFloat(line.getAttribute("data-saved"));

    const savedPercentage = price > 0 ? (saved / price) * 100 : 0;
    const fillElement = line.querySelector(".percentage-fill");
    //@ts-ignore
    fillElement.style.width = `${savedPercentage}%`;
});

const itemContainer = document.querySelectorAll("#items");
itemContainer.forEach(item => {
    const id = item.querySelector(".item").getAttribute("data-id");
    const updateAllocate = document.forms[`updateAllocate${id}`]

    updateAllocate.addEventListener("submit", function (event) {
        event.preventDefault();
        fetch(`/api/item/allocate?id=${id}`, {
            method: "POST",
            body: JSON.stringify({
                //@ts-ignore
                ammountToAllocate: parseFloat(document.getElementById(`newAllocate${id}`).value),
            })
        }).then(async resp => {
            if (resp.ok) {
                window.location.href = "/";
            }
        })
    })
})
