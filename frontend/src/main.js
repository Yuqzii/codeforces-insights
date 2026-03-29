import { solvedTags, solvedRatings, ratingHistory } from "./analytics.js";
import { getColors } from "./charts.js";
import { observeAndAnimate } from "./entrance-anim.js";
import { listenForSearch } from "./search.js";

const root = document.documentElement;
export const main = document.querySelector("main");
const footer = document.querySelector("footer");

const dots = document.getElementById("grid-dots");
const lines = document.getElementById("grid-lines");

const navbar = document.getElementById("nav-bar");
const nav = document.querySelector("#nav-bar nav");
const themeSelect = document.getElementById("theme-select");

let cursorX = window.innerWidth / 2;
let cursorY = window.innerHeight / 2;

document.addEventListener("DOMContentLoaded", () => {
	const savedTheme = localStorage.getItem("theme") || "theme-catppuccin";
	setTheme(savedTheme);
	themeSelect.value = savedTheme;

	themeSelect.addEventListener("change", (e) => {
		const theme = e.target.value;
		setTheme(theme);
	});

	listenForSearch();

	const analyticsCards = document.querySelectorAll(".card");
	observeAndAnimate(analyticsCards);
});

window.addEventListener("mousemove", e => {
	cursorX = e.clientX;
	cursorY = e.clientY;
});

requestAnimationFrame(updateCursorCSS);

window.addEventListener("scroll", () => {
	if (window.scrollY > 50)
		navbar.classList.add("scrolled");
	else
		navbar.classList.remove("scrolled");
});

window.addEventListener("load", () => {
	const hasAnalysed = sessionStorage.getItem("hasAnalysed");
	if (hasAnalysed)
		showMain();
});


function setTheme(theme) {
	root.classList.remove(localStorage.getItem("theme"));
	root.classList.add(theme);
	localStorage.setItem("theme", theme);
	getColors(); // Update chart colors
	solvedRatings.updateChart();
	solvedTags.updateChart();
	ratingHistory.updateChart();
}

function updateCursorCSS() {
	dots.style.maskPosition =
		`${cursorX - window.innerWidth / 2}px ${cursorY - window.innerHeight / 2}px, 0 0`;
	lines.style.maskPosition =
		`${cursorX - window.innerWidth / 2}px ${cursorY - window.innerHeight / 2}px`;

	requestAnimationFrame(updateCursorCSS);
}

export function showMain() {
	// Make initially hidden content appear and give them their CSS defined style.
	main.removeAttribute("style");
	footer.removeAttribute("style");
	nav.removeAttribute("style");
}
