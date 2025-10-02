const bs = require("browser-sync").create();

bs.init({
	server: "dist",
	files: ["dist/**/*"],
	port: 3000,
	open: false,
	notify: false,
	port: 3000,
	ui: {
		port: 3001
	},
	middleware: [
		{
			route: "/api",
			handle: function (req, res, next) {
				const proxy = require("http-proxy").createProxyServer({
					target: "http://server:8080",
					changeOrigin: true,
					ws: true,
					proxyTimeout: 300000, // 5 minutes
				});

				// Rewrite the path to remove /api prefix
				req.url = req.url.replace(/^\/api/, "");

				proxy.web(req, res, {}, function (err) {
					if (err) {
						console.error("Proxy error:", err);
						res.writeHead(502, { "Content-Type": "text/plain" });
						res.end("Bad Gateway");
					}
				});
			}
		}
	],
});

const shutdown = () => {
	console.log("Stopping Browsersync");
	bs.exit();
	process.exit(0);
};

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
