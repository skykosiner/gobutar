const template = document.createElement('template');
template.innerHTML = `
    <style>
    a {
        color: #5899D9;
        transition: 0.3s ease-in-out;
    }

    a:hover {
        color: #3897D3;
    }

    .footer {
        position: fixed;
        bottom: 0rem;
        width: 100%;
        height: 3rem;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        background-color: #F1F1F1;
    }

    .footer * {
        background-color: #F1F1F1;
    }

    .footer>p {
        margin: 0;
        padding: 0.1rem;
    }


    </style>
    <footer class="footer">
        <p id="footer-year">Gobutar ©</p>
        <p>Created by <a href="https://skykosiner.com" target=" _blank">Sky Kosiner</a></p>
    </footer>
    `;

class Footer extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: "open" });
        const content = template.content.cloneNode(true);
        this.shadowRoot.appendChild(content);
        const dateSpan = this.shadowRoot.querySelector("#footer-year")
        dateSpan.textContent += new Date().getFullYear().toString();
    }
}

customElements.define("footer-component", Footer);

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
