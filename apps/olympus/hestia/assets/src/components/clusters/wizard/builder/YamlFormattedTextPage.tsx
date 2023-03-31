import React, {useEffect, useRef} from "react";
import {SelectionText} from "@uiw/react-textarea-code-editor";
import MonacoEditor from "react-monaco-editor/lib/editor";

window.MonacoEnvironment = { getWorkerUrl: () => proxy };

let proxy = URL.createObjectURL(new Blob([`
	self.MonacoEnvironment = {
		baseUrl: 'https://unpkg.com/monaco-editor@latest/min/'
	};
	importScripts('https://unpkg.com/monaco-editor@latest/min/vs/base/worker/workerMain.js');
`], { type: 'text/javascript' }));

export const languageData = [
    'abap', 'aes', 'apex', 'azcli', 'bat', 'c', 'cameligo', 'clojure', 'coffeescript', 'cpp', 'csharp', 'csp', 'css', 'dart', 'dockerfile', 'fsharp', 'go', 'graphql', 'handlebars', 'hcl', 'html', 'ini', 'java', 'javascript', 'json', 'julia', 'kotlin', 'less', 'lex', 'lua', 'markdown', 'mips', 'msdax', 'mysql', 'objective', 'pascal', 'pascaligo', 'perl', 'pgsql', 'php', 'plaintext', 'postiats', 'powerquery', 'powershell', 'pug', 'python', 'r', 'razor', 'redis', 'redshift', 'restructuredtext', 'ruby', 'rust', 'sb', 'scala', 'scheme', 'scss', 'shell', 'sol', 'sql', 'st', 'swift', 'systemverilog', 'tcl', 'twig', 'typescript', 'vb', 'verilog', 'xml', 'yaml'
];

export default function YamlTextField(props: any) {
    const [code, setCode] = React.useState('');
    const onChange = (textInput: string) => {
        setCode(textInput);
    }
    const themeRef = useRef<string>()
    useEffect(() => {
        if (themeRef.current) {
            // @ts-ignore
            const obj = new SelectionText(themeRef.current);
        }
    }, []);
    return (
            <div>
                <MonacoEditor
                    height="800px"
                    width="1000px"
                    language="yaml"
                    theme={'vs-dark'}
                    onChange={(event) => onChange(event)}
                    value={code}
                />
            </div>
    );
}
