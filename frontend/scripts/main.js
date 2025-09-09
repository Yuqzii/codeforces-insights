import { fetchSolvedTagsAndRatings, fetchUserInfo } from "./api.js";
import { hideLoader, showLoader, SolvedRatings, SolvedTags } from "./charts.js";

const userDetails = document.getElementById('user-details');
const solvedRatingsElement = document.getElementById('solved-ratings');
const solvedTagsElement = document.getElementById('solved-tags');

const solvedTags = new SolvedTags();
const solvedRatings = new SolvedRatings();

document.addEventListener('DOMContentLoaded', () => {
	const form = document.getElementById('user-form');
	const input = document.getElementById('handle-input');

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		showLoader(userDetails);
		showLoader(solvedRatingsElement);
		showLoader(solvedTagsElement);

		document.querySelector("main").scrollIntoView({
			behavior: "smooth"
		});

		const tagsRatings = await fetchSolvedTagsAndRatings(handle);
		solvedTags.updateData(tagsRatings.tags);
		solvedTags.updateChart();
		solvedRatings.updateData(tagsRatings.ratings);
		solvedRatings.updateChart();

		updateUserInfo(handle);
	});

	document.getElementById('toggle-other-tags').addEventListener('click', () => {
		solvedTags.toggleOther();
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
