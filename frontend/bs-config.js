const bs = require("browser-sync").create();

bs.init({
	server: "dist",
	files: ["dist/**/*"],
	port: 3000,
	open: false,
	notify: false,
});

const shutdown = () => {
	console.log("Stopping Browsersync");
	bs.exit();
	process.exit(0);
};

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
