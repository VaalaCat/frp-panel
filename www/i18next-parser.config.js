module.exports = {
  input: [
    'api/**/*.{js,jsx,ts,tsx}',
    'config/**/*.{js,jsx,ts,tsx}',
    'components/**/*.{js,jsx,ts,tsx}',
    'pages/**/*.{js,jsx,ts,tsx}',
  ],
  output: 'i18n/locales/$LOCALE.json',
  locales: ['en', 'zh'],
}
