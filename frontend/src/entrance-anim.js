export function observeAndAnimate() {
	const analyticsElements = document.querySelectorAll(".card");
	const observer = new IntersectionObserver(entries => {
		entries.forEach(entry => {
			if (entry.isIntersecting) {
				entry.target.classList.toggle("fade-in", true);
			}
		});
	});

	analyticsElements.forEach(element => {
		observer.observe(element);
	});
}
