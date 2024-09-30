import * as React from 'react';
import {useEffect, useRef, useState} from 'react';
import MonacoEditor from 'react-monaco-editor';

import {uiUrl} from '../base';
import {Button} from './button';
import {PhaseIcon} from './phase-icon';
import {SuspenseMonacoEditor} from './suspense-monaco-editor';

interface Props {
    type?: string;
    value: string;
    buttons?: React.ReactNode;
    lang: string;
    fields?: string[];
    onLangChange: (lang: string) => void;
    onChange?: (value: string) => void;
}

export function ObjectEditorPlain({type, value, buttons, lang, fields, onChange, onLangChange}: Props) {
    const [error, setError] = useState<Error>();
    const editor = useRef<MonacoEditor>(null);

    useEffect(() => {
        if (!editor.current || value === editor.current.editor.getValue()) {
            return;
        }
        editor.current.editor.setValue(value);
    }, [editor, value]);

    useEffect(() => {
        if (!type || lang !== 'json') {
            return;
        }

        (async () => {
            const uri = uiUrl('assets/jsonschema/schema.json');
            try {
                const res = await fetch(uri);
                const swagger = await res.json();
                // lazy load this, otherwise all of monaco-editor gets imported into the main bundle
                const languages = (await import(/* webpackChunkName: "monaco-editor" */ 'monaco-editor/esm/vs/editor/editor.api')).languages;
                // adds auto-completion to JSON only
                languages.json.jsonDefaults.setDiagnosticsOptions({
                    validate: true,
                    schemas: [
                        {
                            uri,
                            fileMatch: ['*'],
                            schema: {
                                $id: 'http://workflows.argoproj.io/' + type + '.json',
                                $ref: '#/definitions/' + type,
                                $schema: 'http://json-schema.org/draft-07/schema',
                                definitions: swagger.definitions
                            }
                        }
                    ]
                });
            } catch (err) {
                setError(err);
            }
        })();
    }, [lang, type]);

    // this calculation is rough, it is probably hard to work for for every case, essentially it is:
    // some pixels above and below for buttons, plus a bit of a buffer/padding
    const height = Math.max(600, window.innerHeight * 0.9 - 250);

    return (
        <>
            <div style={{paddingBottom: '1em'}}>
                <Button outline={true} onClick={() => onLangChange(lang === 'yaml' ? 'json' : 'yaml')}>
                    <span style={{fontWeight: lang === 'json' ? 'bold' : 'normal'}}>JSON</span>/<span style={{fontWeight: lang === 'yaml' ? 'bold' : 'normal'}}>YAML</span>
                </Button>

                {fields.map(x => (
                    <Button
                        key={x}
                        icon='caret-right'
                        outline={true}
                        onClick={() => {
                            // Attempt to move the correct section of the document. Ideally, we'd have the line at the top of the
                            // editor, but Monaco editor does not have method for this (e.g. `revealLineAtTop`).

                            // find the line for the section in either YAML or JSON
                            const index = value.split('\n').findIndex(y => (lang === 'yaml' ? y.startsWith(x + ':') : y.includes('"' + x + '":')));

                            if (index >= 0) {
                                const lineNumber = index + 1;
                                editor.current.editor.revealLineInCenter(lineNumber);
                                editor.current.editor.setPosition({lineNumber, column: 0});
                                editor.current.editor.focus();
                            }
                        }}>
                        {x}
                    </Button>
                ))}
                {buttons}
            </div>
            <div>
                <SuspenseMonacoEditor
                    ref={editor}
                    key='editor'
                    defaultValue={value}
                    language={lang}
                    height={height + 'px'}
                    options={{
                        readOnly: onChange === null,
                        minimap: {enabled: false},
                        guides: {
                            indentation: false
                        },
                        scrollBeyondLastLine: true
                    }}
                    onChange={v => {
                        if (onChange) {
                            try {
                                onChange(v);
                                setError(null);
                            } catch (e) {
                                setError(e);
                            }
                        }
                    }}
                />
            </div>
            {error && (
                <div style={{paddingTop: '1em'}}>
                    <PhaseIcon value='Error' /> {error.message}
                </div>
            )}
            {onChange && (
                <div>
                    <i className='fa fa-info-circle' />{' '}
                    {lang === 'json' ? <>Full auto-completion enabled.</> : <>Basic completion for YAML. Switch to JSON for full auto-completion.</>}{' '}
                    <a href='https://argo-workflows.readthedocs.io/en/latest/ide-setup/'>Learn how to get auto-completion in your IDE.</a>
                </div>
            )}
        </>
    );
}
