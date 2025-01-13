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
        for (let [key, value] of formData.entries()) {
            console.log(`${key}: ${value}`);
        }

        fetch(`/api/section/new-name?id=${id}`, {
            method: "POST",
            body: formData,
        }).then(async resp => {
            if (resp.ok) {
                const newName = await resp.text();
                const sectionTitle = section.querySelector("h2");
                sectionTitle.textContent = newName;
                window.location.href = "/";
            }
        })
    })
})
