esbuild = require("esbuild");
gaze = require("gaze");
fs = require("fs");

gaze("./public/*", function(err, watcher) {
	if (!fs.existsSync('/app/dist')) {
		fs.mkdirSync('/app/dist', { recursive: true });
	}

	// Copy all watched files into dist
	var watched = this.watched();
	for (const filepath of watched['/app/public/']) {
		console.log("Copying", filepath);
		var filename = filepath.replace(/^.*[\\/]/, '')
		fs.copyFileSync(filepath, `/app/dist/${filename}`);
	}

	// On file changed
	this.on('changed', function(filepath) {
		console.log(filepath + ' was changed');
		var filename = filepath.replace(/^.*[\\/]/, '')
		var destPath = `/app/dist/${filename}`;

		function tryCopy(attempts = 0) {
			try {
				if (!fs.existsSync('/app/dist')) {
					fs.mkdirSync('/app/dist', { recursive: true });
				}
				fs.copyFileSync(filepath, destPath);
			} catch (err) {
				if (err.code === 'ENOENT' && attempts < 5) {
					// Retry after a short delay
					setTimeout(() => tryCopy(attempts + 1), attempts * 50);
				} else {
					console.error('Copy failed:', err);
				}
			}
		}

		tryCopy();
	});
});

async function start() {
	const ctx = await esbuild.context({
		entryPoints: ["./src/main.js", "./src/style.css"],
		bundle: true,
		platform: "node",
		outdir: "./dist/",
		sourcemap: true,
		minify: false,
	});

	await ctx.watch();
	console.log("[esbuild] Watching for changes");

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
