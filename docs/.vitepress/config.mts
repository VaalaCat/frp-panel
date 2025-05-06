import { defineConfig } from 'vitepress'
import { zhConfig } from './config/zh'
import { enConfig } from './config/en'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/frp-panel/',
  locales: {
    root: {
      label: '简体中文',
      ...zhConfig
    },
    en: {
      label: 'English',
      ...enConfig
    }
  },
  title: "Frp-Panel WIKI",
  description: "Wiki of vaalacat's wonderful frp-panel",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
  }
})
