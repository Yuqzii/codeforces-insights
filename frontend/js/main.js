import { fetchUserInfo } from "./api.js";
import { updateSolvedTagsAndRatingsCharts } from "./charts.js";

document.addEventListener('DOMContentLoaded', () => {
	const form = document.getElementById('user-form');
	const input = document.getElementById('handle-input');

	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const handle = input.value.trim();
		if (!handle) return;

		await updateSolvedTagsAndRatingsCharts(handle);
		updateUserInfo(handle);
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

	document.getElementById('user-avatar').src = data.avatar;
	document.getElementById('username').textContent = data.handle;
	document.getElementById('user-rating').textContent = data.rating;
	document.getElementById('user-peak-rating').textContent = data.maxRating;
	document.getElementById('user-country').textContent = data.country;
}

