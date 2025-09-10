const url = '/api/'

export async function fetchUserInfo(handle) {
	const endpoint = `users/${handle}`;

	const response = await fetch(url + endpoint);
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchSolvedRatings(handle) {
	const endpoint = `users/solved-ratings/${handle}`;

	const response = await fetch(url + endpoint);
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchSolvedTagsAndRatings(handle) {
	const endpoint = `users/solved-tags-ratings/${handle}`;

	const response = await fetch(url + endpoint);
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchRatingChanges(handle) {
	const endpoint = `users/rating/${handle}`;

	const response = await fetch(url + endpoint);
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}
