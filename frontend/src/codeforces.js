const url = "https://codeforces.com/api/"

async function safeFetch(endpoint, callback, signal) {
	try {
		const resp = await fetch(url + endpoint, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		if (data.status !== "OK") throw new Error(`Codeforces not OK: ${data.comment}`);
		callback(data.result);
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Request failed:", err);
	}
}

export async function getUserInfo(handle, callback, signal) {
	safeFetch(`user.info?handles=${handle}`, (data) => {
		callback(data[0]);
	}, signal);
}

export async function getSubmissions(handle, callback, signal) {
	safeFetch(`user.status?handle=${handle}`, (data) => {
		callback(data);
	}, signal);
}

export async function getRatingHistory(handle, callback, signal) {
	safeFetch(`user.rating?handle=${handle}`, (data) => {
		callback(data);
	}, signal);
}

export async function getPerformance(ratingHistory, callback, signal) {
	try {
		const resp = await fetch("/api/performance", {
			method: "POST",
			body: JSON.stringify(ratingHistory),
			signal: signal,
		});
		if (!resp.ok) throw new Error(`performance response not ok: ${resp.statusText}`);
		const data = await resp.json();
		callback(data);
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Request failed:", err);
	}
}
