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
		throw err;
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

export async function getProblems(tags, signal) {
	let url = "problemset.problems";
	if (tags) {
		url += "?tags=";
		for (const tag of tags) {
			url += tag + ";";
		}
		url = url.slice(0, -1);
		console.log(url)
	}

	return await cfFetch(url, signal);
}

export async function getPerformance(handle, ratingHistory, signal) {
	try {
		const reqData = {
			handle: handle,
			ratingHistory: ratingHistory
		};
		const resp = await fetch(`${process.env.API_URL}/performance`, {
			method: "POST",
			body: JSON.stringify(reqData),
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
		const resp = await fetch(`${process.env.API_URL}/percentile/${rating}`, { signal });
		if (!resp.ok) throw new Error(`response not ok: ${resp.statusText}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Percentile request failed:", err);
	}
}
