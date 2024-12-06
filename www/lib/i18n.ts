import i18next from 'i18next'
import { initReactI18next } from 'react-i18next'
import { atom } from 'nanostores'
import enTranslations from '../i18n/locales/en.json'
import zhTranslations from '../i18n/locales/zh.json'

const LANGUAGE_KEY = 'LANGUAGE'

// Get initial language from localStorage or default to 'zh'
const getInitialLanguage = () => {
  if (typeof window === 'undefined') return 'zh'
  return localStorage.getItem(LANGUAGE_KEY) || 'zh'
}

export const $language = atom(getInitialLanguage())

const i18n = i18next.createInstance()

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: enTranslations,
      },
      zh: {
        translation: zhTranslations,
      },
    },
    lng: getInitialLanguage(),
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false,
    },
  })

export const setLanguage = async (lng: 'en' | 'zh') => {
  await i18n.changeLanguage(lng)
  $language.set(lng)
  if (typeof window !== 'undefined') {
    localStorage.setItem(LANGUAGE_KEY, lng)
  }
}

export default i18n
