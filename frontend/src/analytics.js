import { SolvedTags, SolvedRatings, RatingHistory, hideLoader, showLoader, getRatingColor } from "./charts.js";
import { getUserInfo } from "./codeforces.js";

const apiUrl = '/api/';

const toggleOtherTags = document.getElementById('toggle-other-tags');
const toggle800Probs = document.getElementById('toggle-800-rating');
export const solvedTags = new SolvedTags(toggleOtherTags);
export const solvedRatings = new SolvedRatings(toggle800Probs);
export const ratingHistory = new RatingHistory();

const userDetails = document.getElementById('user-details');

document.addEventListener('DOMContentLoaded', () => {
	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	toggleOtherTags.addEventListener('click', () => {
		solvedTags.toggleOther();
	});

	toggle800Probs.addEventListener('click', () => {
		solvedRatings.toggle800Rating();
	});
});


export async function updateAnalytics(handle, signal) {
	// Set charts to loading
	solvedTags.loading = true;
	toggleOtherTags.style.display = 'none';
	solvedTags.updateChart();
	solvedRatings.loading = true;
	toggle800Probs.style.display = 'none';
	solvedRatings.updateChart();

	ratingHistory.loading = true;
	ratingHistory.updateChart();
	showLoader(userDetails);

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
	updateUserInfo(handle, signal);
}

async function updateTags(handle, signal) {
	return safeUpdate(`users/solved-tags/${handle}`, data => {
		solvedTags.updateData(data);
		solvedTags.loading = false;
		toggleOtherTags.style.display = 'inline';
		solvedTags.updateChart();
	}, signal);
}

async function updateSolvedRatings(handle, signal) {
	return safeUpdate(`users/solved-ratings/${handle}`, data => {
		solvedRatings.updateData(data);
		solvedRatings.loading = false;
		toggle800Probs.style.display = 'inline';
		solvedRatings.updateChart();
	}, signal);
}

async function updateRatingChanges(handle, signal) {
	return safeUpdate(`users/rating/${handle}`, data => {
		ratingHistory.updateRatingData(data);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	}, signal);
}

async function updateSolvedRatingsTime(handle, signal) {
	return safeUpdate(`users/solved-ratings-time/${handle}`, data => {
		ratingHistory.updateSolvedData(data);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	}, signal);
}

async function updatePerformance(handle, signal) {
	return safeUpdate(`users/performance/${handle}`, data => {
		data.sort((a, b) => a.timestamp > b.timestamp);

		ratingHistory.updatePerfomanceData(data);
		ratingHistory.loading = false;
		ratingHistory.updateChart();
	}, signal);
}

async function updateUserInfo(handle, signal) {
	const userInfo = await getUserInfo(handle, signal);

	hideLoader(userDetails);
	document.getElementById('user-title-photo').src = userInfo.titlePhoto;
	document.getElementById('username').textContent = userInfo.handle;
	document.getElementById('user-country').textContent = userInfo.country;

	const rating = document.getElementById('user-rating');
	rating.textContent = userInfo.rating;
	rating.style.setProperty('--text-color', getRatingColor(userInfo.rating));
	const peakRating = document.getElementById('user-peak-rating');
	peakRating.textContent = userInfo.maxRating;
	peakRating.style.setProperty('--text-color', getRatingColor(userInfo.maxRating));
}

async function safeUpdate(endpoint, updater, signal) {
	try {
		const resp = await fetch(apiUrl + endpoint, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		updater(data);
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Request failed:", err);
	}
}
