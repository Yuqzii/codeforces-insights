export function observeAndAnimate(elements) {
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

	elements.forEach(element => {
		observer.observe(element);
	});
}
