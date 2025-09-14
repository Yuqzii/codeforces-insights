import { fetchPerformance, fetchRatingChanges, fetchSolvedTagsAndRatings, fetchUserInfo } from "./api.js";
import { hideLoader, showLoader, SolvedRatings, SolvedTags, RatingHistory, getColors } from "./charts.js";

const root = document.documentElement;

const userDetails = document.getElementById('user-details');
const form = document.getElementById('user-form');
const input = document.getElementById('handle-input');
const perfLoader = document.getElementById('performance-loader');
const toggleOtherTags = document.getElementById('toggle-other-tags')
const themeToggleBtn = document.getElementById('toggle-theme');

const solvedTags = new SolvedTags(toggleOtherTags);
const solvedRatings = new SolvedRatings();
const ratingHistory = new RatingHistory();

document.addEventListener('DOMContentLoaded', () => {
	setTheme(localStorage.getItem('theme') || 'theme-catppuccin');

	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	perfLoader.style.display = 'none';
	toggleOtherTags.style.display = 'none';

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		// Set charts to loading
		showLoader(userDetails);
		solvedRatings.loading = true;
		solvedTags.loading = true;
		ratingHistory.loading = true;
		solvedTags.updateChart();
		solvedRatings.updateChart();
		ratingHistory.updateChart();

		toggleOtherTags.style.display = 'none';

		document.querySelector("main").scrollIntoView({
			behavior: "smooth"
		});

		ratingHistory.updatePerfomanceData([]);
		ratingHistory.updateRatingData([]);

		const tagsRatings = await fetchSolvedTagsAndRatings(handle);
		solvedTags.updateData(tagsRatings.tags);
		solvedTags.loading = false;
		solvedTags.updateChart();
		solvedRatings.updateData(tagsRatings.ratings);
		solvedRatings.loading = false;
		solvedRatings.updateChart();

		const ratingChanges = await fetchRatingChanges(handle);
		ratingHistory.updateRatingData(ratingChanges);
		ratingHistory.loading = false;
		ratingHistory.updateChart();

		perfLoader.style.display = 'flex';
		const performance = await fetchPerformance(handle);
		ratingHistory.updatePerfomanceData(performance);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
		perfLoader.style.display = 'none';

		updateUserInfo(handle);
	});

	toggleOtherTags.addEventListener('click', () => {
		solvedTags.toggleOther();
	});

	themeToggleBtn.addEventListener('click', () => {
		const current = localStorage.getItem('theme') || 'theme-catppuccin';
		const next = current == 'theme-gruvbox' ? 'theme-catppuccin' : 'theme-gruvbox';
		setTheme(next);
	});
});

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
