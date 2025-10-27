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
	const data = await safeFetch(`user.info?handles=${handle}`, signal);
	return data[0];
}

export async function getSubmissions(handle, signal) {
	return await safeFetch(`user.status?handle=${handle}`, signal);
}

export async function getRatingHistory(handle, signal) {
	return await safeFetch(`user.rating?handle=${handle}`, signal);
}

export async function getPerformance(ratingHistory, signal) {
	try {
		const resp = await fetch("/api/performance", {
			method: "POST",
			body: JSON.stringify(ratingHistory),
			signal: signal,
		});
		if (!resp.ok) throw new Error(`performance response not ok: ${resp.statusText}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Request failed:", err);
	}
}
