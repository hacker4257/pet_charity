import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import zhCN from './zh-CN.json';
import enUS from './en-US.json';

i18n
    .use(LanguageDetector)   // 自动检测浏览器语言
    .use(initReactI18next)
    .init({
    resources: {
        'zh-CN': { translation: zhCN },
        'en-US': { translation: enUS },
    },
    fallbackLng: 'zh-CN',
    interpolation: {
        escapeValue: false,   // React 已经做了 XSS 防护
    },
});


export default i18n;
