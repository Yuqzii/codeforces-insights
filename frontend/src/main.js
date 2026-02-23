import { updateAnalytics, solvedTags, solvedRatings, ratingHistory } from "./analytics.js";
import { getColors } from "./charts.js";
import { observeAndAnimate } from "./entrance-anim.js";

const root = document.documentElement;
const main = document.querySelector("main");
const footer = document.querySelector("footer");

const form = document.getElementById("user-form");
const input = document.getElementById("handle-input");
const themeSelect = document.getElementById("theme-select");
const highContrastSlider = document.getElementById("high-contrast-slider");
const navMenuButton = document.getElementById("hamburger");
const navMenu = document.getElementById("nav-menu");

let controller = new AbortController();

let cursorX = window.innerWidth / 2;
let cursorY = window.innerHeight / 2;

let navMenuActive = false;

document.addEventListener("DOMContentLoaded", () => {
	const savedTheme = localStorage.getItem("theme") || "theme-catppuccin";
	setTheme(savedTheme);
	themeSelect.value = savedTheme;

	form.addEventListener("submit", async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		analyzeUser(handle);
	});

	navMenuButton.addEventListener("click", () => {
		navMenuActive = !navMenuActive;
		if (navMenuActive)
			navMenu.classList.add("active");
		else
			navMenu.classList.remove("active");
	});

	themeSelect.addEventListener("change", (e) => {
		const theme = e.target.value;
		setTheme(theme);
	});

	window.addEventListener("click", (e) => {
		if (!navMenu.contains(e.target) && !navMenuButton.contains(e.target)) {
			// Something outside the nav menu was clicked.
			navMenu.classList.remove("active");
			navMenuActive = false;
		}
	});

	highContrastSlider.addEventListener("change", (e) => {
		if (e.target.checked === true)
			root.classList.add("increased-contrast");
		else
			root.classList.remove("increased-contrast");
	});

	const analyticsCards = document.querySelectorAll(".card");
	observeAndAnimate(analyticsCards);
});

window.addEventListener("mousemove", throttle((e) => {
	cursorX = e.clientX;
	cursorY = e.clientY;
	updateCursorCSS();
}, 50));

window.addEventListener("scroll", throttle(() => {
	updateCursorCSS();
}, 50));

window.addEventListener("load", () => {
	const hasAnalysed = sessionStorage.getItem("hasAnalysed");
	if (hasAnalysed)
		showMain();
});

async function analyzeUser(handle) {
	controller.abort();
	controller = new AbortController();

	showMain();
	main.scrollIntoView({
		behavior: "smooth"
	});

	updateAnalytics(handle, controller.signal);
	sessionStorage.setItem("hasAnalysed", "true");
}

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
	root.style.setProperty("--cursor-x", (cursorX + window.scrollX) + "px");
	root.style.setProperty("--cursor-y", (cursorY + window.scrollY) + "px");
}

function throttle(fn, delay) {
	let t = 0;
	return function(...args) {
		const now = Date.now();
		if (now - t >= delay) {
			fn.apply(this, args);
			t = now;
		}
	}
}

function showMain() {
	// Make the main and footer appear and give them their CSS defined style.
	main.removeAttribute("style");
	footer.removeAttribute("style");
}
