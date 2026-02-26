import { updateAnalytics } from "./analytics";
import { main, showMain } from "./main";

const form = document.getElementById("user-form");
const input = document.getElementById("handle-input");

let controller = new AbortController();

export function listenForSearch() {
	form.addEventListener("submit", async (e) => {
		e.preventDefault();

		const handle = input.value.trim();
		if (!handle) return;

		analyzeUser(handle);
	});
}

function analyzeUser(handle) {
	controller.abort();
	controller = new AbortController();

	showMain()
	main.scrollIntoView({
		behavior: "smooth"
	});

	updateAnalytics(handle, controller.signal);
	sessionStorage.setItem("hasAnalysed", "true");
}
