import errorHTML from "./templates/error.html";

const errorTemplate = loadErrorTemplate();

export function toggleContentVisibility(elem, isVisible) {
	const content = elem.querySelector(".analytics-content")
	if (isVisible) content.removeAttribute("style");
	else content.style.display = "none";
}

export function showError(err, parentEl) {
	const errorClone = document.importNode(errorTemplate.content, true);

	const errorMsg = errorClone.querySelector(`[data-field="error-msg"]`);
	errorMsg.textContent = err;

	parentEl.appendChild(errorClone);
}

function loadErrorTemplate() {
	const parser = new DOMParser();
	const doc = parser.parseFromString(errorHTML, "text/html");
	const template = doc.getElementById("error-template");
	return template;
}
