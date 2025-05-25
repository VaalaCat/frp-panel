import type { DefaultTheme, LocaleSpecificConfig } from "vitepress";

export const zhConfig: LocaleSpecificConfig<DefaultTheme.Config> = {
  themeConfig: {
    nav: [
      { text: "首页", link: "/" },
      { text: "源码", link: "https://github.com/vaalacat/frp-panel" },
    ],

    sidebar: [
      {
        text: "快速开始",
        collapsed: false,
        link: "/quick-start",
        items: [
          { text: "Master 部署", link: "/deploy-master" },
          { text: "Server 部署", link: "/deploy-server" },
          { text: "Client 部署", link: "/deploy-client" },
        ],
      },
      {
        text: "配置说明",
        collapsed: false,
        link: "/all-configs",
      },
      {
        text: "贡献指南",
        collapsed: false,
        link: "/contribute",
      },
      {
        text: "截图展示",
        collapsed: false,
        link: "/screenshots",
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/vaalacat/frp-panel" },
    ],
  },
};
