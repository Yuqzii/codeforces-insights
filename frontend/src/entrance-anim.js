export function observeAndAnimate() {
	const analyticsElements = document.querySelectorAll(".card");
	const observer = new IntersectionObserver(entries => {
		entries.forEach(entry => {
			if (!entry.isIntersecting)
				return;

			entry.target.classList.add("fade-in");

			entry.target.addEventListener("animationend", () => {
				entry.target.classList.remove("fade-in");
				entry.target.style.opacity = 1;
			}, { once: true });

			observer.unobserve(entry.target);
		});
	});

	analyticsElements.forEach(element => {
		observer.observe(element);
	});
}
