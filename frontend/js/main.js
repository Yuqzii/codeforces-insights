import { fetchUserInfo } from "./api.js";

document.addEventListener('DOMContentLoaded', () => {
	const form = document.getElementById('user-form');
	const input = document.getElementById('handle-input');

	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const handle = input.value.trim();
		if (!handle) return;

		try {
			await fetchUserInfo(handle);
		} catch (err) {
			alert("Failed to load user stats: " + err.message);
		}
	});
});
