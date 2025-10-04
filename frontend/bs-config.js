const bs = require("browser-sync").create();

bs.init({
	server: "dist",
	files: ["dist/**/*"],
	open: false,
	notify: false,
	port: 3000,
	ui: {
		port: 3001
	},
});

console.log(`
------------------------------------------------------------------------------
IMPORTANT NOTE:
In this project Browsersync is mainly used for hot-reloading.
It is recommended you go to https://localhost:443 instead when nginx starts.
This is to allow HTTPS and HTTP/2 for proper request cancelling on the backend.
------------------------------------------------------------------------------`)

const shutdown = () => {
	console.log("Stopping Browsersync");
	bs.exit();
	process.exit(0);
};

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
