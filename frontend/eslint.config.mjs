import js from "@eslint/js";
import globals from "globals";
import css from "@eslint/css";
import html from "@html-eslint/eslint-plugin";
import { defineConfig } from "eslint/config";

export default defineConfig([
	{
		files: ["**/*.{js,mjs,cjs}"],
		plugins: { js },
		extends: ["js/recommended"],
		languageOptions: {
			globals: {
				...globals.browser,
				process: "readonly",
				Chart: "readonly",
			}
		}
	},
	{
		files: ["**/*.css"],
		plugins: { css },
		language: "css/css",
		extends: ["css/recommended"],
		rules: {
			"css/no-invalid-properties": "off",
			"css/no-unknown-custom-properties": "off",
			"css/use-baseline": "off",
		},
	},
	{
		files: ["**/*.html"],
		plugins: {
			html,
		},
		// When using the recommended rules (or "html/all" for all rules)
		extends: ["html/recommended"],
		language: "html/html",
		rules: {
			"html/no-duplicate-class": "error",
			"html/indent": ["off"],
			"html/attrs-newline": ["off"],
			"html/no-extra-spacing-attrs": ["off"],
		}
	}
]);
