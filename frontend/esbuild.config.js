esbuild = require("esbuild");
fs = require("fs");
path = require("path");

prod = process.env.NODE_ENV === "production";

function copyDir(srcDir, destDir) {
	fs.mkdirSync(destDir, { recursive: true });
	for (const entry of fs.readdirSync(srcDir, { withFileTypes: true })) {
		const srcPath = path.join(srcDir, entry.name);
		const destPath = path.join(destDir, entry.name);
		if (entry.isDirectory())
			copyDir(srcPath, destPath);
		else
			fs.copyFileSync(srcPath, destPath);
	}
}

async function start() {
	const ctx = await esbuild.context({
		entryPoints: ["./src/main.js", "./src/style.css"],
		bundle: true,
		platform: "node",
		outdir: "./dist/",
		sourcemap: !prod,
		minify: prod,
	});

	await ctx.watch();
	console.log("Watching for changes");

	copyDir("./public/", "./dist/");
	const watcher = fs.watch("./public/", { recursive: true }, () => {
		copyDir("./public/", "./dist/");
	});

	const shutdown = async () => {
		console.log("Stopping esbuild");
		watcher.close();
		await ctx.dispose();
		process.exit(0);
	};

	process.on("SIGINT", shutdown);
	process.on("SIGTERM", shutdown);
}

start()
