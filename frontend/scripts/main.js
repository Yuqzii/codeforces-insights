import { fetchPerformance, fetchRatingChanges, fetchSolvedRatingsTime, fetchSolvedTagsAndRatings, fetchUserInfo } from "./api.js";
import { hideLoader, showLoader, SolvedRatings, SolvedTags, RatingHistory, getColors } from "./charts.js";

const root = document.documentElement;

const userDetails = document.getElementById('user-details');
const form = document.getElementById('user-form');
const input = document.getElementById('handle-input');
const perfLoader = document.getElementById('performance-loader');
const toggleOtherTags = document.getElementById('toggle-other-tags')
const themeSelect = document.getElementById('theme-select');

const solvedTags = new SolvedTags(toggleOtherTags);
const solvedRatings = new SolvedRatings();
const ratingHistory = new RatingHistory();

let controller = new AbortController();

document.addEventListener('DOMContentLoaded', () => {
	const savedTheme = localStorage.getItem('theme') || 'theme-catppuccin';
	setTheme(savedTheme);
	themeSelect.value = savedTheme;

	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	perfLoader.style.display = 'none';
	toggleOtherTags.style.display = 'none';

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		analyzeUser(handle);
	});

	toggleOtherTags.addEventListener('click', () => {
		solvedTags.toggleOther();
	});

	themeSelect.addEventListener('change', (e) => {
		const theme = e.target.value;
		setTheme(theme);
	});
});

async function analyzeUser(handle) {
	controller.abort();
	controller = new AbortController();
	const signal = controller.signal;

	// Set charts to loading
	showLoader(userDetails);
	solvedRatings.loading = true;
	solvedTags.loading = true;
	ratingHistory.loading = true;
	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();
	perfLoader.style.display = 'none';

	toggleOtherTags.style.display = 'none';

	document.querySelector("main").scrollIntoView({
		behavior: "smooth"
	});

	ratingHistory.updatePerfomanceData([]);
	ratingHistory.updateRatingData([]);
	ratingHistory.updateSolvedData([]);

	try {
		const tagsRatings = await fetchSolvedTagsAndRatings(handle, { signal });
		solvedTags.updateData(tagsRatings.tags);
		solvedTags.loading = false;
		solvedTags.updateChart();
		solvedRatings.updateData(tagsRatings.ratings);
		solvedRatings.loading = false;
		solvedRatings.updateChart();
	} catch (err) {
		if (err.name == 'AbortError')
			return;
		throw err;
	}

	try {
		const ratingChanges = await fetchRatingChanges(handle, { signal });
		ratingHistory.updateRatingData(ratingChanges);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	} catch (err) {
		if (err.name == 'AbortError')
			return;
		throw err;
	}

	try {
		const solved = await fetchSolvedRatingsTime(handle, { signal });
		ratingHistory.updateSolvedData(solved);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	} catch (err) {
		if (err.name == 'AbortError')
			return;
		throw err;
	}

	perfLoader.style.display = 'flex';
	try {
		const performance = await fetchPerformance(handle, { signal });
		ratingHistory.updatePerfomanceData(performance);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
		perfLoader.style.display = 'none';
	} catch (err) {
		if (err.name == 'AbortError')
			return;
		throw err;
	}

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
