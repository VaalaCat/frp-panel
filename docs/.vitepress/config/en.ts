import type { DefaultTheme, LocaleSpecificConfig } from "vitepress";

export const enConfig: LocaleSpecificConfig<DefaultTheme.Config> = {
  themeConfig: {
    nav: [
      { text: "Home", link: "/en/" },
      { text: "Source Code", link: "https://github.com/vaalacat/frp-panel" },
    ],
    sidebar: [
      {
        text: "Quick Start",
        collapsed: false,
        link: "/en/quick-start",
        items: [
          { text: "Master Deployment", link: "/en/deploy-master" },
          { text: "Server Deployment", link: "/en/deploy-server" },
          { text: "Client Deployment", link: "/en/deploy-client" },
        ],
      },
      {
        text: "Advanced Usage",
        collapsed: false,
        link: "/en/wireguard",
        items: [
          { text: "WireGuard Multi-Hop Networking", link: "/en/wireguard" },
        ],
      },
      {
        text: "Configuration",
        collapsed: false,
        link: "/en/all-configs",
      },
      {
        text: "Contribution Guide",
        collapsed: false,
        link: "/en/contribute",
      },
      {
        text: "FAQ",
        collapsed: false,
        link: "/en/faq",
      },
      {
        text: "Screenshots",
        collapsed: false,
        link: "/en/screenshots",
      },
    ],
    socialLinks: [
      { icon: "github", link: "https://github.com/vaalacat/frp-panel" },
    ],
  },
};
