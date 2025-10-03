import { SolvedTags, SolvedRatings, RatingHistory } from "./charts.js";

const apiUrl = '/api/'

const toggleOtherTags = document.getElementById('toggle-other-tags')
export const solvedTags = new SolvedTags(toggleOtherTags);
export const solvedRatings = new SolvedRatings();
export const ratingHistory = new RatingHistory();

const perfLoader = document.getElementById('performance-loader');

document.addEventListener('DOMContentLoaded', () => {
	toggleOtherTags.style.display = 'none';
	perfLoader.style.display = 'none';

	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	toggleOtherTags.addEventListener('click', () => {
		solvedTags.toggleOther();
	});
});


export async function updateAnalytics(handle, signal) {
	// Set charts to loading
	solvedTags.loading = true;
	solvedTags.updateChart();
	solvedRatings.loading = true;
	solvedRatings.updateChart();
	ratingHistory.loading = true;
	ratingHistory.updateChart();

	// Prevent displaying stale data
	ratingHistory.updatePerfomanceData([]);
	ratingHistory.updateRatingData([]);
	ratingHistory.updateSolvedData([]);

	// Asynchronously update charts
	updateSolvedRatings(handle, signal);
	updateTags(handle, signal);
	updateRatingChanges(handle, signal);
	updateSolvedRatingsTime(handle, signal);
	updatePerformance(handle, signal);
}

async function updateTags(handle, signal) {
	try {
		const data = await fetchHelper(apiUrl + `users/solved-tags/${handle}`, signal);
		solvedTags.updateData(data);
		solvedTags.loading = false;
		solvedTags.updateChart();
	} catch (err) {
		console.error("Failed to update solved tags:", err);
	}
}

async function updateSolvedRatings(handle, signal) {
	try {
		const data = await fetchHelper(apiUrl + `users/solved-ratings/${handle}`, signal);
		solvedRatings.updateData(data);
		solvedRatings.loading = false;
		solvedRatings.updateChart();
	} catch (err) {
		console.error("Failed to update solved ratings:", err);
	}
}

async function updateRatingChanges(handle, signal) {
	try {
		const data = await fetchHelper(apiUrl + `users/rating/${handle}`, signal);
		ratingHistory.updateRatingData(data);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	} catch (err) {
		console.error("Failed to update rating changes:", err);
	}
}

async function updateSolvedRatingsTime(handle, signal) {
	try {
		const data = await fetchHelper(apiUrl + `users/solved-ratings-time/${handle}`, signal);
		ratingHistory.updateSolvedData(data);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	} catch (err) {
		console.error("Failed to update solved ratings over time:", err);
	}
}

async function updatePerformance(handle, signal) {
	try {
		perfLoader.style.display = 'flex';
		const data = await fetchHelper(apiUrl + `users/performance/${handle}`, signal);
		ratingHistory.updatePerfomanceData(data);
		ratingHistory.loading = false;
		perfLoader.style.display = 'none';
		ratingHistory.updateChart();
	} catch (err) {
		console.error("Failed to update performance:", err);
	}
}

async function fetchHelper(url, signal) {
	const resp = await fetch(url, { signal });
	if (!resp.ok)
		throw new Error(`response not ok: ${resp.statusText}`);
	return resp.json();
}
