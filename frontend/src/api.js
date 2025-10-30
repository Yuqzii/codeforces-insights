const url = "https://codeforces.com/api/"

async function cfFetch(endpoint, signal) {
	try {
		const resp = await fetch(url + endpoint, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		if (data.status !== "OK") throw new Error(`Codeforces not OK: ${data.comment}`);
		return data.result;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Codeforces request failed:", err);
	}
}

export async function getUserInfo(handle, signal) {
	const data = await cfFetch(`user.info?handles=${handle}`, signal);
	return data[0];
}

export async function getSubmissions(handle, signal) {
	return await cfFetch(`user.status?handle=${handle}`, signal);
}

export async function getRatingHistory(handle, signal) {
	return await cfFetch(`user.rating?handle=${handle}`, signal);
}

export async function getPerformance(ratingHistory, signal) {
	try {
		const resp = await fetch("/api/performance", {
			method: "POST",
			body: JSON.stringify(ratingHistory),
			signal: signal,
		});
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Performance request failed:", err);
	}
}

export async function getPercentile(rating, signal) {
	try {
		const resp = await fetch(`/api/percentile/${rating}`, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Percentile request failed:", err);
	}
}
