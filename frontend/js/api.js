const url = '/api/'

export async function fetchUserInfo(handle) {
	const endpoint = `users/${handle}`;

		const response = await fetch(url + endpoint);
		if (!response.ok)
			throw new Error(`response not ok: ${response.statusText}`);

		return await response.json();
}
