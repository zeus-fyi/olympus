// @ts-ignore
import React, {useEffect, useRef, useState} from "react";
import {SelectionText} from "@uiw/react-textarea-code-editor";
import MonacoEditor from "react-monaco-editor/lib/editor";
import {editor} from "monaco-editor";
// @ts-ignore
import yaml from 'js-yaml';
import {
    avaxMaxBlockAggReduceExample,
    btcMaxBlockAggReduceExample,
    ethMaxBlockAggReduceExample,
    nearMaxBlockAggReduceExample
} from "./ExampleRequests";
import setTheme = editor.setTheme;

window.MonacoEnvironment = { getWorkerUrl: () => proxy };

let proxy = URL.createObjectURL(new Blob([`
	self.MonacoEnvironment = {
		baseUrl: 'https://unpkg.com/monaco-editor@latest/min/'
	};
	importScripts('https://unpkg.com/monaco-editor@latest/min/vs/base/worker/workerMain.js');
`], { type: 'text/javascript' }));

export const languageData = [
    'abap', 'aes', 'apex', 'azcli', 'bat', 'c', 'cameligo',
    'clojure', 'coffeescript', 'cpp', 'csharp', 'csp', 'css', 'dart',
    'dockerfile', 'fsharp', 'go', 'graphql', 'handlebars', 'hcl', 'html',
    'ini', 'java', 'javascript', 'json', 'julia', 'kotlin', 'less', 'lex',
    'lua', 'markdown', 'mips', 'msdax', 'mysql', 'objective', 'pascal',
    'pascaligo', 'perl', 'pgsql', 'php', 'plaintext', 'postiats', 'powerquery',
    'powershell', 'pug', 'python', 'r', 'razor', 'redis', 'redshift', 'restructuredtext',
    'ruby', 'rust', 'sb', 'scala', 'scheme', 'scss', 'shell', 'sol', 'sql', 'st', 'swift', 'systemverilog', 'tcl', 'twig', 'typescript', 'vb', 'verilog', 'xml', 'yaml'
];

export default function ExamplePageMarkdownText(props: any) {
    const {code, setCode, language, onChange, procedureName} = props;
    const [example, setExample] = useState<string>('');
    const themeRef = useRef<string>()
    function onSelectThemeChange(e: React.ChangeEvent<HTMLSelectElement>) {
        e.persist();
        document.documentElement.setAttribute('data-color-mode', /^vs$/.test(e.target.value) ? 'light' : 'dark');
        themeRef.current = e.target.value;
        setTheme(e.target.value);
    }
    useEffect(() => {
        if (themeRef.current) {
            // @ts-ignore
            const obj = new SelectionText(themeRef.current);
        }
    }, []);

    switch (procedureName) {
        case 'eth_maxBlockAggReduce':
            setCode(ethMaxBlockAggReduceExample);
            break;
        case 'avax_maxBlockAggReduce':
            setCode(avaxMaxBlockAggReduceExample);
            break;
        case 'near_maxBlockAggReduce':
            setCode(nearMaxBlockAggReduceExample);
            break;
        case 'btcMaxBlockAggReduceExample':
            setCode(btcMaxBlockAggReduceExample);
            break;
        default:
            break;
    }
    return (
        <div>
            <MonacoEditor
                height="500px"
                width="800px"
                language={"json"}
                theme={'vs-dark'}
                onChange={(event) => onChange(event)}
                value={ethMaxBlockAggReduceExample}
                options={{
                    wordWrap: "on"
                }}
            />
        </div>
    );
}
