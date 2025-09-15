const url = '/api/'

export async function fetchUserInfo(handle, { signal } = {}) {
	const endpoint = `users/${handle}`;

	const response = await fetch(url + endpoint, { signal });
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchSolvedRatings(handle, { signal } = {}) {
	const endpoint = `users/solved-ratings/${handle}`;

	const response = await fetch(url + endpoint, { signal });
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchSolvedTagsAndRatings(handle, { signal } = {}) {
	const endpoint = `users/solved-tags-ratings/${handle}`;

	const response = await fetch(url + endpoint, { signal });
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchRatingChanges(handle, { signal } = {}) {
	const endpoint = `users/rating/${handle}`;

	const response = await fetch(url + endpoint, { signal });
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}

export async function fetchPerformance(handle, { signal } = {}) {
	const endpoint = `users/performance/${handle}`;

	const response = await fetch(url + endpoint, { signal });
	if (!response.ok)
		throw new Error(`response not ok: ${response.statusText}`);

	return await response.json();
}
