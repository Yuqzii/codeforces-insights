const bs = require("browser-sync").create();
const { createProxyMiddleware } = require("http-proxy-middleware");

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
		createProxyMiddleware({
			target: "http://server:8080",
			changeOrigin: true,
			pathFilter: "/api/**",
			pathRewrite: { "^/api": "" },
			proxyTimeout: 300000, // 5 minutes
			onError: (err, req, res) => {
				console.error("Proxy error:", err);
				res.writeHead(502, { "Content-Type": "text/plain:" });
				res.end("Bad Gateway");
			}
		})
	],
});

const shutdown = () => {
	console.log("Stopping Browsersync");
	bs.exit();
	process.exit(0);
};

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
