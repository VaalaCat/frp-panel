import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { $language } from '@/store/user';
import enTranslation from './locales/en.json';
import zhTranslation from './locales/zh.json';

const savedLanguage = $language.get();

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: enTranslation,
      },
      zh: {
        translation: zhTranslation,
      },
    },
    lng: savedLanguage || 'zh',
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false,
    },
  });

// 监听语言变化并同步到 i18n
$language.subscribe((newLanguage) => {
  if (newLanguage && i18n.language !== newLanguage) {
    i18n.changeLanguage(newLanguage);
  }
});

// 同步初始语言
if (savedLanguage) {
  i18n.changeLanguage(savedLanguage);
}

export default i18n;
