// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).

import { themes as prismThemes } from "prism-react-renderer";

/** @type {import('@docusaurus/types').Config} */
const config = {
	title: "Battlesnake",
	tagline: "Мультиплеерная игра, где соревнуются игровые модели",
	url: "https://docs.battlesnake.com",
	baseUrl: "/",
	favicon: "img/favicon.ico",

	// GitHub pages deployment config.
	// If you aren't using GitHub pages, you don't need these.
	organizationName: "BattlesnakeOfficial", // Usually your GitHub org/user name.
	projectName: "docs", // Usually your repo name.
	trailingSlash: false,

	i18n: {
		defaultLocale: "ru",
		locales: ["ru"],
	},

	presets: [
		[
			"classic",
			/** @type {import('@docusaurus/preset-classic').Options} */
			({
				docs: {
					routeBasePath: "/",
					sidebarPath: require.resolve("./src/sidebars.js"),
					breadcrumbs: false,
					remarkPlugins: [require("remark-math")],
					rehypePlugins: [require("rehype-katex")],
				},
				theme: {
					customCss: require.resolve("./src/css/custom.css"),
				},
			}),
		],
	],

	plugins: [],

	themeConfig:
		/** @type {import('@docusaurus/preset-classic').ThemeConfig} */
		({
			image: "img/logo-social.png",
			announcementBar: require("./src/announcement.js"),
			navbar: {
				logo: {
					alt: "Battlesnake",
					src: "img/logo-light.svg",
					srcDark: "img/logo-dark.svg",
				},
				items: [
					{
						type: "docSidebar",
						sidebarId: "main",
						label: "Доки",
					},
				],
			},
			prism: {
				theme: prismThemes.github,
				darkTheme: prismThemes.dracula,
			},
		}),
	stylesheets: [
		{
			href: "https://cdn.jsdelivr.net/npm/katex@0.13.24/dist/katex.min.css",
			type: "text/css",
			integrity:
				"sha384-odtC+0UGzzFL/6PNoE8rX/SPcQDXBJ+uRepguP4QkPCm2LBxH3FA3y+fKSiJ+AmM",
			crossorigin: "anonymous",
		},
	],
};

export default config;
