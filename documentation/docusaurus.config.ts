import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

const config: Config = {
  title: "Fuego",
  tagline: "The framework for busy Go developers",
  favicon: "img/fuego.ico",

  // Set the production url of your site here
  url: "https://go-fuego.github.io", //TODO
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/fuego",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "go-fuego", // Usually your GitHub org/user name.
  projectName: "fuego", // Usually your repo name.

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  markdown: {
    mermaid: true,
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/go-fuego/fuego",
          sidebarCollapsed: false,
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/go-fuego/fuego/documentation",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themes: ["@docusaurus/theme-mermaid"],
  plugins: ["docusaurus-lunr-search"],

  themeConfig: {
    // Replace with your project's social card
    image: "img/fuego.png",
    navbar: {
      title: "Fuego",
      logo: {
        alt: "Fuego Logo",
        src: "img/logo.svg",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "tutorialSidebar",
          position: "left",
          label: "📖 Docs",
        },
        {
          href: "https://pkg.go.dev/github.com/go-fuego/fuego",
          position: "left",
          label: "📚 Reference",
        },
        {
          href: "https://github.com/go-fuego/fuego/tree/main/examples/",
          position: "left",
          label: "👀 Examples",
        },
        {
          href: "https://github.com/go-fuego/fuego",
          label: "GitHub",
          position: "right",
        },
        {
          position: "left",
          value: `<div class='display-large' style='background-color: light-dark(rgb(235 237 240), rgb(68 73 80));padding-top: 6px; padding-bottom: 6px; padding-left: 6px; padding-right: 6px;border-radius: 12px; color: light-dark(black, white)'><span>Secure, Document, and Monetize your Fuego API for FREE</span><a target='_blank' href='https://zuplo.link/fuego-web' style='background-color: rgb(255 0 189);padding: 6px;border-radius: 8px; margin-left:6px; color: white; text-wrap: nowrap'>Try Zuplo!</a></div>
            <div class='display-medium' style='background-color: light-dark(rgb(235 237 240), rgb(68 73 80));padding-top: 6px; padding-bottom: 6px; padding-left: 6px; padding-right: 6px;border-radius: 12px; color: light-dark(black, white)'><span class='hidden-when-very-small'>Scale your Fuego API</span><a target='_blank' href='https://zuplo.link/fuego-web' style='background-color: rgb(255 0 189);padding: 6px;border-radius: 8px; margin-left:6px; color: white; text-wrap: nowrap'>Try Zuplo!</a></div>
            <a class='display-small' target='_blank' href='https://zuplo.link/fuego-web' style='background-color: rgb(255 0 189);padding: 6px;border-radius: 8px; margin-left:6px; color: white; text-wrap: nowrap'>Try Zuplo!</a>
            <div class='display-mobile' style='background-color: light-dark(rgb(235 237 240), rgb(68 73 80));padding-top: 6px; padding-bottom: 6px; padding-left: 6px; padding-right: 6px;border-radius: 12px; color: light-dark(black, white)'><span class='hidden-when-very-small'>Scale your Fuego API</span><a target='_blank' href='https://zuplo.link/fuego-web' style='background-color: rgb(255 0 189);padding: 6px;border-radius: 8px; margin-left:6px; color: white; text-wrap: nowrap'>Try Zuplo!</a></div>`,
          type: "html",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Docs",
              to: "/docs/",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "Twitter",
              href: "https://twitter.com/FuegoFramework",
            },
            {
              label: "Youtube",
              href: "https://youtube.com/@Golang-Fuego",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/go-fuego/fuego",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Fuego, Inc. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      defaultLanguage: "go",
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
