import esbuild from "esbuild";

esbuild.build({
	entryPoints: ["src/main.js", "src/style/style.css"],
	bundle: true,
	minify: true,
	outdir: "dist",
	loader: {
		".html": "text",
	},
	define: {
		"process.env.API_URL": JSON.stringify(process.env.API_URL || "https://api.cf-insights.org")
	},
}).catch(() => process.exit(1));
