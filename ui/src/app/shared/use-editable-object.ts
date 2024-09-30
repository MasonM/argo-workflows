import {useState, useMemo} from 'react';
import {stringify, parse} from '../shared/components/object-parser';
import {ScopedLocalStorage} from '../shared/scoped-local-storage';

const defaultLang = 'yaml';

export function useEditableObject<T>(): [T, string, string, boolean, (lang: string) => void, (value: T) => void, (value: T) => void] {
    const storage = new ScopedLocalStorage('object-editor');
    const [templateText, setTemplateText] = useState<string>();
    const [initialTemplateText, setInitialTemplateText] = useState<string>();
    const [lang, setLang] = useState<string>(storage.getItem('lang', defaultLang));

    const template = useMemo(() => (templateText ? parse<T>(templateText) : null), [templateText]);
    const edited = templateText !== initialTemplateText;

    function onLangChange(newLang: string) {
        setLang(newLang);
        storage.setItem('lang', newLang, defaultLang);
        setTemplateText(stringify(template, newLang));
    }

    function setTemplate(value: string | T) {
        if (typeof value === 'string') {
            setTemplateText(value);
        } else {
            setTemplateText(stringify(value, lang));
        }
    }

    function resetTemplate(value: T) {
        const val = stringify(value, lang);
        setTemplateText(val);
        setInitialTemplateText(val);
    }

    return [template, templateText, lang, edited, onLangChange, setTemplate, resetTemplate];
}
