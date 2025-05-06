import type { DefaultTheme, LocaleSpecificConfig } from 'vitepress'

export const zhConfig: LocaleSpecificConfig<DefaultTheme.Config> = {
  themeConfig: {
    nav: [
      { text: '首页', link: '/' },
      { text: '源码', link: 'https://github.com/vaalacat/frp-panel' }
    ],

    sidebar: [
      {
        text: '快速开始',
        collapsed: true,
        link: '/quick-start',
        items: [
          { text: 'Master 部署', link: '/deploy-master' },
          { text: 'Server 部署', link: '/deploy-server' },
          { text: 'Client 部署', link: '/deploy-client' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/vaalacat/frp-panel' }
    ]
  }
}