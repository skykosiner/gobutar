export class Introduction {
    /** @type {number} */
    #idx = 0;

    /** @type {string[]} */
    #slides = [
        "slide-1",
        "slide-2",
    ];

    constructor() {
        if (document.title.includes("Introduction")) {
            document.getElementById("slide-1").style.display = "block";
        };

        document.getElementById("next").addEventListener("click", () => {
            if (this.#idx + 1 < 0 || this.#idx + 1 >= this.#slides.length) {
                return;
            }

            this.#idx += 1;
            this.#changeSlide();
        });

        document.getElementById("prev").addEventListener("click", () => {
            if (this.#idx - 1 < 0 || this.#idx - 1 >= this.#slides.length) {
                return;
            }

            this.#idx -= 1;
            this.#changeSlide();
        });
    }

    #changeSlide() {
        this.#slides.map(slide => {
            if (this.#slides.indexOf(slide) === this.#idx) {
                document.getElementById(slide).style.display = "block";
            } else {
                document.getElementById(slide).style.display = "none";
            }
        });
    }
}

export function setCurrency() {
    alert("test");
}
