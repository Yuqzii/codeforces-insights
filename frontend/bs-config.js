const bs = require("browser-sync").create()

bs.init({
	server: "dist",
	files: ["dist/**/*"],
	port: 3000,
	open: false,
	notify: false,
});
