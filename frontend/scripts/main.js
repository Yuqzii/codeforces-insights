import { fetchUserInfo } from "./api.js";
import { updateSolvedTagsAndRatingsCharts } from "./charts.js";
import { toggleShowOtherTags } from "./solvedTags.js";

document.addEventListener('DOMContentLoaded', () => {
	const form = document.getElementById('user-form');
	const input = document.getElementById('handle-input');

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		document.querySelector("main").scrollIntoView({
			behavior: "smooth"
		});

		await updateSolvedTagsAndRatingsCharts(handle);
		updateUserInfo(handle);
	});

	document.getElementById('toggle-other-tags').addEventListener('click', () => {
		toggleShowOtherTags();
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

	document.getElementById('user-title-photo').src = data.titlePhoto;
	document.getElementById('username').textContent = data.handle;
	document.getElementById('user-rating').textContent = data.rating;
	document.getElementById('user-peak-rating').textContent = data.maxRating;
	document.getElementById('user-country').textContent = data.country;
}
