import type * as Preset from "@docusaurus/preset-classic";
import type { Config } from "@docusaurus/types";
import { themes as prismThemes } from "prism-react-renderer";

const config: Config = {
  title: "My Simple Site",
  tagline: "A single page Docusaurus site",
  favicon: "img/favicon.ico",

  future: {
    v4: true,
  },

  url: "https://your-docusaurus-site.example.com",
  baseUrl: "/",

  organizationName: "facebook",
  projectName: "docusaurus",

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          path: "docs",
          routeBasePath: "/", // Serve docs at root
          sidebarPath: require.resolve("./sidebars.ts"),
          // No edit URL
          editUrl: undefined,
          // Show only intro.md as the single page
          include: ["index.md"],
        },
        blog: false, // Disable blog plugin
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Minimal navbar that we'll hide with CSS
    navbar: {
      title: "Simple Site",
      logo: {
        alt: "Logo",
        src: "img/logo.svg",
      },
      items: [], // No navbar items
    },
    // Minimal footer that we'll hide with CSS
    footer: {
      style: "dark",
      links: [], // No links in footer
      copyright: `Copyright © ${new Date().getFullYear()}`,
    },
    // Disable docs sidebar
    docs: {
      sidebar: {
        hideable: true,
        autoCollapseCategories: true,
      },
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
