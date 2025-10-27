const url = "https://codeforces.com/api/"

async function safeFetch(endpoint, signal) {
	try {
		const resp = await fetch(url + endpoint, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		if (data.status !== "OK") throw new Error(`Codeforces not OK: ${data.comment}`);
		return data.result;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Request failed:", err);
	}
}

export async function getUserInfo(handle, signal) {
	const info = await safeFetch(`user.info?handles=${handle}`, signal);
	return info[0];
}
