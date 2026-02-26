import { updateAnalytics } from "./analytics";
import { main, showMain } from "./main";

const form = document.getElementById("user-form");
const heroInput = document.getElementById("handle-input");
const navInput = document.getElementById("handle-input-nav");

let controller = new AbortController();

export function listenForSearch() {
	form.addEventListener("submit", e => {
		e.preventDefault();

		const handle = heroInput.value.trim();
		if (!handle) return;

		navInput.value = handle;

		showMain()
		main.scrollIntoView({
			behavior: "smooth"
		});

		analyzeUser(handle);
	});

	navInput.addEventListener("change", e => {
		e.preventDefault();

		const handle = navInput.value.trim();
		if (!handle) return;

		heroInput.value = handle;

		analyzeUser(handle);
	});
}

function analyzeUser(handle) {
	controller.abort();
	controller = new AbortController();

	updateAnalytics(handle, controller.signal);
	sessionStorage.setItem("hasAnalysed", "true");
}
