import { updateAnalytics, solvedTags, solvedRatings, ratingHistory } from "./analytics.js";
import { hideLoader, getColors } from "./charts.js";

const root = document.documentElement;

const userDetails = document.getElementById('user-details');
const form = document.getElementById('user-form');
const input = document.getElementById('handle-input');
const themeSelect = document.getElementById('theme-select');


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

	updateUserInfo(handle);
}

async function updateUserInfo(username) {
	let data;
	try {
		data = await fetchUserInfo(username);
	} catch (err) {
		console.error(err);
		return;
	}

	hideLoader(userDetails);

	document.getElementById('user-title-photo').src = data.titlePhoto;
	document.getElementById('username').textContent = data.handle;
	document.getElementById('user-rating').textContent = data.rating;
	document.getElementById('user-peak-rating').textContent = data.maxRating;
	document.getElementById('user-country').textContent = data.country;
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
