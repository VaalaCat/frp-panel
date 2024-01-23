import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import { atom } from 'nanostores'

const LANGUAGE_KEY = 'LANGUAGE'

const resources = {
  en: {
    translation: {
      新建: 'New',
    },
  },
  zh: {
    translation: {
      新建: '新建',
    },
  },
} as const

export const $language = atom('zh')

export const setLanguage = async (lng: 'en' | 'zh') => {
  await i18n.changeLanguage(lng)
  $language.set(lng)
  globalThis.localStorage && localStorage.setItem(LANGUAGE_KEY, lng)
}

i18n.use(initReactI18next).init({
  resources,
  lng: $language.get(),

  interpolation: {
    escapeValue: false,
  },
})

export default i18n
