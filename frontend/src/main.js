import { updateAnalytics, solvedTags, solvedRatings, ratingHistory } from "./analytics.js";
import { getColors } from "./charts.js";
import { observeAndAnimate } from "./entrance-anim.js";

const root = document.documentElement;

const form = document.getElementById('user-form');
const input = document.getElementById('handle-input');
const themeSelect = document.getElementById('theme-select');
const highContrastSlider = document.getElementById('high-contrast-slider');

let controller = new AbortController();

let cursorX = window.innerWidth / 2;
let cursorY = window.innerHeight / 2;

document.addEventListener('DOMContentLoaded', () => {
	const savedTheme = localStorage.getItem('theme') || 'theme-catppuccin';
	setTheme(savedTheme);
	themeSelect.value = savedTheme;

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		analyzeUser(handle);
	});

	themeSelect.addEventListener('change', (e) => {
		const theme = e.target.value;
		setTheme(theme);
	});

	highContrastSlider.addEventListener('change', (e) => {
		if (e.target.checked === true)
			root.classList.add('increased-contrast');
		else
			root.classList.remove('increased-contrast');
	});

	observeAndAnimate();
});

window.addEventListener('mousemove', throttle((e) => {
	cursorX = e.clientX;
	cursorY = e.clientY;
	updateCursorCSS();
}, 50));

window.addEventListener('scroll', throttle(() => {
	updateCursorCSS();
}, 50));

async function analyzeUser(handle) {
	controller.abort();
	controller = new AbortController();

	document.querySelector("main").scrollIntoView({
		behavior: "smooth"
	});

	updateAnalytics(handle, controller.signal);
}

function setTheme(theme) {
	root.classList.remove(localStorage.getItem('theme'));
	root.classList.add(theme);
	localStorage.setItem('theme', theme);
	getColors(); // Update chart colors
	solvedRatings.updateChart();
	solvedTags.updateChart();
	ratingHistory.updateChart();
}

function updateCursorCSS() {
	root.style.setProperty('--cursor-x', (cursorX + window.scrollX) + 'px');
	root.style.setProperty('--cursor-y', (cursorY + window.scrollY) + 'px');
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
